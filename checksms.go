package smsapp

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

func (s *SmsOption) CheckSms(phone, code string) (bool, error) {
	key := PhoneRedisPrefix + phone
	dbcode, err := s.R.Get(key).Result()

	if err != nil && err != redis.Nil {
		emsg := fmt.Sprintf("从 redis 读取手机号[%s]对应的 code 失败, %s", phone, err.Error())
		fmt.Println(emsg)
		return false, errors.New(emsg)
	}

	if err == redis.Nil {
		return false, errors.New("验证码已失效")
	}

	if dbcode != code {
		return false, nil
	}

	s.R.Del(key)
	return true, nil
}
