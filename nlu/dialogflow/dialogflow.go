package dialogflow

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"errors"
	"github.com/zhashkevych/goutalk/booking"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"context"
	"fmt"
	"strconv"
)

const (
	IntentTypeBook   = "restaurant.book"
	IntentTypeList   = "restaurant.bookings.list"
	IntentTypeCancel = "restaurant.bookings.cancel"
	IntentTypeChange = "restaurant.bookings.changedate"

	ResponseCancelationSuccess = "All bookings canceled successfully"
	ResponseChangeSuccess      = "Booking's time changed successfully"
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
	repo          booking.Repository
}

func NewDialogflowProcessor(projectID, lang, jsonPath string, bookingRepo booking.Repository) (*DialogflowProcessor, error) {
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
		repo:          bookingRepo,
	}, nil
}

func (dp *DialogflowProcessor) Process(ctx context.Context, message, chatID, userID string) (string, error) {
	session := fmt.Sprintf("projects/%s/agent/sessions/%s", dp.projectID, chatID)

	textInput := dialogflowpb.TextInput{Text: message, LanguageCode: dp.lang}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: session, QueryInput: &queryInput}

	response, err := dp.sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()

	return dp.processQuery(ctx, queryResult, userID)
	//fulfillmentText := queryResult.GetFulfillmentText()

	//return fulfillmentText, nil
}

func (dp *DialogflowProcessor) processQuery(ctx context.Context, q *dialogflowpb.QueryResult, userID string) (string, error) {
	fulfillmentText := q.GetFulfillmentText()

	switch q.Intent.DisplayName {
	case IntentTypeBook:
		item := &booking.BookItem{
			UserID:        userID,
			Date:          q.Parameters.Fields["date-time"].GetStringValue(),
			Location:      q.Parameters.Fields["location"].GetStringValue(),
			VenueType:     q.Parameters.Fields["venue-type"].GetListValue().GetValues()[0].GetStringValue(),
			VenueTitle:    q.Parameters.Fields["venue-title"].GetStringValue(),
			VenueFacility: q.Parameters.Fields["venue-facility"].GetStringValue(),
		}

		return fulfillmentText, dp.repo.Insert(ctx, item)
	case IntentTypeList:
		items, err := dp.repo.GetByUserID(ctx, userID)
		if err != nil {
			return "", err
		}

		return generateResponseListItems(items), nil
	case IntentTypeCancel:
		err := dp.repo.RemoveByUserID(ctx, userID)
		if err != nil {
			return "", err
		}

		return ResponseCancelationSuccess, nil
	case IntentTypeChange:
		id := q.Parameters.Fields["ID"].GetStringValue()
		date := q.Parameters.Fields["date-time"].GetStringValue()

		err := dp.repo.Update(ctx, id, date)
		if err != nil {
			return "", err
		}

		return ResponseChangeSuccess, nil
	}

	return fulfillmentText, nil
}

func generateResponseListItems(items []*booking.BookItem) string {
	res := "Your bookings: \n"
	for i := range items {
		res += strconv.Itoa(i+1) + ". Booking ID: " + items[i].ID.Hex() + " (Use it cancel booking or change it time and date)\n"
		if items[i].Location != "" {
			res += "\tLocation: " + items[i].Location + "\n"
		}
		if items[i].VenueFacility != "" {
			res += "\tVenue Facility: " + items[i].VenueFacility + "\n"
		}
		if items[i].VenueTitle != "" {
			res += "\tVenue Title: " + items[i].VenueTitle + "\n"
		}
		if items[i].VenueType != "" {
			res += "\tVenue Type: " + items[i].VenueType + "\n"
		}
		if items[i].Date != "" {
			res += "\tDate: " + items[i].Date + "\n"
		}
	}

	return res
}
