package banner

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/Qingluan/AllScannerInThis/common"
	"github.com/Qingluan/jupyter/http"
	jupyter "github.com/Qingluan/jupyter/http"
	"github.com/Qingluan/merkur"
)

type Version struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Option  string `json:"option"`
	// Hit     int    `json:"hit"`
}
type BannerRes struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Option  string `json:"option"`
	// Hit     int      `json:"hit"`
	Ver *Version `json:"version"`
}
type Arg struct {
	Name   string
	Target string
	Ban    BannerRes
	Wait   int
}

var (
	Versions = map[string]Version{}
)

func ScanMain(target common.ScanTarget) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	thread := target.Num
	wait := target.Wait
	ch := make(chan int, thread)
	errPage := target.TestErrPage()
	okchan := make(chan string)
	c := 0
	scanRes := LoadRes()
	Versions = LoadVersion()
	ac := len(scanRes)
	for k, bs := range scanRes {
		// SCANLOOP:
		select {
		case FoundPath := <-okchan:
			if target.ScanType != "all" {
				common.InfoOk("Not \"all\" mode so break ! check: ", common.Green(http.UrlJoin(target.Target, FoundPath)))
				os.Exit(0)
			}
		default:
			go scanBanner(k, target.Target, target.Proxy, errPage, target.RandomUA, bs, c, ac, wait, &okchan, &ch)
			ch <- 1
			c++
		}
	}
	// <-ch
}

func scanBanner(name, target, proxy, errPage string, randomua bool, banners []BannerRes, i, ai, wait int, ifokch *chan string, ch *chan int) {
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
	var Found BannerRes

	// paths := []string{}
	// if strings.Contains(banner.Path, "|") {
	// 	for _, path := range strings.Split(banner.Path, "|") {
	// 		paths = append(paths, http.UrlJoin(target, strings.TrimSpace(path)))
	// 	}
	// } else {
	// 	paths = append(paths, http.UrlJoin(target, banner.Path))

	// }
	size := 0
	for _, banner := range banners {
		path := http.UrlJoin(target, banner.Path)

		if res, err := sess.Get(path); err != nil {
			if strings.Contains(err.Error(), "too many open files") {
				common.InfoErr(err, path)
			}

		} else {

			if res.StatusCode/100 > 3 {
				res.Body.Close()
				return
			}
			if text := res.Text(); text == errPage {

				res.Body.Close()
				return
			}

			switch banner.Option {
			case "keyword":
				if res.Search(strings.ToLower(banner.Content), true) {
					Found = banner
					size = len(res.Html())
					res.Body.Close()
					break
				}
			case "re":
				rec := regexp.MustCompile(banner.Content)
				if version := rec.FindString(string(res.Html())); version != "" {
					Found = banner
					size = len(res.Html())
					res.Body.Close()
					break
				}
			case "md5":
				if res.Md5() == banner.Content {
					Found = banner

					size = len(res.Html())
					res.Body.Close()
					break
				}
			}
		}
	}

	if Found.Name != "" {

		common.Info(common.Red("*"), common.Yellow("*"), common.Green("* ", Found.Name, " *"), common.Yellow("*"), common.Red("* "), common.Cyan(Found.Content), " in ", common.Blue(Found.Path), " size:", common.Blue(size), "       ")

		if ver, ok := Versions[Found.Name]; ok {
			ScanVersion(name, target, ver, sess)
		}
		*ifokch <- Found.Path
	} else {
		common.Infor(fmt.Sprintf("[%5d/%5d] Checking : %s", i, ai, name))
	}

}

func ScanVersion(name, target string, ver Version, sess *jupyter.Session) {
	path2 := http.UrlJoin(target, ver.Path)
	if res, err := sess.Get(path2); err == nil {
		switch ver.Option {
		case "re":
			rec := regexp.MustCompile(ver.Content)
			if version := rec.FindString(string(res.Html())); version != "" {
				common.Info(" ", common.Yellow("└"), common.Green("─ "), ver.Path, " Version: ", common.Green(version), "          ")
			}
		}
	}

}
