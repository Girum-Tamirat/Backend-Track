# 🚀 Task Management REST API (Go & Gin Framework)

This repository contains a simple Task Management REST API developed using **Go** and the **Gin** web framework. The API supports basic **CRUD** (Create, Read, Update, Delete) operations for managing tasks using an in-memory data store.

## 🎯 Objective

The primary goal of this task was to implement a fully functional RESTful API backend, focusing on:
* Implementing all required endpoints (`GET`, `POST`, `PUT`, `DELETE`).
* Using an **in-memory database** (`map` protected by `sync.RWMutex`) for data storage.
* Ensuring proper HTTP status codes and error handling.
* Adhering to Go best practices and a clean folder structure.

## 📁 Project Structure

The project follows the requested organizational structure:

\`\`\`
task_manager/
├── main.go               # Application entry point
├── controllers/
│   └── task_controller.go  # HTTP request handlers
├── models/
│   └── task.go             # Task data structure definition
├── data/
│   └── task_service.go     # Business logic and in-memory data store
├── router/
│   └── router.go           # Gin router setup and route registration
└── go.mod                # Dependency management
\`\`\`

## ⚙️ Prerequisites

* **Go** (version 1.18 or higher recommended)
* **Postman** or **cURL** for testing the endpoints

## 🛠️ How to Run

1.  Navigate into the \`task_manager\` directory (you should already be there).
2.  Initialize the Go module and download the required dependencies (Gin):
    \`\`\`bash
    go mod tidy
    \`\`\`
3.  Run the application:
    \`\`\`bash
    go run main.go
    \`\`\`
    The server will start and listen on \`http://localhost:8080\`.

## 📜 API Endpoints Documentation

The API utilizes the base URL: \`http://localhost:8080/api/v1\`.

### 1. Model Definition (\`models/task.go\`)

Tasks require \`title\` and \`status\` fields.

| Field | Type | Description | Required | Example |
| :--- | :--- | :--- | :--- | :--- |
| \`id\` | \`int\` | Unique identifier (auto-generated) | No | \`1\` |
| \`title\` | \`string\` | Short name for the task | **Yes** | \`"Submit Final Code"\` |
| \`description\` | \`string\` | Detailed explanation | No | \`"Review and push to GitHub."\` |
| \`due_date\` | \`string\` | ISO-8601 formatted date/time | No | \`"2025-12-01T23:59:59Z"\` |
| \`status\` | \`string\` | Task state | **Yes** | \`"pending"\`, \`"done"\` |

### 2. Endpoints Summary

| Method | Endpoint | Description | Status Codes |
| :--- | :--- | :--- | :--- |
| \`POST\` | \`/api/v1/tasks\` | Creates a new task. | \`201 Created\`, \`400 Bad Request\` |
| \`GET\` | \`/api/v1/tasks\` | Retrieves a list of all tasks. | \`200 OK\` |
| \`GET\` | \`/api/v1/tasks/:id\` | Retrieves a specific task by ID. | \`200 OK\`, \`400 Bad Request\`, \`404 Not Found\` |
| \`PUT\` | \`/api/v1/tasks/:id\` | Updates an existing task by ID. | \`200 OK\`, \`400 Bad Request\`, \`404 Not Found\` |
| \`DELETE\` | \`/api/v1/tasks/:id\` | Deletes a task by ID. | \`204 No Content\`, \`400 Bad Request\`, \`404 Not Found\` |

### 3. Example Request (cURL)

**Request:** Create a task using cURL.

\`\`\`bash
curl -X POST http://localhost:8080/api/v1/tasks \
 -H "Content-Type: application/json" \
 -d '{"title":"Test Readme","description":"Verify documentation","status":"pending"}'
\`\`\`

**Response (201 Created):**

\`\`\`json
{
  "id": 1,
  "title": "Test Readme",
  "description": "Verify documentation",
  "status": "pending"
}
\`\`\`

## ⚠️ Important Note

This application uses an **in-memory store**. All tasks created will be **lost** when the server is stopped and restarted.
