package helper

import (
	"fmt"
	"library_management/models"
	"time"
)

func ShowWelcomeMessage() {
	fmt.Println("Welcome to the Library Management System!")
	fmt.Println("This system allows you to borrow and return books, and manage your library membership.")
	fmt.Println("Please follow the prompts to interact with the system.")
}

func ShowErrorMessage(err error) {
	if err != nil {
		fmt.Println("❌ Error:", err)
	} else {
		fmt.Println("❌ An unknown error occurred.")
	}
}
func ShowSuccessMessage(message string) {
	fmt.Println("✅ Success:", message)
}

func ShowAvailableBooks(books []models.Book) {
	if len(books) == 0 {
		fmt.Println("No available books at the moment.")
		return
	}
	fmt.Println("Available Books:")
	fmt.Println("No.\t\tID\t\tTitle\t\t\t\tAuthor")
	for i, book := range books {
		fmt.Printf("%v\t\t%v\t\t%v\t\t\t\t%v\n", i+1, book.ID, book.Title, book.Author)
	}
}
func ShowBorrowedBooks(books []models.Book) {
	if len(books) == 0 {
		fmt.Println("You have not borrowed any books.")
		return
	}
	fmt.Println("Your Borrowed Books:")
	fmt.Println("No.\t\tID\t\tTitle\t\t\t\tAuthor")
	for i, book := range books {
		fmt.Printf("%v\t\t%v\t\t%v\t\t\t\t%v\n", i+1, book.ID, book.Title, book.Author)
	}
}
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
	time.Sleep(100 * time.Millisecond)
}
