package network

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestConnect_InvalidURL(t *testing.T) {
	t.Parallel()

	if _, err := Connect("://bad-url", "token"); err == nil {
		t.Fatal("Connect() expected invalid URL error, got nil")
	}
}

func TestWsClient_SendAndReadMessage(t *testing.T) {
	t.Parallel()

	authHeader := make(chan string, 1)
	received := make(chan Message, 1)

	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader <- r.Header.Get("Authorization")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("server read failed: %v", err)
			return
		}
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Errorf("server unmarshal failed: %v", err)
			return
		}
		received <- msg

		replyPayload, _ := json.Marshal(map[string]string{"clientId": "abc"})
		reply := Message{Type: "welcome", Payload: replyPayload}
		replyData, _ := json.Marshal(reply)
		if err := conn.WriteMessage(websocket.TextMessage, replyData); err != nil {
			t.Errorf("server write failed: %v", err)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client, err := Connect(wsURL, "secret")
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer client.Close()

	if err := client.Send("test_type", map[string]string{"k": "v"}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	gotAuth := <-authHeader
	if gotAuth != "Bearer secret" {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, "Bearer secret")
	}

	gotSent := <-received
	if gotSent.Type != "test_type" {
		t.Fatalf("sent type = %q, want %q", gotSent.Type, "test_type")
	}
	var sentPayload map[string]string
	if err := json.Unmarshal(gotSent.Payload, &sentPayload); err != nil {
		t.Fatalf("sent payload unmarshal error = %v", err)
	}
	if sentPayload["k"] != "v" {
		t.Fatalf("sent payload = %+v, want map[k:v]", sentPayload)
	}

	reply, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage() error = %v", err)
	}
	if reply.Type != "welcome" {
		t.Fatalf("reply type = %q, want %q", reply.Type, "welcome")
	}
}

func TestWsClient_CloseNilConnection(t *testing.T) {
	t.Parallel()

	var client WsClient
	if err := client.Close(); err != nil {
		t.Fatalf("Close() on nil conn should be nil, got %v", err)
	}
}
