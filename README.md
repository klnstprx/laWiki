# Manual

## Cómo ejecutar los microservicios de forma simultánea (docker + docker-compose)

1. Tener docker instalado.
2. Abrir docker (para ejecutar la mv de docker).
3. Crear un archivo llamado ".env.docker".
4. Rellenar variables de entorno necesarias (Mirar .env.docker.default).
5. `docker-compose build`
6. `docker-compose up`

**Si no funciona "docker-compose" prueba con `docker compose build` (sin "-")**.

## Cómo ejecutar sin docker:

1. Tener go instalado.
2. Abrir el directorio de microservicio (Por ejemplo: `cd ./gateway`).
3. Crear ".env" y rellenar (Mirar .env.default).
4. `go run main.go`
