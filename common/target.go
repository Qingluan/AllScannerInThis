package common

import (
	"encoding/json"

	"github.com/Qingluan/jupyter/http"
	jupyter "github.com/Qingluan/jupyter/http"
	"github.com/Qingluan/merkur"
)

type ScanTarget struct {
	Target    string            `json:"target"`
	Num       int               `json:"num"`
	Wait      int               `json:"wait"`
	ScanType  string            `json:"scanType"`
	ScanClass string            `json:"class"`
	Kargs     map[string]string `json:"kargs"`
	Proxy     string            `json:"proxy"`
	RandomUA  bool              `json:"random-ua"`
}

func (target ScanTarget) String() string {
	b, _ := json.Marshal(&target)
	return string(b)
}

func String2Target(buf string) ScanTarget {
	t := ScanTarget{}
	json.Unmarshal([]byte(buf), &t)
	return t
}

func (target ScanTarget) TestErrPage() string {
	url := http.UrlJoin(target.Target, "/asomegSomeErrorPageasdgasdgsdg/asdgdg")

	sess := jupyter.NewSession()
	if target.Proxy != "" {
		sess.SetProxyDialer(merkur.NewProxyDialer(target.Proxy))
	}
	if res, err := sess.Get(url); err != nil {
		return ""
	} else {
		return res.Text()
	}

}
