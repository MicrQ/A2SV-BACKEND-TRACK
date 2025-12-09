/**
 * Interface definitions
 */
package services

import (
	"errors"
	"fmt"
	"library_management/models"
	"sync"
	"time"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	ReserveBook(bookID int, memberID int) error
}

type ReservationRequest struct {
	BookID   int
	MemberID int
	Response chan error
}

type Library struct {
	Books           map[int]models.Book
	Members         map[int]models.Member
	Mutex           sync.Mutex
	ReservationChan chan ReservationRequest
	CancelChans     map[int]chan bool
}

func (lib *Library) AddBook(book models.Book) {
	lib.Mutex.Lock()
	defer lib.Mutex.Unlock()
	lib.Books[book.ID] = book
}

func (lib *Library) RemoveBook(bookID int) {
	lib.Mutex.Lock()
	defer lib.Mutex.Unlock()
	delete(lib.Books, bookID)
}

func (lib *Library) BorrowBook(bookID int, memberID int) error {
	lib.Mutex.Lock()
	defer lib.Mutex.Unlock()

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

	if book.Status == "Reserved" && book.ReservedBy != memberID {
		return errors.New("book reserved by another member")
	}

	// If reserved by this member, cancel the auto-cancel
	if book.Status == "Reserved" && book.ReservedBy == memberID {
		if cancel, ok := lib.CancelChans[bookID]; ok {
			cancel <- true
			delete(lib.CancelChans, bookID)
		}
	}

	book.Status = "Borrowed"
	book.ReservedBy = 0
	book.ReservedAt = time.Time{}
	lib.Books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	lib.Members[memberID] = member

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

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
	book.ReservedBy = 0
	book.ReservedAt = time.Time{}
	l.Books[bookID] = book
	l.Members[memberID] = member
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

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
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	member, exists := l.Members[memberID]
	if !exists {
		return nil
	}
	return member.BorrowedBooks
}

func (lib *Library) ReserveBook(bookID int, memberID int) error {
	req := ReservationRequest{
		BookID:   bookID,
		MemberID: memberID,
		Response: make(chan error, 1),
	}
	lib.ReservationChan <- req
	return <-req.Response
}

func (lib *Library) StartWorker() {
	go lib.reservationWorker()
}

func (lib *Library) reservationWorker() {
	for req := range lib.ReservationChan {
		lib.processReservation(req)
	}
}

func (lib *Library) processReservation(req ReservationRequest) {
	lib.Mutex.Lock()
	defer lib.Mutex.Unlock()

	book, exists := lib.Books[req.BookID]
	if !exists {
		req.Response <- errors.New("book not found")
		return
	}

	if book.Status != "Available" {
		req.Response <- errors.New("book not available for reservation")
		return
	}

	// Reserve the book
	book.Status = "Reserved"
	book.ReservedBy = req.MemberID
	book.ReservedAt = time.Now()
	lib.Books[req.BookID] = book

	// Start auto-cancel goroutine
	cancelChan := make(chan bool, 1)
	lib.CancelChans[req.BookID] = cancelChan
	go func(bookID int, memberID int, cancel chan bool) {
		select {
		case <-time.After(5 * time.Second):
			lib.Mutex.Lock()
			if b, ex := lib.Books[bookID]; ex && b.Status == "Reserved" && b.ReservedBy == memberID {
				b.Status = "Available"
				b.ReservedBy = 0
				b.ReservedAt = time.Time{}
				lib.Books[bookID] = b
				fmt.Println("Reservation for book", bookID, "auto-cancelled")
			}
			lib.Mutex.Unlock()
		case <-cancel:
			// Reservation cancelled (borrowed)
		}
		delete(lib.CancelChans, bookID)
	}(req.BookID, req.MemberID, cancelChan)

	req.Response <- nil
}
