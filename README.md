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
