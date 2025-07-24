# 🌍 Task Manager API Documentation

Welcome to the **Task Manager API**, built with Clean Architecture in Go. This RESTful API provides endpoints for user and task management, secured with JWT authentication. The documentation is designed to be clear, developer-friendly.

---

## 📜 Table of Contents

- [Overview](#-overview)
- [Authentication](#-authentication)
- [Endpoints](#-endpoints)
  - [User Endpoints](#-user-endpoints)
  - [Task Endpoints (Admin Only)](#-task-endpoints-admin-only)
- [Error Handling](#-error-handling)
- [Environment Variables](#️-environment-variables)
  - [Example .env](#-example-env)
- [Usage Examples](#-usage-examples)
- [Notes](#-notes)

---

## 🌿 Overview

The Task Manager API follows **Clean Architecture** principles, ensuring separation of concerns, testability, and independence from external systems. It supports user registration, login, task creation, and management, with role-based access (user and admin). The API is built with Go, uses MongoDB for persistence, and is secured with JWT tokens.

---

## 🔒 Authentication

All protected endpoints require a **Bearer JWT token** in the `Authorization` header. Tokens are obtained via the `/api/v1/users/login` endpoint. Admin-only endpoints require a token with the `admin` role.

**Header Format:**

```
Authorization: Bearer <your_jwt_token>
```

**Token Structure:**

- **Access Token:** Short-lived (e.g., 2 hours), used for API requests.
- **Refresh Token:** Long-lived (e.g., 168 hours), used to obtain new access tokens.

---

## 🛠 Endpoints

All endpoints are prefixed with `/api/v1`. Responses are in JSON format, and dates use ISO 8601 (e.g., `2025-07-23T15:17:00Z`).

---

### 🌞 User Endpoints

#### Register a New User

`POST /api/v1/users/register`
Creates a new user account.

**Request Body:**

```json
{
  "username": "abebe",
  "email": "abebe@example.com",
  "password": "selam123",
  "role": "user"
}
```

**Response:**
`201 Created`

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c2",
  "username": "abebe",
  "email": "abebe@example.com",
  "role": "user"
}
```

---

#### Login

`POST /api/v1/users/login`
Authenticates a user and returns a JWT token.

**Request Body:**

```json
{
  "identifier": "abebe", // Username or email
  "password": "selam123"
}
```

**Response:**
`200 OK`

```json
{
  "access_token": "<jwt_access_token>",
  "refresh_token": "<jwt_refresh_token>"
}
```

---

#### Get User Profile

`GET /api/v1/users/:username`
Retrieves a user’s profile. **Requires authentication.**

**Headers:**

```
Authorization: Bearer <user_token>
```

**Response:**
`200 OK`

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c2",
  "username": "abebe",
  "email": "abebe@example.com",
  "role": "user"
}
```

---

#### Get User's Tasks

`GET /api/v1/users/:username/tasks`
Lists all tasks for a specific user. **Requires authentication.**

**Headers:**

```
Authorization: Bearer <user_token>
```

**Response:**
`200 OK`

```json
{
  "tasks": [
    {
      "id": "60f7c2b8e1b2c2a1b8e1b2c3",
      "title": "Visit Merkato",
      "description": "Buy coffee and spices",
      "due_date": "2025-07-24T12:00:00Z",
      "status": "pending",
      "created_by": "abebe"
    }
  ]
}
```

---

#### Create a Task for User

`POST /api/v1/users/:username/tasks`
Creates a new task for a user. **Requires authentication.**

**Headers:**

```
Authorization: Bearer <user_token>
```

**Request Body:**

```json
{
  "title": "Plan Addis Ababa trip",
  "description": "Book hotel and transport",
  "due_date": "2025-07-30T17:00:00Z",
  "status": "pending"
}
```

**Response:**
`201 Created`

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c4",
  "title": "Plan Addis Ababa trip",
  "description": "Book hotel and transport",
  "due_date": "2025-07-30T17:00:00Z",
  "status": "pending",
  "created_by": "abebe"
}
```

---

#### Update a User's Task

`PUT /api/v1/users/:username/tasks/:id`
Updates a specific task. **Requires authentication.**

**Headers:**

```
Authorization: Bearer <user_token>
```

**Request Body:**

```json
{
  "title": "Plan Addis Ababa trip",
  "description": "Book hotel, transport, and tour",
  "due_date": "2025-07-31T17:00:00Z",
  "status": "completed"
}
```

**Response:**
`200 OK`

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c4",
  "title": "Plan Addis Ababa trip",
  "description": "Book hotel, transport, and tour",
  "due_date": "2025-07-31T17:00:00Z",
  "status": "completed",
  "created_by": "abebe"
}
```

---

#### Delete a User's Task

`DELETE /api/v1/users/:username/tasks/:id`
Deletes a specific task. **Requires authentication.**

**Headers:**

```
Authorization: Bearer <user_token>
```

**Response:**
`204 No Content`

```json
{
  "message": "Task deleted successfully"
}
```

---

#### Get User Task Statistics

`GET /api/v1/users/:username/tasks/stats`
Returns task statistics by status for a user. **Requires authentication.**

**Headers:**

```
Authorization: Bearer <user_token>
```

**Response:**
`200 OK`

```json
[
  {
    "status": "pending",
    "count": 3
  },
  {
    "status": "completed",
    "count": 1
  }
]
```

---

### ❤️ Task Endpoints (Admin Only)

These endpoints require a token with the **admin** role.

#### List All Tasks

`GET /api/v1/tasks`
Lists all tasks in the system.

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Response:**
`200 OK`

```json
[
  {
    "id": "60f7c2b8e1b2c2a1b8e1b2c3",
    "title": "Visit Merkato",
    "description": "Buy coffee and spices",
    "due_date": "2025-07-24T12:00:00Z",
    "status": "pending",
    "created_by": "abebe"
  }
]
```

---

#### Get Task by ID

`GET /api/v1/tasks/:id`
Retrieves a specific task.

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Response:**
`200 OK`

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c3",
  "title": "Visit Merkato",
  "description": "Buy coffee and spices",
  "due_date": "2025-07-24T12:00:00Z",
  "status": "pending",
  "created_by": "abebe"
}
```

---

#### Create Task

`POST /api/v1/tasks`
Creates a new task (admin-assigned).

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Request Body:**

```json
{
  "title": "Organize community event",
  "description": "Plan a cultural festival in Addis",
  "due_date": "2025-08-01T10:00:00Z",
  "status": "pending",
  "created_by": "admin"
}
```

**Response:**
`201 Created`

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c5",
  "title": "Organize community event",
  "description": "Plan a cultural festival in Addis",
  "due_date": "2025-08-01T10:00:00Z",
  "status": "pending",
  "created_by": "admin"
}
```

---

#### Update Task

`PUT /api/v1/tasks/:id`
Updates a specific task.

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Request Body:**

```json
{
  "title": "Organize community event",
  "description": "Plan a cultural festival with music",
  "due_date": "2025-08-01T10:00:00Z",
  "status": "completed",
  "created_by": "admin"
}
```

**Response:**
`200 OK`

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c5",
  "title": "Organize community event",
  "description": "Plan a cultural festival with music",
  "due_date": "2025-08-01T10:00:00Z",
  "status": "completed",
  "created_by": "admin"
}
```

---

#### Delete Task

`DELETE /api/v1/tasks/:id`
Deletes a specific task.

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Response:**
`204 No Content`

```json
{
  "message": "Task deleted successfully"
}
```

---

#### Get Task Statistics

`GET /api/v1/tasks/stats`
Returns task statistics across all users.

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Response:**
`200 OK`

```json
[
  {
    "status": "pending",
    "count": 10
  },
  {
    "status": "completed",
    "count": 5
  }
]
```

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

| Variable                  | Description                       | Default/Example                 |
| ------------------------- | --------------------------------- | ------------------------------- |
| DB_HOST_URI               | MongoDB connection URI            | mongodb+srv://user:pass@host/db |
| APP_ENV                   | Application environment           | development                     |
| SERVER_ADDRESS            | Server address and port           | :8080                           |
| CONTEXT_TIMEOUT           | Request context timeout (seconds) | 2                               |
| DB_HOST                   | MongoDB host                      | mongodb                         |
| DB_PORT                   | MongoDB port                      | 27017                           |
| DB_TASK_COLLECTION        | Task collection name              | tasks                           |
| DB_USER_COLLECTION        | User collection name              | users                           |
| DB_PASS                   | MongoDB password                  | password                        |
| DB_NAME                   | MongoDB database name             | go-mongo                        |
| ACCESS_TOKEN_EXPIRY_HOUR  | Access token expiry (hours)       | 2                               |
| REFRESH_TOKEN_EXPIRY_HOUR | Refresh token expiry (hours)      | 168                             |
| ACCESS_TOKEN_SECRET       | JWT secret for access tokens      | access_token_secret             |
| REFRESH_TOKEN_SECRET      | JWT secret for refresh tokens     | refresh_token_secret            |

---

### 📋 Example .env

```env
DB_HOST_URI=mongodb+srv://user:pass@host/go-mongo
APP_ENV=development
SERVER_ADDRESS=:8080
CONTEXT_TIMEOUT=2
DB_HOST=mongodb
DB_PORT=27017
DB_TASK_COLLECTION=tasks
DB_USER_COLLECTION=users
DB_PASS=password
DB_NAME=go-mongo
ACCESS_TOKEN_EXPIRY_HOUR=2
REFRESH_TOKEN_EXPIRY_HOUR=168
ACCESS_TOKEN_SECRET=access_token_secret
REFRESH_TOKEN_SECRET=refresh_token_secret
```

---

## 🛠 Usage Examples

Below are example API calls using `curl`, showcasing common operations.

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
   -H "Content-Type: application/json" \
   -d '{"username":"abebe","email":"abebe@example.com","password":"selam123","role":"user"}'
```

**Response:**

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c2",
  "username": "abebe",
  "email": "abebe@example.com",
  "role": "user"
}
```

---

### Login

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
   -H "Content-Type: application/json" \
   -d '{"identifier":"abebe","password":"selam123"}'
```

**Response:**

```json
{
  "access_token": "<jwt_access_token>",
  "refresh_token": "<jwt_refresh_token>"
}
```

---

### Create a Task

```bash
curl -X POST http://localhost:8080/api/v1/users/abebe/tasks \
   -H "Content-Type: application/json" \
   -H "Authorization: Bearer <jwt_access_token>" \
   -d '{"title":"Plan Addis Ababa trip","description":"Book hotel and transport","due_date":"2025-07-30T17:00:00Z","status":"pending"}'
```

**Response:**

```json
{
  "id": "60f7c2b8e1b2c2a1b8e1b2c4",
  "title": "Plan Addis Ababa trip",
  "description": "Book hotel and transport",
  "due_date": "2025-07-30T17:00:00Z",
  "status": "pending",
  "created_by": "abebe"
}
```

---

## 🌄 Notes

- **Clean Architecture:** The API is structured with layers (`domain`, `usecase`, `interfaces`, `infrastructure`), ensuring maintainability and testability.
- **Security:** JWT tokens are validated by the `internal/infrastructure/security/jwt_service.go` module.
- **Database:** MongoDB is used, with collections specified in `DB_TASK_COLLECTION` and `DB_USER_COLLECTION`.

---
