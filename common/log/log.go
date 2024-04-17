package log

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/chipong/template/common/util"
)

func Printf(c *gin.Context, format string, v ...any) {
	var variables []any

	_, file, line, _ := runtime.Caller(1)
	splitStr := strings.Split(file, "/")

	fileName := splitStr[len(splitStr)-1]
	variables = append(variables, fileName)
	variables = append(variables, strconv.Itoa(line))

	var fBuffer bytes.Buffer
	fBuffer.WriteString("%s:%s")
	if c != nil {
		uid, _, _ := util.GetHeaderUidSessionKey(c)

		fBuffer.WriteString(", uid:'%s'")
		variables = append(variables, uid)
	}

	fBuffer.WriteString(", ")
	fBuffer.WriteString(format)

	variables = append(variables, v...)

	log.Printf(fBuffer.String(), variables...)
}

func Println(c *gin.Context, v ...any) {
	var variables []any

	_, file, line, _ := runtime.Caller(1)
	splitStr := strings.Split(file, "/")
	fileName := splitStr[len(splitStr)-1]

	var fBuffer bytes.Buffer
	fBuffer.WriteString(fmt.Sprintf("%s:%s", fileName, strconv.Itoa(line)))

	if c != nil {
		uid, _, _ := util.GetHeaderUidSessionKey(c)
		fBuffer.WriteString(fmt.Sprintf(", uid:'%s', ", uid))
	}
	variables = append(variables, fBuffer.String())
	variables = append(variables, v...)

	log.Println(variables...)
}
