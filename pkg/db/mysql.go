package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlDB struct {
	Database
	DB *sql.DB
}

func (m *MysqlDB) Connect() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Pass, m.Host, m.Port, m.Name)
	m.DB, err = sql.Open("mysql", dsn)
	return
}

func (m *MysqlDB) Select(query string) []interface{} {
	rows, err := m.DB.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	col, _ := rows.Columns()
	size := len(col)
	var data []interface{}
	for rows.Next() {
		tmp := make([]interface{}, size)
		t := make([]interface{}, size)
		for k := range tmp {
			t[k] = &tmp[k]
		}
		rows.Scan(t...)
		data = append(data, tmp)
	}
	rows.Close()
	return data
}

func (m *MysqlDB) Insert(query string, data interface{}) {
	stmt, _ := m.DB.Prepare(query)
	val, _ := data.([]interface{})
	for _, v := range val {
		val := v.([]interface{})
		stmt.Exec(val...)
	}
	stmt.Close()
}
