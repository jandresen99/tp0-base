# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1

### Ejercicio 1

#### Explicación

Se pide generar un script en bash llamado `generar-compose.sh` que permita generar un nuevo archivo `.yaml` que contenga una definición de Docker Compose con un servidor y una N cantidad de clientes.

El script `generar-compose.sh` llama a `mi-generador.py`, el cual se encarga de generar el archivo `.yaml`. Este consistira un servidor, N cantidad de clientes y una network.

Si es necesario, este script se actualizara en futuros ejercicios.

#### Ejecución

> Antes de ejecutar el script, debemos asegurarnos que nuestro usuario tenga permisos suficientes para ejecutarlo, utilizando `chmod`.

Para ejecutar el script debemos utilizar el siguiente comando:

```
./generar-compose.sh <archivo_de_salida> N
```

Por ejemplo:

```
./generar-compose.sh docker-compose-dev.yaml 5
```

Luego de generar la definición de Docker Compose, podemos levantar todo utilizando el siguiente comando:

```
make docker-compose-up 
```

### Ejercicio 2

#### Explicación

Se pide modificar el cliente y el servidor para poder persistir los archivos de configuración fuera de la imagen.

Se utilizaron volumenes dentro de la definición de Docker Compose, para indicar que archivos debian persistirse.

Caso Servidor:
```
volumes:
      - ./server/config.ini:/config.ini
```

Caso Cliente:
```
volumes:
      - ./client/config.yaml:/config.yaml
```

Dado que los programas priorizan las variables de entorno por sobre las variables de configuración, se decidio eliminar las variables de entorno de la definición de Docker Compose para este ejercicio.

#### Ejecución

Primero generamos la nueva definición de Docker compose:

```
./generar-compose.sh docker-compose-dev.yaml 5
```

Luego levantamos todo utilizando el siguiente comando:

```
make docker-compose-up 
```

### Ejercicio 3

#### Explicación

Para este ejercicio, se genero un script en bash llamado `validar-echo-server.sh` que verifica la comunicación con el servidor utilizando netcat.

Para esto, el script levanta un nuevo container en Docker dentro de la misma red que el servidor, y ejecuta el siguiente comando:

```
echo 'hello' | nc server 12345
```

Luego, comprueba que la respuesta sea igual al mensaje envia y muestra el log correspondiente.

#### Ejecución

Generamos la definición de Docker compose:

```
./generar-compose.sh docker-compose-dev.yaml 5
```

Levantamos todo utilizando el siguiente comando:

```
make docker-compose-up 
```

Finalmente ejecutamos el script validador:
```
./validar-echo-server.sh
```

> Antes de ejecutar el script, debemos asegurarnos que nuestro usuario tenga permisos suficientes para ejecutarlo, utilizando `chmod`.


### Ejercicio 4

#### Explicación

Se modificó el código para que los clientes y el servidor finalicen al recibir la señal SIGTERM.

Para los clientes se utilzaron las librería `os/signal` y `syscall` que nos permiten recibir la señal SIGTERM y almacenarla en una variable. Luego, dentro del loop se comprueba si la señal fue recibida y finaliza el programa.

#### Ejecución

Generamos la definición de Docker compose:

```
./generar-compose.sh docker-compose-dev.yaml 5
```

Levantamos todo utilizando el siguiente comando:

```
make docker-compose-up 
```

Para enviar la señal SIGTERM podemos utilizar el siguiente comando:

```
docker compose -f docker-compose-dev.yaml stop -t 10
```

## Parte 2

### Ejercicio 5

#### Explicación

#### Ejecución

#### Protocolo

### Ejercicio 6

#### Explicación

#### Ejecución

#### Protocolo

### Ejercicio 7

#### Explicación

#### Ejecución

#### Protocolo

## Parte 3

### Ejercicio 8

#### Explicación

#### Ejecución

#### Mecanismo de sincronización