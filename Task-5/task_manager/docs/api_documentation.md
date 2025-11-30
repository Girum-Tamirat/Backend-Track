# Task Manager API — Authentication & Authorization

Base URL: http://localhost:8080

## Auth

### POST /register
Body:
```json
{
  "username":"alice",
  "password":"secret"
}
```
- The **first user registered automatically becomes an admin**.

### POST /login
Body:
```json
{
  "username":"alice",
  "password":"secret"
}
```
Response: `{ "token": "<JWT>", "username": "...", "role": "admin" }`

**Usage for Protected Endpoints:**
Use the header `Authorization: Bearer <JWT>` for all protected requests.

## Endpoints (Protected)

| Method | Route | Description | Required Role |
| :--- | :--- | :--- | :--- |
| `GET` | `/tasks` | List all tasks | Authenticated User |
| `GET` | `/tasks/:id` | Get single task by ID | Authenticated User |
| `POST` | `/tasks` | Create a new task | **Admin Only** |
| `PUT` | `/tasks/:id` | Update a task | **Admin Only** |
| `DELETE` | `/tasks/:id` | Delete a task | **Admin Only** |
| `POST` | `/users/:username/promote` | Promote a user to admin | **Admin Only** |

## Security Notes
- Passwords are securely hashed using **bcrypt**.
- JWT is signed using the `JWT_SECRET` environment variable (defaults to a placeholder).
- Tokens expire in 24 hours.
