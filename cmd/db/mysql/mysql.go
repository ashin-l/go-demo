package main

import (
	"fmt"
	"math/rand"

	"github.com/ashin-l/go-demo/pkg/db"
)

const (
	START = 1572567621
	END   = 1572600201
	DAYS  = 27
)

func main() {
	db := &db.MysqlDB{
		Database: db.Database{
			Host: "192.168.152.181",
			Port: "3306",
			User: "root",
			Pass: "1234.Com",
			Name: "tianhe",
		},
	}
	err := db.Connect()
	if err != nil {
		fmt.Println("connect db error: ", err)
	}
	empids := db.Select("select EmpID from employee")
	ids := make([]string, 0, 10)
	for _, v := range empids {
		value, _ := v.([]interface{})
		for _, val := range value {
			t, _ := val.([]byte)
			ids = append(ids, string(t))
		}
	}
	data := make([]interface{}, 0, 10)
	for i := 0; i != DAYS; i++ {
		start := START + 24*60*60*i
		end := END + 24*60*60*i
		for _, v := range ids {
			fmt.Println(v)
			r := rand.Intn(70)
			if r == 37 {
				continue
			}
			tmp := make([]interface{}, 0, 2)
			tmp = append(tmp, v)
			if r == 13 {
				tmp = append(tmp, start+3600)
			} else {
				tmp = append(tmp, start)
			}
			tmp1 := make([]interface{}, 0, 2)
			tmp1 = append(tmp1, v)
			if r == 17 {
				tmp1 = append(tmp1, end-3600)
			} else {
				tmp1 = append(tmp1, end)
			}
			data = append(data, tmp, tmp1)
		}
	}
	sqlstr := "insert into attendanceNow (EmpID, Time) values(?, ?)"
	db.Insert(sqlstr, data)
}
