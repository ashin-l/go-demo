package util

import (
	"io/ioutil"
	"net/http"
)

func ReadResponse(resp *http.Response) ([]byte, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp.Body.Close()
	return b, err
}

// func UploadImage(str string) (path string, err error) {
// 	file, err := ioutil.TempFile("", "*.jpg")
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()

// 	decodebyte, err := base64.StdEncoding.DecodeString(str)
// 	if err != nil {
// 		return
// 	}

// 	file.Write(decodebyte)

// 	var obj map[string]interface{}
// 	req := httplib.Post(uri)
// 	req.PostFile("file", file.Name())
// 	req.Param("output", "json")
// 	err = req.ToJSON(&obj)
// 	if err != nil {
// 		return
// 	}
// 	path = prefix + obj["path"].(string)
// 	return
// }

// func AnalyzeBImage(uri, bimage string) (target string, err error) {
// 	var result map[string]interface{}
// 	req := httplib.Post(uri)
// 	param := make(map[string]interface{})
// 	param["bimage"] = bimage
// 	req.JSONBody(param)
// 	err = req.ToJSON(&result)
// 	if err != nil {
// 		pzap.Sugar.Error("调用ai接口报错: ", err)
// 		return
// 	}
// 	status := result["status"].(float64)
// 	if status == 1 {
// 		b, err := json.Marshal(result["target"])
// 		if err != nil {
// 			pzap.Sugar.Error("序列化 target 报错: ", err)
// 		} else {
// 			target = string(b)
// 		}
// 	}
// 	return
// }

// func AnalyzePImage(url, pimage string) ([]byte, error) {
// 	param := make(map[string]interface{})
// 	param["imgPath"] = pimage
// 	jstr, err := json.Marshal(&param)
// 	if err != nil {
// 		log.Warn(err)
// 		return nil, err
// 	}

// 	// reader := strings.NewReader(string(jstr))
// 	// log.Info("param: ", string(jstr))
// 	// log.Info("url: ", url)

// 	req, err := http.NewRequest("POST", url, bytes.NewReader(jstr))
// 	if err != nil {
// 		log.Warn(err)
// 		return nil, err
// 	}

// 	req.ContentLength = int64(len(jstr))
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	// resp, err := http.Post(algserv+path, "application/json", bytes.NewReader(jstr))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return nil, errors.New("算法服务器报错" + resp.Status)
// 	}

// 	return readRespse(resp)
// }
