import socket
import logging
import signal
from common import utils


class Server:
    def __init__(self, port, listen_backlog, clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.running = False
        self.clients = clients
        self. client_sockets = []
        self.finished_clients = 0
        self.winners = []

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
        for client_sock in self.client_sockets:
            client_sock.close()
        logging.info(f'action: shutdown | result: success')

    def __handle_client_connection(self, client_sock):
        try:
            message = utils.receive_message(client_sock)
            if message == "BET":
                bet_count = 0
                while True:
                    bets, finish = utils.decode_bets(client_sock, bet_count)
                    if finish:
                        utils.acknowledge_bets(client_sock, bet_count)
                        logging.info(f'action: apuesta_recibida | result: success | cantidad: {bet_count}')
                        self.finished_clients += 1
                        break

                    utils.store_bets(bets)
                    bet_count += len(bets)
            if message == "RESULTS":
                agency_id = utils.receive_message(client_sock)
                if int(self.finished_clients) == int(self.clients):
                    if not self.winners:
                        bets = utils.load_bets()
                        winners = [bet for bet in bets if utils.has_won(bet)]
                        self.winners = winners
                        logging.info(f'action: sorteo | result: success')
                    
                    agency_winners = [bet for bet in self.winners if bet.agency == int(agency_id)]

                    utils.send_results(client_sock, agency_winners)
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
