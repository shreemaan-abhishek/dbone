package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	ID_SIZE       = 4
	USERNAME_SIZE = 32
	EMAIL_SIZE    = 255

	ID_OFFSET       = 0
	USERNAME_OFFSET = ID_OFFSET + ID_SIZE
	EMAIL_OFFSET    = USERNAME_OFFSET + USERNAME_SIZE
	ROW_SIZE        = EMAIL_OFFSET + EMAIL_SIZE

	PAGE_SIZE       = 4096
	ROWS_PER_PAGE   = PAGE_SIZE / ROW_SIZE
	TABLE_MAX_PAGES = 100
	TABLE_SIZE      = PAGE_SIZE * TABLE_MAX_PAGES
)

type Row struct {
	id       uint32
	username string
	email    string
}

type Table struct {
	numRows uint32
	pages   []byte
	dbFile  string
}

func main() {
	var table Table
	table.dbFile = "dbone.db"
	data, err := os.ReadFile(table.dbFile)
	if err != nil {
		log.Printf("failed to read db file: %s", err.Error())
	} else {
		table.pages = make([]byte, len(data))
		copy(table.pages[:], data)
		table.numRows = uint32(len(data)) / ROW_SIZE
	}

	for {
		fmt.Print("dbone > ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("failed to read from stdin")
		}
		if input[0] == '.' {
			if doMetaCommand(input) != 0 {
				err := os.WriteFile(table.dbFile, table.pages[:], 0644)
				if err != nil {
					log.Fatalf("failed to read db file: %s", err.Error())
				}

				return
			}
			continue
		}

		table.PrepareCommand(input)

	}
}

func (table *Table) PrepareCommand(statement []byte) *Row {
	// TODO: validate the format of statement
	statement_str := string(statement)
	var row Row
	var command string

	fmt.Sscanf(statement_str, "%s %d %s %s", &command, &row.id, &row.email, &row.username)

	switch command {
	case "select":
		for i := 0; i < int(table.numRows); i++ {
			row, bound := rowSlot(uint32(i))
			data := deserialise(table.pages[row:bound])
			fmt.Println(data)
		}

	case "insert":
		slot, bound := rowSlot(table.numRows)
		fmt.Println("inserting @ ", slot)
		serialisedRow := serialise(row)
		table.pages = append(table.pages, make([]byte, ROW_SIZE)...)
		copy(table.pages[slot:bound], serialisedRow)
		table.numRows = table.numRows + 1
		return &row

	default:
		fmt.Println("unrecognized command")
	}
	return &row
}

func serialise(data Row) []byte {
	bin := make([]byte, ROW_SIZE)
	binary.LittleEndian.PutUint32(bin, data.id)
	copy(bin[USERNAME_OFFSET:], data.username[:])
	copy(bin[EMAIL_OFFSET:], data.email[:])

	return bin
}

func deserialise(bin []byte) Row {
	var row Row
	row.id = binary.LittleEndian.Uint32(bin[ID_OFFSET:])
	row.username = string(bin[USERNAME_OFFSET : USERNAME_OFFSET+USERNAME_SIZE])
	row.email = string(bin[EMAIL_OFFSET : EMAIL_OFFSET+EMAIL_SIZE])

	return row
}

func doMetaCommand(input []byte) int {
	inputStr := string(input)
	switch inputStr {
	case ".exit\n":
		return -1
	default:
		fmt.Println("unrecognized command")
		return 0
	}
}

func rowSlot(rowNum uint32) (uint32, uint32) {
	pageNum := rowNum / ROWS_PER_PAGE
	rowOffset := rowNum % ROWS_PER_PAGE
	byteOffset := rowOffset * ROW_SIZE

	slot := pageNum + byteOffset
	bound := slot + ROW_SIZE

	if bound > TABLE_SIZE {
		log.Fatalf("row %d out of bounds", rowNum)
	}

	return slot, bound
}
