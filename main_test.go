package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPrepareCommandInsert(t *testing.T) {
	tests := []struct {
		name          string
		statement     string
		expectedRows  uint32
		expectedID    uint32
		expectedEmail string
		expectedUser  string
	}{
		{
			name:          "insert first row",
			statement:     "insert 1 alice@example.com alice",
			expectedRows:  1,
			expectedID:    1,
			expectedEmail: "alice@example.com",
			expectedUser:  "alice",
		},
		{
			name:          "insert row with longer values",
			statement:     "insert 42 bob.smith@example.com bobsmith",
			expectedRows:  1,
			expectedID:    42,
			expectedEmail: "bob.smith@example.com",
			expectedUser:  "bobsmith",
		},
		{
			name:          "insert row with numbers in username",
			statement:     "insert 100 user123@test.com user123",
			expectedRows:  1,
			expectedID:    100,
			expectedEmail: "user123@test.com",
			expectedUser:  "user123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			old := os.Stdout
			_, w, _ := os.Pipe()
			os.Stdout = w

			table := Table{
				numRows: 0,
				pages:   make([]byte, 0),
				dbFile:  "test.db",
			}

			row := table.PrepareCommand([]byte(tt.statement))

			w.Close()
			os.Stdout = old

			if table.numRows != tt.expectedRows {
				t.Errorf("numRows = %d, want %d", table.numRows, tt.expectedRows)
			}

			if row.id != tt.expectedID {
				t.Errorf("row.id = %d, want %d", row.id, tt.expectedID)
			}

			email := bytes.TrimRight([]byte(row.email), "\x00")
			if string(email) != tt.expectedEmail {
				t.Errorf("row.email = %q, want %q", string(email), tt.expectedEmail)
			}

			username := bytes.TrimRight([]byte(row.username), "\x00")
			if string(username) != tt.expectedUser {
				t.Errorf("row.username = %q, want %q", string(username), tt.expectedUser)
			}

			if len(table.pages) < ROW_SIZE {
				t.Errorf("pages length = %d, want at least %d", len(table.pages), ROW_SIZE)
			}

			slot, bound := rowSlot(0)
			storedRow := deserialise(table.pages[slot:bound])
			if storedRow.id != tt.expectedID {
				t.Errorf("stored row.id = %d, want %d", storedRow.id, tt.expectedID)
			}
		})
	}
}

func TestPrepareCommandMultipleInserts(t *testing.T) {

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	table := Table{
		numRows: 0,
		pages:   make([]byte, 0),
		dbFile:  "test.db",
	}

	statements := []struct {
		statement string
		id        uint32
	}{
		{"insert 1 alice@example.com alice", 1},
		{"insert 2 bob@example.com bob", 2},
		{"insert 3 charlie@example.com charlie", 3},
	}

	for _, stmt := range statements {
		row := table.PrepareCommand([]byte(stmt.statement))
		if row.id != stmt.id {
			t.Errorf("row.id = %d, want %d", row.id, stmt.id)
		}
	}

	w.Close()
	os.Stdout = old
	io.Copy(io.Discard, r)

	if table.numRows != 3 {
		t.Errorf("numRows = %d, want 3", table.numRows)
	}

	for i := uint32(0); i < 3; i++ {
		slot, bound := rowSlot(i)
		row := deserialise(table.pages[slot:bound])
		if row.id != i+1 {
			t.Errorf("row %d: id = %d, want %d", i, row.id, i+1)
		}
	}
}

func TestPrepareCommandSelect(t *testing.T) {

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	table := Table{
		numRows: 0,
		pages:   make([]byte, 0),
		dbFile:  "test.db",
	}

	table.PrepareCommand([]byte("insert 1 alice@example.com alice"))
	table.PrepareCommand([]byte("insert 2 bob@example.com bob"))

	table.PrepareCommand([]byte("select"))

	w.Close()
	os.Stdout = old

	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	output := buf.String()

	if !bytes.Contains([]byte(output), []byte("{1")) {
		t.Error("select output should contain first row with id 1")
	}
	if !bytes.Contains([]byte(output), []byte("{2")) {
		t.Error("select output should contain second row with id 2")
	}
}

func TestPrepareCommandSelectEmpty(t *testing.T) {

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	table := Table{
		numRows: 0,
		pages:   make([]byte, 0),
		dbFile:  "test.db",
	}

	table.PrepareCommand([]byte("select"))

	w.Close()
	os.Stdout = old

	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	output := buf.String()

	_ = output
}

func TestPrepareCommandUnrecognized(t *testing.T) {

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	table := Table{
		numRows: 0,
		pages:   make([]byte, 0),
		dbFile:  "test.db",
	}

	table.PrepareCommand([]byte("delete 1"))

	w.Close()
	os.Stdout = old

	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	output := buf.String()

	if !bytes.Contains([]byte(output), []byte("unrecognized command")) {
		t.Error("unrecognized command should produce error message")
	}
}
