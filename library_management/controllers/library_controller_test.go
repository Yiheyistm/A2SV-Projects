package controllers

import (
	"library_management/models"
	"library_management/services"
	"testing"
)

func TestShowMenu(t *testing.T) {
	// This function is interactive and requires user input.
	// Testing it would require mocking user input and output.
	// Consider using a library like "os" or "bufio" for advanced testing.
}

func TestAddBook(t *testing.T) {
	service := services.NewLibraryService()
	book := models.Book{ID: 3, Title: "New Book", Author: "Author Name", Status: "available"}
	service.AddBook(book)
	if _, exists := service.Books[book.ID]; !exists {
		t.Errorf("Book was not added successfully")
	}
}

func TestRemoveBook(t *testing.T) {
	service := services.NewLibraryService()
	message, err := service.RemoveBook(1)
	if err != nil || message == "" {
		t.Errorf("Failed to remove book: %v", err)
	}
}

func TestBorrowBook(t *testing.T) {
	service := services.NewLibraryService()
	message, err := service.BorrowBook(1, 1)
	if err != nil || message == "" {
		t.Errorf("Failed to borrow book: %v", err)
	}
}

func TestReturnBook(t *testing.T) {
	service := services.NewLibraryService()
	service.BorrowBook(1, 1)
	message, err := service.ReturnBook(1, 1)
	if err != nil || message == "" {
		t.Errorf("Failed to return book: %v", err)
	}
}

func TestListAvailableBooks(t *testing.T) {
	service := services.NewLibraryService()
	books, err := service.ListAvailableBooks()
	if err != nil || len(books) == 0 {
		t.Errorf("Failed to list available books: %v", err)
	}
}

func TestListBorrowedBooks(t *testing.T) {
	service := services.NewLibraryService()
	service.BorrowBook(1, 1)
	books, err := service.ListBorrowedBooks(1)
	if err != nil || len(books) == 0 {
		t.Errorf("Failed to list borrowed books: %v", err)
	}
}

func TestRegisterMember(t *testing.T) {
	service := services.NewLibraryService()
	member := models.Member{ID: 3, Name: "New Member"}
	service.RegisterMember(member)
	if _, exists := service.Members[member.ID]; !exists {
		t.Errorf("Member was not registered successfully")
	}
}
