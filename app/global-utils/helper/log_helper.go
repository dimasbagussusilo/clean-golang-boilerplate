package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"appsku-golang/app/global-utils/constants"
	"appsku-golang/app/global-utils/model"

	"github.com/sirupsen/logrus"
)

// var DefaultStatusText = map[int]string{
// 	http.StatusInternalServerError: "Something went wrong, please try again later",
// 	http.StatusNotFound:            "data not found",
// }

var DefaultStatusText = map[int]string{
	http.StatusInternalServerError: "Terjadi Kesalahan, Silahkan Coba lagi Nanti",
	http.StatusNotFound:            "Data tidak Ditemukan",
	http.StatusBadRequest:          "Ada kesalahan pada request data, silahkan dicek kembali",
}

func WriteLog(err error, errorCode int, message interface{}) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		output := &model.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
		}
		outputForPrint := &model.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
			Line:       fmt.Sprintf("%d", line),
			Filename:   file,
			Function:   funcName,
		}

		output.SystemMessage = err.Error()
		if message == nil {
			output.Message = DefaultStatusText[errorCode]
			if output.Message == "" {
				output.Message = http.StatusText(errorCode)
				outputForPrint.Message = http.StatusText(errorCode)
			}
		} else {
			output.Message = message
			outputForPrint.Message = message
		}
		if errorCode == http.StatusInternalServerError {
			output.Line = fmt.Sprintf("%d", line)
			output.Filename = file
			output.Function = funcName
		}

		logForPrint := map[string]interface{}{}
		_ = DecodeMapType(outputForPrint, &logForPrint)

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(logForPrint).Error(err)
		return output
	}

	return nil
}

func WriteLogWithContext(ctx context.Context, err error, errorCode int, message interface{}) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		output := &model.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
		}
		outputForPrint := &model.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
			Line:       fmt.Sprintf("%d", line),
			Filename:   fmt.Sprintf("%s:%d", file, line),
			Function:   funcName,
		}

		output.SystemMessage = err.Error()
		if message == nil {
			output.Message = DefaultStatusText[errorCode]
			if output.Message == "" {
				output.Message = http.StatusText(errorCode)
				outputForPrint.Message = http.StatusText(errorCode)
			}
		} else {
			output.Message = message
			outputForPrint.Message = message
		}
		if errorCode == http.StatusInternalServerError {
			output.Line = fmt.Sprintf("%d", line)
			output.Filename = file
			output.Function = funcName
		}

		logForPrint := map[string]interface{}{}
		_ = DecodeMapType(outputForPrint, &logForPrint)
		_, logForPrint[constants.RequestId] = GetRequestIDContext(ctx)

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(logForPrint).Error(err)
		return output
	}

	return nil
}

func NewWriteLog(errorLog model.ErrorLog) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {

		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		var output *model.ErrorLog

		if errorLog.StatusCode == 500 {
			output = &model.ErrorLog{
				Line:          fmt.Sprintf("%d", line),
				Filename:      file,
				Function:      funcName,
				Message:       errorLog.Message,
				SystemMessage: errorLog.Err.Error(),
				Err:           errorLog.Err,
				StatusCode:    errorLog.StatusCode,
			}
		} else if errorLog.StatusCode == http.StatusUnauthorized || errorLog.StatusCode == http.StatusForbidden || errorLog.StatusCode == http.StatusNotFound || errorLog.StatusCode == http.StatusConflict || errorLog.StatusCode == http.StatusBadRequest || errorLog.StatusCode == http.StatusUnprocessableEntity || errorLog.StatusCode == http.StatusExpectationFailed {
			output = &model.ErrorLog{
				Message:       errorLog.Message,
				SystemMessage: errorLog.SystemMessage,
				Err:           errorLog.Err,
				StatusCode:    errorLog.StatusCode,
			}

		} else {
			output = &model.ErrorLog{
				Message:       []string{"Error tidak dikenali"},
				SystemMessage: []string{"Undefined error"},
				Err:           errorLog.Err,
				StatusCode:    errorLog.StatusCode,
			}
		}

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(log).Error(errorLog.Message)
		return output
	}

	return nil
}

func WriteLogConsumer(consumerName string, consumerTopic string, consumerPartition int, consumerOffset int64, consumerKey string, err error, errorCode int, message interface{}) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		var output *model.ErrorLog

		if errorCode == 500 {
			output = &model.ErrorLog{
				Line:              fmt.Sprintf("%d", line),
				Filename:          file,
				Function:          funcName,
				Message:           message,
				SystemMessage:     err.Error(),
				ConsumerName:      consumerName,
				ConsumerTopic:     consumerTopic,
				ConsumerPartition: consumerPartition,
				ConsumerOffset:    consumerOffset,
				ConsumerKey:       consumerKey,
				Err:               err,
				StatusCode:        errorCode,
			}
		}

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		logrus.WithFields(log).Error(message)
		return output
	}

	return nil
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
