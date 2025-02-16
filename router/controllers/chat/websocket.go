package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"bhbdev/jam/lib/page"

	"github.com/olahol/melody"
	"github.com/ollama/ollama/api"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type WebSocketMessage struct {
	UserMessage string `json:"user-message"`
}

type ChatResponse struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message,omitempty"`
	Done      bool    `json:"done"`
}

type Message struct {
	ID      string        `json:"id"`
	Role    string        `json:"role"` // system, user, assistant
	Content template.HTML `json:"content"`
	Append  bool
}

var History []api.Message

type WebSocketHandler struct {
	WebSocket *melody.Melody
	*api.Client
}

func NewWebSocketHandler() *WebSocketHandler {
	ollamaClient, err := api.ClientFromEnvironment()
	if err != nil {
		slog.Error("Failed to create ollama client", "error", err)
		panic(err)
	}

	handler := &WebSocketHandler{
		melody.New(),
		ollamaClient,
	}
	handler.ApplyWebSocketHandlers()
	return handler
}

func (h *WebSocketHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	slog.Info("Upgrading to websocket connection")
	h.WebSocket.HandleRequest(w, r)
}

func (h *WebSocketHandler) ApplyWebSocketHandlers() {

	h.WebSocket.HandleConnect(func(s *melody.Session) {
		slog.Info("New websocket connection")
		ctx := s.Request.Context()
		// assistant greeting

		stream := false
		greetChatRequest := &api.ChatRequest{
			Model: "deepseek-r1:1.5b",
			Messages: []api.Message{
				{
					Role:    "user",
					Content: "Hello. You are a friendly assistant.",
				},
			},
			Stream: &stream,
		}
		err := h.Client.Chat(ctx, greetChatRequest, func(botChatResponse api.ChatResponse) error {

			if botChatResponse.Done {

				slog.Info("Chat greeting done", "response", botChatResponse)

				History = append(History, botChatResponse.Message)

				var botMessageWSResponse bytes.Buffer
				p := page.New(ctx)
				p.Data["Chat"] = ChatResponse{
					CreatedAt: time.Now().Format(time.TimeOnly),
					Message: Message{
						ID:      "anonymous",
						Role:    "assistant",
						Content: template.HTML(filterHTML(botChatResponse)),
					},
				}

				p.Partial(&botMessageWSResponse, "chat-list")
				slog.Info("Broadcasting message to websocket connections", "message", botMessageWSResponse.String())
				h.WebSocket.Broadcast(botMessageWSResponse.Bytes())

			}

			return nil
		})

		if err != nil {
			slog.Error("Failure requesting ollama chat api", "error", err)
		}
	})

	h.WebSocket.HandleMessage(func(s *melody.Session, msg []byte) {

		ctx := s.Request.Context()
		slog.Info("Received message from websocket connection", "message", string(msg))

		var wsMessage WebSocketMessage
		if err := json.Unmarshal(msg, &wsMessage); err != nil {
			slog.Error("Failed to unmarshal message", "error", err)
		}

		userMessage := api.Message{
			Role:    "user",
			Content: wsMessage.UserMessage,
		}

		var userMessageWSResponse bytes.Buffer
		p := page.New(ctx)
		p.Data["Chat"] = ChatResponse{
			CreatedAt: time.Now().Format(time.TimeOnly),
			Message: Message{
				ID:      "anonymous",
				Role:    "user",
				Content: template.HTML(wsMessage.UserMessage),
			},
		}

		p.Partial(&userMessageWSResponse, "chat-list")
		slog.Info("Broadcasting message to websocket connections", "message", userMessageWSResponse.String())
		h.WebSocket.Broadcast(userMessageWSResponse.Bytes())

		History = append(History, userMessage)

		// assistant chat request
		stream := false
		botChatRequest := &api.ChatRequest{
			Model:    "deepseek-r1:1.5b",
			Messages: History,
			Stream:   &stream,
		}

		slog.Info("Requesting ollama chat", "request", botChatRequest)

		err := h.Client.Chat(ctx, botChatRequest, func(botChatResponse api.ChatResponse) error {

			if botChatResponse.Done {

				slog.Info("Chat request done", "response", botChatResponse)

				History = append(History, botChatResponse.Message)

				var botMessageWSResponse bytes.Buffer
				p := page.New(ctx)
				p.Data["Chat"] = ChatResponse{
					CreatedAt: time.Now().Format(time.TimeOnly),
					Message: Message{
						ID:      "anonymous",
						Role:    "assistant",
						Content: template.HTML(filterHTML(botChatResponse)),
					},
				}

				p.Partial(&botMessageWSResponse, "chat-list")
				slog.Info("Broadcasting message to websocket connections", "message", botMessageWSResponse.String())
				h.WebSocket.Broadcast(botMessageWSResponse.Bytes())

			}

			return nil
		})

		if err != nil {
			slog.Error("Failure requesting ollama chat api", "error", err)
		}
	})
}

func filterHTML(response api.ChatResponse) string {
	// remove everything up to the end of the think tag
	content := regexp.MustCompile(`(?s).*?<\/think>`).ReplaceAllString(response.Message.Content, "")
	slog.Info(fmt.Sprintf("Content: %s", content))
	unsafe := blackfriday.Run([]byte(content))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(html)
}
