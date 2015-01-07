package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/cihub/seelog"
	"github.com/codegangsta/negroni"
)

var (
	logFilename = "./logs/ducking-ninja.log"
)

func init() {
	config := `
<seelog>
	<formats>
		<format id="ducking-ninja" format="%Date %Time;%Msg%n"/>
	</formats>
	<outputs formatid="ducking-ninja">
		<rollingfile type="size" filename="` + logFilename + `" maxsize="100000" maxrolls="10" />
	</outputs>
</seelog>
`
	logger, err := seelog.LoggerFromConfigAsString(config)
	if err != nil {
		log.Fatal(err)
	}
	seelog.ReplaceLogger(logger)
}

func logging() negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		haddr := req.Header.Get("X-Real-IP")
		if haddr == "" {
			haddr = req.Header.Get("X-Forwarded-For")
		}
		seelog.Infof("%v;%v;%v;%v", req.Method, req.URL.Path, req.RemoteAddr, haddr)
		next(w, req)
	}
}

func getLogs() []string {
	bytes, err := ioutil.ReadFile(logFilename)
	if err != nil {
		log.Fatal(err)
	}

	data := strings.Split(string(bytes), "\n")
	sort.Sort(sort.Reverse(sort.StringSlice(data)))
	return data
}
