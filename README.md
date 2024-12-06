# Manual

## Cómo ejecutar los microservicios de forma simultánea (docker + docker-compose)

1. Tener docker instalado.
2. Abrir docker (para ejecutar la mv de docker).
3. Copiar el contenido `default_config.docker.toml` y pegarlo en un archivo llamado `config.docker.toml`.
4. Rellenar `config.docker.toml` con los datos necesarios. Por ejemplo: MONGODB_URI="mongodb+srv://username:password@cluster0.rfz8f.mongodb.net/".
5. `docker-compose build`
6. `docker-compose up`

## Cómo ejecutar el frontend con docker

1. Tener docker instalado
2. Abrir docker (para ejecutar la mv de docker).
3. Navegar hacia el directorio del frontend del proyecto en un CLI.
4. `docker-compose build`
5. `docker-compose up`
6. Abrir la app en localhost:5173

**Si no funciona "docker-compose" prueba con `docker compose build` (sin "-")**.
   
## Cómo ejecutar backend sin docker:

1. Tener go instalado.
2. Copiar el contenido `default_config.toml` y pegarlo en un archivo llamado `config.toml`.
3. Abrir el directorio de microservicio (Por ejemplo: `cd ./gateway`).
4. Rellenar `config.toml` con los datos necesarios. Por ejemplo: MONGODB_URI="mongodb+srv://username:password@cluster0.rfz8f.mongodb.net/".
5. `go run main.go`

**La diferencia entre config.toml y config.docker.toml es que en uno llamamos el host de cada microservio por "localhost" y en el otro por el nombre del servicio en docker-compose.**

## Cómo ejecutar frontend sin docker:

1. Tener Node.js instalado
2. Ejecutar 'npm install' en el directorio del frontend
3. 'npm run dev'
4. Navegar a http://localhost:5173/

Estructura del repositorio:

```bash
├── clientexample
├── docs
│   ├── component__Diagramas_de_Componentes.png
│   ├── deployment__Diagrama_de_Despliegue.png
│   ├── ifml.json
│   ├── UML.mdzip
│   └── UML.mdzip.bak
├── README.md
└── src
    ├── backend
    │   ├── auth
    │   │   ├── config
    │   │   │   └── auth_config.go
    │   │   ├── Dockerfile
    │   │   ├── docs
    │   │   │   ├── docs.go
    │   │   │   ├── swagger.json
    │   │   │   └── swagger.yaml
    │   │   ├── go.mod
    │   │   ├── go.sum
    │   │   ├── handler
    │   │   │   └── auth_handler.go
    │   │   ├── main.go
    │   │   ├── model
    │   │   │   └── auth_model.go
    │   │   ├── README.md
    │   │   └── router
    │   │       └── auth_router.go
    │   ├── comment
    │   │   ├── comment-service
    │   │   ├── comment-service.log
    │   │   ├── config
    │   │   │   └── comment_config.go
    │   │   ├── database
    │   │   │   └── comment_database.go
    │   │   ├── Dockerfile
    │   │   ├── docs
    │   │   │   ├── docs.go
    │   │   │   ├── swagger.json
    │   │   │   └── swagger.yaml
    │   │   ├── go.mod
    │   │   ├── go.sum
    │   │   ├── handler
    │   │   │   └── comment_handler.go
    │   │   ├── main.go
    │   │   ├── model
    │   │   │   └── comment_model.go
    │   │   └── router
    │   │       └── comment_router.go
    │   ├── config.docker.toml
    │   ├── config.toml
    │   ├── default_config.docker.toml
    │   ├── default_config.toml
    │   ├── docker-compose.yml
    │   ├── entry
    │   │   ├── config
    │   │   │   └── entry_config.go
    │   │   ├── database
    │   │   │   └── entry_database.go
    │   │   ├── Dockerfile
    │   │   ├── docs
    │   │   │   ├── docs.go
    │   │   │   ├── swagger.json
    │   │   │   └── swagger.yaml
    │   │   ├── entry-service
    │   │   ├── entry-service.log
    │   │   ├── entry-service.pid
    │   │   ├── go.mod
    │   │   ├── go.sum
    │   │   ├── handler
    │   │   │   └── entry_handler.go
    │   │   ├── main.go
    │   │   ├── model
    │   │   │   └── entry_model.go
    │   │   └── router
    │   │       └── entry_router.go
    │   ├── gateway
    │   │   ├── api-gateway
    │   │   ├── api-gateway.log
    │   │   ├── config
    │   │   │   └── gateway_config.go
    │   │   ├── Dockerfile
    │   │   ├── docs
    │   │   │   ├── docs.go
    │   │   │   ├── swagger.json
    │   │   │   └── swagger.yaml
    │   │   ├── go.mod
    │   │   ├── go.sum
    │   │   ├── handler
    │   │   │   ├── gateway_health.go
    │   │   │   └── gateway_proxy.go
    │   │   ├── main.go
    │   │   ├── middleware
    │   │   │   └── gateway_middleware.go
    │   │   └── router
    │   │       └── gateway_router.go
    │   ├── Makefile
    │   ├── media
    │   │   ├── config
    │   │   │   └── media_config.go
    │   │   ├── database
    │   │   │   └── media_database.go
    │   │   ├── Dockerfile
    │   │   ├── docs
    │   │   │   ├── docs.go
    │   │   │   ├── swagger.json
    │   │   │   └── swagger.yaml
    │   │   ├── go.mod
    │   │   ├── go.sum
    │   │   ├── handler
    │   │   │   └── media_handler.go
    │   │   ├── main.go
    │   │   ├── media-service
    │   │   ├── media-service.log
    │   │   ├── model
    │   │   │   └── model_media.go
    │   │   └── router
    │   │       └── media_router.go
    │   ├── run-services.ps1
    │   ├── swagger-config.json
    │   ├── version
    │   │   ├── config
    │   │   │   └── version_config.go
    │   │   ├── database
    │   │   │   └── version_database.go
    │   │   ├── Dockerfile
    │   │   ├── docs
    │   │   │   ├── docs.go
    │   │   │   ├── swagger.json
    │   │   │   └── swagger.yaml
    │   │   ├── go.mod
    │   │   ├── go.sum
    │   │   ├── handler
    │   │   │   └── version_handler.go
    │   │   ├── main.go
    │   │   ├── model
    │   │   │   └── version_model.go
    │   │   ├── router
    │   │   │   └── version_router.go
    │   │   ├── version-service
    │   │   └── version-service.log
    │   └── wiki
    │       ├── config
    │       │   └── wiki_config.go
    │       ├── database
    │       │   └── wiki_database.go
    │       ├── Dockerfile
    │       ├── docs
    │       │   ├── docs.go
    │       │   ├── swagger.json
    │       │   └── swagger.yaml
    │       ├── go.mod
    │       ├── go.sum
    │       ├── handler
    │       │   └── wiki_handler.go
    │       ├── main.go
    │       ├── model
    │       │   └── wiki_model.go
    │       ├── router
    │       │   └── wiki_router.go
    │       ├── wiki-service
    │       └── wiki-service.log
    └── frontend
        ├── deno.lock
        ├── dist
        │   ├── assets
        │   │   ├── index-BPvgi06w.css
        │   │   └── index-D1NZhrJd.js
        │   ├── index.html
        │   └── vite.svg
        ├── eslint.config.js
        ├── index.html
        ├── package.json
        ├── package-lock.json
        ├── public
        │   └── vite.svg
        ├── src
        │   ├── api.js
        │   ├── App.css
        │   ├── App.jsx
        │   ├── index.css
        │   ├── main.jsx
        │   └── pages
        │       └── Home.jsx
        └── vite.config.js

82 directories, 199 files
```
