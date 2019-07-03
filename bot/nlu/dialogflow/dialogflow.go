package dialogflow

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"errors"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"context"
	"fmt"
	"log"
)

type Response struct {
	Intent     string            `json:"intent"`
	Confidence float32           `json:"confidence"`
	Entities   map[string]string `json:"entities"`
}

type DialogflowProcessor struct {
	projectID     string
	lang          string
	sessionClient *dialogflow.SessionsClient
}

func NewDialogflowProcessor(projectID, lang, jsonPath string) (*DialogflowProcessor, error) {
	sessionClient, err := dialogflow.NewSessionsClient(context.Background(), option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, err
	}

	if projectID == "" {
		return nil, errors.New("empty projectID")
	}

	return &DialogflowProcessor{
		sessionClient: sessionClient,
		projectID:     projectID,
		lang:          lang,
	}, nil
}

func (dp *DialogflowProcessor) Process(ctx context.Context, message, userID, chatID string) (string, error) {
	session := fmt.Sprintf("projects/%s/agent/sessions/%s", dp.projectID, chatID+userID)

	textInput := dialogflowpb.TextInput{Text: message, LanguageCode: dp.lang}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: session, QueryInput: &queryInput}

	response, err := dp.sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	log.Printf("recieved response: %+v", response)

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()

	return fulfillmentText, nil
}
