package common

import (
	"github.com/Qingluan/jupyter/http"
	jupyter "github.com/Qingluan/jupyter/http"
	"github.com/Qingluan/merkur"
)

type Source interface {
	Iter() chan []string
}
type Scaner struct {
	target ScanTarget
	Source Source

	// ResNamePath
}

func NewScanner(target ScanTarget, source Source) *Scaner {
	return &Scaner{
		target: target,
		Source: source,
	}
}

func (scan *Scaner) Scan(callback func(name string, scan ScanTarget, resp *jupyter.SmartResponse) bool) {
	errPage := scan.target.TestErrPage()
	chs := make(chan int, scan.target.Num)
	breakChans := make(chan int)
	for args := range scan.Source.Iter() {
	LOOP:
		select {
		case <-breakChans:
			break LOOP
		default:
			if len(args) == 2 {
				go scan.scan(args[0], args[1], errPage, scan.target, &chs, &breakChans, callback)
				chs <- 1
			} else if len(args) == 1 {
				go scan.scan(args[0], args[0], errPage, scan.target, &chs, &breakChans, callback)
				chs <- 1
			}
		}
	}
}

func (scan *Scaner) scan(name, path, errPage string, target ScanTarget, chs, breakChans *chan int, callBack func(name string, target ScanTarget, resp *jupyter.SmartResponse) bool) {
	defer func() { <-*chs }()
	url := http.UrlJoin(target.Target, path)
	sess := jupyter.NewSession()
	if target.Proxy != "" {
		sess.SetProxyDialer(merkur.NewProxyDialer(target.Proxy))
	}
	sess.RandomeUA = target.RandomUA
	if resp, err := sess.Get(url); err == nil {
		if resp.Text() != errPage {
			if callBack(name, scan.target, resp) {
				*breakChans <- 1
				return
			}
		}
		Infor("Check: ", url)
	} else {
		InfoErr(err, url)
	}
}
