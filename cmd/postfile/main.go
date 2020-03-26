package main

import (
	"fmt"

	"github.com/astaxie/beego/httplib"
)

func main() {
	var obj interface{}
	req := httplib.Post("http://192.168.152.44:30002/group1/upload")
	req.PostFile("file", "111.jpg") //注意不是全路径
	req.Param("output", "json")
	req.Param("scene", "")
	req.Param("path", "/h5/111.jpg")
	req.ToJSON(&obj)
	fmt.Println(obj)
}
