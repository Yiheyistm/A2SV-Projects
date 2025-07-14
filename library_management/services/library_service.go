package services

import (
	"errors"
	"fmt"
	"library_management/models"
	"strconv"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID string) (string, error)
	BorrowBook(memberID, bookID string) (string, error)
	ReturnBook(memberID, bookID string) (string, error)
	ListAvailableBooks() ([]models.Book, error)
	ListBorrowedBooks(memberID string) ([]models.Book, error)
}

type LibraryService struct {
	Books   map[int]models.Book
	Members map[int]models.Member
}

func NewLibraryService() *LibraryService {
	// Initialize the library service with empty maps for books and members
	// sample data can be added here if needed
	// For example:
	books := map[int]models.Book{
		1: {ID: 1, Title: "1984", Author: "George Orwell", Status: "available"},
		2: {ID: 2, Title: "Mockingbird", Author: "Harper Lee", Status: "available"},
	}
	members := map[int]models.Member{
		1: {ID: 1, Name: "John Doe", BorrowedBooks: []models.Book{}},
		2: {ID: 2, Name: "Jane Smith", BorrowedBooks: []models.Book{}},
	}
	return &LibraryService{
		Books:   books,
		Members: members,
	}
}

func (s *LibraryService) AddBook(book models.Book) {
	s.Books[book.ID] = book
}

func (s *LibraryService) RemoveBook(bookID int) (string, error) {

	if _, exists := s.Books[bookID]; exists {
		delete(s.Books, bookID)
		return fmt.Sprintf("Book with ID: %v removed.\n", bookID), nil
	}
	return "", errors.New("Book with ID: " + strconv.Itoa(bookID) + " not found")
}

func (s *LibraryService) BorrowBook(memberID, bookID int) (string, error) {
	if book, exist := s.Books[bookID]; exist && book.Status == "available" {
		book.Status = "borrowed"
		s.Books[bookID] = book
		if member, exist := s.Members[memberID]; exist {
			member.BorrowedBooks = append(member.BorrowedBooks, book)
			s.Members[memberID] = member
			// Update the book status in the member's borrowed list
			for i, borrowedBook := range member.BorrowedBooks {
				if borrowedBook.ID == bookID {
					member.BorrowedBooks[i].Status = "borrowed"
					break
				}
			}
			return fmt.Sprintln("Borrowed Successfully"), nil
		} else {
			return "", errors.New("your member ID is not found, Please register first")
		}
	}
	return "", errors.New("Book with ID: " + strconv.Itoa(bookID) + " not found")
}

func (s *LibraryService) ReturnBook(memberID, bookID int) (string, error) {

	if book, exist := s.Books[bookID]; exist && book.Status == "borrowed" {
		book.Status = "available"
		s.Books[bookID] = book
		if member, exist := s.Members[memberID]; exist {

			for k, borrowedBook := range member.BorrowedBooks {
				if borrowedBook.ID == bookID {
					member.BorrowedBooks = append(member.BorrowedBooks[:k], member.BorrowedBooks[k+1:]...)
					for i, borrowedBook := range member.BorrowedBooks {
						if borrowedBook.ID == bookID {
							member.BorrowedBooks[i].Status = "available"
							break
						}
					}
					s.Members[memberID] = member
					return fmt.Sprintln("Successfully returned, Thank You!!"), nil
				}
			}
			return "", errors.New("book with ID: " + strconv.Itoa(bookID) + " is not found from your borrowed list, Please Check Again!!!")

		}
		return "", errors.New("your member ID is not found")
	}
	return "", errors.New("Book with ID: " + strconv.Itoa(bookID) + " not borrowed")
}

func (s *LibraryService) ListAvailableBooks() ([]models.Book, error) {
	var availableBooks []models.Book
	for _, book := range s.Books {
		if book.Status == "available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks, nil
}

func (s *LibraryService) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	if member, exist := s.Members[memberID]; exist {
		if len(member.BorrowedBooks) > 0 {
			return member.BorrowedBooks, nil
		}
		return nil, errors.New("you are not borrowed anything, Thank You")
	}
	return nil, errors.New("your Member ID is not found")
}
func (s *LibraryService) RegisterMember(member models.Member) {
	s.Members[member.ID] = member
}
