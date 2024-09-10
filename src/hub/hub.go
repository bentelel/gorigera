package hub

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Hub represents a Smart Home Hub
type Hub struct {
	Token            string
	APIBaseURL       string
	WebSocketBaseURL string
	Conn             *websocket.Conn
}

// NewHub creates a new instance of Hub
func NewHub(token, ipAddress, port, apiVersion string) *Hub {
	if port == "" {
		port = "8443"
	}
	if apiVersion == "" {
		apiVersion = "v1"
	}

	return &Hub{
		Token:            token,
		APIBaseURL:       fmt.Sprintf("https://%s:%s/%s", ipAddress, port, apiVersion),
		WebSocketBaseURL: fmt.Sprintf("wss://%s:%s/%s", ipAddress, port, apiVersion),
	}
}

// Headers returns the HTTP headers required for authentication
func (h *Hub) Headers() http.Header {
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	return headers
}

// CreateEventListener sets up a WebSocket connection
func (h *Hub) CreateEventListener() error {
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	conn, _, err := dialer.Dial(h.WebSocketBaseURL, h.Headers())
	if err != nil {
		return err
	}
	h.Conn = conn
	go h.listen()
	return nil
}

// listen handles WebSocket events
func (h *Hub) listen() {
	defer h.Conn.Close()

	for {
		_, message, err := h.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)
	}
}

// StopEventListener stops the WebSocket connection
func (h *Hub) StopEventListener() {
	if h.Conn != nil {
		h.Conn.Close()
		h.Conn = nil
	}
}

// request makes HTTP requests with error handling
func (h *Hub) request(method, route string, data interface{}) ([]byte, error) {
	url := h.APIBaseURL + route

	var jsonData []byte
	var err error
	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header = h.Headers()
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	return io.ReadAll(resp.Body)
}

// Patch sends a PATCH request to the specified route
func (h *Hub) Patch(route string, data interface{}) (string, error) {
	response, err := h.request(http.MethodPatch, route, data)
	return string(response), err
}

// Get sends a GET request to the specified route
func (h *Hub) Get(route string) (map[string]interface{}, error) {
	response, err := h.request(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	return result, err
}

// Post sends a POST request to the specified route
func (h *Hub) Post(route string, data interface{}) (map[string]interface{}, error) {
	response, err := h.request(http.MethodPost, route, data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	return result, err
}

// Delete sends a DELETE request to the specified route
func (h *Hub) Delete(route string, data interface{}) error {
	_, err := h.request(http.MethodDelete, route, data)
	return err
}

// GetDeviceDataByID fetches device data by its ID
func (h *Hub) GetDeviceDataByID(id string) (map[string]interface{}, error) {
	deviceData, err := h.Get("/devices/" + id)
	if err != nil {
		return nil, err
	}
	return deviceData, nil
}

// GetAirPurifiers fetches all air purifiers
func (h *Hub) GetAirPurifiers() ([]map[string]interface{}, error) {
	devices, err := h.Get("/devices")
	if err != nil {
		return nil, err
	}

	airPurifiers := filterDevices(devices, "airPurifier")
	return airPurifiers, nil
}

// filterDevices filters devices by type
func filterDevices(devices map[string]interface{}, deviceType string) []map[string]interface{} {
	var result []map[string]interface{}
	if deviceList, ok := devices["devices"].([]interface{}); ok {
		for _, device := range deviceList {
			deviceMap := device.(map[string]interface{})
			if deviceMap["type"] == deviceType {
				result = append(result, deviceMap)
			}
		}
	}
	return result
}
