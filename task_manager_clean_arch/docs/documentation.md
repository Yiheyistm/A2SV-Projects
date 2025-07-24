# 🌍 Task Manager Clean Architecture API Documentation

Welcome to the **Task Manager Clean Architecture API**. This RESTful API is built in Go, following Clean Architecture principles for maintainability, testability, and scalability. It provides endpoints for user and task management, secured with JWT authentication, and uses MongoDB for persistence.

---

## 📁 Project Structure

```
task_manager_clean_arch/
├── cmd/
│   └── api/
│       └── main.go                # Application entry point
├── config/
│   ├── config.go                  # Configuration loading
│   └── env.go                     # Environment variable helpers
├── docs/
│   └── documentation.md           # This documentation
├── internal/
│   ├── domain/                    # Core business entities and interfaces
│   │   ├── db.go
│   │   ├── refresh_token.go
│   │   ├── task.go
│   │   └── user.go
│   ├── infrastructure/            # External tech (DB, JWT, etc.)
│   │   ├── database/
│   │   │   ├── mongo_config.go
│   │   │   ├── user_entity.go
│   │   │   └── user_mapper.go
│   │   ├── persistence/
│   │   │   ├── task_repo.go
│   │   │   └── user_repo.go
│   │   └── security/
│   │       ├── jwt_service.go
│   │       └── password_service.go
│   ├── interfaces/
│   │   ├── http/
│   │   │   ├── dto/
│   │   │   │   ├── refresh_token_dto.go
│   │   │   │   ├── refresh_token_mapper.go
│   │   │   │   ├── task_dto.go
│   │   │   │   ├── task_mapper.go
│   │   │   │   ├── user_dto.go
│   │   │   │   └── user_mapper.go
│   │   │   ├── handler/
│   │   │   │   ├── refresh_token_handler.go
│   │   │   │   ├── task_handler.go
│   │   │   │   └── user_handler.go
│   │   │   └── router/
│   │   │       ├── auth_route.go
│   │   │       ├── refresh_token_route.go
│   │   │       ├── route.go
│   │   │       ├── task_route.go
│   │   │       └── user_route.go
│   │   └── middleware/
│   │       └── auth.go
│   └── usecase/
│       ├── refresh_token_usecase.go
│       ├── task_usecase.go
│       └── user_usercase.go
├── tmp/                           # Temporary build files
├── go.mod                         # Go module definition
├── go.sum                         # Go dependencies checksum
├── .air.toml                      # Air live reload configuration
├── .env                           # Environment variables
└── README.md                      # Project documentation
```

---

## 🏗️ Clean Architecture Layers

- **Domain:** Core business entities and repository interfaces.
- **Usecase:** Application-specific business rules.
- **Interfaces:** Adapters for HTTP handlers, middleware, DTOs, and routers.
- **Infrastructure:** External technologies (DB, JWT, password hashing, etc.).

---

## 🔒 Authentication

All protected endpoints require a **Bearer JWT token** in the `Authorization` header.
Tokens are obtained via the `/api/v1/users/login` endpoint.
Admin-only endpoints require a token with the `admin` role.

**Header Example:**

```
Authorization: Bearer <your_jwt_token>
```

- **Access Token:** Short-lived (e.g., 2 hours), used for API requests.
- **Refresh Token:** Long-lived (e.g., 168 hours), used to obtain new access tokens.

---

## 🛠 Endpoints

All endpoints are prefixed with `/api/v1`. Responses are in JSON format, and dates use ISO 8601.

### User Endpoints

#### Register a New User

- **POST** `/api/v1/users/register`
- **Body:**
  ```json
  {
    "username": "abebe",
    "email": "abebe@example.com",
    "password": "selam123",
    "role": "user"
  }
  ```
- **Response:** `201 Created`

#### Login

- **POST** `/api/v1/users/login`
- **Body:**
  ```json
  {
    "identifier": "abebe",
    "password": "selam123"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "access_token": "<jwt_access_token>",
    "refresh_token": "<jwt_refresh_token>"
  }
  ```

#### Refresh Token

- **POST** `/api/v1/users/refresh`
- **Body:**
  ```json
  {
    "refreshToken": "<jwt_refresh_token>"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "accessToken": "<new_access_token>",
    "refreshToken": "<new_refresh_token>"
  }
  ```

#### Get User Profile

- **GET** `/api/v1/users/:username`
- **Headers:** `Authorization: Bearer <user_token>`
- **Response:** `200 OK`

#### Get User's Tasks

- **GET** `/api/v1/users/:username/tasks`
- **Headers:** `Authorization: Bearer <user_token>`
- **Response:** `200 OK`

#### Get a User's Task by ID

- **GET** `/api/v1/users/:username/tasks/:id`
- **Headers:** `Authorization: Bearer <user_token>`
- **Response:** `200 OK`

#### Create a Task for User

- **POST** `/api/v1/users/:username/tasks`
- **Headers:** `Authorization: Bearer <user_token>`
- **Body:** (see TaskRequest in code)
- **Response:** `201 Created`

#### Update a User's Task

- **PUT** `/api/v1/users/:username/tasks/:id`
- **Headers:** `Authorization: Bearer <user_token>`
- **Body:** (see TaskRequest in code)
- **Response:** `200 OK`

#### Delete a User's Task

- **DELETE** `/api/v1/users/:username/tasks/:id`
- **Headers:** `Authorization: Bearer <user_token>`
- **Response:** `204 No Content`

#### Get User Task Statistics

- **GET** `/api/v1/users/:username/tasks/stats`
- **Headers:** `Authorization: Bearer <user_token>`
- **Response:** `200 OK`

---

### Task Endpoints (Admin Only)

All `/tasks` endpoints require admin privileges.

#### List All Tasks

- **GET** `/api/v1/tasks`
- **Headers:** `Authorization: Bearer <admin_token>`
- **Response:** `200 OK`

#### Get Task by ID

- **GET** `/api/v1/tasks/:id`
- **Headers:** `Authorization: Bearer <admin_token>`
- **Response:** `200 OK`

#### Create Task

- **POST** `/api/v1/tasks`
- **Headers:** `Authorization: Bearer <admin_token>`
- **Body:** (see TaskRequest in code)
- **Response:** `201 Created`

#### Update Task

- **PUT** `/api/v1/tasks/:id`
- **Headers:** `Authorization: Bearer <admin_token>`
- **Body:** (see TaskRequest in code)
- **Response:** `200 OK`

#### Delete Task

- **DELETE** `/api/v1/tasks/:id`
- **Headers:** `Authorization: Bearer <admin_token>`
- **Response:** `204 No Content`

#### Get Task Statistics

- **GET** `/api/v1/tasks/stats`
- **Headers:** `Authorization: Bearer <admin_token>`
- **Response:** `200 OK`

---

## 🚨 Error Handling

Errors are returned in JSON format with appropriate HTTP status codes.

**Error Format:**

```json
{
  "error": "Error message"
}
```

**Common Errors:**

| Status Code | Description                                 |
| ----------- | ------------------------------------------- |
| 400         | Bad Request: Invalid input data.            |
| 401         | Unauthorized: Missing or invalid JWT token. |
| 403         | Forbidden: Insufficient permissions.        |
| 404         | Not Found: Resource not found.              |
| 500         | Internal Server Error: Server-side issue.   |

**Example:**

```json
{
  "error": "User not found"
}
```

---

## ⚙️ Environment Variables

The API relies on environment variables for configuration, loaded via the `config` package. These are stored in a `.env` file or set in the environment.

| Variable                  | Description                       | Example/Default                 |
| ------------------------- | --------------------------------- | ------------------------------- |
| DB_HOST_URI               | MongoDB connection URI            | mongodb+srv://user:pass@host/db |
| APP_ENV                   | Application environment           | development                     |
| SERVER_ADDRESS            | Server address and port           | :8080                           |
| CONTEXT_TIMEOUT           | Request context timeout (seconds) | 2                               |
| DB_USER                   | MongoDB user                      | nicko                           |
| DB_HOST                   | MongoDB host                      | go-mongo                        |
| DB_PORT                   | MongoDB port                      | 27017                           |
| DB_TASK_COLLECTION        | Task collection name              | tasks                           |
| DB_USER_COLLECTION        | User collection name              | users                           |
| DB_PASS                   | MongoDB password                  | 123456                          |
| DB_NAME                   | MongoDB database name             | task_manager                    |
| ACCESS_TOKEN_EXPIRY_HOUR  | Access token expiry (hours)       | 2                               |
| REFRESH_TOKEN_EXPIRY_HOUR | Refresh token expiry (hours)      | 168                             |
| ACCESS_TOKEN_SECRET       | JWT secret for access tokens      | your_access_token_secret        |
| REFRESH_TOKEN_SECRET      | JWT secret for refresh tokens     | your_refresh_token_secret       |

### Example .env

```env
DB_HOST_URI=mongodb+srv://user:pass@host/task_manager
APP_ENV=development
SERVER_ADDRESS=:8080
CONTEXT_TIMEOUT=2
DB_USER=abeto
DB_HOST=go-mongo
DB_PORT=27017
DB_TASK_COLLECTION=tasks
DB_USER_COLLECTION=users
DB_PASS=qwe123
DB_NAME=task_manager
ACCESS_TOKEN_EXPIRY_HOUR=2
REFRESH_TOKEN_EXPIRY_HOUR=168
ACCESS_TOKEN_SECRET=your_access_token_secret
REFRESH_TOKEN_SECRET=your_refresh_token_secret
```

---

## 🛠 Usage Examples

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
   -H "Content-Type: application/json" \
   -d '{"username":"abebe","email":"abebe@example.com","password":"selam123","role":"user"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
   -H "Content-Type: application/json" \
   -d '{"identifier":"abebe","password":"selam123"}'
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/users/refresh \
   -H "Content-Type: application/json" \
   -d '{"refreshToken":"<jwt_refresh_token>"}'
```

### Create a Task

```bash
curl -X POST http://localhost:8080/api/v1/users/abebe/tasks \
   -H "Content-Type: application/json" \
   -H "Authorization: Bearer <jwt_access_token>" \
   -d '{"title":"Plan Addis Ababa trip","description":"Book hotel and transport","due_date":"2025-07-30T17:00:00Z","status":"pending"}'
```

---

## 🌄 Notes

- **Clean Architecture:** The API is structured with layers (`domain`, `usecase`, `interfaces`, `infrastructure`), ensuring maintainability and testability.
- **Security:** JWT tokens are validated by the `internal/infrastructure/security/jwt_service.go` module.
- **Database:** MongoDB is used, with collections specified in `DB_TASK_COLLECTION` and `DB_USER_COLLECTION`.
- **Live Reload:** Use `.air.toml` and [Air](https://github.com/cosmtrek/air) for development live reloading.
