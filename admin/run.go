package admin

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/Qingluan/AllScannerInThis/common"
	jupyter "github.com/Qingluan/jupyter/http"
	"github.com/Qingluan/merkur"
)

var result []string

func ScanMain(target common.ScanTarget) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// threa := os.Args[2]
	thread := target.Num
	url := target.Target
	scan_type := strings.ToUpper(target.ScanType)
	var file string
	switch scan_type {
	case "PHP":
		file = php
	case "ASP":
		file = asp
	case "JSP":
		file = jsp
	case "ASPX":
		file = aspx
	case "DIR":
		file = dirs
	case "MDB":
		file = mdb
	default:
		file = dirs + "\n" + mdb
		os.Exit(1)
	}

	arr := strings.Split(file, "\n")
	//字典长度
	lens := len(arr)
	//每个线程分配任务数

	task := lens / thread

	ch := make(chan int)

	for i := 0; i < thread; i++ {
		go run(url, arr, i, task, ch, target.Proxy)
	}
	<-ch

}

func testErrorPage(url, proxy string) (code int, errorPage []byte) {
	url = url + "/someErrorPage"
	c, _, page := scandir(url, proxy)
	return c, page
}

func run(urls string, dir []string, tnum int, task int, ch chan int, proxy string) {
	_, ErrPage := testErrorPage(urls, proxy)
	for i := tnum*task + 1; i < (tnum*task)+task; i++ {
		dir[i-1] = strings.TrimSpace(dir[i-1])
		url := urls + dir[i-1]
		if strings.TrimSpace(url) == "" {
			continue
		}
		code, err, buf := scandir(url, proxy)
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(buf)) == strings.TrimSpace(string(ErrPage)) {
			continue
		}
		if code == 403 || code == 404 {
			common.Infor("Checking: ", dir[i-1])
		} else {
			common.Info(fmt.Sprintf("Found: %s [%d]!!!", dir[i-1], code))
			result = append(result, dir[i-1])
		}
	}
	ch <- 1
}

func scandir(url string, proxy string) (int, error, []byte) {
	session := jupyter.NewSession()
	if proxy != "" {
		dialer := merkur.NewProxyDialer(proxy)
		if dialer != nil {
			session.SetProxyDialer(dialer)
		}
	}
	resp, err := session.Get(url)
	var status int
	var page []byte
	if err != nil {
		status = 404
	} else {
		status = resp.StatusCode
		page, err = ioutil.ReadAll(resp.Body)
	}
	return status, err, page
}
