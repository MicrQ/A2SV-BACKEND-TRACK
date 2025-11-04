/**
 * Interface definitions
 */
package services

import (
	"errors"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}


type Library struct {
	Books 	map[int]models.Book
	Members map[int]models.Member
}


func (lib *Library) AddBook(book models.Book) {
	lib.Books[book.ID] = book
}


func (lib *Library) RemoveBook(bookID int) {
	delete(lib.Books, bookID)
}


func (lib *Library) BorrowBook(bookID int, memberID int) error {
	// borrowing a book
	book, exists := lib.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	member, exists := lib.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}

	book.Status = "Borrowed"
	lib.Books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	lib.Members[memberID] = member

	return nil
}


func (l *Library) ReturnBook(bookID int, memberID int) error {
	// returns a book from a member
	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	member, exists := l.Members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	// Checkinigg if the book was borrowed
	found := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("member did not borrow this book")
	}

	book.Status = "Available"
	l.Books[bookID] = book
	l.Members[memberID] = member
	return nil
}


func (l *Library) ListAvailableBooks() []models.Book {
	// returns a list of available books
	var available []models.Book
	for _, book := range l.Books {
		if book.Status == "Available" {
			available = append(available, book)
		}
	}
	return available
}


func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.Members[memberID]
	if !exists {
		return nil
	}
	return member.BorrowedBooks
}
