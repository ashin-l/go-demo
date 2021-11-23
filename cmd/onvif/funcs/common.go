package funcs

import (
	"io/ioutil"
	"net/http"

	"github.com/ashin-l/go-demo/pkg/config"
)

type Configuration struct {
	Ip       string
	Username string
	Password string
}

var Conf = &Configuration{}

func init() {
	err := config.LoadConf("app.yml", "yaml", Conf)
	if err != nil {
		panic(err)
	}
}

func readResponse(resp *http.Response) string {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	return string(b)
}
