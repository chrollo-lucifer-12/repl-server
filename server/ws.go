package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type wsWriter struct {
	conn *websocket.Conn
}

func (w *wsWriter) Write(p []byte) (int, error) {
	message := map[string]interface{}{
		"type": "output",
		"data": string(p),
	}
	err := w.conn.WriteJSON(message)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

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
		s.l.Error("ws upgrade failed:", err)
		return
	}

	defer conn.Close()

	ctx := context.Background()

	writer := &wsWriter{conn: conn}

	pr, pw := io.Pipe()
	defer pw.Close()

	var replStarted bool

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.l.Error("ws read error:", err)
			return
		}

		var msgData map[string]string
		if err := json.Unmarshal(msg, &msgData); err != nil {
			writer.Write([]byte("invalid JSON\n"))
			continue
		}

		userId := msgData["userId"]

		switch msgData["type"] {

		case "init_project":
			{
				userIdUint, _ := strconv.ParseUint(userId, 10, 64)
				user, err := s.db.FindUser(uint(userIdUint))
				if err != nil || user == nil {
					writer.Write([]byte("user not found"))
					return
				}
				s.d.StartContainer(ctx, writer, msgData["userId"])
			}

		case "react_project":
			{
				if replStarted {
					continue
				}
				replStarted = true
				go func() {
					err := s.d.StartInteractiveRepl(ctx, userId, pr, writer)
					if err != nil {
						writer.Write([]byte(err.Error()))
					}
				}()

				pw.Write([]byte("npm create vite@latest my-app -- --template react\n"))
			}

		case "input":
			if data, ok := msgData["data"]; ok {
				pw.Write([]byte(data))
			}

		case "write_file":
			_ = s.d.WriteFile(
				ctx,
				userId,
				msgData["path"],
				msgData["content"],
				writer,
			)

		case "read_file":
			_ = s.d.ReadFile(
				ctx,
				userId,
				msgData["path"],
				writer,
			)

		case "list_files":
			_ = s.d.ListFiles(
				ctx,
				userId,
				msgData["path"],
				writer,
			)

		case "remove_file":
			_ = s.d.RemoveFile(
				ctx,
				msgData["path"],
				userId,
				writer,
			)

		case "stat_file":
			_ = s.d.StatFile(
				ctx,
				userId,
				msgData["path"],
				writer,
			)

		case "search_file":
			_ = s.d.SearchInFile(
				ctx,
				userId,
				msgData["path"],
				msgData["search"],
				writer,
			)

		case "rename_file":
			_ = s.d.RenameFileDir(
				ctx,
				userId,
				msgData["path"],
				msgData["new_name"],
				writer,
			)
		case "create_dir":
			_ = s.d.CreateDir(ctx, userId, msgData["path"], writer)

		case "resize_terminal":
			_ = s.d.ResizeTerminal(ctx, userId, msgData["rows"], msgData["cols"])

		default:
			writer.Write([]byte("unknown message type\n"))
		}
	}
}
