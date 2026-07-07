# go-server

A minimal HTTP/HTTPS server written in Go that exposes a small key/value style API backed by a JSON file "database". Data is stored per numeric ID and persisted to disk as JSON, with no external database dependency.

## Features

- Simultaneous HTTP (`:80`) and HTTPS (`:443`) listeners
- Simple file-backed JSON table implementation (`DataBase.JsonTable`) with thread-safe reads/writes via mutex
- Generic, reusable table type (`JsonTable[T Indexable]`) — any type implementing `GetId() uint32` can be persisted
- Small REST-style API for storing arbitrary string payloads keyed by ID

## Project Structure

```
go-server/
├── main.go                     # HTTP handlers, server bootstrap
├── go.mod
├── DataBase/
│   ├── Models.go                # Indexable interface, IndexedType, AppData model
│   ├── JsonTable.go             # Generic JSON-file-backed table with CRUD helpers
│   └── Tables/
│       └── Tables.go            # Table instances (e.g. AppData table)
└── Helpers/
    ├── Convert/
    │   └── Convert.go           # String/number conversion helpers
    └── File/
        └── File.go              # Filesystem helpers (exists, read, write, mkdir)
```

## Requirements

- Go 1.24.3 or newer

## Running

```bash
cd go-server
go run .
```

By default the server listens on:

- `http://localhost:80`
- `https://localhost:443` (only if `_cert.pem` and `_key.pem` are present in the working directory)

If the TLS certificate/key files are not found next to the binary, the HTTPS listener is skipped and a message is logged — the HTTP listener still starts normally.

### Generating a self-signed certificate (for local testing)

```bash
openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout _key.pem -out _cert.pem -days 365 \
  -subj "/CN=localhost"
```

### Building a binary

```bash
go build -o go-server .
./go-server
```

## Data Storage

On first run, a `db/` directory is created next to the executable's working directory. Each table is persisted as a JSON array in its own file, e.g. `db/AppData.json`.

## API Reference

All endpoints are served under `/api/data`.

| Method | Endpoint          | Query Params | Description                                      |
|--------|-------------------|--------------|---------------------------------------------------|
| GET    | `/`                | —            | Health check, returns `It Works!`                 |
| POST   | `/api/data/add`    | `id` (uint32) | Stores the request body as data for the given ID (creates or overwrites) |
| GET    | `/api/data/get`    | `id` (uint32) | Returns the stored data for the given ID, or an error message if not found |
| GET    | `/api/data/clear`  | —            | Removes all stored entries                        |
| GET    | `/api/data/ids`    | —            | Returns a list of all stored IDs                  |

### Examples

```bash
# Store data under id=1
curl -X POST "http://localhost:80/api/data/add?id=1" -d "hello world"

# Retrieve it
curl "http://localhost:80/api/data/get?id=1"

# List all stored ids
curl "http://localhost:80/api/data/ids"

# Clear all data
curl "http://localhost:80/api/data/clear"
```

## Notes

- `id` values are parsed as unsigned 32-bit integers; invalid or missing values default to `0`.
- The JSON table implementation is generic and can be reused for other data models by implementing the `DataBase.Indexable` interface (a `GetId() uint32` method).
