/*
 * @Author: lqc
 * @Date: 2021-08-31 09:03:33
 * @Description: 基于onvif实现对摄像机的控制
 */
package main

import (
	"fmt"
	"os"

	"github.com/ashin-l/go-demo/cmd/onvif/funcs"
	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/option"
	"github.com/use-go/onvif"
)

const (
	login    = "admin"
	password = "ai123456"
)

func main() {
	opt := option.New()
	err := opt.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	logger.Init(opt)
	dev, err := onvif.NewDevice(opt.Camera.Ip)
	if err != nil {
		panic(err)
	}
	dev.Authenticate(opt.Camera.Username, opt.Camera.Password)
	// 抓拍
	// funcs.Snapshot(dev)

	// ptz控制:连续移动
	// funcs.ContiMove(dev)

	// funcs.GetStatus(dev)
	// ptz控制:绝对移动
	funcs.AbsoluteMove(dev)

	// 获取预置点
	// funcs.GetPreset(dev)

	// 获取ptz状态
	// funcs.GetStatus(dev)

	// 获取 media.profile
	// funcs.GetProfile(dev)

	// resp, err := dev.CallMethod(snapurl)
	// sm := gosoap.SoapMessage(readResponse(resp))
	// fmt.Println(sm.Body())
	// data := media.GetSnapshotUriResponse{}
	// err = xml.Unmarshal([]byte(sm.Body()), &data)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(data)
	//Preparing commands
	// systemDateAndTyme := device.GetSystemDateAndTime{}
	// getCapabilities := device.GetCapabilities{Category: "All"}
	// createUser := device.CreateUsers{User: onvif.User{
	// 	Username:  "TestUser",
	// 	Password:  "TestPassword",
	// 	UserLevel: "User",
	// },
	// }

	// //Commands execution
	// systemDateAndTymeResponse, err := dev.CallMethod(systemDateAndTyme)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	fmt.Println(readResponse(systemDateAndTymeResponse))
	// }
	// getCapabilitiesResponse, err := dev.CallMethod(getCapabilities)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	fmt.Println(readResponse(getCapabilitiesResponse))
	// }
	// createUserResponse, err := dev.CallMethod(createUser)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	/*
	// 		You could use https://github.com/use-go/onvif/gosoap for pretty printing response
	// 	*/
	// 	fmt.Println(gosoap.SoapMessage(readResponse(createUserResponse)).StringIndent())
	// }
}
