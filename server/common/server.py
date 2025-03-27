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
        self.client_sockets = []

    def run(self):
        self.running = True
        signal.signal(signal.SIGTERM, self.shutdown)

        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                if client_sock:
                    self.client_sockets.append(client_sock)
                    self.__handle_client_connection(client_sock)
            except:
                if not self.running:
                    break
    
    def shutdown(self, signum=None, frame=None):
        self.running = False
        logging.info(f'action: shutdown | result: in_progress')
        if self._server_socket:
            self._server_socket.close()
        for client_socket in self.client_sockets:
            client_socket.close()
        logging.info(f'action: shutdown | result: success')

    def __handle_client_connection(self, client_sock):
        try:
            bet_count = 0
            while True:
                bets, finish = utils.decode_bets(client_sock, bet_count)
                if finish:
                    utils.acknowledge_bets(client_sock, bet_count)
                    logging.info(f'action: apuesta_recibida | result: success | cantidad: {bet_count}')
                    break
                
                utils.store_bets(bets)
                bet_count += len(bets)
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            self.client_sockets.remove(client_sock)

    def __accept_new_connection(self):
        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
