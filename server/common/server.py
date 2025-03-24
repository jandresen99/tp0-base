import socket
import logging
import signal
import threading
from common import utils


class Server:
    def __init__(self, port, listen_backlog, clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.running = False
        self.clients = clients
        self.finished_clients = 0
        self.winners = []
        self.lock = threading.Lock()
        self.threads = []

    def run(self):
        self.running = True
        signal.signal(signal.SIGTERM, self.shutdown)

        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                if client_sock:
                    client_thread = threading.Thread(target=self.__handle_client_connection, args=(client_sock,))
                    client_thread.daemon = True
                    client_thread.start()
                    self.threads.append(client_thread)
            except:
                if not self.running:
                    break
        
        self.__wait_for_threads()
    
    def __wait_for_threads(self):
        for thread in self.threads:
            thread.join()
    
    def shutdown(self, signum=None, frame=None):
        logging.info("action: shutdown | result: in_progress")
        self.running = False
        if self._server_socket:
            self._server_socket.close()
        self.__wait_for_threads()
        logging.info("action: shutdown | result: success")

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            message = utils.receive_message(client_sock)
            if message == "BET":
                bet_count = 0
                #logging.info(f"action: apuesta_recibida | result: in_progress | cantidad: {bet_count}")
                while True:
                    bets, _, finish = utils.decode_bets(client_sock, bet_count)
                    if finish:
                        utils.acknowledge_bets(client_sock, bet_count)
                        logging.info(f'action: apuesta_recibida | result: success | cantidad: {bet_count}')
                        with self.lock:
                            self.finished_clients += 1
                        break

                    utils.store_bets(bets)
                    bet_count += len(bets)
                    #logging.info(f'action: apuesta_recibida | result: in_progress | cantidad: {bet_count}')
            if message == "RESULTS":
                agency_id = utils.receive_message(client_sock)
                #logging.info(f'action: sorteo | result: pending | finished_clients: {self.finished_clients} | clients: {self.clients} | agency_id: {agency_id}')
                with self.lock:
                    if int(self.finished_clients) == int(self.clients):
                        #logging.info(f'action: sorteo | result: in_progress | finished_clients: {self.finished_clients} | clients: {self.clients} | agency_id: {agency_id}')
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
