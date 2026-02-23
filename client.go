package bombahead

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/N3moAhead/bombahead-go/internal/network"
)

const (
	msgWelcome            = "welcome"
	msgBackToLobby        = "back_to_lobby"
	msgUpdateLobby        = "update_lobby"
	msgPlayerStatusUpdate = "player_status_update"
	msgServerError        = "error"
	msgClassicInput       = "classic_input"
	msgClassicState       = "classic_state"
	msgGameStart          = "game_start"
)

type welcomePayload struct {
	ClientID string `json:"clientId"`
}

type errorPayload struct {
	Message string `json:"message"`
}

type playerStatusUpdatePayload struct {
	IsReady   bool   `json:"isReady"`
	AuthToken string `json:"authToken,omitempty"`
}

type classicInputPayload struct {
	Move Action `json:"move"`
}

type classicStatePayload struct {
	Players    []Player   `json:"players"`
	Field      fieldWire  `json:"field"`
	Bombs      []Bomb     `json:"bombs"`
	Explosions []Position `json:"explosions"`
}

type fieldWire struct {
	Width  int               `json:"width"`
	Height int               `json:"height"`
	Field  []json.RawMessage `json:"field"`
}

// Run starts the bot and connects to the game server
// It blocks until the connection closes or an unrecoverable error occurs
func Run(userBot Bot) {
	wsURL := os.Getenv("BOMBAHEAD_WS_URL")
	if wsURL == "" {
		wsURL = "ws://localhost:8038/ws"
	}

	token := os.Getenv("BOMBAHEAD_TOKEN")
	if token == "" {
		token = os.Getenv("BOMBERMAN_CLIENT_AUTH_TOKEN")
	}
	if token == "" {
		token = "dev-token-local"
	}

	log.Printf("Connecting to %s...", wsURL)

	client, err := network.Connect(wsURL, token)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	if err := client.Send(msgPlayerStatusUpdate, playerStatusUpdatePayload{IsReady: true, AuthToken: token}); err != nil {
		log.Fatalf("Failed to send initial ready state: %v", err)
	}

	var myID string

	for {
		msg, err := client.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		switch msg.Type {
		case msgWelcome:
			var payload welcomePayload
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				log.Printf("Failed to parse welcome payload: %v", err)
				continue
			}
			myID = payload.ClientID
			log.Printf("Connected as %s", myID)

		case msgUpdateLobby:
			// Intentionally ignored for bot decision logic
			continue

		case msgServerError:
			var payload errorPayload
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				log.Printf("Server error (unparsed payload): %s", string(msg.Payload))
				continue
			}
			log.Printf("Server error: %s", payload.Message)

		case msgGameStart:
			log.Printf("Game started")

		case msgBackToLobby:
			if err := client.Send(msgPlayerStatusUpdate, playerStatusUpdatePayload{IsReady: true}); err != nil {
				log.Printf("Failed to re-ready in lobby: %v", err)
				return
			}

		case msgClassicState:
			state, err := parseClassicState(msg.Payload, myID)
			if err != nil {
				log.Printf("Failed to parse classic state: %v", err)
				continue
			}

			helpers := NewGameHelpers(state)
			action := userBot.GetNextMove(state, helpers)

			if err := client.Send(msgClassicInput, classicInputPayload{Move: action}); err != nil {
				log.Printf("Failed to send action: %v", err)
				return
			}

		default:
			log.Printf("Ignoring message type %q", msg.Type)
		}
	}
}

func parseClassicState(data []byte, myID string) (*GameState, error) {
	var payload classicStatePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal classic state: %w", err)
	}

	cells := make([]CellType, payload.Field.Width*payload.Field.Height)
	for i := 0; i < len(cells) && i < len(payload.Field.Field); i++ {
		cell, err := decodeCell(payload.Field.Field[i])
		if err != nil {
			return nil, fmt.Errorf("decode field cell %d: %w", i, err)
		}
		cells[i] = cell
	}

	state := &GameState{
		Players:    payload.Players,
		Field:      Field{Width: payload.Field.Width, Height: payload.Field.Height, Cells: cells},
		Bombs:      payload.Bombs,
		Explosions: payload.Explosions,
	}

	for _, p := range payload.Players {
		if p.ID == myID {
			player := p
			state.Me = &player
			continue
		}
		state.Opponents = append(state.Opponents, p)
	}

	if state.Me == nil && len(payload.Players) > 0 {
		// Fallback when welcome wasn't received yet.
		player := payload.Players[0]
		state.Me = &player
		if len(payload.Players) > 1 {
			state.Opponents = payload.Players[1:]
		}
	}

	return state, nil
}

func decodeCell(raw json.RawMessage) (CellType, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return CellType(s), nil
	}

	var n int
	if err := json.Unmarshal(raw, &n); err == nil {
		switch n {
		case 0:
			return Wall, nil
		case 1:
			return Air, nil
		case 2:
			return Box, nil
		default:
			return Air, nil
		}
	}

	return Air, fmt.Errorf("unsupported cell encoding: %s", string(raw))
}
