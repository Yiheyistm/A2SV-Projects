# Library Management System Documentation

This document provides an overview of the Library Management System, outlining its functionality, setup, usage, and testing procedures.

## Overview

The Library Management System is a Go-based application that allows users to manage a library's books and members. Key features include:
- Viewing available books.
- Adding and removing books.
- Borrowing and returning books.
- Registering library members.

## Project Structure

```
library_management/
├── controllers/
│   └── library_controller_test.go  # Contains unit tests for library actions
├── docs/
│   └── documentation.md            # This documentation file
├── models/
│   └── (model files for Book, Member, etc.)
├── services/
│   └── (service implementation for library operations)
└── main.go                         # Application entry point
```

## Setup

1. **Prerequisites:**
    - Go installed in your system.
    - Proper environment setup for your Go workspace.

2. **Build the Application:**
    - Navigate to the project directory.
    - Run `go build` to compile the project.
    - Alternatively, run the project with `go run main.go`.

## Usage

### Main Application

The application starts in `main.go` by calling the `ShowMenu` function from the controllers package:
```go
func main() {
    controllers.ShowMenu()
    fmt.Println("Thank you for using the Library Management System!")
}
```
Users can interact with the system via the console menu which allows selection of operations like adding, removing, borrowing, or returning books.

### Library Operations

- **Add Book:** Use the `AddBook` function in the library service to add new books.
- **Remove Book:** Use the `RemoveBook` function to remove a book by its ID.
- **Borrow & Return Book:** Functions `BorrowBook` and `ReturnBook` manage the book lending process.
- **List Books:** Functions `ListAvailableBooks` and `ListBorrowedBooks` display the current status of the library's inventory.

## Testing

The project includes unit tests within the `controllers/library_controller_test.go` file. These tests cover:
- Adding a book
- Removing a book
- Borrowing and returning a book
- Listing available and borrowed books
- Registering a member

Run the tests using:
```
go test ./...
```

## Additional Information

For complex interactions involving user input, consider reviewing testing strategies using Go packages like `os` or `bufio` to simulate console input and output.

This documentation will be updated as new features or changes are implemented in the Library Management System.
