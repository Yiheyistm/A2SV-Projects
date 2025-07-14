package helper

import (
	"library_management/models"
	"testing"
)

func TestShowAvailableBooks(t *testing.T) {
	books := []models.Book{
		{ID: 1, Title: "1984", Author: "George Orwell"},
		{ID: 2, Title: "Mockingbird", Author: "Harper Lee"},
	}
	ShowAvailableBooks(books)
}

func TestShowBorrowedBooks(t *testing.T) {
	books := []models.Book{
		{ID: 1, Title: "1984", Author: "George Orwell"},
		{ID: 2, Title: "Mockingbird", Author: "Harper Lee"},
	}
	ShowBorrowedBooks(books)
}

func TestClearScreen(t *testing.T) {
	ClearScreen()
}
