package dialogflow

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"context"
	"fmt"
	"google.golang.org/api/option"
)

type Response struct {
	Intent     string            `json:"intent"`
	Confidence float32           `json:"confidence"`
	Entities   map[string]string `json:"entities"`
}

type DialogflowProcessor struct {
	projectID     string
	lang          string
	timeZone      string
	sessionClient *dialogflow.SessionsClient
}

func NewDialogflowProcessor(projectID, jsonPath, lang, timezone string) (*DialogflowProcessor, error) {
	sessionClient, err := dialogflow.NewSessionsClient(context.Background(), option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, err
	}

	return &DialogflowProcessor{
		sessionClient: sessionClient,
		projectID:     projectID,
		lang:          lang,
		timeZone:      timezone,
	}, nil
}

func (dp *DialogflowProcessor) Process(ctx context.Context, message, userID, chatID string) error {
	session := fmt.Sprintf("projects/%s/agent/sessions/%s/%s", dp.projectID, chatID, userID)
	request := dialogflowpb.DetectIntentRequest{
		Session: session,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         message,
					LanguageCode: dp.lang,
				},
			},
		},
		QueryParams: &dialogflowpb.QueryParameters{
			TimeZone: dp.timeZone,
		},
	}

	response, err := dp.sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return err
	}

	//queryResult := response.GetQueryResult()
	//if queryResult.Intent != nil {
	//	r.Intent = queryResult.Intent.DisplayName
	//	r.Confidence = float32(queryResult.IntentDetectionConfidence)
	//}
	//
	//r.Entities = make(map[string]string)
	//params := queryResult.Parameters.GetFields()
	//if len(params) > 0 {
	//	for paramName, p := range params {
	//		fmt.Printf("Param %s: %s (%s)", paramName, p.GetStringValue(), p.String())
	//		extractedValue := extractDialogflowEntities(p)
	//		r.Entities[paramName] = extractedValue
	//	}
	//}
	//return
}
