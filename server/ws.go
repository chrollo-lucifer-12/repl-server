package server

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansi.ReplaceAllString(s, "")
}

func (s *Server) wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	s.l.Info("new connection :", conn.RemoteAddr().String())
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.l.Error("read error:", err)
			return
		}

		var msgData map[string]interface{}
		if err := json.Unmarshal(msg, &msgData); err != nil {
			s.l.Error("invalid JSON message:", err)
			continue
		}

		if msgData["type"] == "execute" {
			cmd, ok := msgData["command"].(string)
			if !ok {
				s.l.Error("invalid command format")
				continue
			}

			s.l.Info("executing command:", cmd)
			res, err := s.t.Run(cmd)
			if err != nil {
				s.l.Error("execution error:", err)
			}

			cleanRes := stripANSI(res)

			response := map[string]interface{}{
				"type":   "response",
				"result": cleanRes,
			}

			formatted, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				s.l.Error("json marshal error:", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, formatted); err != nil {
				s.l.Error("write error:", err)
			}
		}
	}
}
