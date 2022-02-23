package utils

import (
	"fmt"
	"path/filepath"
	"runtime"

	"go.uber.org/zap"
)

// 增加数据的错误
func InsertErrors(err error, count int64, datalen int64) error {
	where := caller(1, false)
	var e string
	if err != nil {
		e = where + err.Error() + ","
	} else {
		e = where
	}
	if count != datalen {
		return fmt.Errorf(e + fmt.Sprintf("count error, count: %d, datalen: %d", count, datalen))
	}
	return fmt.Errorf(e)
}

// 增加数据的错误
func UpdateErrors(err error, count int64, datalen int64) error {
	where := caller(1, false)
	var e string
	if err != nil {
		e = where + err.Error() + ","
	} else {
		e = where
	}
	if count != datalen {
		return fmt.Errorf(e + fmt.Sprintf("count error, count: %d, datalen: %d", count, datalen))
	}
	return fmt.Errorf(e)
}

// 格式化错误
func Errors(err error, msg string) error {
	where := caller(1, false)
	zap.L().Error(msg + err.Error())
	return fmt.Errorf(where + msg + ": " + err.Error())
}

// 获取源代码行数
func caller(calldepth int, short bool) string {
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file = "???"
		line = 0
	} else if short {
		file = filepath.Base(file)
	}

	return fmt.Sprintf("%s:%d: ", file, line)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
