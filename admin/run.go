package admin

import (
	"log"

	"github.com/Qingluan/AllScannerInThis/asset"
	"github.com/Qingluan/AllScannerInThis/common"
	"github.com/Qingluan/jupyter/http"
)

func ScanMain(target common.ScanTarget) {
	buf, err := asset.Asset("Res/admin/admin.txt")
	if err != nil {
		log.Fatal("Load err:", err)
	}
	bSource := common.NewBufSouce(buf)
	scanner := common.NewScanner(target, bSource)
	scanner.Scan(func(name string, target common.ScanTarget, resp *http.SmartResponse) bool {
		if resp.StatusCode/100 < 4 {
			common.InfoOk("Found :", common.Green(name), "Size:", common.Blue(len(resp.Html())), "MD5:", common.Yellow(resp.Md5()))
		}
		return false
	})
}
