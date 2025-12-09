/**
 * Member model for the library management system
 */
package models

type Member struct {
	ID   int
	Name string
	BorrowedBooks []Book
}
