# Task Manager API Documentation

A RESTful API for managing tasks, built with Go and Gin. This API supports basic CRUD operations using an in-memory database.

## Base URL
```
http://localhost:8080
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
    "id": 1,
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "due_date": "2025-11-30",
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
  "id": 1,
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "due_date": "2025-11-30",
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
  "due_date": "2025-12-28",
  "status": "pending"
}
```
- **Response:**
```json
201 Created
{
  "id": 2,
  "title": "New Task",
  "description": "Random description",
  "due_date": "2025-12-28",
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
  "due_date": "2025-11-30",
  "status": "in_progress"
}
```
- **Response:**
```json
200 OK
{
  "id": 1,
  "title": "Buy groceries (updated)",
  "description": "Milk, eggs, bread, bananas",
  "due_date": "2025-11-30",
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
- Common HTTP status codes: 200 (OK), 201 (Created), 204 (No Content), 400 (Bad Request), 404 (Not Found)

## Notes
- All date/time fields use ISO 8601 format (e.g., `2025-11-30`).
- The API uses an in-memory database; data will be lost when the server restarts.

## Example cURL Requests

- **Create Task:**
```bash
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"title":"New Task","description":"Test","due_date":"2025-12-28","status":"pending"}'
```

- **Get All Tasks:**
```bash
curl http://localhost:8080/tasks
```

- **Update Task:**
```bash
curl -X PUT http://localhost:8080/tasks/1 -H "Content-Type: application/json" -d '{"title":"Updated Task","description":"Updated","due_date":"2025-11-30","status":"completed"}'
```

- **Delete Task:**
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

### copyright© 2025

