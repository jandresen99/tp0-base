import sys

def generate_file(filename, clients):
    compose = f"""name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
    - PYTHONUNBUFFERED=1
    - CLIENTS={clients}
    networks:
    - testing_net
    volumes:
      - ./server/config.ini:/config.ini
"""
    
    for i in range(1, clients + 1):
        compose += f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={i}
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{i}.csv:/agency.csv
    depends_on:
      - server
"""
    
    compose += """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""
    
    with open(filename, "w") as f:
        f.write(compose)
    
    print(f"Archivo {filename} generado con {clients} clientes.")

if __name__ == "__main__":
    filename = sys.argv[1]
    clients = int(sys.argv[2])
    
    generate_file(filename, clients)