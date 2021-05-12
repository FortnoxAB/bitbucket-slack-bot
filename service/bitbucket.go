package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fortnoxab/bitbucket-slack-bot/config"
	"github.com/fortnoxab/bitbucket-slack-bot/models"
)

type Bitbucket struct {
	Config *config.Config
}

func NewBitbucket(config *config.Config) *Bitbucket {
	return &Bitbucket{
		Config: config,
	}
}

func (b *Bitbucket) CanMerge(project, repoSlug string, prID int) (models.Merge, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/merge", project, repoSlug, prID)
	merge := models.Merge{}
	err := b.do("GET", u, nil, &merge)
	/*
		if err != nil {
			if strings.Contains(err.Error(), "already been merged") {
				return merge, nil
			}
		}
	*/
	return merge, err
}

type ErrorResponse struct {
	Errors []struct {
		Context       interface{} `json:"context"`
		Message       string      `json:"message"`
		ExceptionName string      `json:"exceptionName"`
	} `json:"errors"`
}

func (e ErrorResponse) Error() string {
	errs := []string{}
	for _, v := range e.Errors {
		errs = append(errs, v.Message)
	}
	return strings.Join(errs, ",")
}

func (b *Bitbucket) do(method, uri string, body io.Reader, response interface{}) error {
	client := &http.Client{}
	u := fmt.Sprintf("%s/%s", b.Config.BitbucketURL, uri)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(b.Config.BitbucketUser, b.Config.BitbucketPassword)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		decoder := json.NewDecoder(resp.Body)
		errorResponse := &ErrorResponse{}
		err := decoder.Decode(errorResponse)
		if err != nil {
			return err
		}
		return errorResponse
	}

	if response != nil {
		decoder := json.NewDecoder(resp.Body)
		return decoder.Decode(response)
	}

	return nil
}
