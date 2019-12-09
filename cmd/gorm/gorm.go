package main

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	NORMAL  = 0
	LATE    = 1
	LEAVE   = 2
	ABSENCE = 3
)

const (
	START = 1572567621
	END   = 1572600201
	DAYS  = 27
)

type Employee struct {
	EmpID string
	DepID string
}

type Attendance struct {
	Day       string
	EmpID     string `gorm:"column:EmpID"`
	DepID     string `gorm:"column:DepID"`
	StartTime string `gorm:"column:StartTime"`
	EndTime   string `gorm:"column:EndTime"`
	Status    int
}

func main() {
	db, err := gorm.Open("mysql", "root:1234.Com@tcp(192.168.152.181:3306)/tianhe")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	rows, err := db.Raw("select EmpID, DepID from employee").Rows() // (*sql.Rows, error)
	defer rows.Close()
	emps := []Employee{}
	for rows.Next() {
		tmp := Employee{}
		rows.Scan(&tmp.EmpID, &tmp.DepID)
		emps = append(emps, tmp)
	}
	data := []Attendance{}
	for i := 0; i != DAYS; i++ {
		start := int64(START + 24*60*60*i)
		end := int64(END + 24*60*60*i)
		day := time.Unix(start, 0).Format("2006-01-02")
		for _, v := range emps {
			ad := Attendance{
				EmpID: v.EmpID,
				DepID: v.DepID,
				Day:   day,
			}
			r := rand.Intn(10)
			if r == 7 {
				ad.Status = ABSENCE
			} else if r == 3 {
				ad.StartTime = time.Unix(start+3600, 0).Format("2006-01-02 15:04:05")
				ad.EndTime = time.Unix(end, 0).Format("2006-01-02 15:04:05")
				ad.Status = LATE
			} else if r == 1 {
				ad.StartTime = time.Unix(start, 0).Format("2006-01-02 15:04:05")
				ad.EndTime = time.Unix(end-3600, 0).Format("2006-01-02 15:04:05")
				ad.Status = LEAVE
			} else {
				ad.StartTime = time.Unix(start, 0).Format("2006-01-02 15:04:05")
				ad.EndTime = time.Unix(end, 0).Format("2006-01-02 15:04:05")
				ad.Status = NORMAL
			}
			data = append(data, ad)
		}
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, v := range data {
		if err := tx.Table("attendance").Create(&v).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()

	//	// 更新 - 更新product的price为2000
	//	db.Model(&product).Update("Price", 2000)
	//
	//	// 删除 - 删除product
	//	db.Delete(&product)
}
