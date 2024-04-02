package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"body"`
}

type GitHubWebhook struct {
    Issue struct {
        HTMLURL string `json:"html_url"`
    } `json:"issue"`
}

func handler(ctx context.Context, i json.RawMessage) (string, error) {
	
	var event MyEvent
	if err := json.Unmarshal(i, &event); err != nil { 
		return "", err
	}

	var githubAction GitHubWebhook
	if err := json.Unmarshal(i, &githubAction); err != nil { 
		return "", err
	}

	slackMessage := map[string]string{
		"text": fmt.Sprintf("Issue created: %s", githubAction.Issue.HTMLURL),
	}

	slackMessageBytes, err := json.Marshal(slackMessage)
	if err != nil { 
		return "", err
	}

	slackChannelURL := os.Getenv("SLACK_URL")
	if slackChannelURL == "" { 
		return "", fmt.Errorf("SLACK_URL environment variable not set")
	}

	response, err := http.Post(slackChannelURL, "application/json", bytes.NewBuffer(slackMessageBytes))
	if err != nil { 
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK { 
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	responseMessage, err := ioutil.ReadAll(response.Body)
	if err != nil { 
		return "", err
	}


	return string(responseMessage), nil
}

func main() {
	// Start the Lambda function handler
	lambda.Start(handler)
}
