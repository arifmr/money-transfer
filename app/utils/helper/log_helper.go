package helper

import (
	"fmt"
	"money-transfer/app/models"
	"net/http"
	"runtime"
	"strings"
)

var DefaultStatusText = map[int]string{
	http.StatusInternalServerError: "Terjadi Kesalahan, Silahkan Coba lagi Nanti",
	http.StatusNotFound:            "Data tidak Ditemukan",
	http.StatusBadRequest:          "Ada kesalahan pada request data, silahkan dicek kembali",
}

func WriteLog(err error, errorCode int, message interface{}) *models.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		output := &models.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
		}

		output.SystemMessage = err.Error()
		if message == nil {
			output.Message = DefaultStatusText[errorCode]
			if output.Message == "" || output.Message == nil {
				output.Message = http.StatusText(errorCode)
			}
		} else {
			output.Message = message
		}

		output.Line = fmt.Sprintf("%d", line)
		output.Filename = fmt.Sprintf("%s:%d", file, line)
		output.Function = runtime.FuncForPC(pc).Name()
		return output
	}

	return nil
}
