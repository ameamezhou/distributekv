package dbpkg

import (
	"bytes"
	"os"
	"testing"
)

func TestDataBase_SetKey(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "kv.db")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer f.Close()
	name := f.Name()
	defer os.Remove(name)

	db, err := NewDataBase(&name)
	if err != nil {
		t.Fatalf("could not create a new database: %v", err)
	}
	defer db.Close()

	if err = db.SetKey("party", []byte("xiaoqizho")); err != nil {
		t.Fatalf("could not write key: %v", err)
	}

	value, err := db.GetValue("party")
	if err != nil {
		t.Fatalf("could not get key: %v", err)
	}
	if !bytes.Equal(value, []byte("xiaoqizhou")) {
		t.Fatalf("unexpected value for key party: got %q, want %q", value, "xiaoqizhou")
	}
}