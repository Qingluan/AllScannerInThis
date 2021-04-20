package admin

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/Qingluan/AllScannerInThis/common"
	jupyter "github.com/Qingluan/jupyter/http"
	"github.com/Qingluan/merkur"
)

var result []string

func ScanMain(target common.ScanTarget) {
	common.Info("thread:", target.Num)
	common.Info("type  :", target.ScanType)
	common.Info("target:", target.Target)

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
		// os.Exit(1)
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

func testErrorPage(url, proxy string) (code int, errorPage string) {
	// url = url + "/someErrorPageAbabababababababababbababa"

	c, _, page := scandir(url, proxy)
	common.Info("tset Err Pge:", url, " len:", len(page))
	return c, page
}

func J(u, e string) string {
	e = strings.TrimSpace(e)
	if strings.HasSuffix(u, "/") {
		if strings.HasPrefix(e, "/") {
			return u + e[1:]
		} else {
			return u + e
		}
	} else {
		if strings.HasPrefix(e, "/") {
			return u + e
		} else {
			return u + "/" + e

		}
	}
}

func ErrPage(url, e string) string {
	if strings.HasPrefix(e, "/") {
		e = "asdfdsfasdfaeroerrreea" + e
	} else {
		e = "asdfasdffasdgasdgdasgasdgsd/" + e
	}
	return J(url, e)
}

func run(urls string, dir []string, tnum int, task int, ch chan int, proxy string) {
	_, ErrHash := testErrorPage(ErrPage(urls, dir[0]), proxy)
	ai := len(dir)
	// common.Info("Err Hash:", ErrHash)
	for i := tnum*task + 1; i < (tnum*task)+task; i++ {
		dir[i-1] = strings.TrimSpace(dir[i-1])
		url := J(urls, dir[i-1])
		if strings.TrimSpace(url) == "" {
			continue
		}
		code, err, hash := scandir(url, proxy)
		if err != nil {
			continue
		}
		if hash == ErrHash || len(hash) == 0 {
			continue
		}
		if code == 403 || code == 404 {
			common.Infor(fmt.Sprintf("[%d/%d]", i, ai), "Checking: ", dir[i-1])
		} else {
			common.Info(fmt.Sprintf("Found: %s length:%d [%d] !!!", common.Green(J(urls, dir[i-1])), len(hash), code))
			result = append(result, dir[i-1])
		}
	}
	ch <- 1
}

func scandir(url string, proxy string) (int, error, string) {
	session := jupyter.NewSession()
	if proxy != "" {
		dialer := merkur.NewProxyDialer(proxy)
		if dialer != nil {
			session.SetProxyDialer(dialer)
		}
	}
	resp, err := session.Get(url)
	var status int
	var page string
	if err != nil {
		status = 404
		return status, err, ""
	} else {
		status = resp.StatusCode
		page = resp.Text()
		// page, err = ioutil.ReadAll(resp.Body)

	}
	// fmt.Println(url, "----\n---", page)
	return status, err, page
}
