# RealTime Messenger

![Generic badge](https://img.shields.io/badge/maintainer-devworlds-blue.svg)
[![codecov](https://codecov.io/gh/devworlds/eda-message-go/branch/main/graph/badge.svg)](https://codecov.io/gh/devworlds/eda-message-go)
[![Test](https://github.com/devworlds/eda-message-go/actions/workflows/build.yml/badge.svg)](https://github.com/devworlds/eda-message-go/actions/workflows/build.yml)
![Generic badge](https://img.shields.io/badge/version-v0.1.0-green.svg)


### Kubernetes Project: Local Deployment with Minikube, Helm, Docker, and Go

This guide explains how to set up and run a local environment using Minikube, Helm, Docker, Go, and the required services (PostgreSQL, Kafka, Auth, WebSocket, and Persistence).

---

#### **Prerequisites**

Make sure you have the following tools installed:

- [Minikube](https://minikube.sigs.k8s.io/docs/)
- [Helm](https://helm.sh/docs/intro/install/)
- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Go 1.23+](https://go.dev/doc/install)

---

#### **1. Starting Minikube**

```sh
minikube start
```

---

#### **2. Setting Up Helm**

Add the Bitnami repository and update:

```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

---

#### **3. Installing Dependencies with Helm**

**PostgreSQL:**

```sh
helm install postgres bitnami/postgresql -f charts/postgresql/values.yaml
```

**Kafka:**

```sh
helm install kafka bitnami/kafka -f charts/kafka/values.yaml
```

---

#### **4. Creating Kubernetes Secrets**

```sh
kubectl create secret generic auth-service-secret --from-literal=jwtSecret=secret-token-here
```

---

#### **5. Building and Deploying Services**

Before each build, run:

```sh
eval $(minikube docker-env)
```

**WebSocket Service:**

```sh
docker build -f websocket/Dockerfile -t websocket-service .
kubectl delete deployment websocket-service --ignore-not-found
kubectl apply -f deployments/websocket/websocket-deployment.yaml
```

**Auth Service:**

```sh
docker build -f auth/Dockerfile -t auth-service .
kubectl delete deployment auth-service --ignore-not-found
kubectl apply -f deployments/auth/auth-deployment.yaml
```

**Persistence Service:**

```sh
docker build -f persistence/Dockerfile -t persistence-service .
kubectl delete deployment persistence-service --ignore-not-found
kubectl apply -f deployments/persistence/persistence-deployment.yaml
```

---

#### **6. Port Forwarding**

**WebSocket Service:**

```sh
kubectl port-forward service/websocket-service 30080:80
```

**PostgreSQL:**

```sh
kubectl port-forward svc/postgres-postgresql 5432:5432
```

**Auth Service:**

```sh
kubectl port-forward deployment/auth-service 8081:8081
```

---
### Services Overview

Below is a brief explanation of the main services in this project and their responsibilities:

---

### Auth Service
| Method | Route            | Description                        |
|--------|------------------|------------------------------------|
| POST   | `localhost:8081/login`            | Login for access token     |

### Websocket Service
| Method | Route            | Description                        |
|--------|------------------|------------------------------------|
|  GET  | `ws://localhost:30080/ws`            | Websocket server    |

---

#### **Auth Service**

Responsible for creating and validating JWT tokens. It provides authentication endpoints that other services use to verify if a client is authorized.

---

#### **Authentication (Token Generation)**

The client must send a JSON payload containing `username` and `password`. The Auth Service checks if the user exists in the database and if the credentials are correct. If authentication is successful, a JWT token is generated and returned in the response.

**Example request:**

```json
POST /login
{
  "username": "your_username",
  "password": "your_password"
}
```

```bash
curl -X POST http://localhost:8081/login -H "Content-Type: application/json" -d '{"username":"your_username","password":"your_password"}'
```
**Example response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
---

#### **Token Validation (WebSocket Connection)**

When a client connects to the WebSocket Service, the **first message** sent by the client must be a valid JWT token. The WebSocket Service forwards this token to the Auth Service for validation.

- If the token is valid, the client is authenticated and allowed to send and receive broadcast messages.
- If the first message is not a valid JWT token, the connection is immediately closed.

**Important:**
While the client is not authenticated, it cannot send or receive any broadcast messages. Only after successful authentication (valid token) does the client gain access to the messaging system.

---

This flow ensures that only authenticated users can participate in the real-time messaging system, providing both security and control over client communications.

#### **WebSocket Service**

Manages client connections via WebSocket. For each client, it checks authentication by communicating with the Auth Service. Only authenticated clients can send and receive messages.
Messages sent by clients are published to the Kafka topic `websocket-messages`. The WebSocket Service also consumes this topic to broadcast messages to all connected and authenticated clients.

---

#### **Persistence Service**

Consumes messages from the `websocket-messages` Kafka topic. Its main responsibility is to persist all messages sent by clients into the PostgreSQL database. This ensures that messages are not lost, even if there are issues with broadcasting to clients.

---

#### **Kafka**

Acts as the message broker between services. All client messages are sent to the `websocket-messages` topic, which is consumed by both the WebSocket Service (for broadcasting) and the Persistence Service (for storage).

---

#### **PostgreSQL**

Stores all messages received from the Persistence Service, providing durability and reliability for message history and recovery.

---

This architecture ensures secure authentication, reliable message delivery, and data persistence for all client communications.

---

### Usage

Below are example usage scenarios for interacting with the application:

---

#### **Scenario 1: Successful Authentication and Real-Time Messaging**

After completing the port forwarding steps, open `index.html` in your browser. This page provides a simple interface to connect to the WebSocket server.

1. **Generate a JWT Token:**
   Use the `/login` endpoint to authenticate with your username and password. The response will include a JWT token.

2. **Connect to the WebSocket Server:**
   Use the interface in `index.html` to establish a WebSocket connection.

3. **Authenticate:**
   As the **first message**, send your JWT token through the WebSocket connection.
   If the token is valid, you will be authenticated and added to the hub of authenticated clients.

4. **Send and Receive Messages:**
   Once authenticated, you can send and receive real-time messages.
   Example message format:

   ```json
   {
     "id": "023e4567-e89b-12d3-a456-426614174000",
     "content": "Player 7 pick one cards.",
     "timestamp": "2024-05-11T15:30:00Z"
   }
   ```

   When you send a message:
   - It is published to the Kafka topic `websocket-messages`.
   - The Persistence Service consumes the message and saves it to PostgreSQL.
   - The WebSocket Service also consumes the message and broadcasts it to all authenticated clients in real time.

---

#### **Scenario 2: Invalid Authentication**

If you open `index.html` and connect to the WebSocket server, but the **first message** you send is **not** a valid JWT token, the server will immediately close your connection.

**Note:**
Clients must be authenticated before they can send or receive any broadcast messages. Only after successful authentication (by sending a valid JWT token as the first message) will the client be able to participate in real-time messaging.

---

This flow ensures secure, real-time communication where only authenticated users can interact with the system.

#### **Notes**

- Make sure your Docker context is set to Minikube (`eval $(minikube docker-env)`) before building images.
- To update images, repeat the build process and re-apply the deployment.
- The configuration files (`values.yaml` and `deployment.yaml`) should be properly set up according to your needs.

---

#### **References**

- [Minikube Documentation](https://minikube.sigs.k8s.io/docs/)
- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [Docker Documentation](https://docs.docker.com/)

---
