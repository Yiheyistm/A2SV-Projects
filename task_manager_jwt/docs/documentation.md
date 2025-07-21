# Task Manager Mongo

## Overview

Task Manager Mongo is a RESTful API built with Go and Gin, using MongoDB as the database. It provides endpoints for user registration, authentication, and task management, supporting both user and admin roles.

## Features

- User registration and login with JWT authentication
- Role-based access control (user/admin)
- CRUD operations for tasks
- User-specific and admin-level task management
- Task statistics and filtering

## Folder Structure

```
├── config/         # Database and environment configuration
├── controllers/    # HTTP handlers for users and tasks
├── middleware/     # Authentication and authorization middleware
├── models/         # Data models for users and tasks
├── routes/         # API route definitions
├── services/       # Business logic for users and tasks
├── utils/          # Utility functions (e.g., JWT parsing)
```

## Installation

1. Clone the repository:
   ```bash
   git clone <your-repo-url>
   cd task_manager_mongo
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set environment variables (or use a `.env` file):

   - `MONGO_URI` (default: mongodb://localhost:27017)
   - `MONGO_DB_NAME` (your database name)
   - `MONGO_COLLECTION_NAME` (default: users)
   - `JWT_SECRET` (your secret key)

4. Run the server:
   ```bash
   go run main.go
   ```

## API Endpoints

### Public

- `POST /api/v1/users/register` — Register a new user
- `POST /api/v1/users/login` — Login and receive JWT

### Authenticated (User)

- `GET /api/v1/users/:username/tasks` — List user's tasks
- `GET /api/v1/users/:username/tasks/:id` — Get a specific user task
- `POST /api/v1/users/:username/tasks` — Create a user task
- `PUT /api/v1/users/:username/tasks/:id` — Update a user task
- `DELETE /api/v1/users/:username/tasks/:id` — Delete a user task
- `GET /api/v1/users/:username/tasks/stats` — Get user's task stats

### Admin Only

- `GET /api/v1/tasks` — List all tasks
- `GET /api/v1/tasks/:id` — Get any task by ID
- `POST /api/v1/tasks` — Create a task
- `PUT /api/v1/tasks/:id` — Update a task
- `DELETE /api/v1/tasks/:id` — Delete a task
- `GET /api/v1/tasks/stats` — Get task statistics
- `GET /api/v1/users` — List all users
- `GET /api/v1/users/:username` — Get user by username

## Authentication & Authorization

- JWT-based authentication is required for all endpoints except registration and login.
- Admin endpoints require the user to have a `role` of `admin` in their JWT claims.
