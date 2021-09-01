/*
 * @Author: lqc
 * @Date: 2021-08-17 13:14:50
 * @Description: 实现摄像机抓拍图片并保存到本地
 */
package funcs

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/gosoap"
	"github.com/use-go/onvif/media"
)

func Snapshot(dev *goonvif.Device) {
	snapurl := media.GetSnapshotUri{}
	fmt.Println(snapurl)

	pkgPath := strings.Split(reflect.TypeOf(snapurl).PkgPath(), "/")
	pkg := strings.ToLower(pkgPath[len(pkgPath)-1])
	fmt.Println(pkg)

	resp, err := dev.CallMethod(snapurl)
	sm := gosoap.SoapMessage(readResponse(resp))
	data := media.GetSnapshotUriResponse{}
	err = xml.Unmarshal([]byte(sm.Body()), &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	uri := string(data.MediaUri.Uri)
	url := uri[:7] + Conf.Username + ":" + Conf.Password + "@" + uri[7:]
	resp, err = http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create("test.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
