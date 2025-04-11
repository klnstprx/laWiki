# LaWiki

LaWiki is a knowledge-sharing platform built using a microservices architecture. It allows users to create, edit, and translate wiki entries, comment on versions, and manage media.
## Project Description

LaWiki is composed of the following microservices:

*   **Auth Service:** Handles user authentication and authorization using Google OAuth 2.0 and JWT (JSON Web Tokens).
*   **Wiki Service:** Manages the creation, retrieval, updating, and deletion of wiki entries.
*   **Entry Service:** Manages individual entries within a wiki, including their translated fields.
*   **Version Service:** Tracks and manages different versions of entries, allowing for content history and rollback capabilities.
*   **Comment Service:** Enables users to comment on specific versions of entries, fostering discussion and collaboration.
*   **Media Service:** Handles the uploading, storage, and retrieval of media files (e.g., images) using Cloudinary.
*   **Translation Service:** Provides translation capabilities for wiki entries using the DeepL API.
*   **API Gateway:** Acts as a single entry point for all client requests, routing them to the appropriate microservice. Also handles authentication and security concerns.
*   **Frontend:** A React-based user interface that provides a user-friendly way to interact with the platform.

## Prerequisites

Before running the project, ensure that you have the following installed:

*   **Go:** Version 1.22 or later.
*   **Node.js:** Version 18 or later.
*   **Docker:**  Required for containerization.
*   **Docker Compose:**  Required for orchestrating multi-container Docker applications.
*   **MongoDB:** Ensure you have a MongoDB instance running or accessible.

You may need to install other dependencies like `npx` for swagger combining, you can check the Makefile for the required dependencies.

## Configuration

Each backend service utilizes a `config.toml` file to manage its settings.  Example configurations are provided in `src/backend/default_config.toml` and `src/backend/default_config.docker.toml`.
  * `default_config.toml` is for local development
  * `default_config.docker.toml` is for running using docker

Modify these files to match your environment.  Key configuration parameters include:

*   **MongoDB URI:**  The connection string for your MongoDB instance.
*   **API Gateway URL:** The URL where the API Gateway service is running.
*   **JWT Secret:**  A secret key used for signing JWTs.  Keep this secure!
*   **Service URLs:** The URLs of the other microservices (used by the API Gateway).
*   **Cloudinary Credentials:** Required for the Media Service if using Cloudinary for media storage.
*	**MailSender Credentials**: Required for MailSender API
*   **DeepL API Key:** Required for the Translation Service.

**Important:** Store secrets (like API keys and the JWT secret) securely, especially in production environments. Do not commit them directly to your repository.

## Running the Project

There are two primary ways to run the project: using Docker Compose (recommended) or running the services locally.

### Method 1: Using Docker Compose

This method simplifies the setup process by containerizing all the services.

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd <repository_directory>
    ```

2.  **Configure Docker environment variables**:
    If you are using Docker, you can set environment variables directly in the `docker-compose.yml` or within each service's environment section.

3.  **Build and Run the Services:**

    ```bash
    cd src/backend
    docker-compose up --build
    ```

    This command builds the Docker images for all services defined in `docker-compose.yml` and starts them. The `-d` flag can be added to run the containers in detached mode (background).

4.  **Access the Application:**

    *   The frontend will be accessible at `http://localhost:5173` (or the port you configured in `src/frontend/vite.config.js` and `src/backend/default_config.toml`).
    *   The API Gateway will be accessible at `http://localhost:8000` (or the port you configured in `src/backend/default_config.toml`).
    *   Swagger documentation for the API can be found at `http://localhost:8000/api/swagger/index.html`.

### Method 2: Running Services Locally

This method requires you to build and run each service individually.

1.  **Navigate to a service directory:**

    ```bash
    cd src/backend/<service_name>
    ```

    where `<service_name>` is the name of the service you want to run (e.g., `auth`, `wiki`, `gateway`).

2.  **Build the service:**
   Use the provided Makefile to build the service.

    ```bash
    make build-<service_name>
    ```

    For example, to build the auth service, you'd run:

    ```bash
    make build-auth
    ```

3.  **Run the service:**
   Use the provided Makefile to run the service. This command will build and execute the service.

    ```bash
    make run-<service_name>-service
    ```

    For example, to run the auth service, you'd run:

    ```bash
    make run-auth-service
    ```

    This compiles the Go code and runs the service.  It also redirects output to a log file and saves the process ID to a `.pid` file for later stopping.

4.  **Repeat steps 1-3 for each backend service.** Make sure to configure each service with the correct ports and service URLs in the config.toml files.

5.  **Frontend setup:**
    First, navigate to the frontend directory.

    ```bash
    cd src/frontend
    ```

6.  **Install dependencies:**

    ```bash
    npm install
    ```

7.  **Run the frontend:**

    ```bash
    npm run dev
    ```

    This starts the development server, and the frontend will be accessible at `http://localhost:5173` (or the port you configured).
   **IMPORTANT**: set the `VITE_API_BASE_URL` environment variable in the frontend to the correct `API_GATEWAY_URL`

### Generating and Combining Swagger Documentation

The backend services include Swagger documentation. To generate and combine these documents:

1.  **Navigate to the `src/backend` directory:**

    ```bash
    cd src/backend
    ```

2.  **Run the `combine-swagger` target in the Makefile:**

    ```bash
    make combine-swagger
    ```

    This command uses `swag` to generate the Swagger documentation for each service and then combines them into a single `swagger.json` file located in the `gateway/docs/` directory. The API Gateway serves this combined documentation.

## Makefile Commands

The `src/backend/Makefile` provides several convenient commands:

*   `make all`: Builds and pushes all backend services to Docker Hub (requires Docker Hub credentials and a suitable Docker Hub repository name).
*   `make build`: Builds all backend services.
*   `make push`: Pushes all backend services to Docker Hub.
*   `make build-<service_name>`: Builds a specific service (e.g., `make build-auth`).
*   `make push-<service_name>`: Pushes a specific service to Docker Hub.
*   `make run-all`: Runs all backend services locally (builds them first).
*   `make run-<service_name>-service`: Runs a specific service locally (e.g., `make run-auth-service`).
*   `make clean`: Stops all running services and removes their PID files.
*   `make combine-swagger`: Generates and combines Swagger documentation for all services.

## Frontend Configuration

The frontend's base URL is set by the `VITE_API_BASE_URL` environment variable. This variable points to the API Gateway's address.


## License

This project is licensed under the MIT License.
