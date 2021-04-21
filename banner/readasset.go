package banner

import (
	"encoding/json"
	"log"

	"github.com/Qingluan/AllScannerInThis/asset"
)

func LoadRes() map[string][]BannerRes {
	buf, err := asset.Asset("Res/banner/banner.json")
	if err != nil {
		log.Fatal("load fail:", err)
	}
	array := make(map[string]interface{})
	if err := json.Unmarshal(buf, &array); err != nil {
		log.Fatal("load json err:", err)
	}
	data := make(map[string][]BannerRes)

	for path, bf := range array {
		bfs := bf.([]interface{})
		bans := []BannerRes{}
		for _, bff := range bfs {
			buf, _ := json.Marshal(bff)
			e := BannerRes{}
			json.Unmarshal(buf, &e)
			e.Path = path
			bans = append(bans, e)
		}
		data[path] = bans
	}
	return data
}
