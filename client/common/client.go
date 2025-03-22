package common

import (
	"bufio"
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
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	bet    Bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet Bet) *Client {
	client := &Client{
		config: config,
		bet:    bet,
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
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	select {
	case <-sigChan:
		log.Infof("action: shutdown | result: success")
		return
	default:

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		err := sendBet(c.conn, c.bet)
		if err != nil {
			return
		}

		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		response_data := strings.Split(msg, ",")
		rsp_doc, _ := strconv.Atoi(strings.TrimSpace(response_data[0]))
		rsp_num, _ := strconv.Atoi(strings.TrimSpace(response_data[1]))
		bet_doc, _ := strconv.Atoi(strings.TrimSpace(c.bet.Document))
		bet_num, _ := strconv.Atoi(strings.TrimSpace(c.bet.Number))
		if rsp_doc == bet_doc && rsp_num == bet_num {
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", c.bet.Document, c.bet.Number)
		} else {
			log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v", c.bet.Document, c.bet.Number)
		}
	}
}
