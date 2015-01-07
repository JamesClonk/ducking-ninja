package main

import (
	"log"
	"net/http"

	"github.com/cihub/seelog"
	"github.com/codegangsta/negroni"
)

func init() {
	testConfig := `
<seelog>
	<formats>
		<format id="ducking-ninja" format="%Date %Time;%Msg%n"/>
	</formats>
	<outputs formatid="ducking-ninja">
		<rollingfile type="size" filename="./logs/ducking-ninja.log" maxsize="1048576" maxrolls="10" />
	</outputs>
</seelog>
`
	logger, err := seelog.LoggerFromConfigAsString(testConfig)
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
