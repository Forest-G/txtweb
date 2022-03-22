package redisdata

import (
	"errors"
	"fmt"
	"txt/pkg/config"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

var (
	conn redis.Conn
	err  error
)

func Init() {
	go func() {
		setdb := redis.DialDatabase(config.C.Redisdb)
		setpassword := redis.DialPassword(config.C.Redispassword)
		// 建立连接
		conn, err = redis.Dial("tcp", config.C.Redisport, setdb, setpassword)

		if err != nil {
			fmt.Println("redis.Dial err=", err)
			logrus.Error(conn.Err())
			return
		}
		fmt.Println("redis连接成功")
	}()
}

//把tocken 和邮箱作为键值对存起来
func Ruserlogin(username string, token string) error {
	// 通过go向redis写入数据 string [key - value]
	_, err = conn.Do("set", username, token)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = conn.Do("expire", username, 86400*15) //一天是86400秒
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

//查询用户toker
func FindRuserlogin(username string, token string) error {
	// 读取数据 获取名字
	r, err := redis.String(conn.Do("GET", username))
	//err就代表着值是空的******************
	if err != nil {
		logrus.Error(err)
		return err
	}
	if r == token {
		return nil
	} else {
		return errors.New("用户身份已过期")
	}
}

//把验证码 和邮箱作为键值对存起来
func EmailCode(email string, icode string) error {
	// 通过go向redis写入数据 string [key - value]
	_, err = conn.Do("SET", email, icode)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = conn.Do("expire", email, 60) //一天是86400秒
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

//查询邮箱验证码
func FindVcode(email string) string {
	// 读取数据 获取名字
	r, err := redis.String(conn.Do("Get", email))
	//err就代表着值是空的******************
	if err != nil {
		logrus.Error(err)
		return "验证码已过期"
	}
	return r
}

//查询邮箱验证码
func FindVcode_Validation(email string, icode string) error {
	// 读取数据 获取名字
	r, err := redis.String(conn.Do("Get", email))
	//err就代表着值是空的******************
	if err != nil {
		logrus.Error(err)
		return errors.New("验证码已过期")
	}
	if r == icode {
		return nil
	} else {
		return errors.New("验证失败")
	}
}
