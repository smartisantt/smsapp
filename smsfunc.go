package smsapp

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func CanSend(phone string) (bool, error) {
	if len(phone) != PHONE_LEN {
		return false, errors.New("手机号长度不对")
	}
	reg := `^1[3456789]\d{9}$`
	rgx := regexp.MustCompile(reg)

	if !rgx.MatchString(phone) {
		return false, errors.New("手机号码格式不对")
	}
	return true, nil
}

func GenerateSmsCode(codeLen int) string {
	rand.Seed(time.Now().UnixNano())
	var code string
	for i := 0; i < codeLen; i++ {
		c := rand.Intn(10)
		code += strconv.Itoa(c)
	}
	return code
}
