# Manual

## Como ejecutar cada microservico a la vez (docker + docker-compose)

1. Tener docker instalado
2. Abrir docker (para ejecutar la mv de docker)
3. `docker-compose build`
4. `docker-compose up`

**Si no funciona "docker-compose" prueba con `docker compose build` (sin "-")**

## Como ejecutar sin docker

1. Tener go instalado
2. Abrir el directorio de microservicio (Por ejemplo: `cd ./gateway`)
3. `go run main.go`
