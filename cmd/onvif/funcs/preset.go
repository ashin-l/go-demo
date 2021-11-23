/*
 * @Author: lqc
 * @Date: 2021-09-24 13:08:59
 * @Description: 预置点
 */
package funcs

import (
	"fmt"

	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/gosoap"
	"github.com/use-go/onvif/ptz"
)

func GetPreset(dev *goonvif.Device) {
	pgp := ptz.GetPresets{ProfileToken: "Profile_1"}
	resp, err := dev.CallMethod(&pgp)
	if err != nil {
		panic(err)
	}
	sm := gosoap.SoapMessage(readResponse(resp))
	fmt.Println(sm)
	// data := ptz.GetPresetsResponse{}
	// err = xml.Unmarshal([]byte(sm.Body()), &data)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(data)
}
