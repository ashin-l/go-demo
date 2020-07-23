package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresDB struct {
	DataSource
	DB *gorm.DB
}

func (m *PostgresDB) Connect() (err error) {
	str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", m.Host, m.Port, m.UserName, m.PassWord, m.DbName, m.Sslmode)
	m.DB, err = gorm.Open("postgres", str)
	return
}

func (m *PostgresDB) Save(data interface{}) {
	m.DB.Save(data)
}

func (m *PostgresDB) First(data interface{}) error {
	return m.DB.First(data).Error
}

func (m *PostgresDB) Find(data interface{}) {
	m.DB.Find(data)
}
