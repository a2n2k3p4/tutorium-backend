package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Minimal login request/response structs
type LoginRequest struct {
	Action   string `json:"Action"`
	ID       string `json:"ID"`
	Password string `json:"Password"`
}

type LoginResponse struct {
	Status     string `json:"status"`
	ID         string `json:"id"`
	Token      string `json:"token"`
	ImgProfile string `json:"img_profile"`
	NameTH     string `json:"name_th"`
	SurnameTH  string `json:"surname_th"`
	NameEN     string `json:"name_en"`
	ErrDesc    string `json:"err_desc,omitempty"`
}

// NisitKUClient handles API communication
type NisitKUClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewNisitKUClient creates a new minimal client
func NewNisitKUClient(baseURL string) *NisitKUClient {
	return &NisitKUClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Login authenticates a user and returns profile info
func (c *NisitKUClient) Login(username, password string) (*LoginResponse, error) {
	request := LoginRequest{
		Action:   "login",
		ID:       username,
		Password: password,
	}

	var response LoginResponse
	var responses []LoginResponse

	err := c.doRequest(c.BaseURL, request, &response)
	if err != nil {
		// Try as array
		err = c.doRequest(c.BaseURL, request, &responses)
		if err != nil {
			return nil, err
		}
		if len(responses) == 0 {
			return nil, fmt.Errorf("empty response")
		}
		response = responses[0]
	}

	if response.Status != "true" {
		return nil, fmt.Errorf("login failed: %s", response.ErrDesc)
	}

	return &response, nil
}

// Internal function to send POST request and parse response
func (c *NisitKUClient) doRequest(url string, request interface{}, responsePtr interface{}) error {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "NisitKU-Go-Client/1.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	trimmed := strings.TrimSpace(string(body))
	if len(trimmed) > 0 && trimmed[0] == '[' {
		var tempArray []json.RawMessage
		if err := json.Unmarshal(body, &tempArray); err != nil {
			return fmt.Errorf("error parsing array response: %v", err)
		}
		if len(tempArray) > 0 {
			if err := json.Unmarshal(tempArray[0], responsePtr); err != nil {
				return fmt.Errorf("error parsing first array element: %v", err)
			}
		} else {
			return fmt.Errorf("empty array response")
		}
	} else {
		if err := json.Unmarshal(body, responsePtr); err != nil {
			return fmt.Errorf("error parsing object response: %v", err)
		}
	}

	return nil
}
