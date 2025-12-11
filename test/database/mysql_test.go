package database

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlRootPassword = "dev123"
	mysqlDatabaseName = "gonoweb"
)

// go test -v -count=1 -timeout 30s -run ^TestPingMysql$ github.com/sfshf/gonoweb/test/database
func TestPingMysql(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@/%s", mysqlRootPassword, mysqlDatabaseName))
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	t.Log("Mysql is OK !")
}
