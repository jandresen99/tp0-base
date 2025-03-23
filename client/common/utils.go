package common

import (
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
	messageLenght := len(message)
	if messageLenght > 8192 {
		return errors.New("message exceeds 8kb")
	}

	messageBytes := []byte(message)
	messageSize := uint16(len(messageBytes))

	sizeBuffer := make([]byte, 2)
	binary.BigEndian.PutUint16(sizeBuffer, messageSize)

	_, err := conn.Write(sizeBuffer)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | cantidad: %v | error: %v", betCount, err)
		return err
	}

	_, err = conn.Write(messageBytes)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | cantidad: %v | error: %v", betCount, err)
		return err
	}

	log.Infof("action: apuesta_enviada | result: success | cantidad: %v", betCount)

	return nil
}

func sendFinishMessage(conn net.Conn) error {
	message := "FINISH"
	messageBytes := []byte(message)
	messageSize := uint16(len(messageBytes))

	sizeBuffer := make([]byte, 2)
	binary.BigEndian.PutUint16(sizeBuffer, messageSize)

	_, err := conn.Write(sizeBuffer)
	if err != nil {
		log.Errorf("action: finalizar_envio | result: fail | error: %v", err)
		return err
	}

	_, err = conn.Write(messageBytes)
	if err != nil {
		log.Errorf("action: finalizar_envio | result: fail | error: %v", err)
		return err
	}

	return nil
}
