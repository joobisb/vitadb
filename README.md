# VitaDB

**VitaDB** is an experimental, distributed, fault-tolerant database designed for both learning and potential production use. 

The focus of VitaDB is to explore and implement core concepts of distributed systems, such as durability, consistency, and fault tolerance, starting with basic features and progressively building toward a robust, scalable solution.

## Current Status

VitaDB is in the early stages of development. The initial focus is on implementing a Write-Ahead Log (WAL) for durability and ensuring data integrity during operations. 

As development progresses, additional features like transactional support, concurrency control, and performance optimizations will be introduced.

## Getting Started

Since VitaDB is still under active development, it is not yet ready for production use. However, you are welcome to explore the code, contribute, or use it as a learning resource.

## Running VitaDB Locally

To run VitaDB on your local machine, follow these steps:

1. **Prerequisites**
   - Go 1.21 or later
   - Make (optional, but recommended)

2. **Clone the Repository**
```bash
git clone https://github.com/yourusername/vitadb.git
cd vitadb
```

3. **Build the Project**
If you have Make installed:
```bash
make build
```
Otherwise, use Go directly:
```bash
go build ./...
```

4. **Run VitaDB**
If you have Make installed:
```bash
make run
```
Otherwise, use Go directly:
```bash
go run cmd/main.go
```

5. **Using the CLI**
Once VitaDB is running, you can interact with it using the built-in CLI. Here are some basic commands:
- Set a key-value pair: `set <key> <value>`
- Get a value: `get <key>`
- Delete a key: `delete <key>`
- Exit the CLI: `exit`

6. **Running Tests**
To run the test suite:
`make test`
Or without Make:
`go test ./...`

## License

VitaDB is released under the [MIT License](LICENSE).
