/*
 * @Author: lqc
 * @Date: 2021-08-31 08:53:04
 * @Description: 实现摄像机连续移动
 */
package funcs

import (
	"time"

	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/util"
	goonvif "github.com/use-go/onvif"
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

	_, err := dev.CallMethod(&pcm)
	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
	ps := ptz.Stop{}
	ps.ProfileToken = "Profile_1"
	_, err = dev.CallMethod(&ps)
	if err != nil {
		panic(err)
	}
}

func AbsoluteMove(dev *goonvif.Device) {
	pf := GetProfiles(dev)
	pam := ptz.AbsoluteMove{ProfileToken: pf.Token}
	pam.Position.PanTilt.X = 0.185778
	pam.Position.PanTilt.Y = 1.00
	pam.Position.Zoom.X = 0
	// pam.Position.PanTilt.Space = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/PositionGenericSpace"
	// pam.Position.Zoom.Space = "http://www.onvif.org/ver10/tptz/ZoomSpaces/PositionGenericSpace"
	pam.Position.PanTilt.Space = pf.PTZConfiguration.DefaultAbsolutePantTiltPositionSpace
	pam.Position.Zoom.Space = pf.PTZConfiguration.DefaultAbsoluteZoomPositionSpace
	pam.Speed.PanTilt.X = 1.0
	pam.Speed.PanTilt.Y = 1.0
	// pam.Speed.PanTilt.Space = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/GenericSpeedSpace"
	logger.Logger().Info("pam:", pam)
	resp, err := dev.CallMethod(&pam)
	if err != nil {
		panic(err)
	}

	bs, _ := util.ReadResponse(resp)
	logger.Logger().Info(string(bs))
	time.Sleep(3 * time.Second)
}
