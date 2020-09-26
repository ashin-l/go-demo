package main

import (
	"encoding/json"
	"fmt"
)

type Wheel struct {
	Size int
}

type Engine struct {
	Type  string
	Power int
}

type Car struct {
	Wheel
	Engine
}

func main() {
	car := Car{
		Wheel: Wheel{
			Size: 1,
		},
		Engine: Engine{
			Type:  "a",
			Power: 300,
		},
	}

	jstr := `{"Size":33,"Type":"bb","Power":500}`
	car1 := Car{}
	json.Unmarshal([]byte(jstr), &car1)
	fmt.Println(car, car1)
}
