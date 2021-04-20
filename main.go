package main

import (
	"flag"

	"github.com/Qingluan/AllScannerInThis/admin"
	"github.com/Qingluan/AllScannerInThis/common"
)

var (
	ScanClass = ""
	Target    = ""
)

func main() {
	target := common.ScanTarget{}
	flag.StringVar(&target.Target, "t", "", "set target [url/ip]")
	flag.IntVar(&target.Num, "c", 10, "set channel async num")
	flag.StringVar(&target.ScanType, "type", "", "set scan type")
	flag.StringVar(&target.ScanClass, "scan", "admin", "set scan class type : [admin/dir/banner]")
	flag.StringVar(&target.Proxy, "proxy", "", "set proxy value")

	flag.Parse()

	switch target.ScanClass {
	case "admin":
		admin.ScanMain(target)
	}
}
