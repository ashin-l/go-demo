/*
 * @Author: lqc
 * @Date: 2021-09-24 15:59:04
 * @Description: 获取当前ptz状态
 */

package funcs

import (
	"encoding/xml"

	"github.com/ashin-l/go-demo/cmd/onvif/model"
	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/util"
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

	bs, err := util.ReadResponse(resp)
	if err != nil {
		panic(err)
	}

	sm := gosoap.SoapMessage(bs)
	data := model.GetStatusResponse{}
	err = xml.Unmarshal([]byte(sm.Body()), &data)
	if err != nil {
		panic(err)
	}

	logger.Logger().Infof("status: %+v\n\n", data)
}
