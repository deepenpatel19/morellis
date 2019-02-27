package mysql

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"
)

var dsn string

func init() {
	if os.Getenv("TEST_DSN") != "" {
		dsn = os.Getenv("TEST_DSN")
	} else {
		panic("No test DSN defined")
	}
}

func newTestDB(t *testing.T) (*sql.DB, func()) {

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile("./testdata/setup.sql")

	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}
