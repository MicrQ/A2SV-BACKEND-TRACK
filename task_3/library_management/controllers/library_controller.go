package controllers

import (
	"fmt"
	"library_management/models"
	"library_management/services"
)



func Console(lib *services.Library) {
	for {

		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books by Member")
		fmt.Println("7. Exit")

		var choice int
		fmt.Print("\nEnter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var id int
			var title, author string
			fmt.Print("Enter Book ID: ")
			fmt.Scan(&id)
			fmt.Print("Enter Book Title: ")
			fmt.Scan(&title)
			fmt.Print("Enter Book Author: ")
			fmt.Scan(&author)
			book := models.Book{ID: id, Title: title, Author: author, Status: "Available"}
			lib.AddBook(book)
			fmt.Println("Book added successfully!")

		case 2:
			var id int
			fmt.Print("Enter Book ID to remove: ")
			fmt.Scan(&id)
			lib.RemoveBook(id)
			fmt.Println("Book removed successfully!")
			
		case 3:
			var bookID, memberID int
			fmt.Print("Enter Book ID to borrow: ")
			fmt.Scan(&bookID)
			fmt.Print("Enter Member ID: ")
			fmt.Scan(&memberID)
			err := lib.BorrowBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book borrowed successfully!")
			}

		case 4:
			var bookID, memberID int
			fmt.Print("Enter Book ID to return: ")
			fmt.Scan(&bookID)
			fmt.Print("Enter Member ID: ")
			fmt.Scan(&memberID)
			err := lib.ReturnBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book returned successfully!")
			}
			
		case 5:
			books := lib.ListAvailableBooks()
			fmt.Println("Available Books:")
			for _, book := range books {
				fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
			}
			if len(books) == 0 {
				fmt.Println("No available books.")
			}

		case 6:
			var memberID int
			fmt.Print("Enter Member ID: ")
			fmt.Scan(&memberID)
			books := lib.ListBorrowedBooks(memberID)
			fmt.Println("Borrowed Books:")
			for _, book := range books {
				fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
			}
			if len(books) == 0 {
				fmt.Println("No books borrowed by this member.")
			}

		case 7:
			fmt.Println("Exiting From MicrQ Library... Goodbye!")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
