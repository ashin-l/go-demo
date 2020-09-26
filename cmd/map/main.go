package main

import (
	"fmt"
)

func main() {
	m := make(map[string]interface{})
	m["key"] = "value"
	fmt.Println(m)
	fmt.Println(m["hi"])
	v, ok := m["hi"].(string)
	if ok {
		fmt.Println(v)
	} else {
		fmt.Println("not ok")
	}
}
