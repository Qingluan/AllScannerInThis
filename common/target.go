package common

import "encoding/json"

type ScanTarget struct {
	Target    string            `json:"target"`
	Num       int               `json:"num"`
	Wait      int               `json:"wait"`
	ScanType  string            `json:"scanType"`
	ScanClass string            `json:"class"`
	Kargs     map[string]string `json:"kargs"`
	Proxy     string            `json:"proxy"`
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
