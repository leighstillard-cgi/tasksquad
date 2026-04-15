package handlers

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHub_RegisterUnregister(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	hub := NewHub(logger)

	go hub.Run()

	// Create mock client
	client := &Client{
		hub:  hub,
		send: make(chan []byte, 256),
	}

	// Register
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client, got %d", hub.ClientCount())
	}

	// Unregister
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("expected 0 clients, got %d", hub.ClientCount())
	}
}

func TestHub_Broadcast(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	hub := NewHub(logger)

	go hub.Run()

	// Create mock clients
	client1 := &Client{hub: hub, send: make(chan []byte, 256)}
	client2 := &Client{hub: hub, send: make(chan []byte, 256)}

	hub.register <- client1
	hub.register <- client2
	time.Sleep(10 * time.Millisecond)

	// Broadcast a message
	message := []byte(`{"test": "data"}`)
	hub.Broadcast(message)
	time.Sleep(10 * time.Millisecond)

	// Check both clients received the message
	select {
	case msg := <-client1.send:
		if string(msg) != string(message) {
			t.Errorf("client1 received wrong message: %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("client1 did not receive message")
	}

	select {
	case msg := <-client2.send:
		if string(msg) != string(message) {
			t.Errorf("client2 received wrong message: %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("client2 did not receive message")
	}
}

func TestHub_ServeWs(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	hub := NewHub(logger)

	go hub.Run()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(hub.ServeWs))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect WebSocket client
	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("expected status 101, got %d", resp.StatusCode)
	}

	// Give hub time to register the client
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client after connection, got %d", hub.ClientCount())
	}
}
