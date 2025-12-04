# Task Manager API Documentation

A RESTful API for managing tasks, built with Go and Gin following **Clean Architecture** principles. This API supports CRUD operations with MongoDB persistence, JWT authentication, and role-based access control.

## Architecture Overview

The project follows Clean Architecture with these layers:

```
task_manager/
├── Delivery/           # HTTP handlers and routing
│   ├── main.go
│   ├── controllers/
│   └── routers/
├── Domain/             # Core business entities
│   └── domain.go
├── Infrastructure/     # External services (JWT, password hashing)
│   ├── auth_middleWare.go
│   ├── jwt_service.go
│   └── password_service.go
├── Repositories/       # Data access interfaces and implementations
│   ├── task_repository.go
│   └── user_repository.go
├── Usecases/           # Business logic
│   ├── task_usecases.go
│   └── user_usecases.go
└── Tests/              # Test suites
    ├── mocks/
    ├── infrastructure/
    ├── usecases/
    ├── middleware/
    ├── controllers/
    └── repositories_integration/
```

## Base URL
```
http://localhost:8080
```

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection URI | `mongodb://localhost:27017` |
| `DB_NAME` | Database name | `taskmanager` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-secret-key` |

Example setup:
```bash
export MONGODB_URI="mongodb://localhost:27017"
export DB_NAME="taskmanager"
export JWT_SECRET="your-super-secret-key"
```

---

## Authentication

The API uses JWT (JSON Web Tokens) for authentication.

### Roles
- **admin**: Full access (create, update, delete tasks; promote users)
- **user**: Read-only access to tasks

> **Note**: The first registered user automatically becomes an admin.

### Auth Endpoints

#### Register
- **POST /register**
- **Description:** Create a new user account
- **Request Body:**
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```
- **Response:**
```json
201 Created
{
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "john_doe",
    "role": "user",
    "created_at": "2025-12-04T10:00:00Z"
  }
}
```

#### Login
- **POST /login**
- **Description:** Authenticate and receive a JWT token
- **Request Body:**
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```
- **Response:**
```json
200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Promote User (Admin only)
- **POST /promote/:id**
- **Auth:** Required (Admin)
- **Description:** Promote a user to admin role
- **Response:** `204 No Content`

---

## Task Endpoints

All task endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### 1. Get All Tasks
- **GET /tasks**
- **Auth:** Required
- **Description:** Retrieve a list of all tasks.
- **Response:**
```json
200 OK
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Buy groceries",
      "description": "Milk, eggs, bread",
      "due_date": "2025-11-30T00:00:00Z",
      "status": "pending"
    }
  ]
}
```

---

### 2. Get Task By ID
- **GET /tasks/:id**
- **Auth:** Required
- **Description:** Retrieve details of a single task by its ID.
- **Response:**
```json
200 OK
{
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "due_date": "2025-11-30T00:00:00Z",
    "status": "pending"
  }
}
```
- **Error Response:**
```json
404 Not Found
{
  "error": "task not found"
}
```

---

### 3. Create Task
- **POST /tasks**
- **Auth:** Required
- **Description:** Create a new task.
- **Valid Status Values:** `pending`, `in_progress`, `completed`
- **Request Body:**
```json
{
  "title": "New Task",
  "description": "Random description",
  "due_date": "2025-12-28T00:00:00Z",
  "status": "pending"
}
```
- **Response:**
```json
201 Created
{
  "id": "507f1f77bcf86cd799439012",
  "title": "New Task",
  "description": "Random description",
  "due_date": "2025-12-28T00:00:00Z",
  "status": "pending"
}
```
- **Error Response:**
```json
400 Bad Request
{
  "error": "invalid request body"
}
```

---

### 4. Update Task
- **PUT /tasks/:id**
- **Auth:** Required
- **Description:** Update an existing task.
- **Request Body:**
```json
{
  "title": "Buy groceries (updated)",
  "description": "Milk, eggs, bread, bananas",
  "due_date": "2025-11-30T00:00:00Z",
  "status": "in_progress"
}
```
- **Response:**
```json
200 OK
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Buy groceries (updated)",
  "description": "Milk, eggs, bread, bananas",
  "due_date": "2025-11-30T00:00:00Z",
  "status": "in_progress"
}
```
- **Error Response:**
```json
404 Not Found
{
  "error": "task not found"
}
```

---

### 5. Delete Task
- **DELETE /tasks/:id**
- **Auth:** Required (Admin only)
- **Description:** Delete a specific task.
- **Response:**
```
204 No Content
```
- **Error Response:**
```json
404 Not Found
{
  "error": "task not found"
}
```

---

## Error Handling
- All error responses are returned in JSON format with an `error` field.
- Common HTTP status codes:

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Success |
| 201 | Created - Resource created |
| 204 | No Content - Success with no response body |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error |

## Notes
- All date/time fields use RFC3339 format (e.g., `2025-11-30T00:00:00Z`).
- The API uses MongoDB for persistent data storage; data persists across server restarts.
- Task IDs are MongoDB ObjectIDs represented as hexadecimal strings.
- JWT tokens expire after 24 hours.
- First registered user is automatically assigned the `admin` role.

---

## Testing

The project includes comprehensive tests using [Testify](https://github.com/stretchr/testify).

### Test Structure
```
Tests/
├── mocks/                      # Mock repositories for unit tests
├── infrastructure/             # JWT and password service tests
├── usecases/                   # Business logic tests
├── middleware/                 # Auth middleware tests
├── controllers/                # HTTP handler tests
└── repositories_integration/   # MongoDB integration tests
```

### Run All Unit Tests
```bash
SKIP_INTEGRATION=true go test ./Tests/... -v
```

### Run Integration Tests
Requires MongoDB running:
```bash
go test ./Tests/repositories_integration -v
```

### Run Tests with Coverage
```bash
# Generate coverage report
go test ./Tests/... -coverprofile=coverage.out -covermode=atomic

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Test Categories

| Category | Description | Command |
|----------|-------------|---------|
| Infrastructure | JWT, password hashing | `go test ./Tests/infrastructure -v` |
| Usecases | Business logic with mocks | `go test ./Tests/usecases -v` |
| Middleware | Auth and role checks | `go test ./Tests/middleware -v` |
| Controllers | HTTP handlers | `go test ./Tests/controllers -v` |
| Integration | MongoDB operations | `go test ./Tests/repositories_integration -v` |

---

## Example cURL Requests

### Register a User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Create Task (with auth)
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"New Task","description":"Test","due_date":"2025-12-28T00:00:00Z","status":"pending"}'
```

### Get All Tasks (with auth)
```bash
curl http://localhost:8080/tasks \
  -H "Authorization: Bearer <your-token>"
```

### Update Task
```bash
curl -X PUT http://localhost:8080/tasks/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Task","description":"Updated","due_date":"2025-11-30T00:00:00Z","status":"completed"}'
```

### Delete Task (admin only)
```bash
curl -X DELETE http://localhost:8080/tasks/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <admin-token>"
```

### Copyright © 2025

