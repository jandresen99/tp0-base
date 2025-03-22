package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Bet struct {
	AgencyId  string
	Name      string
	LastName  string
	Document  string
	BirthDate string
	Number    string
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

	messageBytes := []byte(message)
	messageSize := uint32(len(messageBytes))

	sizeBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBuffer, messageSize)

	_, err := conn.Write(sizeBuffer)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v | error: %v", bet.Document, bet.Number, err)
		return err
	}

	_, err = conn.Write(messageBytes)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v | error: %v", bet.Document, bet.Number, err)
		return err
	}

	return nil
}
