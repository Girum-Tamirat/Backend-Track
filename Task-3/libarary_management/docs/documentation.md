# Concurrent Book Reservation System

## Overview
This project extends the Library Management System by adding concurrency for handling book reservations.

## Key Features
- **Goroutines** process reservation requests simultaneously.
- **Channels** queue reservation requests.
- **Mutex** ensures safe updates to shared data (books and members).
- **Auto-cancellation** after 5 seconds if not borrowed.

## How It Works
1. `main.go` initializes the library and starts the `LibraryController`.
2. The `controller` launches multiple goroutines simulating different members reserving books at the same time.
3. Each reservation request is sent to a **channel**.
4. A **worker goroutine** continuously reads from this channel and processes reservations.
5. The `Library` struct uses a **sync.Mutex** to lock shared data during updates.
6. If a book remains "Reserved" for 5 seconds, a goroutine automatically resets its status to "Available".

## Run
```bash
cd library_management
go run main.go
