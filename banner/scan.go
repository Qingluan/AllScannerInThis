package banner

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"

	"github.com/Qingluan/AllScannerInThis/asset"
	"github.com/Qingluan/AllScannerInThis/common"
	"github.com/Qingluan/jupyter/http"
	jupyter "github.com/Qingluan/jupyter/http"
	"github.com/Qingluan/merkur"
)

type Version struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Option  string `json:"option"`
	Hit     int    `json:"hit"`
}
type BannerRes struct {
	Path    string   `json:"path"`
	Content string   `json:"content"`
	Option  string   `json:"option"`
	Hit     int      `json:"hit"`
	Ver     *Version `json:"version"`
}
type Arg struct {
	Name   string
	Target string
	Ban    BannerRes
	Wait   int
}

func ScanMain(target common.ScanTarget) {
	buf, err := asset.Asset("Res/banner.json")
	if err != nil {
		log.Fatal("Load Banner.json failed", err)
	}
	scanRes := make(map[string]interface{})
	if err := json.Unmarshal(buf, &scanRes); err != nil {
		log.Fatal(err)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	thread := target.Num
	// ch := make(chan Arg, thread)
	// url := target.Target
	wait := target.Wait
	ch := make(chan int, thread)
	errPage := target.TestErrPage()
	okchan := make(chan int)
	c := 0
	ac := len(scanRes)
	for k, bi := range scanRes {
	SCANLOOP:
		select {
		case <-okchan:
			if target.ScanType != "all" {
				break SCANLOOP
			}
		default:

			buf, _ := json.Marshal(bi)
			b := BannerRes{}
			err := json.Unmarshal(buf, &b)
			if err != nil {
				log.Fatal("Trans buf -> Banner :", err, k)
			}
			go scanBanner(k, target.Target, target.Proxy, errPage, target.RandomUA, b, c, ac, wait, &okchan, &ch)
			ch <- 1
			c++
		}
	}
	// <-ch
}

func scanBanner(name, target, proxy, errPage string, randomua bool, banner BannerRes, i, ai, wait int, ifokch, ch *chan int) {
	defer func() { <-*ch }()
	sess := jupyter.NewSession()
	if proxy != "" {
		if dailer := merkur.NewProxyDialer(proxy); dailer != nil {
			sess.SetProxyDialer(dailer)
		}
	}
	if randomua {
		sess.RandomeUA = true
	}
	Found := false
	paths := []string{}
	if strings.Contains(banner.Path, "|") {
		for _, path := range strings.Split(banner.Path, "|") {
			paths = append(paths, http.UrlJoin(target, strings.TrimSpace(path)))
		}
	} else {
		paths = append(paths, http.UrlJoin(target, banner.Path))

	}
	for _, path := range paths {

		if res, err := sess.Get(path); err != nil {

		} else {
			if text := res.Text(); text == errPage {
				return
			}

			switch banner.Option {
			case "keyword":
				if res.Search(strings.ToLower(banner.Content), true) {
					Found = true
					break
				}
			case "md5":
				if res.Md5() == banner.Content {
					Found = true
					break
				}
			}
		}
	}

	if Found {

		common.Info(common.Green(name), " in ", common.Blue(banner.Path), "       ")
		if banner.Ver != nil {
			ScanVersion(name, target, banner, sess)
		}
	} else {
		common.Infor(fmt.Sprintf("[%5d/%5d] Checking : %s", i, ai, name))
	}

}

func ScanVersion(name, target string, banner BannerRes, sess *jupyter.Session) {

	path2 := http.UrlJoin(target, banner.Ver.Path)
	if res, err := sess.Get(path2); err == nil {
		switch banner.Ver.Option {
		case "re":
			rec := regexp.MustCompile(banner.Ver.Content)
			if version := rec.FindString(string(res.Html())); version != "" {
				common.Info("  ", name, "Version:", common.Green(version), "          ")
			}
		}
	}
}
