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
