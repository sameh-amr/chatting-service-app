# Chatting Service App

A real-time chat application built with Go (backend) and React (frontend). This project demonstrates a full-stack architecture with authentication, user management, and real-time messaging using WebSockets.

## Architecture Overview (basic)

```
+-------------------+         WebSocket/API         +-------------------+
|                   | <--------------------------> |                   |
|   React Frontend  |                             |     Go Backend     |
|  (Vite, TS, UI)   |                             | (REST, WS, Auth)   |
|                   | <------ REST API ----------> |                   |
+-------------------+                             +-------------------+
         |                                                    |
         |                                                    |
         |                                                    |
         v                                                    v
+-------------------+                             +-------------------+
|                   |                             |                   |
|   User Browser    |                             |   PostgreSQL DB   |
|                   |                             | (users, messages) |
+-------------------+                             +-------------------+
```


## Clean Architecture

This project follows a Clean Architecture approach for the Go backend, which separates concerns and makes the codebase easier to maintain and extend. The main layers are:

- **Models:** Core business entities (e.g., User, Message) in `backend/models/`.
- **Repository:** Data access logic (e.g., database queries) in `backend/repository/`.
- **Service:** Business logic and use cases in `backend/service/`.
- **HTTP Handlers:** API endpoints and request/response handling in `backend/httphandlers/`.
- **WebSocket:** Real-time communication logic in `backend/websocket/`.
- **Utils/DTO:** Utility functions and data transfer objects for validation, JWT, etc.

This separation allows for:
- Easier testing (mock repositories/services)
- Swappable infrastructure (e.g., change DB or transport)
- Clear boundaries between business logic and delivery (API/UI)

The frontend is organized by feature (contexts, components, pages) for clarity and maintainability.


- **Frontend:**
  - React (Vite) SPA
  - Communicates with backend via REST API for auth, user, and message operations
  - Uses WebSocket for real-time messaging
  - State managed with React Context
- **Backend:**
  - Go HTTP server (Gorilla Mux)
  - REST API for CRUD/auth
  - WebSocket server for real-time events
  - JWT authentication
  - GORM ORM for DB access
- **Database:**
  - PostgreSQL (or SQLite for dev)
  - Stores users, messages, sessions

## Tech Stack

- **Backend:** Go, Gorilla Mux, GORM, JWT, WebSockets
- **Frontend:** React, Vite, TypeScript, Tailwind CSS
- **Database:** PostgreSQL (or SQLite for dev)
- **Other:** Docker, Swagger (OpenAPI) for API docs

## Setup Instructions

### Prerequisites
- Go 1.20+
- Node.js 18+
- Docker (optional, for easy setup)

### Local Development

#### 1. Backend
```sh
cd backend
# Copy .env.example to .env and set DB connection if needed
# go run main.go
```

#### 2. Frontend
```sh
cd client
npm install
npm run dev
```

#### 3. Docker (Full stack)
```sh
docker-compose up --build
```

#### 4. API Docs
- Open `backend/swagger.yaml` in [Swagger Editor](https://editor.swagger.io/) or use Swagger UI.

## API Usage

- All endpoints are documented in [`backend/swagger.yaml`](backend/swagger.yaml).
- Example endpoints:
  - `POST /auth/signup` — Register a new user
  - `POST /auth/login` — Login and get JWT
  - `POST /auth/logout` — Logout
  - `GET /auth/users` — List all users (except self)
  - `POST /messages` — Send a message
  - `GET /messages?user1=...&user2=...` — Get messages between users
  - `POST /upload` — Upload a file (multipart/form-data)

You can import the Swagger file into Postman or use Swagger UI for interactive API testing.



## Known Limitations

- No group chat support (1:1 and broadcast only)
- No message search or advanced filtering
- No push notifications
- Minimal error handling on the frontend
- No production-grade security hardening (for demo/dev use)

---

For questions or contributions, open an issue or PR.
