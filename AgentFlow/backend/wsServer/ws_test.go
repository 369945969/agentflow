package main

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketForwarding(t *testing.T) {
	// 1. 连接到正在运行的服务器 (请确保执行 verify_ws.sh 时服务器已启动)
	u := url.URL{Scheme: "ws", Host: "localhost:3002", Path: "/ws"}
	t.Logf("🔌 Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Dial error: %v (Ensure server is running on port 3002)", err)
	}
	defer c.Close()

	// 2. 发送测试请求
	testMsg := map[string]string{
		"userId":  "test_user",
		"message": "你好",
	}
	payload, _ := json.Marshal(testMsg)
	if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// 3. 接收并验证响应
	receivedDelta := false
	receivedDone := false

	// 设置超时，防止死等
	c.SetReadDeadline(time.Now().Add(10 * time.Second))

	for i := 0; i < 20; i++ { // 最多读取 20 条消息
		_, message, err := c.ReadMessage()
		if err != nil {
			t.Logf("Read loop ended: %v", err)
			break
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(message, &resp); err != nil {
			t.Errorf("Invalid JSON received: %s", string(message))
			continue
		}

		msgType, _ := resp["type"].(string)
		t.Logf("Received message type: %s", msgType)

		switch msgType {
		case "delta":
			receivedDelta = true
		case "done":
			receivedDone = true
			goto EndLoop
		case "error":
			t.Fatalf("Server returned error: %v", resp["error"])
		}
	}

EndLoop:
	if !receivedDelta && !receivedDone {
		t.Error("Did not receive any delta or done messages from server")
	}
	
	if receivedDelta {
		t.Log("✅ Successfully received stream delta")
	}
	if receivedDone {
		t.Log("✅ Successfully received done signal")
	}
}
