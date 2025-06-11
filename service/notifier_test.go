package service

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/fortnoxab/bitbucket-slack-bot/models"
	"github.com/nlopes/slack"
)

type Author struct {
	email   string
	slackId string
}

var FallbackUser = Author{
	email:   "fallback.user@example.com",
	slackId: "slack-fallback-user",
}
var PrAuthor = Author{
	email:   "pr.author@example.com",
	slackId: "slack-id-author",
}
var PreviousThreadCommenter1 = Author{
	email:   "thread.commenter1@example.com",
	slackId: "slack-id-thread-commenter1",
}
var PreviousThreadCommenter2 = Author{
	email:   "thread.commenter2@example.com",
	slackId: "slack-id-thread-commenter2",
}

type MockSlackGerUserByEmail struct {
	CalledWith                     []string
	ReturnFallbackUser             *slack.User
	ReturnErr                      error
	ReturnPrAuthor                 *slack.User
	ReturnPreviousThreadCommenter1 *slack.User
	ReturnPreviousThreadCommenter2 *slack.User
	ReturnLastCommenter            *slack.User
}

func (m *MockSlackGerUserByEmail) call(email string) (*slack.User, error) {
	m.CalledWith = append(m.CalledWith, email)
	switch email {
	case PrAuthor.email:
		return m.ReturnPrAuthor, m.ReturnErr
	case PreviousThreadCommenter1.email:
		return m.ReturnPreviousThreadCommenter1, m.ReturnErr
	case PreviousThreadCommenter2.email:
		return m.ReturnPreviousThreadCommenter2, m.ReturnErr
	default:
		return m.ReturnFallbackUser, m.ReturnErr
	}
}

type MockSlackPostMessage struct {
	CalledWith []struct {
		ChannelID string
		Options   []slack.MsgOption
	}
	ReturnTimestamp string
	ReturnMsgID     string
	ReturnErr       error
}

func (m *MockSlackPostMessage) call(channelID string, options ...slack.MsgOption) (string, string, error) {
	m.CalledWith = append(m.CalledWith, struct {
		ChannelID string
		Options   []slack.MsgOption
	}{channelID, options})

	return m.ReturnTimestamp, m.ReturnMsgID, m.ReturnErr
}

type MockBitbucketGetPrActivity struct {
	CalledWith []struct {
		project  string
		repoSlug string
		prId     int
	}
	ReturnActivities models.Activities
	ReturnErr        error
}

func (m *MockBitbucketGetPrActivity) call(project string, repoSlug string, prId int) (models.Activities, error) {
	m.CalledWith = append(m.CalledWith, struct {
		project  string
		repoSlug string
		prId     int
	}{project, repoSlug, prId})
	return m.ReturnActivities, m.ReturnErr
}

/*
This test verifies the logic for notifying bitbucket PR thread participants through slack
The mock webhooks and pr_activities responses contain the following case:

	[thread] comment 1 (author: PrAuthor)
		[thread] comment 2 (author: ThreadCommenter1)
			[thread] comment 3 (author: ThreadCommenter2)
				[thread] comment 4 (author: ThreadCommenter3) <- new comment that triggers the webhook

3 users must be notified once (not ThreadCommenter3, since he posted the comment):
  - PrAuthor
  - ThreadCommenter1
  - ThreadCommenter2
*/
func TestPrCommentAdded(t *testing.T) {
	expect := func(msg string, received any, expected any) {
		if expected != received {
			t.Errorf("❌ %s [actual/expect] \"%s\" / \"%s\"", msg, received, expected)
		} else {
			// fmt.Printf("✅ %s\n", msg)
		}
	}

	prActivitiesJSON, err1 := os.ReadFile("../testdata/pr_activities_mock.json")
	webhookBodyJSON, err2 := os.ReadFile("../testdata/webhook_body_mock.json")
	if err1 != nil || err2 != nil {
		t.Fatalf("failed to read fixtures")
	}
	var prActivitiesMock models.Activities
	var webhookBodyMock models.WebhookBody
	if err := json.Unmarshal(prActivitiesJSON, &prActivitiesMock); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if err := json.Unmarshal(webhookBodyJSON, &webhookBodyMock); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	mockSlackGerUserByEmail := &MockSlackGerUserByEmail{
		ReturnFallbackUser: &slack.User{
			ID: FallbackUser.slackId,
		},
		ReturnPrAuthor: &slack.User{
			ID: PrAuthor.slackId,
		},
		ReturnPreviousThreadCommenter1: &slack.User{
			ID: PreviousThreadCommenter1.slackId,
		},
		ReturnPreviousThreadCommenter2: &slack.User{
			ID: PreviousThreadCommenter2.slackId,
		},
		ReturnErr: nil,
	}

	mockSlackPostMessage := &MockSlackPostMessage{
		ReturnTimestamp: "0",
		ReturnMsgID:     "1",
		ReturnErr:       nil,
	}

	mockBitbucketGetPrActivity := &MockBitbucketGetPrActivity{
		ReturnActivities: prActivitiesMock,
		ReturnErr:        nil,
	}

	// should return different users based on their email
	f, _ := mockSlackGerUserByEmail.call(FallbackUser.email)
	expect("should load fallback user ", FallbackUser.slackId, f.ID)
	p, _ := mockSlackGerUserByEmail.call(PrAuthor.email)
	expect("should load PR author user ", PrAuthor.slackId, p.ID)
	t1, _ := mockSlackGerUserByEmail.call(PreviousThreadCommenter1.email)
	expect("should load ThreadCommenter1 user ", PreviousThreadCommenter1.slackId, t1.ID)
	t2, _ := mockSlackGerUserByEmail.call(PreviousThreadCommenter2.email)
	expect("should load ThreadCommenter2 user ", PreviousThreadCommenter2.slackId, t2.ID)

	prCommentAdded(
		mockSlackGerUserByEmail.call,
		mockSlackPostMessage.call,
		mockBitbucketGetPrActivity.call,
		&webhookBodyMock)

	expect("should notify PR author", mockSlackPostMessage.CalledWith[0].ChannelID, PrAuthor.slackId)
	expect("should notify previous thread commenter 1", mockSlackPostMessage.CalledWith[1].ChannelID, PreviousThreadCommenter1.slackId)
	expect("should notify previous thread commenter 2", mockSlackPostMessage.CalledWith[2].ChannelID, PreviousThreadCommenter2.slackId)
	expect("should NOT double notify or notify the author of last comment (thread commenter 3)", len(mockSlackPostMessage.CalledWith), 3)
}
