# WebSocket Chat API Backend (Learning Project)

This is a **learning project** designed to explore and demonstrate a basic chat backend written in Go. It provides WebSocket support for real-time communication, JWT-based authentication, and minimal persistence. Most of the code is **placeholder** or **boilerplate** and is not production-ready unless you properly configure it, particularly for secret management using **HashiCorp Vault**.

---

## Features

* Real-time chat using WebSockets
* JWT-based authentication
* Token generation with custom claims (Access & Refresh tokens)
* Integration-ready with HashiCorp Vault for secure key management
* Basic separation of concerns between services, handlers, and routes

---

## Requirements

* Go 1.23+
* [HashiCorp Vault](https://www.vaultproject.io/)
* A Unix-based shell (tested on Linux/macOS)

---

## Installation & Setup

1. **Clone the repository**

```bash
git clone https://github.com/your-username/websocket-chat-backend.git
cd websocket-chat-backend
```

2. **Initialize Go modules**

```bash
go mod tidy
```

3. **Install HashiCorp Vault**

Follow the instructions at [https://developer.hashicorp.com/vault/docs/install](https://developer.hashicorp.com/vault/docs/install)

Example for Linux:

```bash
wget https://releases.hashicorp.com/vault/1.15.2/vault_1.15.2_linux_amd64.zip
unzip vault_1.15.2_linux_amd64.zip
sudo mv vault /usr/local/bin/
```

4. **Start Vault (in development mode)**

```bash
vault server -dev
```

Copy the **Root Token** printed in the terminal.

5. **Export the Vault token in your shell**

```bash
export VAULT_TOKEN="<your-root-token>"
```

6. **(Optional) Store a private key in Vault for signing JWTs**

```bash
vault kv put secret/jwt/private key=@private.pem
```

7. **Run the server**

```bash
make build
make run
```

---

## Notes

* This is a **learning-only** codebase. The cryptography, authentication logic, and error handling are **simplified**.
* For serious use, implement:

  * Input validation
  * Logging
  * Better session handling
  * Persistent storage (PostgreSQL, Redis, etc.)
  * Proper TLS and secure WebSocket (`wss://`) support

* My online server has all these features that this mock codebase does not have currently

---

## License

This project is licensed under the MIT License.

