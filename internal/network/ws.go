package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// Message represents the old bombahead websocket envelope.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// WsClient handles the websocket connection.
type WsClient struct {
	conn *websocket.Conn
}

// Connect establishes a connection to the game server.
func Connect(serverURL string, token string) (*WsClient, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	header := http.Header{}
	if token != "" {
		header.Add("Authorization", "Bearer "+token)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &WsClient{conn: c}, nil
}

// ReadMessage reads one envelope message from the websocket.
func (w *WsClient) ReadMessage() (*Message, error) {
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var m Message
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, fmt.Errorf("invalid message envelope: %w", err)
	}

	return &m, nil
}

// Send sends a typed envelope message.
func (w *WsClient) Send(msgType string, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	msg := Message{
		Type:    msgType,
		Payload: payloadBytes,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := w.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

// Close closes the websocket connection.
func (w *WsClient) Close() error {
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}
