package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	c "github.com/e-gloo/orlog/internal/commands"
	"github.com/e-gloo/orlog/internal/pkg/logging"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func NewServer(ctx context.Context, addr, port string) *http.Server {
	logger := logging.GetLogger()

	// Create http server
	router := http.NewServeMux()
	router.HandleFunc("/connect", MessageHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", addr, port),
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	return srv
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade", "err", err)
		return
	}
	defer conn.Close()
	ch := NewCommandHandler(conn)
	if err := c.SendPacket(conn, c.CreateOrJoin, &c.CreateOrJoinMessage{Welcome: "Welcome to Orlog!"}); err != nil {
		slog.Error("Error sending packet", "err", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if err, ok := err.(*websocket.CloseError); ok {
				ch.HandleRagequit()
			} else {
				slog.Error("read", "err", err)
			}
			break
		}

		packet := &c.Packet{}
		err = json.Unmarshal(message, packet)
		if err != nil {
			slog.Error("Error unmarshalling packet", "err", err)
			return
		}

		slog.Debug("New message", "packet", packet)

		err = ch.Handle(packet)
		if err != nil {
			slog.Error("Error handling packet", "err", err)
			return
		}
	}
}
