package helper

import (
	"appsku-golang/app/global-utils/constants"
	"context"
	"regexp"
	"strconv"
	"strings"
)

func GetRequestIDContext(ctx context.Context) (string, interface{}) {
	val := ctx.Value(constants.RequestId)
	if val == nil {
		val = ""
	}

	return constants.RequestId, val
}

func SetRequestIDToContext(c context.Context, reqId string) context.Context {
	ctx := context.WithValue(c, constants.RequestId, reqId)
	return ctx
}

func DefineIDCountryCode(number string) (string, bool) {
	pattern := regexp.MustCompile(`^(\+\d{1,3})?\d+$`)

	number = strings.ReplaceAll(number, "+", "")

	if len(number) >= 8 {
		if number[:1] == "0" {
			number = "62" + number[1:]
		}

		if pattern.MatchString(number) {
			return number, true
		}
	}

	return number, false
}

func DefineMobileNumber(number string) (string, bool) {

	if len(number) > 8 {
		firstLine := number[:2]

		if firstLine[:2] == "62" {
			number = "0" + number[2:]
		}

		if firstLine[:1] == "8" {
			numbx, _ := strconv.Atoi(number[:1])
			if numbx > 0 {
				number = "0" + number
			}
		}
	}

	return number, false
}
