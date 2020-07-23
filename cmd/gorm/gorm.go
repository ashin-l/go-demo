package main

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	_ int = iota
	ABSENCE
	NORMAL
	LATE
	LEAVE
	OVERTM
)

const (
	START = 1589156415
	END   = 1589188995
	DAYS  = 30
)

type Employee struct {
	EmpId string
	DepId int
	Name  string
}

type Attendance struct {
	Day       string
	EmpId     string
	DepId     int
	Name      string
	StartTime int64
	EndTime   int64
	Status    int
}

func main() {
	db, err := gorm.Open("postgres", "host=192.168.152.37 port=5432 user=postgres dbname=smart_park password=postgres sslmode=disable")
	if err != nil {
		panic("failed to connect database " + err.Error())
	}
	defer db.Close()

	rows, err := db.Raw("select emp_id, dep_id, name from cfg_employee").Rows() // (*sql.Rows, error)
	defer rows.Close()
	emps := []Employee{}
	for rows.Next() {
		tmp := Employee{}
		rows.Scan(&tmp.EmpId, &tmp.DepId, &tmp.Name)
		emps = append(emps, tmp)
	}
	data := []Attendance{}
	for i := 0; i != DAYS; i++ {
		start := int64(START + 24*60*60*i)
		end := int64(END + 24*60*60*i)
		dt := time.Unix(start, 0)
		if dt.Weekday() == time.Saturday || dt.Weekday() == time.Sunday {
			continue
		}
		day := time.Unix(start, 0).Format("2006-01-02")
		for _, v := range emps {
			ad := Attendance{
				EmpId: v.EmpId,
				DepId: v.DepId,
				Day:   day,
				Name:  v.Name,
			}
			r := rand.Intn(10)
			if r == 7 {
				ad.Status = ABSENCE
			} else if r == 3 {
				ad.StartTime = (start + 3600) * 1000
				ad.EndTime = end * 1000
				//ad.StartTime = time.Unix(start+3600, 0).Format("2006-01-02 15:04:05")
				//ad.EndTime = time.Unix(end, 0).Format("2006-01-02 15:04:05")
				ad.Status = LATE
			} else if r == 1 {
				ad.StartTime = start * 1000
				ad.EndTime = (end - 3600) * 1000
				//ad.StartTime = time.Unix(start, 0).Format("2006-01-02 15:04:05")
				//ad.EndTime = time.Unix(end-3600, 0).Format("2006-01-02 15:04:05")
				ad.Status = LEAVE
			} else {
				ad.StartTime = start * 1000
				ad.EndTime = end * 1000
				//ad.StartTime = time.Unix(start, 0).Format("2006-01-02 15:04:05")
				//ad.EndTime = time.Unix(end, 0).Format("2006-01-02 15:04:05")
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
		if err := tx.Table("sdr_attendance_statistics").Create(&v).Error; err != nil {
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
