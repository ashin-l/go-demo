/*
 * @Author: lqc
 * @Date: 2021-09-24 15:59:04
 * @Description: 获取当前ptz状态
 */

package funcs

import (
	"encoding/xml"
	"fmt"

	"github.com/ashin-l/go-demo/cmd/onvif/model"
	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/gosoap"
	"github.com/use-go/onvif/ptz"
)

func GetStatus(dev *goonvif.Device) {
	GetProfiles(dev)
	pgs := ptz.GetStatus{ProfileToken: "Profile_1"}
	resp, err := dev.CallMethod(&pgs)
	if err != nil {
		panic(err)
	}
	sm := gosoap.SoapMessage(readResponse(resp))

	data := model.GetStatusResponse{}
	err = xml.Unmarshal([]byte(sm.Body()), &data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("status: %+v\n\n", data)
}
