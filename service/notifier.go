package service

import (
	"fmt"
	"regexp"

	"github.com/fortnoxab/bitbucket-slack-bot/models"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

type Notifier struct {
	RTM       *slack.RTM
	Slack     *slack.Client
	Bitbucket *Bitbucket
}

func NewNotifier(s *slack.Client, rtm *slack.RTM, bitbucket *Bitbucket) *Notifier {
	return &Notifier{
		RTM:       rtm,
		Slack:     s,
		Bitbucket: bitbucket,
	}
}

var rex = regexp.MustCompile(`@"(.+?)"`)

// ProcessWebhook is called when we get a new webhook request from bitbucket with a body.
func (n *Notifier) ProcessWebhook(b *models.WebhookBody) error {
	switch b.EventKey {
	case "pr:opened": // Notify all current reviewers
		for _, reviewer := range b.PullRequest.Reviewers {
			u, err := n.Slack.GetUserByEmail(reviewer.User.EmailAddress)
			if err != nil {
				logrus.Errorf("error fetching %s from slack: %s", reviewer.User.EmailAddress, err)
				continue
			}
			_, _, err = n.Slack.PostMessage(u.ID, b.FormatMessage("is waiting for your review.", "opened pull request")...)
			if err != nil {
				logrus.Error(err)
				return nil
			}
		}
	case "pr:reviewer:updated": // AddedReviewers AND RemovedReviewers have changed notify added ones
		for _, user := range b.AddedReviewers {
			u, err := n.Slack.GetUserByEmail(user.EmailAddress)
			if err != nil {
				logrus.Errorf("error fetching %s from slack: %s", user.EmailAddress, err)
				continue
			}
			_, _, err = n.Slack.PostMessage(u.ID, b.FormatMessage(fmt.Sprintf("has %d/%d approvals", b.ApprovedCount(), len(b.PullRequest.Reviewers)), "added you as reviewer")...)
			if err != nil {
				logrus.Error(err)
				return nil
			}
		}

	case "pr:reviewer:needs_work":
		user, err := n.Slack.GetUserByEmail(b.PullRequest.Author.User.EmailAddress)
		if err != nil {
			return err
		}
		_, _, err = n.Slack.PostMessage(user.ID, b.FormatMessage(fmt.Sprintf("has %d/%d approvals", b.ApprovedCount(), len(b.PullRequest.Reviewers)), "said needs work")...)
		return err
	// message to AUTHOR if it can be merged
	case "pr:reviewer:approved":
		merge, err := n.Bitbucket.CanMerge(b.ID())
		if err != nil {
			logrus.Error(err)
		}
		if !merge.CanMerge && err == nil {
			return nil
		}
		user, err := n.Slack.GetUserByEmail(b.PullRequest.Author.User.EmailAddress)
		if err != nil {
			return err
		}
		if merge.CanMerge {
			_, _, err = n.Slack.PostMessage(user.ID, b.FormatMessage(fmt.Sprintf("has %d/%d approvals and can be merged now.", b.ApprovedCount(), len(b.PullRequest.Reviewers)), "approved")...)
		} else {
			_, _, err = n.Slack.PostMessage(user.ID, b.FormatMessage(fmt.Sprintf("has %d/%d approvals", b.ApprovedCount(), len(b.PullRequest.Reviewers)), "approved")...)
		}
		return err
	case "pr:comment:added":
		return prCommentAdded(n.Slack.GetUserByEmail, n.Slack.PostMessage, n.Bitbucket.GetPrActivity, b)
	}
	return nil
}

func prCommentAdded(
	slackGetUserByEmail func(email string) (*slack.User, error),
	slackPostMessage func(channelID string, options ...slack.MsgOption) (string, string, error),
	bitbucketGetPrActivity func(project string, repoSlug string, prID int) (models.Activities, error),
	b *models.WebhookBody,
) error {
	// If mention. Also notify the person mentioned
	matches := rex.FindAllStringSubmatch(b.Comment.Text, -1)
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		mentionedUsername := match[1]

		if mentionedUsername == b.Actor.Name { // Skip notifying yourself
			logrus.Debug("Skip notifying yourself")
			continue
		}
		if mentionedUsername == b.PullRequest.Author.User.Name { // Skip notifying author (notified below anyways)
			logrus.Debug("Skip notifying author")
			continue
		}

		// TODO make this domain configurable
		user, err := slackGetUserByEmail(mentionedUsername + "@fortnox.se")
		if err != nil {
			return err
		}
		_, _, err = slackPostMessage(user.ID, b.FormatMessage(b.Comment.Text, "mentioned you in comment")...)
		if err != nil {
			return err
		}
	}

	// Notify PR author
	if b.Actor.Name == b.PullRequest.Author.User.Name { // dont notify yourself
		logrus.Debug("Skip notifying author is same as actor")
		return nil
	}
	user, err := slackGetUserByEmail(b.PullRequest.Author.User.EmailAddress)
	if err != nil {
		return err
	}
	_, _, err = slackPostMessage(user.ID, b.FormatMessage(b.Comment.Text, "commented")...)
	if err != nil {
		return err
	}

	// Notify comment thread
	threads := []models.Comment{}
	activities, err := bitbucketGetPrActivity(b.ID())
	if err != nil {
		return err
	}
	for _, v := range activities.Values {
		if v.Action == "COMMENTED" {
			threads = append(threads, v.Comment)
		}
	}
	logrus.Debug("Number of comment threads: ", len(threads))

	for _, thread := range threads {
		authorEmails := map[string]bool{}
		found := false
		traverseThread(&thread, authorEmails, b.Comment.ID, &found)
		if found {
			delete(authorEmails, b.Comment.Author.EmailAddress)  // never notify the comment author
			delete(authorEmails, b.PullRequest.Author.User.Name) // never notify the PR author twice (already done before)
			logrus.Debug("Authors: " + fmt.Sprint(len(authorEmails)))

			for a := range authorEmails {
				user, err := slackGetUserByEmail(a)
				if err != nil {
					return err
				}

				commentUrl := fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d/overview?commentId=%d",
					b.BitbucketURL,
					b.PullRequest.ToRef.Repository.Project.Key,
					b.PullRequest.ToRef.Repository.Slug,
					b.PullRequest.ID,
					b.Comment.ID,
				)
				message := fmt.Sprintf("\n👉 %s\n\n%s", commentUrl, b.Comment.Text)
				_, _, err = slackPostMessage(user.ID, b.FormatMessage(message, "commented on thread")...)
				if err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}

func traverseThread(c *models.Comment, authors map[string]bool, commentId int, found *bool) {
	logrus.Debug("Comment by '" + c.Author.DisplayName + "' | ID: " + fmt.Sprint(c.ID))
	authors[c.Author.EmailAddress] = true
	if c.ID == commentId {
		*found = true
	}

	for _, child := range c.Comments {
		traverseThread(child, authors, commentId, found)
	}
}
