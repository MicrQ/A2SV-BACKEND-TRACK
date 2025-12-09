/**
 * Book model for the library management system
 */
package models

import "time"

type Book struct {
	ID         int
	Title      string
	Author     string
	Status     string
	ReservedBy int       // Member ID who reserved the book
	ReservedAt time.Time // Time when reserved
}
