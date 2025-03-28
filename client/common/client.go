package common

import (
	"encoding/csv"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	file   *os.File
	reader *csv.Reader
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, file *os.File, reader *csv.Reader) *Client {
	client := &Client{
		config: config,
		file:   file,
		reader: reader,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(sigChan chan os.Signal) {
	// Create the connection the server in every loop iteration. Send an
	c.createClientSocket()

	betCount := 0

	log.Infof("action: comenzar_envio | result: in_progress")
	err := sendStartBets(c.conn)
	if err != nil {
		log.Errorf("action: comenzar_envio | result: fail | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	log.Infof("action: comenzar_envio | result: success")

	for {
		batch, err := getAgencyData(c.reader, c.config.ID, c.config.BatchMaxAmount)
		if err != nil {
			return
		}

		if len(batch) == 0 {
			break
		}

		err = sendBetBatch(c.conn, batch)
		if err != nil {
			return
		}

		betCount += len(batch)

		select {
		case <-sigChan:
			log.Infof("action: shutdown | result: success")
			c.conn.Close()
			c.file.Close()
			return
		default:
		}
	}

	c.file.Close()

	log.Infof("action: apuesta_enviada | result: success | cantidad: %v", betCount)
	log.Infof("action: finalizar_envio | result: in_progress")

	err = sendFinishMessage(c.conn)
	if err != nil {
		log.Errorf("action: finalizar_envio | result: fail | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	msg, err := receiveMessage(c.conn)
	if err != nil {
		log.Errorf("action: finalizar_envio | result: fail | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	responseCount, _ := strconv.Atoi(strings.TrimSpace(msg))

	if responseCount != betCount {
		log.Errorf("action: finalizar_envio | result: fail | msg: %v | error: response does not match",
			msg,
		)
		return
	}

	log.Infof("action: finalizar_envio | result: success")
	c.conn.Close()

	sleepMultiplier := 0
	for {
		select {
		case <-sigChan:
			log.Infof("action: shutdown | result: success")
			return
		default:
			log.Infof("action: consulta_ganadores | result: in_progress | attempt: %v", sleepMultiplier)
			c.createClientSocket()

			results, err := sendAskForResults(c.conn, c.config.ID)
			if err == nil {
				log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(results))
				c.conn.Close()
				time.Sleep(100 * time.Millisecond)
				return
			} else {
				c.conn.Close()
				sleepMultiplier++
				time.Sleep(time.Duration(sleepMultiplier*1000) * time.Millisecond)
			}
		}
	}
}
