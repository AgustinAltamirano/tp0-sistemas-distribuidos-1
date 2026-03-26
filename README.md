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
