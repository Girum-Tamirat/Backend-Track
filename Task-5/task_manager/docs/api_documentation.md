# Task Manager (Clean Architecture)

## Overview
This project implements a Task Management API using Clean Architecture:
- Domain: core entities (Task, User)
- Repositories: interfaces + Mongo implementations
- Usecases: business logic
- Infrastructure: JWT, password hashing, middleware
- Delivery: Gin controllers and routers

## Quick Run
1. Start MongoDB locally or set MONGO_URI.
2. (Optional) export JWT_SECRET
3. go mod tidy
4. go run Delivery/main.go

## Endpoints
- POST /register
- POST /login
- GET /tasks (auth)
- GET /tasks/:id (auth)
- POST /tasks (admin)
- PUT /tasks/:id (admin)
- DELETE /tasks/:id (admin)
- POST /users/:username/promote (admin)

Auth: Authorization header `Bearer <token>` returned from /login.

