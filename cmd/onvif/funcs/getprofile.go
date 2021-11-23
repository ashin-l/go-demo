/*
 * @Author: lqc
 * @Date: 2021-10-19 15:28:25
 * @Description:
 */

package funcs

import (
	"encoding/xml"
	"fmt"

	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/gosoap"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/xsd/onvif"
)

// func GetProfile(dev *goonvif.Device) media.GetProfileResponse {
// 	mgp := media.GetProfile{}
// 	resp, err := dev.CallMethod(&mgp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	sm := gosoap.SoapMessage(readResponse(resp))
// 	fmt.Println(sm.Body())

// 	data := media.GetProfileResponse{}
// 	err = xml.Unmarshal([]byte(sm.Body()), &data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("=====================")
// 	fmt.Printf("profile: %+v\n", data)
// 	fmt.Println("=====================")
// 	fmt.Println("token: ", data.Profile.Token)
// 	return data
// }

func GetProfiles(dev *goonvif.Device) onvif.Profile {
	mgp := media.GetProfiles{}
	resp, err := dev.CallMethod(&mgp)
	if err != nil {
		panic(err)
	}
	sm := gosoap.SoapMessage(readResponse(resp))

	data := media.GetProfilesResponse{}
	err = xml.Unmarshal([]byte(sm.Body()), &data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("profiles: %+v\n", data)
	fmt.Println("profileToken:", data.Profiles[0].Token)
	return data.Profiles[0]
}
