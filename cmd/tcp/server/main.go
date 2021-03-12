package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

type SensorData struct {
	Temperature   float32
	Humidity      float32
	Airpressure   float32
	Windspeed     float32
	Winddirection float32
	Timestamp     int64
}

func main() {
	fmt.Println("server start...")
	l, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		fmt.Println("server error: ", err.Error())
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("connect error: ", err.Error())
			continue
		}
		fmt.Println("connect!!!")
		go process(conn)
	}
}

func process(conn net.Conn) {
	fmt.Println("1111111111111")
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err.Error())
			return
		}
		fmt.Println("size:", n)
		fmt.Println("data:", string(buf[:n]))
		// jstr, _ := json.Marshal(`{"msgType":1,"code":0}`)
		jstr := []byte(`{"msgType":1,"code":0}`)
		fmt.Println("len:", len(jstr))
		jstr = append(jstr, '\n')
		fmt.Println("jstr", jstr)
		conn.Write(jstr)
		fmt.Println("len:", len(jstr))
		// 气象传感器
		// sData := SensorData{}

		// sData.Temperature = ByteToFloat32(buf, 6)
		// sData.Humidity = ByteToFloat32(buf, 10)
		// sData.Airpressure = ByteToFloat32(buf, 14)
		// sData.Windspeed = ByteToFloat32(buf, 18)
		// sData.Winddirection = ByteToFloat32(buf, 22)

		// sData.Timestamp = time.Now().UnixNano() / 1e6
		// fmt.Println(sData)
	}

}

func ByteToFloat32(buf []byte, index int) float32 {
	bits := binary.LittleEndian.Uint32(buf[index : index+4])
	return math.Float32frombits(bits)
}
