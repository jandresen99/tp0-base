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

	n, err := conn.Write(sizeBuffer)
	if err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}
	if n != len(sizeBuffer) {
		log.Error("action: send_message | result: fail | error: incomplete write")
		return errors.New("incomplete write")
	}

	unsentData := messageBytes
	for len(unsentData) > 0 {
		n, err = conn.Write(unsentData)
		if err != nil {
			log.Errorf("action: send_message | result: fail | error: %v", err)
			return err
		}
		unsentData = unsentData[n:]
	}

	return nil
}

func receiveMessage(conn net.Conn) (string, error) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	msg = strings.TrimSpace(msg)
	if err == nil && msg == "ERROR" {
		log.Error("action: receive_message | result: fail | error: server failed")
		return msg, errors.New("server failed")
	}

	return msg, err
}

func sendBetBatch(conn net.Conn, batch []Bet, betCount int) error {
	bets_str := make([]string, 0, len(batch))

	for _, bet := range batch {
		bet_str := fmt.Sprintf(
			"%s,%s,%s,%s,%s,%s",
			bet.AgencyId,
			bet.Name,
			bet.LastName,
			bet.Document,
			bet.BirthDate,
			bet.Number,
		)
		bets_str = append(bets_str, bet_str)
	}

	message := strings.Join(bets_str, ";")

	return sendMessage(conn, message)
}

func sendFinishMessage(conn net.Conn) error {
	return sendMessage(conn, "FINISH")
}

func sendStartBets(conn net.Conn) error {
	return sendMessage(conn, "BET")
}

func sendAskForResults(conn net.Conn, agencyId string) ([]string, error) {
	message1 := "RESULTS"
	message2 := agencyId

	err := sendMessage(conn, message1)
	if err != nil {
		return nil, err
	}

	err = sendMessage(conn, message2)
	if err != nil {
		return nil, err
	}

	msg, err := receiveMessage(conn)
	if err != nil {
		return nil, err
	}

	if msg == "NOWINNERS" {
		return []string{}, nil
	}

	return strings.Split(msg, ","), nil
}
