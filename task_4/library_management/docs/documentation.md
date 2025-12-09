# Library Management System (Go)

## Overview

The Library Management System is a console-based application built in Go that demonstrates fundamental programming concepts such as structs, interfaces, methods, maps, and slices. It provides a simple yet effective way to manage a library's book inventory and member borrowing activities through a text-based user interface.

This project serves as an educational example for backend development in Go, showcasing object-oriented principles, error handling, data management using built-in Go data structures, and concurrent programming with Goroutines, Channels, and Mutexes.

## Features

- **Book Management**: Add and remove books from the library inventory
- **Member Operations**: Borrow and return books with member tracking
- **Reservation System**: Reserve books concurrently with auto-cancellation
- **Inventory Queries**: List available books and view borrowed books by member
- **Console Interface**: Interactive menu-driven user experience
- **Concurrency**: Handles multiple reservation requests safely using Goroutines and Channels
- **Error Handling**: Robust validation for operations like borrowing non-existent books

## Architecture

### Key Components

- **Models** (`models/`):
  - `Book`: Represents a book with ID, title, author, status, reservation info
  - `Member`: Represents a library member with ID, name, and list of borrowed books

- **Services** (`services/`):
  - `LibraryManager` interface: Defines the contract for library operations
  - `Library` struct: Implements the library management logic using maps for efficient data storage, with concurrency support

- **Controllers** (`controllers/`):
  - `Console` function: Handles user interaction through a command-line menu

### Data Structures

- Books are stored in a `map[int]models.Book` for O(1) access by ID
- Members are stored in a `map[int]models.Member` for efficient lookup
- Borrowed books are tracked as slices within each member
- Reservation requests are queued using a buffered channel
- Mutex protects shared data from race conditions

## Concurrency Implementation

The system uses Go's concurrency primitives to handle multiple reservation requests:

- **Goroutines**: Used to process reservation requests asynchronously and handle auto-cancellation timers
- **Channels**: A buffered channel queues incoming reservation requests to the worker goroutine
- **Mutex**: Protects the library's shared state (books and members maps) from concurrent access

### Reservation Flow

1. User requests to reserve a book
2. Request is sent to the reservation channel
3. Worker goroutine processes the request, reserving the book if available
4. A timer goroutine starts, auto-cancelling the reservation after 5 seconds if not borrowed
5. If the user borrows the reserved book within 5 seconds, the timer is cancelled

This ensures thread-safe operations and prevents race conditions when multiple users try to reserve the same book simultaneously.

## Prerequisites

- Go 1.16 or later installed on your system
- Basic understanding of command-line interfaces

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/MicrQ/A2SV-BACKEND-TRACK.git
   cd A2SV-BACKEND-TRACK/task_3/library_management
   ```

2. Initialize Go modules:
   ```bash
   go mod tidy
   ```

## Build Commands

### Build the Application

To compile the application into an executable:

```bash
go build -o library_app main.go
```

This creates a binary named `library_app` in the current directory.


## Usage

### Running the Application

Start the application using one of the following methods:

**Option 1: Direct run (recommended for development)**
```bash
go run main.go
```

**Option 2: Build and run**
```bash
go build -o library_app main.go
./library_app
```

### Console Menu Options

Upon running, you'll see a welcome message and the main menu:

1. **Add Book**: Enter book ID, title, and author to add a new book
2. **Remove Book**: Enter book ID to remove from inventory
3. **Reserve Book**: Enter book ID and member ID to reserve a book (auto-cancels in 5 seconds)
4. **Borrow Book**: Enter book ID and member ID to borrow a book (can borrow reserved books)
5. **Return Book**: Enter book ID and member ID to return a borrowed book
6. **List Available Books**: Display all books currently available
7. **List Borrowed Books by Member**: Enter member ID to see their borrowed books
8. **Exit**: Quit the application

### Example Usage

```
===== Welcome to MicrQ Library =====

1. Add Book
2. Remove Book
3. Borrow Book
4. Return Book
5. List Available Books
6. List Borrowed Books by Member
7. Exit

Enter your choice: 1
Enter Book ID: 101
Enter Book Title: The MicrQ
Enter Book Author: Abenet Gebre
Book added successfully!
```

## License

This project is part of the A2SV Backend Track and is intended for educational purposes.
---
**December 2025.**
