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

Se pide modificar el `Client` y el `Server` para poder persistir los archivos de configuración fuera de la imagen.

Se utilizaron volumenes dentro de la definición de Docker Compose, para indicar que archivos debian persistirse.

Caso `Server`:
```
volumes:
      - ./server/config.ini:/config.ini
```

Caso `Client`:
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

Para este ejercicio, se genero un script en bash llamado `validar-echo-server.sh` que verifica la comunicación con el `Server` utilizando netcat.

Para esto, el script levanta un nuevo container en Docker dentro de la misma red que el `Server`, y ejecuta el siguiente comando:

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

Se modificó el código para que los `Client` y el `Server` finalicen al recibir la señal SIGTERM.

Para los `Client` se utilzaron las librería `os/signal` y `syscall` que nos permiten recibir la señal SIGTERM y almacenarla en una variable. Luego, dentro del loop se comprueba si la señal fue recibida y finaliza el programa.

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

Se modificaron el `Client` y el `Server` para implementar el envío y recepción de apuestas.

El `Client` tendrá los datos de la apuesta definidos como variables de entorno:

- NOMBRE
- APELLIDO
- DOCUMENTO
- NACIMIENTO
- NUMERO

Al ejecutarse, se tomaran estos y se enviaran al `Server` utilizando el protocolo que se explicara en la siguiente sección.

El `Server` recibira la apuesta y la almacenara en un archivo csv utilizando la función `store_bets`

#### Protocolo

Si el `Cliente` quiere enviar un mensaje al `Server`, debera enviarle dos mensajes. El primero indicando la longitud del siguiente, y el segundo con el contenido que se desea enviar.

Los mensajes que contengan apuestas tendrán el siguiente formato

```
ID AGENCIA,NOMBRE,APELLIDO,DOCUMENTO,NACIMIENTO,NUMERO
```

La respuesta del `Server` en caso de recibir correctamente el mensaje será reenviar los siguientes datos:

```
DOCUMENTO,NACIMIENTO
```

De esta manera el `Client` podrá confirmar que el `Server` recibir los datos correctamente.

Si el `Client` no recibe una respuesta o los datos no coinciden, arrojara un error y finalizara.


#### Ejecución

Generamos la definición de Docker compose:

```
./generar-compose.sh docker-compose-dev.yaml 5
```

Levantamos todo utilizando el siguiente comando:

```
make docker-compose-up 
```

### Ejercicio 6

#### Explicación

Se modificaron los `Client` para que puedan enviar varias apuestas a la vez. Estas apuestas se leen a partir de un archivo `.csv` que es persistido utilizando volumenes.

Los `Client` enviaran una cierta cantidad de apuestas por mensaje (batch). La cantidad de mensajes estara determinada por el valor `maxAmount` dentro del archivo de configuración.

El `Server` recibira estas apuestas y las guardará utilizando la función `store_bets` al igual que en el ejercicio anterior.

#### Protocolo

Se realizaron las siguientes modificaciones al protocolo para lograr un correcto funcionamiento:

- Las apuestas se enviaran con el mismo formato del ejercicio 5, con el agregado de que dentro de un batch estarán separadas por un `;`

```
ID AGENCIA,NOMBRE1,APELLIDO1,DOCUMENTO1,NACIMIENTO1,NUMERO1;
ID AGENCIA,NOMBRE2,APELLIDO2,DOCUMENTO2,NACIMIENTO2,NUMERO2;
ID AGENCIA,NOMBRE3,APELLIDO3,DOCUMENTO3,NACIMIENTO3,NUMERO3
```

- Cuando un `Client` termina de enviar todas sus apuestas, envía un último mensaje con la palabra `FINISH` para indicarle al `Server` que ya no va a recibir más apuestas

- Una vez procesadas todas las apuestas y recibido el mensaje `FINISH`, el `Server` le enviara un último mensaje a cada `Client` con la cantidad de apuestas que recibió y este último podrá utilizar ese dato para determinar si todas las apuestas fueron enviadas o si se perdió alguna

#### Ejecución

Generamos la definición de Docker compose:

```
./generar-compose.sh docker-compose-dev.yaml 5
```

Levantamos todo utilizando el siguiente comando:

```
make docker-compose-up 
```

### Ejercicio 7

#### Explicación

Se realizan modificaciones para que al finalizar el envia de las apuestas de todos los `Client`, estos puedan consultar los ganadores.

Cuando el `Server` recibe el mensaje `FINISH` de cada `Client`, se habilita la realización del sorteo utilizando las funciones `load_bets` y `has_won`.

El sorteo se realiza una sola vez y los resultados se guardan en una variable del `Server` para que pueda ser consultado por varios `Client` sin tener que realizar el proceso del sorteo múltiples veces.

#### Protocolo

Se realizaron las siguientes modificaciones al protocolo para lograr un correcto funcionamiento:

- Se agregan dos mensajes para inicializar procesos. 
- `BET` lo envia un `Client` para indicarle al `Server` que va a empezar a recibir las apuestas. 
- `RESULTS` lo envia un `Client` para indicarle al `Server` que quiere obtener los resultados del sorteo, seguido de su ID para que el `Server` sepa que agencia es
- Si todos los `Client` enviaron el mensaje `FINISH`, al recibir el primer `RESULTS` el `Server` realizara el sorteo y le enviara los documentos de los ganadores de su agencia al `Client` con el siguiente formato

```
DOCUMENTO1,DOCUMENTO2,DOCUMENTO3,DOCUMENTO4
```

- Si nunguna de las apuestas de la agencia resultó ganadores, el `Server` le enviará el mensaje `NOWINNERS`

- Si el `Server` todavía no recibió el mensaje `FINISH` de todos los `Client`, entonces no va a responder ningun mensaje de `RESULTS`. Por lo cual, los `Client` deberán reintentar hasta conseguir una respuesta, realizando un sleep que se incrementara en cada iteración.


#### Ejecución

Generamos la definición de Docker compose:

```
./generar-compose.sh docker-compose-dev.yaml 5
```

Levantamos todo utilizando el siguiente comando:

```
make docker-compose-up 
```

## Parte 3

### Ejercicio 8

#### Explicación

Se modificó el `Server` para que pueda procesar multiples mensajes al mismo tiempo. El protocolo no sufrió ningún cambio.

#### Mecanismo de sincronización

Se utilizo la librería `threading` para generar un `thread` por cada connexión.

Para evitar condiciones de carrera se implementaron `locks` para el contador de `Client` que enviaron el mensaje `FINISH` y otro para utilizar la función `store_bets`.

#### Ejecución

Generamos la definición de Docker compose:

```
./generar-compose.sh docker-compose-dev.yaml 5
```

Levantamos todo utilizando el siguiente comando:

```
make docker-compose-up 
```