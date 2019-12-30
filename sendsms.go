package smsapp

import (
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const PhoneRedisPrefix = "SMS:PHONE:"

var PHONE_LEN = 11
var MAX_CODE_LEN = 6

func (s *SmsOption) SendSms(phone, code, content string) error {
	ok, err := CanSend(phone)
	if !ok {
		return err
	}

	if s.Default {
		// 发送默认短信数据
		code = GenerateSmsCode(MAX_CODE_LEN)
		content = fmt.Sprintf("验证码：【%s】，此验证码只用于登录您的账户，请勿提供给别人。", code)
	}

	if !s.Debug {
		err := s.sendSms(phone, content)
		if err != nil {
			return nil
		}
	}

	// 存储到redis当中
	key := PhoneRedisPrefix + phone
	s.R.Set(key, code, 10*time.Minute)
	return nil
}

func (s *SmsOption) sendSms(phone, content string) error {
	v := url.Values{}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	v.Set("account", s.Account)
	v.Set("password", s.Passwd)
	v.Set("mobile", phone)
	v.Set("content", content)
	v.Set("time", now)
	v.Set("format", "json")

	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, s.Url, body)
	if err != nil {
		fmt.Println("生成发送短信请求失败 ", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("给手机号 ", phone, "发送短信失败 ", err.Error())
		return err
	}

	fmt.Println("给手机号: ", phone, "发送的信息是:", content)

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取发送短信正码的返回消息失败 ", err.Error())
		return err
	}

	var smsRet SendSmsResponse
	err = jsoniter.Unmarshal(ret, &smsRet)
	if err != nil {
		fmt.Println("反序列化短信发送商返回的数据时失败 ", err.Error())
		return err
	}

	if smsRet.Code != SMS_RESPONSE_OK {
		fmt.Println("发送短信失败 ", smsRet.Code, smsRet.Msg)
		return fmt.Errorf("发送短信失败 %s", smsRet.Msg)
	}

	return nil
}
