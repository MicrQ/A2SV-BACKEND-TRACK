# Library Management System (Go)

## Overview

The Library Management System is a console-based application built in Go that demonstrates fundamental programming concepts such as structs, interfaces, methods, maps, and slices. It provides a simple yet effective way to manage a library's book inventory and member borrowing activities through a text-based user interface.

This project serves as an educational example for backend development in Go, showcasing object-oriented principles, error handling, and data management using built-in Go data structures.

## Features

- **Book Management**: Add and remove books from the library inventory
- **Member Operations**: Borrow and return books with member tracking
- **Inventory Queries**: List available books and view borrowed books by member
- **Console Interface**: Interactive menu-driven user experience
- **Error Handling**: Robust validation for operations like borrowing non-existent books

## Architecture

### Key Components

- **Models** (`models/`):
  - `Book`: Represents a book with ID, title, author, and status
  - `Member`: Represents a library member with ID, name, and list of borrowed books

- **Services** (`services/`):
  - `LibraryManager` interface: Defines the contract for library operations
  - `Library` struct: Implements the library management logic using maps for efficient data storage

- **Controllers** (`controllers/`):
  - `Console` function: Handles user interaction through a command-line menu

### Data Structures

- Books are stored in a `map[int]models.Book` for O(1) access by ID
- Members are stored in a `map[int]models.Member` for efficient lookup
- Borrowed books are tracked as slices within each member

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
3. **Borrow Book**: Enter book ID and member ID to borrow a book
4. **Return Book**: Enter book ID and member ID to return a borrowed book
5. **List Available Books**: Display all books currently available
6. **List Borrowed Books by Member**: Enter member ID to see their borrowed books
7. **Exit**: Quit the application

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
**November 2025.**
