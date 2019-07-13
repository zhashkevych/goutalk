package dialogflow

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/goutalk/alert"
	"github.com/zhashkevych/goutalk/booking"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"strconv"
)

const (
	IntentTypeBook     = "restaurant.book"
	IntentTypeList     = "restaurant.bookings.list"
	IntentTypeCancel   = "restaurant.bookings.cancel"
	IntentTypeChange   = "restaurant.bookings.changedate"
	IntentTypeInterval = "restaurant.bookings.setinterval"

	ResponseCancelationOneSuccess = "Booking canceled successfully"
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
	alerter       *alert.Alerter
}

func NewDialogflowProcessor(projectID, lang, jsonPath string, bookingRepo booking.Repository, a *alert.Alerter) (*DialogflowProcessor, error) {
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
		alerter:       a,
	}, nil
}

func (dp *DialogflowProcessor) Process(ctx context.Context, message, roomID, userID string) (string, error) {
	session := fmt.Sprintf("projects/%s/agent/sessions/%s", dp.projectID, roomID)

	textInput := dialogflowpb.TextInput{Text: message, LanguageCode: dp.lang}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: session, QueryInput: &queryInput}

	response, err := dp.sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()

	return dp.processQuery(ctx, queryResult, roomID, userID)
}

func (dp *DialogflowProcessor) processQuery(ctx context.Context, q *dialogflowpb.QueryResult, roomID, userID string) (string, error) {
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

		if err := dp.repo.Insert(ctx, item); err != nil {
			return "Can't make booking on this time, it's already booked.", nil
		}

		//dp.alerter.AddTask(roomID, userID, item.ID.Hex(), timeDiff, time.After(timeDiff))

		return fulfillmentText, nil
	case IntentTypeList:
		items, err := dp.repo.GetByUserID(ctx, userID)
		if err != nil {
			return "", err
		}

		return generateResponseListItems(items), nil
	case IntentTypeCancel:
		var err error
		id := q.Parameters.Fields["ID"].GetStringValue()
		if id != "" {
			err = dp.repo.RemoveByID(ctx, userID, id)
			if err != nil {
				return "", err
			}

			return ResponseCancelationOneSuccess, nil
		}

		err = dp.repo.RemoveByUserID(ctx, userID)
		if err != nil {
			return "", err
		}

		return fulfillmentText, nil
	case IntentTypeChange:
		id := q.Parameters.Fields["ID"].GetStringValue()
		date := q.Parameters.Fields["date-time"].GetStringValue()

		logrus.Printf("change id: %s", id)

		err := dp.repo.Update(ctx, id, date)
		if err != nil {
			return "", err
		}

		return fulfillmentText, nil
	default:
		return fulfillmentText, nil
	}
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
