package client

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/e-gloo/orlog/internal/commands"
	og "github.com/e-gloo/orlog/internal/orlog"
	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	ioh  IOHandler
	conn *websocket.Conn
	game *ClientGame
}

func NewCommandHandler(ioh IOHandler, conn *websocket.Conn) *CommandHandler {
	return &CommandHandler{
		ioh:  ioh,
		conn: conn,
	}
}

func (ch *CommandHandler) Handle(conn *websocket.Conn, packet *commands.Packet) error {
	switch packet.Command {
	case commands.CreateOrJoin:
		return ch.handleCreateOrJoin()
	case commands.AddPlayer:
		return ch.handleAddPlayer(conn)
	case commands.GameStarting:
		return ch.handleGameStarting(packet)
	case commands.CommandOK:
		if packet.Data != "" {
			ch.ioh.DisplayMessage(fmt.Sprintf("%s\n", packet.Data))
		}
		return nil
	case commands.CommandError:
		slog.Debug("Oops désolé :D")
		return nil
	default:
		return ch.handleDefaultCase(packet.Command)
	}
}

func (ch *CommandHandler) handleCreateOrJoin() error {
	gameUuid := ""
	ch.ioh.DisplayMessage("Enter the game UUID (empty for new): ")
	err := ch.ioh.ReadInput(&gameUuid)
	if err != nil {
		return err
	}

	var command commands.Command
	if gameUuid != "" {
		command = commands.JoinGame
	} else {
		command = commands.CreateGame
	}

	err = commands.SendPacket(ch.conn, &commands.Packet{
		Command: command,
		Data:    gameUuid,
	})
	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}
	return nil
}

func (ch *CommandHandler) handleAddPlayer(conn *websocket.Conn) error {
	username := "Player"
	ch.ioh.DisplayMessage("Enter your name : ")

	err := ch.ioh.ReadInput(&username)
	if err != nil {
		return err
	}

	err = commands.SendPacket(conn, &commands.Packet{
		Command: commands.AddPlayer,
		Data:    username,
	})
	if err != nil {
		return fmt.Errorf("error sending packet: %w", err)
	}

	return nil
}

func (ch *CommandHandler) handleGameStarting(packet *commands.Packet) error {
	var game *og.Game

	if err := json.Unmarshal([]byte(packet.Data), game); err != nil {
		return fmt.Errorf("error unmarshalling game: %w", err)
	}

	ch.game = NewClientGame(game)

	return nil
}

func (ch *CommandHandler) handleDefaultCase(command commands.Command) error {
	slog.Debug("Unknown command", "command", command)

	return nil
}
