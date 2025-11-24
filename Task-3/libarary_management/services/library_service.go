package services

import (
	"errors"
	"sync"
	"time"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	ReserveBook(bookID int, memberID int) error
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
	mu      sync.Mutex // Protects Books and Members
}

func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()

	book.Status = "Available"
	l.Books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.Books[bookID]; !exists {
		return errors.New("book not found")
	}
	delete(l.Books, bookID)
	return nil
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if book.Status != "Available" && book.Status != "Reserved" {
		return errors.New("book not available for borrowing")
	}

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Books[bookID] = book
	l.Members[memberID] = member
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("member not found")
	}
	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	// Remove from borrowed list
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			break
		}
	}

	book.Status = "Available"
	l.Members[memberID] = member
	l.Books[bookID] = book
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	available := []models.Book{}
	for _, book := range l.Books {
		if book.Status == "Available" {
			available = append(available, book)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	if member, ok := l.Members[memberID]; ok {
		return member.BorrowedBooks
	}
	return nil
}

// ReserveBook handles concurrent reservation logic
func (l *Library) ReserveBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}
	if book.Status != "Available" {
		return errors.New("book already borrowed or reserved")
	}

	book.Status = "Reserved"
	l.Books[bookID] = book
	go func() {
		time.Sleep(5 * time.Second)
		l.mu.Lock()
		defer l.mu.Unlock()

		if l.Books[bookID].Status == "Reserved" {
			book := l.Books[bookID]
			book.Status = "Available"
			l.Books[bookID] = book
		}
	}()
	return nil
}
