package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name     string
	Age      int
	GenderId int8 `json:"xxx"` // json注解解决字段名称不一致问题
}

func main() {
	jstr := `{"name":"h","age":12,"genderId":31}`
	p := Person{}
	json.Unmarshal([]byte(jstr), &p)
	fmt.Println(p)
}
