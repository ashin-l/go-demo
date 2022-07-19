/*
 * @Author: lqc
 * @Date: 2021-08-17 13:14:50
 * @Description: 实现摄像机抓拍图片并保存到本地
 */
package funcs

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"

	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/option"
	"github.com/ashin-l/go-demo/pkg/util"
	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/gosoap"
	"github.com/use-go/onvif/media"
)

func Snapshot(opt *option.Options, dev *goonvif.Device) {
	snapurl := media.GetSnapshotUri{ProfileToken: GetProfiles(dev).Token}
	resp, err := dev.CallMethod(snapurl)
	bs, err := util.ReadResponse(resp)
	if err != nil {
		logger.Logger().Error("read resp error: ", err)
		return
	}

	sm := gosoap.SoapMessage(bs)
	data := media.GetSnapshotUriResponse{}
	err = xml.Unmarshal([]byte(sm.Body()), &data)
	if err != nil {
		panic(err)
	}

	logger.Logger().Info(data)
	uri := string(data.MediaUri.Uri)
	url := uri[:7] + opt.Camera.Username + ":" + opt.Camera.Password + "@" + uri[7:]
	logger.Logger().Info("url:", url)
	resp, err = http.Get(url)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.Logger().Info("http请求报错:", resp.Status)
		return
	}

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
