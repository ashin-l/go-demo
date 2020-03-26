package main

import (
	"fmt"
	"reflect"
)

func main() {
	type Person struct {
		Name string
		Age  int
	}

	var a interface{}
	a = Person{"haha", 12}

	st := reflect.TypeOf(a)
	fmt.Println(st.Name())
	val := reflect.ValueOf(a)
	num := val.NumField()

	for i := 0; i != num; i++ {
		v := val.Field(i)
		f := st.Field(i)
		fmt.Println(f.Name, v)
	}

}
