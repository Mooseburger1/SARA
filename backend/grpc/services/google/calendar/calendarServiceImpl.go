package calendar

import (
	"backend/grpc/proto/api/POGO"
	"backend/grpc/proto/api/calendar"
	"backend/grpc/services/google/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	LIST_CALENDAR_LIST_ENDPOINT = "https://www.googleapis.com/calendar/v3/users/me/calendarList"
	LIST_CALENDAR_GET_ENDPOINT = "https://www.googleapis.com/calendar/v3/users/me/calendarList/calendarId"
	GET                         = "GET"
)

// listCalendarList is a package private funciton utilized to make
// an http request to the google Calendar API server, specifically the endpoint
// which returns a list of all calendars for authorized calling client. The
// response is unmarshaled and converted into an CalendarListResponse protobuf
func listCalendarList(rpc *calendar.CalendarListRequest, logger *log.Logger) (*calendar.CalendarListResponse, error) {
	client, err := utils.CreateClient(rpc.GetClientInfo())

	if err != nil {
		logger.Printf("Error creating client: %v", err)
		st := utils.CreateClientCreationError(err)
		return nil, st.Err()
	}

	req, err := http.NewRequest(GET, LIST_CALENDAR_LIST_ENDPOINT, nil)
	if err != nil {
		panic(err)
	}

	query := req.URL.Query()
	if rpc.PageToken != "" {
		query.Add("pageToken", rpc.GetPageToken())
	}
	if rpc.MaxResults != 0 {
		query.Add("maxResults", strconv.Itoa(int(rpc.GetMaxResults())))
	}
	if rpc.ShowDeleted != false {
		query.Add("showDeleted", strconv.FormatBool(rpc.GetShowDeleted()))
	}
	if rpc.ShowHidden != false {
		query.Add("showHidden", strconv.FormatBool(rpc.GetShowHidden()))
	}
	if rpc.SyncToken != "" {
		query.Add("syncToken", rpc.GetSyncToken())
	}

	logger.Printf("query: %v", query.Encode())

	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		logger.Printf("Call to List Calendar List returned status code %v, not %v", resp.StatusCode, http.StatusOK)
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}
		st := utils.CreateErrorResponseError(resp.StatusCode, bodyBytes)
		return nil, st.Err()
	}

	result := calendarListResponseDecoder(resp.Body)

	return listCalendarResponsePogo2Proto(&result), nil
}

// calendarListResponseDecoder takes in the response body from the API request made to the
// server and unmarshals it into a POGO (plain old golang object). This is needed to be done so
// the response can be easily converted to a response protobuf and returned to the caller
func calendarListResponseDecoder(body io.ReadCloser) POGO.CalendarListResponse {
	var result POGO.CalendarListResponse
	json.NewDecoder(body).Decode(&result)
	return result
}

// listCalendarResponsePogo2Proto converts a golang struct containing the data of a ListCalendarResponse
// and converts it to an actual instance of a ListCalendarResponse protobuf
func listCalendarResponsePogo2Proto(result *POGO.CalendarListResponse) *calendar.CalendarListResponse {

	var items []*calendar.CalendarList

	for _, item := range result.Items {
		var reminders []*calendar.Reminders
		for _, reminder := range item.DefaultReminders {
			reminders = append(reminders, &calendar.Reminders{Method: reminder.Method, Minutes: int64(reminder.Minutes)})
		}

		var notifications []*calendar.NotificationSettings_Notifications
		for _, notification := range item.NotificationSettings.Notifications {
			notifications = append(notifications, &calendar.NotificationSettings_Notifications{Type: notification.Type, Method: notification.Method})
		}

		items = append(items, &calendar.CalendarList{Id: item.Id,
			Summary:              item.Summary,
			Description:          item.Description,
			Location:             item.Location,
			Timezone:             item.Timezone,
			ColorId:              item.ColorId,
			BackgroundColor:      item.BackgroundColor,
			ForegroundColor:      item.ForegroundColor,
			Hidden:               item.Hidden,
			Selected:             item.Selected,
			AccessRole:           item.AccessRole,
			DefaultReminders:     reminders,
			NotificationSettings: &calendar.NotificationSettings{Notifications: notifications},
			Primary:              item.Primary,
			Deleted:              item.Deleted})

	}

	return &calendar.CalendarListResponse{NextPageToken: result.NextPageToken, NextSyncToken: result.NextSyncToken, Items: items}

}
