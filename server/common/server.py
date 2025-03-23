import socket
import logging
import signal
from common import utils


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.running = False

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        self.running = True
        signal.signal(signal.SIGTERM, self.shutdown)

        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                if client_sock:
                    self.__handle_client_connection(client_sock)
            except:
                if not self.running:
                    break
    
    def shutdown(self, signum=None, frame=None):
        self.running = False
        logging.info(f'action: shutdown | result: in_progress')
        if self._server_socket:
            self._server_socket.close()
            logging.info(f'action: shutdown | result: success')

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bet_count = 0
            while True:
                bets, _, finish = utils.decode_bets(client_sock, bet_count)
                if finish:
                    utils.acknowledge_bets(client_sock, bet_count)
                    logging.info(f'action: apuesta_recibida | result: success | cantidad: {bet_count}')
                    break
                
                utils.store_bets(bets)
                bet_count += len(bets)
                logging.info(f'action: apuesta_almacenada | result: success | cantidad: {bet_count}')
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
