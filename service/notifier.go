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

// ProcessWebhook is called when we get a new webhook request from bitbucket with a body
func (n *Notifier) ProcessWebhook(b *models.WebhookBody) error {

	switch b.EventKey {
	case "pr:opened": //Notify all current reviewers
		for _, reviewer := range b.PullRequest.Reviewers {
			u, err := n.Slack.GetUserByEmail(reviewer.User.EmailAddress)
			if err != nil {
				logrus.Error(err)
				return nil
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
				logrus.Error(err)
				return nil
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
		return n.prCommentAdded(b)

	}
	return nil
}

func (n *Notifier) prCommentAdded(b *models.WebhookBody) error {
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

		//TODO make this domain configurable
		user, err := n.Slack.GetUserByEmail(mentionedUsername + "@fortnox.se")
		if err != nil {
			return err
		}
		_, _, err = n.Slack.PostMessage(user.ID, b.FormatMessage(fmt.Sprintf("%s", b.Comment.Text), "mentioned you in comment")...)
		if err != nil {
			return err
		}
	}

	if b.Actor.Name == b.PullRequest.Author.User.Name { // dont notify yourself
		logrus.Debug("Skip notifying author is same as actor")
		return nil
	}
	user, err := n.Slack.GetUserByEmail(b.PullRequest.Author.User.EmailAddress)
	if err != nil {
		return err
	}
	_, _, err = n.Slack.PostMessage(user.ID, b.FormatMessage(fmt.Sprintf("%s", b.Comment.Text), "commented")...)
	return err

}

func (n *Notifier) getIMByEmail(email string) (string, error) {
	user, err := n.Slack.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	_, _, channelID, err := n.Slack.OpenIMChannel(user.ID)
	return channelID, err
}
