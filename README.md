# TP0: Docker + Comunicaciones + Concurrencia

- **Alumno**: Agustín Altamirano
- **Padrón**: 110237

## Ejercicio 1

El script `generar-compose.sh` permite generar un archivo docker compose con una cantidad de clientes variable. Se usa
de la siguiente forma:

```bash
./generar-compose.sh <nombre_del_archivo> <cantidad_clientes>
```

Este script simplemente ejecuta un script de Python: `generate-compose.py`, el cual hace uso del módulo `pyyaml` para
generar el archivo según la cantidad de clientes.

## Ejercicio 2

Para permitir cambiar los archivos de configuración `config.yml` y `config.ini` sin necesidad de reconstruir la imagen, se montan estos archivos como volúmenes en el contenedor. De esta forma, cualquier cambio que se realice en estos archivos en el host se reflejará automáticamente dentro el contenedor.

Los volúmenes se definen en el archivo `docker-compose.yml` de la siguiente manera:

```yaml
name: tp0
services:
  server:
    volumes:
      - ./server/config.ini:/config.ini:ro

  client1:
    volumes:
      - ./client/config.yaml:/config.yaml:ro
```

Como observación, el sufijo `:ro` indica que el volumen es de solo lectura, lo que significa que el contenedor no podrá modificar los archivos montados.

Adicionalmente, se actualizó el script `generate-compose.py` para incluir estos volúmenes.

## Ejercicio 3

En este ejercicio se implementa un validador del servidor mediante el script `validar-echo-server.sh`. El mismo envía
un mensaje de prueba al echo server y verifica que la respuesta sea correcta. Para ejecutar el validador, se puede usar
el siguiente comando:

```bash
./validar-echo-server.sh
```

Como es imposible garantizar que la computadora host tenga instalado `netcat` (o mismo que posea un sistema operativo
UNIX), se implementó el validador dentro de un container de Docker, levantado mediante docker compose desde el script.
Además, para que el validador pueda comunicarse con el servidor sin que este último tenga que exponer un puerto, ambos
contenedores deben compartir la misma red. Esto se logra utilizando la opción `network_mode: "container:server"` en la definición del servicio del validador, lo que permite que el validador utilice la misma red local del contenedor del
servidor.

Para ello, se utiliza el archivo `docker-compose-validator.yml`, el cual define un servicio llamado `validator` que a
su vez ejecuta otro script llamado `run.sh`. Dentro de este script, se realiza el envío y recepción del mensaje de
prueba utilizando `netcat`.

Una aclaración: en `run.sh` se definen las variables `SERVER_HOST` y `SERVER_PORT` para especificar la dirección y el
puerto del servidor al que se desea conectar. Estas variables se pueden modificar según sea necesario para apuntar al
servidor correcto.

## Ejercicio 4

En este ejercicio se implementa el manejo de la señal `SIGTERM` tanto en el cliente como en el servidor, permitiendo un
_graceful shutdown_ de ambos procesos, cerrando todos sus recursos (sus sockets).

### Cliente

En `main.go`, la función `handleSignal()` crea un channel de señales y registra la notificación de `SIGTERM` mediante
`signal.Notify`. Este channel se pasa al cliente en su constructor.

En `client.go`, dentro de `StartClientLoop()`, en cada iteración del loop se utiliza un `select` con dos cases:

- El channel de señales: si se recibe una señal `SIGTERM`, se loguea el evento y se retorna inmediatamente, finalizando
  el loop.
- El case `default`: ejecuta el envío y recepción normal de mensajes.

Además, el método `closeClientSocket()` se encarga de cerrar la conexión TCP de forma segura.

### Servidor

En `main.py`, se registra un handler para `SIGTERM` usando la librería `signal`. Cuando se recibe la señal, se invoca a
la función handler de la señal, que llama a `server.request_shutdown()`. Este método simplemente setea el flag
`_stop_requested = True`.

En `server.py`, el loop principal verifica en cada iteración si `_stop_requested` es `True` para salir del
bucle. El socket del servidor tiene un timeout de 0.2 segundos, lo que evita que `accept()` bloquee
indefinidamente y permite que el servidor chequee periódicamente si se solicitó el shutdown. Cuando el loop termina, el
context manager (`__exit__`) llama a `close()`, que cierra el socket del servidor de forma ordenada.

En ambos casos, se eligió por un cierre _polite_ de los procesos. Al recibir la señal, tanto cliente como servidor
esperan a que finalice el ciclo actual del echo antes de cerrar los sockets.

## Ejercicio 5

En este ejercicio se implementa un sistema de batchría distribuido, donde los clientes (agencias de batchría) registran
apuestas en un servidor central (batchría central) mediante un protocolo binario sobre TCP.

### SafeSocket

Tanto el cliente como el servidor implementan un `SafeSocket` que encapsula al socket TCP nativo y garantiza lecturas y
escrituras completas, manejando los problemas de _short reads_ y _short writes_ inherentes a TCP.

Las implementaciones de cliente y servidor son análogas. Ambas disponibilizan dos métodos principales: `Read(n)` para
leer exactamente `n` bytes, y `Write(data)` para enviar todos los bytes de `data`. Ambos métodos manejan casos de
lectura o escritura parcial, acumulando o enviando en loop hasta completar la operación. Si se detecta que la conexión
fue cerrada (lectura de 0 bytes o escritura de 0 bytes), se retorna un error.

### Capas de abstracción

La arquitectura tanto del cliente como del servidor se organiza en tres capas, de abajo hacia arriba:

1. **SafeSocket**: capa de transporte que garantiza lecturas y escrituras completas sobre TCP, abstrayendo los short
   reads/writes. Es la única capa que interactúa directamente con el socket.

2. **Protocolo**: capa de serialización/deserialización que define el formato binario de los mensajes. Se encarga de
   convertir estructuras de datos del dominio del problema (como las apuestas) a secuencias de bytes o viceversa.
   Utiliza el `SafeSocket` para enviar y recibir estos bytes.

3. **LotteryAgency / LotteryCentral** : capa de lógica de negocio. La agencia arma las apuestas y las envía mediante el
   protocolo, y procesa las confirmaciones recibidas. La batchría central recibe las apuestas, las almacena y envía la
   confirmación. Ninguna de estas clases conoce detalles de serialización ni de sockets.

### Mensajes del protocolo

Para todo el protocolo, se considera **big-endian** como el endianness de la red.

#### RegisterBetBatch (código de mensaje: 1)

Enviado por el cliente al servidor para registrar un batch de apuestas.

**Header:**

| Campo        | Tamaño  | Tipo   | Descripción                                        |
| ------------ | ------- | ------ | -------------------------------------------------- |
| messageCode  | 2 bytes | uint16 | Código del mensaje (`1` = REGISTER_BET_BATCH)      |
| agencyId     | 2 bytes | uint16 | Id de la agencia que envía las apuestas            |
| betAmount    | 4 bytes | uint32 | Cantidad de apuestas en el batch                   |
| betBatchSize | 4 bytes | uint32 | Tamaño total en bytes de los datos de las apuestas |

**Por cada apuesta** (repetido `betAmount` veces):

| Campo           | Tamaño   | Tipo   | Descripción                                 |
| --------------- | -------- | ------ | ------------------------------------------- |
| document        | 4 bytes  | uint32 | DNI del participante                        |
| birthdate       | 10 bytes | string | Fecha de nacimiento en formato `YYYY-MM-DD` |
| number          | 4 bytes  | uint32 | Número apostado                             |
| firstNameLength | 2 bytes  | uint16 | Longitud del nombre en bytes                |
| lastNameLength  | 2 bytes  | uint16 | Longitud del apellido en bytes              |
| firstName       | variable | string | Nombre del participante (UTF-8)             |
| lastName        | variable | string | Apellido del participante (UTF-8)           |

Cada apuesta tiene un tamaño base fijo de 22 bytes (`BET_BASE_SIZE`) más los bytes variables del nombre y apellido.

#### ConfirmBetBatch (código de mensaje: 2)

Enviado por el servidor al cliente como respuesta a un `RegisterBetBatch`.

| Campo       | Tamaño  | Tipo   | Descripción                                            |
| ----------- | ------- | ------ | ------------------------------------------------------ |
| messageCode | 2 bytes | uint16 | Código del mensaje (`2` = CONFIRM_BET_BATCH)           |
| agencyId    | 2 bytes | uint16 | Id de la agencia                                       |
| betAmount   | 4 bytes | uint32 | Cantidad de apuestas registradas                       |
| resultCode  | 1 byte  | uint8  | Resultado de la operación (`0` = SUCCESS, `1` = ERROR) |

### Secuencia de comunicación

La comunicación sigue un esquema simple:

1. El cliente construye un batch de apuestas (para este ejericio solo cuenta con una apuesta) y envía un mensaje
   **RegisterBetBatch** al servidor con las apuestas serializadas en formato binario.
2. El servidor recibe el mensaje, deserializa las apuestas, las almacena y responde con un mensaje
   **ConfirmBetBatch** indicando el resultado de la operación (éxito o error) y la cantidad de apuestas procesadas.

### Persistencia de apuestas

El servidor almacena las apuestas recibidas en el archivo `bets.csv`, utilizando formato CSV con las columnas:
`agency`, `first_name`, `last_name`, `document`, `birthdate`, `number`.

Para que este archivo persista más allá del ciclo de vida del contenedor y sea accesible desde el host, se monta un
**Docker volume** en el `docker-compose-dev.yaml`:

```yaml
server:
  volumes:
    - ./.data_server/bets.csv:/bets.csv
```

De esta forma, las apuestas se guardan en `.data_server/bets.csv` en el host.

### Datos de la apuesta desde variables de entorno

Los datos de la apuesta del cliente se toman desde variables de entorno definidas en el archivo `client/.env`:

```env
CLI_BET_FIRSTNAME=Santiago Lionel
CLI_BET_LASTNAME=Lorca
CLI_BET_DOCUMENT=30904465
CLI_BET_BIRTHDATE=1999-03-17
CLI_BET_NUMBER=7574
```

Este archivo se referencia en el `docker-compose-dev.yaml` mediante la directiva `env_file`, lo que permite modificar
los datos de la apuesta sin reconstruir la imagen del cliente.

### Ejercicio 6

En este ejercicio se reemplaza la lectura de una única apuesta desde variables de entorno por la lectura de múltiples
apuestas desde un archivo CSV, enviándolas al servidor en batches.

#### Lectura del CSV

Cada cliente lee sus apuestas desde un archivo CSV ubicado en `/data/agency-{ID}.csv` dentro del contenedor. Este
archivo se monta como volumen de solo lectura en el `docker-compose-dev.yaml`:

```yaml
client1:
  volumes:
    - ./.data/agency-1.csv:/data/agency-1.csv:ro
```

El componente `BetCSVReader` abre el archivo y expone un método `Next()` que lee una línea por invocación, parsea los
campos (nombre, apellido, documento, fecha de nacimiento y número) y retorna un `Bet`. Cuando no quedan más registros,
retorna `io.EOF`.

#### Construcción de batches

El `BatchBuilder` acumula las apuestas leídas por el `BetCSVReader` en lotes, respetando dos límites:

1. **Cantidad máxima de apuestas por batch** (`CLI_BATCH_MAXAMOUNT`): configurable mediante variable de entorno.
2. **Tamaño máximo del payload** (`MAX_BET_BATCH_SIZE` = 8 KB): límite fijo en bytes para el payload serializado del
   batch.

Para cada apuesta, el builder calcula su tamaño serializado (22 bytes base + longitud del nombre + longitud del
apellido) y verifica si agregarla al batch actual excedería alguno de los dos límites. Si es así, envía el batch
acumulado al servidor y comienza uno nuevo con la apuesta pendiente.

#### Flujo de envío

El método `Run()` de la agencia orquesta el proceso completo:

1. Abre el `BetCSVReader` con la ruta al CSV correspondiente.
2. Crea un `BatchBuilder` con los límites configurados.
3. Invoca `BuildFromReader`, que itera sobre todas las apuestas del CSV y, por cada batch completo, ejecuta un
   callback que:
   - Envía el mensaje `RegisterBetBatch` con el batch serializado.
   - Espera la respuesta `ConfirmBetBatch` del servidor.
   - Loguea el resultado y la cantidad de apuestas procesadas.
4. Al finalizar el CSV, envía el último batch parcial (si lo hay).
