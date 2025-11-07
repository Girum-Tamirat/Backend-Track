package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

type LibraryController struct {
	Service *services.Library
}

func NewLibraryController(service *services.Library) *LibraryController {
	return &LibraryController{Service: service}
}

func (c *LibraryController) Start() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- Library Management System ---")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books by Member")
		fmt.Println("7. Exit")
		fmt.Print("Enter your choice: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, _ := strconv.Atoi(choiceStr)

		switch choice {
		case 1:
			c.addBook(reader)
		case 2:
			c.removeBook(reader)
		case 3:
			c.borrowBook(reader)
		case 4:
			c.returnBook(reader)
		case 5:
			c.listAvailableBooks()
		case 6:
			c.listBorrowedBooks(reader)
		case 7:
			fmt.Println("Exiting system...")
			return
		default:
			fmt.Println("Invalid choice, try again.")
		}
	}
}

func (c *LibraryController) addBook(reader *bufio.Reader) {
	fmt.Print("Enter Book ID: ")
	id, _ := strconv.Atoi(readLine(reader))
	fmt.Print("Enter Title: ")
	title := readLine(reader)
	fmt.Print("Enter Author: ")
	author := readLine(reader)

	book := models.Book{ID: id, Title: title, Author: author}
	c.Service.AddBook(book)
	fmt.Println("Book added successfully!")
}

func (c *LibraryController) removeBook(reader *bufio.Reader) {
	fmt.Print("Enter Book ID to remove: ")
	id, _ := strconv.Atoi(readLine(reader))
	if err := c.Service.RemoveBook(id); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Book removed successfully!")
}

func (c *LibraryController) borrowBook(reader *bufio.Reader) {
	fmt.Print("Enter Member ID: ")
	memberID, _ := strconv.Atoi(readLine(reader))
	fmt.Print("Enter Member Name: ")
	name := readLine(reader)

	if _, exists := c.Service.Members[memberID]; !exists {
		c.Service.Members[memberID] = models.Member{ID: memberID, Name: name}
	}

	fmt.Print("Enter Book ID to borrow: ")
	bookID, _ := strconv.Atoi(readLine(reader))

	if err := c.Service.BorrowBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Book borrowed successfully!")
}

func (c *LibraryController) returnBook(reader *bufio.Reader) {
	fmt.Print("Enter Member ID: ")
	memberID, _ := strconv.Atoi(readLine(reader))
	fmt.Print("Enter Book ID to return: ")
	bookID, _ := strconv.Atoi(readLine(reader))

	if err := c.Service.ReturnBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Book returned successfully!")
}

func (c *LibraryController) listAvailableBooks() {
	books := c.Service.ListAvailableBooks()
	fmt.Println("\nAvailable Books:")
	for _, b := range books {
		fmt.Printf("ID: %d | Title: %s | Author: %s\n", b.ID, b.Title, b.Author)
	}
}

func (c *LibraryController) listBorrowedBooks(reader *bufio.Reader) {
	fmt.Print("Enter Member ID: ")
	memberID, _ := strconv.Atoi(readLine(reader))
	books := c.Service.ListBorrowedBooks(memberID)
	fmt.Printf("\nBorrowed Books by Member %d:\n", memberID)
	for _, b := range books {
		fmt.Printf("ID: %d | Title: %s | Author: %s\n", b.ID, b.Title, b.Author)
	}
}

func readLine(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
