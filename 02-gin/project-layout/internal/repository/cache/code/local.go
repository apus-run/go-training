package code

import (
	"context"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"strconv"
)

type codeMemoryCache struct {
	local *freecache.Cache
}

func NewCodeMemoryCache() CodeCache {
	return &codeMemoryCache{
		local: freecache.NewCache(1024 * 1024 * 10),
	}
}

func (m *codeMemoryCache) Set(ctx context.Context, biz, phone, code string) error {
	res := m.setCode(biz, phone, code)

	switch res {
	case 0:
		return nil
	case -1:
		// 最近一次发送的验证码还没有过期, 不允许重发
		return ErrCodeSendTooMany
	default:
		// 发送验证码遇到未知错误, 系统错误, 比如说 -2, 是key冲突了
		// 其它响应码, 不知道是什么错误

		// TODO:: 这里应该有一个监控, 监控到这里的错误, 然后报警

		return ErrUnknownForCode
	}
}

func (m *codeMemoryCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res := m.verifyCode(biz, phone, code)

	switch res {
	case 0:
		return true, nil
	case -1:
		// 验证次数耗尽, 一般有可能是恶意攻击, 比如说有人在不停的尝试验证码
		return false, ErrCodeVerifyTooManyTimes
	default:
		// 验证码不一致
		return false, nil
	}
}

func (m *codeMemoryCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (m *codeMemoryCache) setCode(biz, phone, code string) int {
	// 设置缓存过期时间为 10 分钟
	expire := 10 * 60

	// key的规则, code:业务:手机号码
	// phone_code:login:152xxxxxxxx
	key := m.key(biz, phone)

	// 验证次数，我们一个验证码，最多重复三次，这个记录还可以验证几次
	// phone_code:login:152xxxxxxxx:cnt
	cntKey := key + ":cnt"

	// 过期时间
	ttl, err := m.local.TTL([]byte(key))
	// 如果为空, 说明是过期的 或者 key 不存在
	if ttl == uint32(0) && errors.Is(err, freecache.ErrNotFound) {
		err := m.local.Set([]byte(key), []byte(code), expire)
		if err != nil {
			return -2
		}
		err = m.local.Set([]byte(cntKey), []byte(strconv.Itoa(3)), expire)
		if err != nil {
			return -2
		}

		// 设置成功
		return 0
	}

	// ttl < 540, ttl < 540 是发了一个验证码，已经超过一分钟了，可以重新发送
	if ttl < uint32(540) {
		err := m.local.Set([]byte(key), []byte(code), expire)
		if err != nil {
			return -2
		}

		err = m.local.Set([]byte(cntKey), []byte(strconv.Itoa(3)), expire)
		if err != nil {
			return -2
		}

		// 设置成功
		return 0
	} else {
		// 发送太频繁, 已经发送了一个验证码，但是还不到一分钟
		return -1
	}
}

func (m *codeMemoryCache) verifyCode(biz, phone, expectedCode string) int {
	// 设置缓存过期时间为 10 分钟
	expire := 10 * 60

	// key的规则, code:业务:手机号码
	// phone_code:login:152xxxxxxxx
	key := m.key(biz, phone)

	// 验证次数，我们一个验证码，最多重复三次，这个记录还可以验证几次
	// phone_code:login:152xxxxxxxx:cnt
	cntKey := key + ":cnt"

	// 获取验证码验证次数
	cntBytes, err := m.local.Get([]byte(cntKey))
	if err != nil {
		return -2
	}
	cnt, err := strconv.Atoi(string(cntBytes))
	if err != nil {
		return 0
	}
	// 验证次数已经耗尽
	if cnt <= 0 {
		return -1
	}

	// 获取库里的验证码
	codeBytes, err := m.local.Get([]byte(key))
	if err != nil {
		return -2
	}
	// 如果验证码不存在
	if len(codeBytes) == 0 {
		return -2
	}

	if string(codeBytes) == expectedCode {
		// 把次数标记位 -1，认为验证码不可用
		err := m.local.Set([]byte(cntKey), []byte(strconv.Itoa(-1)), expire)
		if err != nil {
			return -2
		}
		return 0
	} else {
		// 有可能用户输错了验证码
		// 验证次数 -1
		cntStr := strconv.Itoa(cnt - 1)
		err = m.local.Set([]byte(cntKey), []byte(cntStr), expire)
		if err != nil {
			return -2
		}
		return -2
	}
}
