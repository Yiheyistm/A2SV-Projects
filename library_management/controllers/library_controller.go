package controllers

import (
	"bufio"
	"errors"
	"fmt"
	"library_management/helper"
	"library_management/models"
	"library_management/services"
	"os"
	"strings"
)

func ShowMenu() {
	defer fmt.Println("Exiting the Library Management System. Goodbye!")
	helper.ShowWelcomeMessage()
	service := services.NewLibraryService()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n************** Library Management System Menu **************")

		fmt.Println("1. Add a Book")
		fmt.Println("2. Remove a Book")
		fmt.Println("3. Borrow a Book")
		fmt.Println("4. Return a Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Member Registration")
		fmt.Println("8. Clear Screen")
		fmt.Println("0. Exit")
		fmt.Print("Please select an option: \n")
		var choice int
		fmt.Scan(&choice)
		switch choice {
		case 1:
			reader.ReadString('\n')
			fmt.Println("Please enter the book details:")
			var book models.Book
			book.ID = len(service.Books) + 1
			fmt.Print("Title: ")
			book.Title, _ = reader.ReadString('\n')
			book.Title = strings.TrimSpace(book.Title)
			fmt.Print("Author: ")
			book.Author, _ = reader.ReadString('\n')
			book.Author = strings.TrimSpace(book.Author)
			book.Status = "available"
			service.AddBook(book)
			helper.ShowSuccessMessage("Book added successfully!")
		case 2:
			var bookID int
			fmt.Print("Enter the book ID to remove: ")
			fmt.Scan(&bookID)
			message, err := service.RemoveBook(bookID)
			if err != nil {
				helper.ShowErrorMessage(err)
			} else {
				helper.ShowSuccessMessage(message)
			}
		case 3:
			var memberID, bookID int
			fmt.Print("Enter your member ID: ")
			fmt.Scan(&memberID)
			fmt.Print("Enter the book ID to borrow: ")
			fmt.Scan(&bookID)
			message, err := service.BorrowBook(memberID, bookID)
			if err != nil {
				helper.ShowErrorMessage(err)
			} else {
				helper.ShowSuccessMessage(message)
			}
		case 4:
			var memberID, bookID int
			fmt.Print("Enter your member ID: ")
			fmt.Scan(&memberID)
			fmt.Print("Enter the book ID to return: ")
			fmt.Scan(&bookID)
			message, err := service.ReturnBook(memberID, bookID)
			if err != nil {
				helper.ShowErrorMessage(err)
			} else {
				helper.ShowSuccessMessage(message)
			}
		case 5:
			availableBooks, err := service.ListAvailableBooks()
			if err != nil {
				helper.ShowErrorMessage(err)
			} else {
				helper.ShowAvailableBooks(availableBooks)
			}
		case 6:
			var memberID int
			fmt.Print("Enter your member ID: ")
			fmt.Scan(&memberID)
			if memberID <= 0 {
				helper.ShowErrorMessage(errors.New("invalid member ID"))
				continue
			}
			borrowedBooks, err := service.ListBorrowedBooks(memberID)
			if err != nil {
				helper.ShowErrorMessage(err)
			} else {
				helper.ShowBorrowedBooks(borrowedBooks)
			}
		case 7:
			reader.ReadString('\n')
			var member models.Member
			fmt.Println("Please enter your details for registration:")
			member.ID = len(service.Members) + 1
			fmt.Print("Name: ")
			member.Name, _ = reader.ReadString('\n')
			member.Name = strings.TrimSpace(member.Name)
			service.RegisterMember(member)
			helper.ShowSuccessMessage("Member registered successfully!")
		case 8:
			helper.ClearScreen()
		case 0:
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

