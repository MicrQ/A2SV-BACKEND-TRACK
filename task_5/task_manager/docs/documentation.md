# Task Manager API Documentation

A RESTful API for managing tasks, built with Go and Gin. This API supports basic CRUD operations using MongoDB for persistent data storage.

## Base URL
```
http://localhost:8080
```

---

## MongoDB Configuration

The API uses MongoDB for data persistence. Configure the following environment variables:

- `MONGODB_URI`: MongoDB connection URI (default: `mongodb://localhost:27017`)
- `DB_NAME`: Database name (default: `taskmanager`)
- `COLLECTION_NAME`: Collection name (default: `tasks`)

Example:
```bash
export MONGODB_URI="mongodb://localhost:27017"
export DB_NAME="taskmanager"
export COLLECTION_NAME="tasks"
```

---

## Endpoints

### 1. Get All Tasks
- **GET /tasks**
- **Description:** Retrieve a list of all tasks.
- **Response:**
```json
200 OK
[
  {
    "id": "507f1f77bcf86cd799439011",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "due_date": "2025-11-30T00:00:00Z",
    "status": "pending"
  }
]
```

---

### 2. Get Task By ID
- **GET /tasks/:id**
- **Description:** Retrieve details of a single task by its ID.
- **Response:**
```json
200 OK
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "due_date": "2025-11-30T00:00:00Z",
  "status": "pending"
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
- **Description:** Create a new task.
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
- Common HTTP status codes: 200 (OK), 201 (Created), 204 (No Content), 400 (Bad Request), 404 (Not Found), 500 (Internal Server Error)

## Notes
- All date/time fields use RFC3339 format (e.g., `2025-11-30T00:00:00Z`).
- The API uses MongoDB for persistent data storage; data persists across server restarts.
- Task IDs are MongoDB ObjectIDs represented as hexadecimal strings.

## Example cURL Requests

- **Create Task:**
```bash
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"title":"New Task","description":"Test","due_date":"2025-12-28T00:00:00Z","status":"pending"}'
```

- **Get All Tasks:**
```bash
curl http://localhost:8080/tasks
```

- **Update Task:**
```bash
curl -X PUT http://localhost:8080/tasks/507f1f77bcf86cd799439011 -H "Content-Type: application/json" -d '{"title":"Updated Task","description":"Updated","due_date":"2025-11-30T00:00:00Z","status":"completed"}'
```

- **Delete Task:**
```bash
curl -X DELETE http://localhost:8080/tasks/507f1f77bcf86cd799439011
```

### copyright© 2025

