package common

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Bet struct {
	AgencyId  string
	Name      string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

func sendMessage(conn net.Conn, message string) error {
	messageBytes := []byte(message)
	messageLenght := len(message)
	if messageLenght > 8192 {
		log.Error("action: send_message | result: fail | error: message exceeds 8kb")
		return errors.New("message exceeds 8kb")
	}

	messageSize := uint16(messageLenght)

	sizeBuffer := make([]byte, 2)
	binary.BigEndian.PutUint16(sizeBuffer, messageSize)

	_, err := conn.Write(sizeBuffer)
	if err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}

	_, err = conn.Write(messageBytes)
	if err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}

	return nil
}

func receiveMessage(conn net.Conn) (string, error) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	msg = strings.TrimSpace(msg)

	return msg, err
}

func sendBet(conn net.Conn, bet Bet) error {
	message := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s\n",
		bet.AgencyId,
		bet.Name,
		bet.LastName,
		bet.Document,
		bet.BirthDate,
		bet.Number,
	)

	return sendMessage(conn, message)
}
