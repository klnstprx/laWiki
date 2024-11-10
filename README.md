# Manual

## Cómo ejecutar los microservicios de forma simultánea (docker + docker-compose)

1. Tener docker instalado.
2. Abrir docker (para ejecutar la mv de docker).
3. Copiar el contenido `default_config.docker.toml` y pegarlo en un archivo llamado `config.docker.toml`.
4. Rellenar `config.docker.toml` con los datos necesarios. Por ejemplo: MONGODB_URI="mongodb+srv://username:password@cluster0.rfz8f.mongodb.net/".
5. `docker-compose build`
6. `docker-compose up`

**Si no funciona "docker-compose" prueba con `docker compose build` (sin "-")**.

## Cómo ejecutar sin docker:

1. Tener go instalado.
2. Copiar el contenido `default_config.toml` y pegarlo en un archivo llamado `config.toml`.
3. Abrir el directorio de microservicio (Por ejemplo: `cd ./gateway`).
4. Rellenar `config.toml` con los datos necesarios. Por ejemplo: MONGODB_URI="mongodb+srv://username:password@cluster0.rfz8f.mongodb.net/".
5. `go run main.go`

**La diferencia entre config.toml y config.docker.toml es que en uno llamamos el host de cada microservio por "localhost" y en el otro por el nombre del servicio en docker-compose.**
