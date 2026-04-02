package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"wsServer/pkg/team"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
)

type WSMessage struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"` // 新增 GroupID
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	port := "3002"
	fmt.Printf("🚀 Agent Team WebSocket Server starting on ws://localhost:%s/ws\n", port)
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
		
		// 如果没有 GroupID，默认使用 'default'
		if wsMsg.GroupID == "" {
			wsMsg.GroupID = "default"
		}

		log.Printf("[%s@%s] -> Team: %s", wsMsg.UserID, wsMsg.GroupID, wsMsg.Message)

		// 异步处理消息编排
		go handleTeamOrchestration(conn, wsMsg)
	}
}

func handleTeamOrchestration(conn *websocket.Conn, msg WSMessage) {
	// 1. 初始化 Team Context
	ctx, err := team.Manager.InitTeam(msg.GroupID)
	if err != nil {
		sendError(conn, "Failed to initialize team: "+err.Error())
		return
	}

	// 2. Main Agent 决策 (Router)
	// 通知前端正在思考调度
	sendThinking(conn, "Main Agent 正在分析任务并指派专家...", true)
	
	// 这里为了演示流畅性，我们先做一个简单的 mock 决策，或者直接调用 DecideNextAgent
	// 在实际复杂场景中，这里应该是一个真实的 LLM 调用
	// 简单起见，我们先假设如果是自由模式，且没有指定，就直接用 Main Agent 回复（或随机选一个）
	// 或者我们实现一个简化版的 DecideNextAgent，直接返回第一个可用的 Agent，或者根据规则
	
	// 简化版逻辑：如果 Group 有成员，选第一个成员；否则选 Main Agent
	var targetSessionID string
	var targetModelID string
	var targetAgentName string
	
	ctx.Mutex.RLock()
	// 简单策略：如果有成员，选第一个成员回答；否则 Main Agent 回答
	// 真实的 RouterDecision 逻辑可以在 pkg/team/manager.go 中完善
	if len(ctx.MemberSessions) > 0 {
		for uid, sess := range ctx.MemberSessions {
			targetSessionID = sess
			profile := ctx.Profiles[uid]
			targetModelID = profile.ModelID
			targetAgentName = profile.Name
			break // Just pick one for now
		}
		sendThinking(conn, fmt.Sprintf("任务已指派给专家: %s", targetAgentName), false)
	} else {
		targetSessionID = ctx.MainSessionID
		sendThinking(conn, "由 Team Leader 直接处理", false)
	}
	ctx.Mutex.RUnlock()

	// 3. 执行任务 (Worker Agent)
	// 构建 Prompt: 注入其他成员信息，实现“协同感知”
	// 这里我们直接发送用户消息，但在高级模式下，应该封装成 BuildWorkerPrompt
	
	req, err := team.BuildRequest(targetSessionID, targetModelID, msg.Message)
	if err != nil {
		sendError(conn, "Failed to build request: "+err.Error())
		return
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

	// 4. 流式转发结果
	streamResponse(conn, resp.Body)
}

func streamResponse(conn *websocket.Conn, body io.ReadCloser) {
	reader := bufio.NewReader(body)
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

func sendThinking(conn *websocket.Conn, content string, isStreaming bool) {
	// 模拟 OpenCode 的 Thinking 格式发送给前端
	conn.WriteJSON(map[string]interface{}{
		"type": "delta",
		"raw": map[string]interface{}{
			"content":     content,
			"is_thinking": true,
			// isStreaming 为 true 时前端会显示“思考中...”，false 为完成
		},
	})
}
