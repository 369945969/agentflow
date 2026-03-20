package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
	openCodeURL = "http://localhost:3001/api/sse/chat"
)

type WSMessage struct {
	UserID  string `json:"userId"`
	ModelID string `json:"modelId"`
	Message string `json:"message"`
}

type SSEData struct {
	Content string `json:"content"`
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	port := "3002"
	fmt.Printf("🚀 WebSocket Server starting on ws://localhost:%s/ws\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Println("🔌 Client connected")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			sendError(conn, "Invalid JSON format")
			continue
		}

		if wsMsg.Message == "" {
			sendError(conn, "Message is required")
			continue
		}

		log.Printf("[%s] -> OpenCode: %s", wsMsg.UserID, wsMsg.Message)

		// Call OpenCode SSE
		go forwardToOpenCode(conn, wsMsg)
	}
}

func forwardToOpenCode(conn *websocket.Conn, msg WSMessage) {
	payload, _ := json.Marshal(map[string]string{"message": msg.Message})
	req, err := http.NewRequest("POST", openCodeURL, bytes.NewBuffer(payload))
	if err != nil {
		sendError(conn, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if msg.UserID != "" {
		req.Header.Set("X-User-ID", msg.UserID)
	}
	if msg.ModelID != "" {
		req.Header.Set("X-Model-ID", msg.ModelID)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendError(conn, "Failed to connect to OpenCode: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		sendError(conn, fmt.Sprintf("OpenCode error (%d): %s", resp.StatusCode, string(body)))
		return
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Stream read error: %v", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				conn.WriteJSON(map[string]string{"type": "done"})
				break
			}

			// Forward data to client
			var sseData interface{}
			if err := json.Unmarshal([]byte(data), &sseData); err == nil {
				conn.WriteJSON(map[string]interface{}{
					"type": "delta",
					"raw":  sseData,
				})
			} else {
				conn.WriteJSON(map[string]string{
					"type": "raw",
					"data": data,
				})
			}
		}
	}
	log.Println("✅ Stream finished")
}

func sendError(conn *websocket.Conn, errMsg string) {
	conn.WriteJSON(map[string]string{
		"type":  "error",
		"error": errMsg,
	})
}
