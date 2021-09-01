/*
 * @Author: lqc
 * @Date: 2021-08-31 08:53:04
 * @Description: 实现摄像机连续移动
 */
package funcs

import (
	"fmt"
	"time"

	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/gosoap"
	"github.com/use-go/onvif/ptz"
)

// 连续移动
func ContiMove(dev *goonvif.Device) {

	// pf := media.GetProfile{}
	// resp, err := dev.CallMethod(&pf)
	// if err != nil {
	// 	panic(err)
	// }
	// sm := gosoap.SoapMessage(readResponse(resp))
	// fmt.Println(sm)
	// data := media.GetProfilesResponse{}
	// err = xml.Unmarshal([]byte(sm.Body()), &data)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(data)

	pcm := ptz.ContinuousMove{}
	pcm.Velocity.PanTilt.X = 0.2
	pcm.ProfileToken = "Profile_1"

	resp, err := dev.CallMethod(&pcm)
	if err != nil {
		panic(err)
	}
	sm := gosoap.SoapMessage(readResponse(resp))
	fmt.Println(sm)
	time.Sleep(3 * time.Second)
	ps := ptz.Stop{}
	ps.ProfileToken = "Profile_1"
	_, err = dev.CallMethod(&ps)
	if err != nil {
		panic(err)
	}
}
