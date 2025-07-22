# Task Manager Clean Architecture

## Overview

Task Manager Clean Architecture is a Go application that helps users manage tasks efficiently. The project is structured using the Clean Architecture pattern, ensuring separation of concerns, testability, and scalability.

## Project Structure

```
task_manager_clean_arch/
├── cmd/
│   └── api/
│       └── main.go                # Application entry point
├── config/
│   ├── config.go                  # Configuration loading
│   └── env.go                     # Environment variable helpers
├── internal/
│   ├── app/
│   │   └── app.go                 # Application setup
│   ├── domain/
│   │   ├── task.go                # Task entity
│   │   ├── user.go                # User entity
│   │   ├── task_repo.go           # Task repository interface
│   │   └── user_repo.go           # User repository interface
│   ├── infrastructure/
│   │   ├── database/
│   │   │   └── db.go              # Database connection
│   │   ├── persistence/
│   │   │   ├── task_repo.go       # Task repository implementation
│   │   │   └── user_repo.go       # User repository implementation
│   │   └── security/
│   │       ├── jwt_service.go     # JWT handling
│   │       └── password_service.go# Password hashing
│   ├── interfaces/
│   │   ├── http/
│   │   │   ├── task_handler.go    # Task HTTP handlers
│   │   │   ├── user_handler.go    # User HTTP handlers
│   │   │   └── router.go         # HTTP router setup
│   │   └── middleware/
│   │       └── auth.go           # Authentication middleware
│   └── usecase/
│       ├── task_usecase.go        # Task business logic
│       └── user_usercase.go       # User business logic
├── tmp/                           # Temporary build files
├── go.mod                         # Go module definition
├── go.sum                         # Go dependencies checksum
├── .air.toml                      # Air live reload configuration
└── README.md                      # Project documentation
```

## Clean Architecture Layers

- **Domain**: Core business entities and repository interfaces.
- **Usecase**: Application-specific business rules.
- **Interfaces**: Adapters for HTTP handlers, middleware, etc.
- **Infrastructure**: External technologies (DB, JWT, etc.).

## Setup Instructions

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yiheyistm/task_manager_clean_arch.git
   cd task_manager_clean_arch
   ```

2. **Install Dependencies**

   Ensure you have Go installed (1.20+ recommended), then run:

   ```bash
   go mod tidy
   ```

3. **Run the Application**

   To start the application:

   ```bash
   go run ./cmd/api/main.go
   ```

   Or use [Air](https://github.com/cosmtrek/air) for live reloading during development:

   ```bash
   go install github.com/cosmtrek/air@latest
   air
   ```

## Usage

- The application exposes RESTful HTTP endpoints for managing tasks and users.
- Use tools like Postman or curl to interact with the API.
- See the `docs/documentation.md` file for detailed API documentation and endpoint descriptions.

## Configuration

- Environment variables and configuration files are managed in the `config/` directory.
- Edit `.env` or use environment variables as needed for your setup.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
