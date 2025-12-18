package zoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
)


type ZoomID string

func (id *ZoomID) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case float64:
		*id = ZoomID(fmt.Sprintf("%.0f", val))
	case string:
		*id = ZoomID(val)
	default:
		return fmt.Errorf("cannot unmarshal %T into ZoomID", v)
	}
	return nil
}

type ZoomClient struct {
	auth       *ZoomAuth
	config     *config.ZoomConfig
	httpClient *http.Client
}

type CreateMeetingRequest struct {
	Topic     string    `json:"topic"`
	Type      int       `json:"type"`
	StartTime *string   `json:"start_time,omitempty"`
	Duration  *int      `json:"duration,omitempty"`
	Password  *string   `json:"password,omitempty"`
	Settings  *MeetingSettings `json:"settings,omitempty"`
}

type MeetingSettings struct {
	HostVideo        bool `json:"host_video"`
	ParticipantVideo bool `json:"participant_video"`
	JoinBeforeHost   bool `json:"join_before_host"`
	MuteUponEntry    bool `json:"mute_upon_entry"`
	WaitingRoom      bool `json:"waiting_room"`
}

type CreateMeetingResponse struct {
	ID        ZoomID    `json:"id"`
	Topic     string    `json:"topic"`
	Type      int       `json:"type"`
	StartTime string    `json:"start_time"`
	Duration  int       `json:"duration"`
	Timezone  string    `json:"timezone"`
	Password  string    `json:"password"`
	JoinURL   string    `json:"join_url"`
	StartURL  string    `json:"start_url"`
	Settings  MeetingSettings `json:"settings"`
}

type UpdateMeetingRequest struct {
	Topic     string    `json:"topic,omitempty"`
	Type      int       `json:"type,omitempty"`
	StartTime *string   `json:"start_time,omitempty"`
	Duration  *int      `json:"duration,omitempty"`
	Password  *string   `json:"password,omitempty"`
	Settings  *MeetingSettings `json:"settings,omitempty"`
}

type GetMeetingResponse struct {
	ID        ZoomID    `json:"id"`
	Topic     string    `json:"topic"`
	Type      int       `json:"type"`
	StartTime string    `json:"start_time"`
	Duration  int       `json:"duration"`
	Timezone  string    `json:"timezone"`
	Password  string    `json:"password"`
	JoinURL   string    `json:"join_url"`
	StartURL  string    `json:"start_url"`
	Settings  MeetingSettings `json:"settings"`
}

func NewZoomClient(cfg *config.ZoomConfig) *ZoomClient {
	return &ZoomClient{
		auth:   NewZoomAuth(cfg),
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (zc *ZoomClient) CreateMeeting(userID string, req CreateMeetingRequest) (*CreateMeetingResponse, error) {
	accessToken, err := zc.auth.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/users/%s/meetings", zc.config.BaseURL, userID)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		"Content-Type":  "application/json",
	}

	var response CreateMeetingResponse
	err = makeHTTPRequest("POST", url, headers, json.RawMessage(reqBody), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create meeting: %w", err)
	}

	return &response, nil
}

func (zc *ZoomClient) UpdateMeeting(meetingID string, req UpdateMeetingRequest) error {
	accessToken, err := zc.auth.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/meetings/%s", zc.config.BaseURL, meetingID)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		"Content-Type":  "application/json",
	}

	err = makeHTTPRequest("PATCH", url, headers, json.RawMessage(reqBody), nil)
	if err != nil {
		return fmt.Errorf("failed to update meeting: %w", err)
	}

	return nil
}

func (zc *ZoomClient) GetMeeting(meetingID string) (*GetMeetingResponse, error) {
	accessToken, err := zc.auth.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/meetings/%s", zc.config.BaseURL, meetingID)
	
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	var response GetMeetingResponse
	err = makeHTTPRequest("GET", url, headers, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get meeting: %w", err)
	}

	return &response, nil
}

func (zc *ZoomClient) DeleteMeeting(meetingID string) error {
	accessToken, err := zc.auth.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/meetings/%s", zc.config.BaseURL, meetingID)
	
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	err = makeHTTPRequest("DELETE", url, headers, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete meeting: %w", err)
	}

	return nil
}

func makeHTTPRequest(method, url string, headers map[string]string, body interface{}, response interface{}) error {
	var reqBody io.Reader
	if body != nil {
		if rawMsg, ok := body.(json.RawMessage); ok {
			reqBody = bytes.NewBuffer(rawMsg)
		} else {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return err
			}
			reqBody = bytes.NewBuffer(bodyBytes)
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return err
		}
	}

	return nil
}

