import csv
import datetime
import time
import logging


""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574


""" A lottery bet registry. """
class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)

""" Checks whether a bet won the prize or not. """
def has_won(bet: Bet) -> bool:
    return bet.number == LOTTERY_WINNER_NUMBER

"""
Persist the information of each bet in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def store_bets(bets: list[Bet]) -> None:
    with open(STORAGE_FILEPATH, 'a+') as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow([bet.agency, bet.first_name, bet.last_name,
                             bet.document, bet.birthdate, bet.number])

"""
Loads the information all the bets in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def load_bets() -> list[Bet]:
    with open(STORAGE_FILEPATH, 'r') as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])

"""
Receives a message from a client socket.
"""
def receive_message(client_sock):
    size = int.from_bytes(client_sock.recv(2), byteorder='big')

    data = b""
    while len(data) < size:
        packet = client_sock.recv(size - len(data))
        if not packet:
            raise ConnectionError("Connection closed unexpectedly")
        data += packet
    
    msg = data.decode('utf-8').strip()

    return msg

"""
Decodes a bet from a client socket.
"""
def decode_bet(client_sock):
    msg = receive_message(client_sock)
    addr = client_sock.getpeername()

    bet_data = msg.split(',')
    if len(bet_data) != 6:
        logging.error(f"action: receive_message | result: fail | ip: {addr[0]} | msg: {msg} | error: Invalid bet data")
        raise ValueError("Invalid bet data")
    
    bet = Bet(bet_data[0], bet_data[1], bet_data[2], bet_data[3], bet_data[4], bet_data[5])
    return bet, addr, msg

"""
Acknowledges a bet to a client socket.
"""
def acknowledge_bet(client_sock, document, number):
    client_sock.send("{},{}\n".format(document,number).encode('utf-8'))