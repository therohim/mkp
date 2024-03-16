package utils

import (
	"fmt"
	"strings"
	"time"
)

func GenerateTrxCode(prefix string, sufix string) string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t := time.Now().In(loc)

	return fmt.Sprintf("%s-%d-%s", prefix, t.Unix(), strings.ToUpper(sufix))
}

func JakartaTime(t time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return t.UTC()
	}

	return t.UTC().In(loc)
}

func CheckPointerValue(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}
