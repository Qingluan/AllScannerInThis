package main

import (
	"flag"
	"log"

	"github.com/Qingluan/AllScannerInThis/admin"
	"github.com/Qingluan/AllScannerInThis/banner"
	"github.com/Qingluan/AllScannerInThis/common"
	"github.com/Qingluan/AllScannerInThis/dir"
	"github.com/Qingluan/jupyter/http"
	"github.com/Qingluan/merkur"
)

var (
	ScanClass = ""
	Target    = ""
)

func main() {
	target := common.ScanTarget{}
	flag.StringVar(&target.Target, "u", "", "set target [url/ip]")
	flag.IntVar(&target.Num, "c", 10, "set channel async num")
	flag.StringVar(&target.ScanType, "t", "", "set scan type \n\tadmin:[php,jsp,dir,cgi,mdb,asp,aspx] default:dir+cgi \n\tbanner:[all] default: ")
	flag.StringVar(&target.ScanClass, "scan", "admin", "set scan class type : [admin/dir/banner]")
	flag.StringVar(&target.Proxy, "proxy", "", "set proxy value")
	flag.BoolVar(&target.RandomUA, "ua", false, "true : random user-agent for scanner")
	flag.Parse()
	common.Info("Scan:", common.Yellow(target.ScanClass))
	switch target.ScanClass {
	case "admin":
		admin.ScanMain(target)
	case "dir":
		dir.ScanMain(target)
	case "banner":
		banner.ScanMain(target)
	case "md5":
		sess := http.NewSession()
		if target.Proxy != "" {
			sess.SetProxyDialer(merkur.NewProxyDialer(target.Proxy))
		}
		sess.RandomeUA = target.RandomUA
		if res, err := sess.Get(target.Target); err != nil {
			log.Fatal(err)
		} else {
			common.InfoOk("md5: ", common.Green(res.Md5()))
		}

	}
}
