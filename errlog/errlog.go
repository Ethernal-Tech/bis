package errlog

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func Println(err error) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	funcNameEls := strings.Split(funcName, ".")
	fmt.Println(funcName, file, funcNameEls)
	if len(funcNameEls) == 2 {
		funcName = funcNameEls[1] + "()"
	} else {
		funcName = funcNameEls[1] + "." + funcNameEls[2] + "()"
	}

	dateTime := strings.Split(time.Now().String(), ".")[0]

	// (if) for docker container
	// (else) for general purposes
	if strings.Count(file, "app") == 2 {
		file = "./app" + strings.Split(file, "/app")[2]
	} else {
		file = "./app" + strings.Split(file, "/app")[1]
	}

	fmt.Printf("\033[31m%v ERROR: \033[0m%v:%v [%v]: \033[32m%v\033[0m\n", dateTime, file, line, funcName, err)
}
