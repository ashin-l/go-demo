package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format"
)

func init() {
	format.RegisterAll()
}

func main() {
	srcfile := flag.String("src", "rtsp://admin:ai123456@192.168.160.101:554/cam/realmonitor?channel=1&subtype=0", "Source file")
	dstfile := flag.String("dst", "output.flv", "Output file")
	max := flag.Int("max", 5, "Max seconds")
	flag.Parse()

	src, err := avutil.Open(*srcfile)
	if err != nil {
		panic(err)
	}
	dst, err := avutil.Create(*dstfile)
	if err != nil {
		panic(err)
	}
	// same as calling avutil.CopyFile(dst, src) but added
	// max duration in case the src is live and never ends
	err = CopyFileMax(dst, src, time.Duration(*max)*time.Second)
	if err != nil {
		panic(err)
	}
}

func CopyPacketsMax(dst av.PacketWriter, src av.PacketReader, max time.Duration) (err error) {
	for {
		var pkt av.Packet
		if pkt, err = src.ReadPacket(); err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		// break when max time has been reached
		if max > 0 && pkt.Time >= max {
			return
		}

		if err = dst.WritePacket(pkt); err != nil {
			return
		}
	}
	return
}

func CopyFileMax(dst av.Muxer, src av.Demuxer, max time.Duration) (err error) {
	var streams []av.CodecData
	if streams, err = src.Streams(); err != nil {
		fmt.Println("111")
		return
	}
	if err = dst.WriteHeader(streams); err != nil {
		fmt.Println("222")
		return
	}
	if err = CopyPacketsMax(dst, src, max); err != nil {
		fmt.Println("333")
		return
	}
	if err = dst.WriteTrailer(); err != nil {
		fmt.Println("444")
		return
	}
	return
}
