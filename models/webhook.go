package models

import (
	"fmt"

	"github.com/nlopes/slack"
)

type WebhookBody struct {
	BitbucketURL string `json:"bitbucketURL"`
	EventKey     string `json:"eventKey"`
	Date         string `json:"date"`
	Actor        User   `json:"actor"`
	PullRequest  struct {
		ID          int    `json:"id"`
		Version     int    `json:"version"`
		Title       string `json:"title"`
		Description string `json:"description"`
		State       string `json:"state"`
		Open        bool   `json:"open"`
		Closed      bool   `json:"closed"`
		CreatedDate int64  `json:"createdDate"`
		UpdatedDate int64  `json:"updatedDate"`
		Locked      bool   `json:"locked"`
		ToRef       Ref    `json:"toRef"`
		FromRef     Ref    `json:"fromRef"`
		Author      struct {
			User     User   `json:"user"`
			Role     string `json:"role"`
			Approved bool   `json:"approved"`
			Status   string `json:"status"`
		} `json:"author"`
		Reviewers []struct {
			User               User   `json:"user"`
			LastReviewedCommit string `json:"lastReviewedCommit"`
			Role               string `json:"role"`
			Approved           bool   `json:"approved"`
			Status             string `json:"status"`
		} `json:"reviewers"`
		Participants []interface{} `json:"participants"`
	} `json:"pullRequest"`
	// Comment is present if eventKey is pr:comment...
	Comment *Comment `json:"comment"`
	// Participant is present if eventKey is pr:reviewer...
	Participant *Participant `json:"participant"`
	// PreviousStatus is present if eventKey is pr:reviewer:approved , pr:reviewer:unapproved, pr:reviewer:needs_work
	PreviousStatus   string `json:"previousStatus"`
	AddedReviewers   []User `json:"addedReviewers"`
	RemovedReviewers []User `json:"removedReviewers"`
}

type Comment struct {
	Properties struct {
		RepositoryID int `json:"repositoryId"`
	} `json:"properties"`
	ID          int           `json:"id"`
	Version     int           `json:"version"`
	Text        string        `json:"text"`
	Author      User          `json:"author"`
	CreatedDate int64         `json:"createdDate"`
	UpdatedDate int64         `json:"updatedDate"`
	Comments    []interface{} `json:"comments"`
	Tasks       []interface{} `json:"tasks"`
}

type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type Participant struct {
	User               User   `json:"user"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
	Role               string `json:"role"`
	Approved           bool   `json:"approved"`
	Status             string `json:"status"`
}

type Ref struct {
	ID           string `json:"id"`
	DisplayID    string `json:"displayId"`
	LatestCommit string `json:"latestCommit"`
	Repository   struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
		} `json:"project"`
		Public bool `json:"public"`
	} `json:"repository"`
}

func (w WebhookBody) ApprovedCount() int {
	approvedCount := 0
	for _, v := range w.PullRequest.Reviewers {
		if v.Approved {
			approvedCount++
		}
	}
	return approvedCount
}

func (w WebhookBody) GetPrURL() string {
	u := fmt.Sprintf("%s/projects/%s/repos/%s/pull-requests/%d/overview",
		w.BitbucketURL,
		w.PullRequest.ToRef.Repository.Project.Key,
		w.PullRequest.ToRef.Repository.Slug,
		w.PullRequest.ID,
	)
	return u
}

// ID can be used to send into functions that needs the uniqu path to the exact pull request.
func (w WebhookBody) ID() (string, string, int) {
	return w.PullRequest.ToRef.Repository.Project.Key,
		w.PullRequest.ToRef.Repository.Slug,
		w.PullRequest.ID
}

func (w WebhookBody) FormatMessage(msg string, action string) []slack.MsgOption {
	color := "#36a64f" // green

	switch w.EventKey {
	case "pr:reviewer:unapproved":
		color = "#F34343"
	case "pr:reviewer:needs_work":
		color = "#B8C043"
	}

	authorName := fmt.Sprintf(
		"%s (%s)\n%s %s",
		w.PullRequest.ToRef.Repository.Project.Key,
		w.PullRequest.ToRef.Repository.Slug,
		w.Actor.DisplayName,
		action,
	)

	title := fmt.Sprintf("PR #%d: %s", w.PullRequest.ID, w.PullRequest.Title)

	attachment := slack.Attachment{
		AuthorName: authorName,
		Color:      color,
		Text:       msg,
		Title:      title,
		TitleLink:  w.GetPrURL(),
	}

	return []slack.MsgOption{
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	}
}
