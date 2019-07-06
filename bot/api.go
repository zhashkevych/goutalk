package bot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type loginData struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type loginResponse struct {
	Credentials string `json:"credentials"`
}

type sendMessageData struct {
	RoomID  string `json:"room_id"`
	Message string `json:"message"`
}

// login sends POST request to goutalk's /login endpoint
func (c *ChatBot) login() (string, error) {
	data := &loginData{UserName: c.username, Password: c.password}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	loginURL := url.URL{Scheme: "http", Host: c.serverHost, Path: "/login"}
	resp, err := http.Post(loginURL.String(), "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	response := new(loginResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return "", err
	}

	return response.Credentials, nil
}

// sendMessage sends POST request to goutalk's /message endpoint
func (c *ChatBot) sendMessage(roomID, message string) error {
	data := &sendMessageData{RoomID: roomID, Message: message}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	messageURL := url.URL{Scheme: "http", Host: c.serverHost, Path: "/message"}
	req, err := http.NewRequest(http.MethodPost, messageURL.String(), bytes.NewBuffer(dataBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
