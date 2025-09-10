# Flat-File Order Records

This project is a simple Go application for managing order records. It uses a binary flat-file (`orders.db`) as its database and provides both a Command-Line Interface (CLI) and a Graphical User Interface (GUI) for interaction.

## Features

*   **Insert**: Add new order records.
*   **Update**: Modify existing records by their `OrderID`.
*   **Read All**: Display all records in the database.
*   **Search**: Find specific records by their `OrderID`.
*   **Dual Interface**: Can be run as a standard CLI application or as a Fyne-based GUI application.

## How to Run

First, ensure you have Go installed.

### Build the Application

Run the following command in the project's root directory to build the executable:

```bash
go build .
```

### Run the CLI

To run the command-line interface:

```bash
./flatfile-order-records
```

### Run the GUI

To run the graphical user interface, use the `-gui` flag:

```bash
./flatfile-order-records -gui
```

## Concurrency and Race Conditions

This section addresses how the application handles concurrent operations.

### Within a Single Application Instance

If multiple threads (goroutines) within a single running application call the database functions simultaneously, race conditions are prevented by a `sync.RWMutex`.

*   `d.mu.Lock()` and `d.mu.Unlock()` in `Insert` and `UpdateById` ensure that only one goroutine can write to the file at a time.
*   `d.mu.RLock()` and `d.mu.RUnlock()` in `ReadAll` and `SearchById` allow multiple goroutines to read the file concurrently but will wait if a write lock is held.

This provides a robust in-process concurrency control mechanism.

### Running Multiple Application Instances

**There is a major race condition if you run multiple instances of the application (e.g., two GUI windows) simultaneously.**

The `sync.RWMutex` only works for goroutines *within the same process*. It cannot coordinate access between separate, independent processes. If two instances run at the same time, they can read and write to the `orders.db` file concurrently, leading to **lost updates** and potential **data corruption**.

To solve this, a more advanced inter-process locking mechanism would be required, such as using file locks provided by the operating system (e.g., `flock` on Unix-like systems).