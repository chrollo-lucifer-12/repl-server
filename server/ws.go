package server

import (
	"log"
	"net/http"

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

func (s *Server) wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	s.l.Info("new connection :", conn.RemoteAddr().String())
	defer conn.Close()
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			s.l.Error("read error:", err)
			return
		}

		s.l.Info("received message:", string(msg), "type:", messageType)
	}
}
