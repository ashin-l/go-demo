package main

import (
	"fmt"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	fi, _ := ioutil.ReadFile("noir.aif")

	runFFMPEGFromStdin(populateStdin(fi))
}

func populateStdin(file []byte) func(io.WriteCloser) {
	return func(stdin io.WriteCloser) {
		defer stdin.Close()
		io.Copy(stdin, bytes.NewReader(file))
	}
}

func runFFMPEGFromStdin(populate_stdin_func func(io.WriteCloser)) {
	//cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-ab", "128k", "-f", "mp3", "pipe:1")
	// ffmpeg -rtsp_transport tcp -buffer_size 1024000 -max_delay 500000 -stimeout 20000000 -analyzeduration 100 -max_delay 100 -i rtsp://admin:1234.com@192.168.165.123:554/h265/ch1/main/av_stream -an -vcodec copy -f flv pipe:1
	cmd := exec.Command("ffmpeg", "-rtsp_transport", "tcp", "-buffer_size", "1024000", "-max_delay", "500000", "-stimeout", "20000000", "-analyzeduration", "100", "-max_delay", "100", "-i", "rtsp://admin:1234.com@192.168.165.123:554/h265/ch1/main/av_stream", "-an", "-vcodec", "copy", "-f", "flv", "pipe:1")
	//stdin, err := cmd.StdinPipe()
	//if err != nil {
	//	log.Panic(err)
	//}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Panic(err)
	}
	//populate_stdin_func(stdin)
	//fo, _ := os.Create("output.mp3")
	//io.Copy(fo, stdout)
	fo, _ := os.Create("output.flv")
	io.Copy(fo, stdout)
	fmt.Println("====")

	err = cmd.Wait()
	if err != nil {
		log.Panic(err)
	}
}
