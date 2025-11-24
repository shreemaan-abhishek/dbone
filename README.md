# dbone

A lightweight, file-based database implementation written in Go. This project demonstrates fundamental database concepts including data persistence, serialization, and a simple REPL interface.

## Features

- **Simple REPL Interface**: Interactive command-line interface for database operations
- **File-based Persistence**: Data is stored in a binary file (`dbone.db`)
- **Fixed-size Rows**: Efficient storage with predetermined row sizes
- **Page-based Architecture**: Data organized in 4KB pages for better memory management
- **Basic CRUD Operations**: Support for insert and select operations

## Architecture

The database uses a page-based storage system with the following characteristics:

- **Row Structure**: Each row contains:
  - `id` (uint32, 4 bytes)
  - `username` (string, 32 bytes)
  - `email` (string, 255 bytes)
  - Total row size: 291 bytes

- **Storage Layout**:
  - Page size: 4096 bytes
  - Rows per page: 14
  - Maximum pages: 100
  - Maximum capacity: 1,400 rows

## Installation

```bash
git clone https://github.com/shreemaan-abhishek/dbone.git
cd dbone
go build
```

## Usage

Start the database:

```bash
./dbone
```

### Available Commands

#### Insert Data
```
insert <id> <email> <username>
```
Example:
```
dbone > insert 1 user@example.com johndoe
```

#### Select All Data
```
select
```
Example:
```
dbone > select
```

#### Meta Commands
- `.exit` - Save and exit the database

## Example Session

```
dbone > insert 1 alice@example.com alice
inserting @ 0
dbone > insert 2 bob@example.com bob
inserting @ 291
dbone > select
{1 alice alice@example.com}
{2 bob bob@example.com}
dbone > .exit
```

## Technical Details

### Data Serialization
- Uses little-endian byte order for numeric data
- Fixed-width fields for predictable memory layout
- Binary encoding for efficient storage

### Persistence
- Data is automatically loaded from `dbone.db` on startup
- Changes are written to disk when exiting with `.exit`
- Binary file format for compact storage

## Roadmap

- [x] Simple REPL interface
- [x] INSERT operations
- [x] SELECT operations (full table scan)
- [x] File-based persistence
- [x] Page-based storage architecture
- [ ] DELETE operations
- [ ] UPDATE operations
- [ ] WHERE clause support for SELECT queries
- [ ] B-tree indexing
- [ ] Query optimization
- [ ] Dynamic schema support
- [ ] Transaction support
- [ ] Concurrent access handling

## License

MIT License - see [LICENSE](LICENSE) file for details

## Author

Shreemaan Abhishek

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.
