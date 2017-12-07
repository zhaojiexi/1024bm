package com

import (
	"fmt"
	"log"
	"flag"
	"time"
	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
	host = flag.String("106.75.13.30", "106.75.13.30:6379", "")
	/*
	redis设置了AUTH安全验证才需要passd,一般普通写法为c, err := redis.Dial("tcp", "127.0.0.1:6379")
	*/
	password = flag.String("passwd", "Zxcasdqwe123!@#", "")
)

const redigoSession string = "redigoSessionId"

//测试struct
type MyUser struct {
	UserName  string
	UserPhone string
}

//连接池
type Pool struct {
	//Dial 是创建链接的方法
	Dial func() (redis.Conn, error)

	//TestOnBorrow 是一个测试链接可用性的方法
	TestOnBorrow func(c redis.Conn, t time.Time) error

	//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
	MaxIdle int

	//最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
	MaxActive int

	/*
	最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
	应该设置一个比redis服务端超时时间更短的时间
	*/
	IdleTimeout time.Duration

	/*
	当链接数达到最大后是否阻塞，如果不的话，达到最大后返回错误
	如果Wait被设置成true，则Get()方法将会阻塞
	*/
	Wait bool
}

/*
初始化一个pool
host:ip、端口
passwd：密码
返回连接池指针
*/
func RedigoPool(host, passwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   5,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				fmt.Println(err)
				log.Fatalf("redigo->RedigoPool->redis.Dial()初始化连接池时报错: %s\n", err)
			}
			if _, err := c.Do("AUTH", passwd); err != nil {
				c.Close()

			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

//累加计数,每执行一次加(自增长)
func SetCount(countKey string){
	flag.Parse()
	pool = RedigoPool(*host, *password)

	conn := pool.Get()//从连接池获取连接
	defer conn.Close()//用完后放回连接池

	//redis操作
	value, err := conn.Do("INCR", countKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	value1, err := conn.Do("EXPIRE", countKey,30)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)//返回ok
	fmt.Println(value1)//返回ok
}

func SetMap() {
	flag.Parse()
	pool = RedigoPool(*host, *password)

	conn := pool.Get() //从连接池获取连接
	defer conn.Close() //用完后放回连接池

	myuser := map[string]*MyUser{
		"CCC": &MyUser{UserName: "caimin", UserPhone: "13162578783"},
		"DDD": &MyUser{UserName: "caimin", UserPhone: "13162578783"},
	}

	//保存Map
	for sym, row := range myuser {
		if _, err := conn.Do("HMSET", redis.Args{sym}.AddFlat(row)...); err != nil {
			log.Fatal(err)
		}
	}

	//获取所有Map下Key的值（命令 hgetall AAA）
	for sym := range myuser {
		values, err := redis.Values(conn.Do("HGETALL", sym))
		if err != nil {
			log.Fatal(err)
		}

		//打印出myuser下所有Key值
		var value MyUser
		if err := redis.ScanStruct(values, &value); err != nil {
			log.Fatal(err)
		}
		log.Printf("%s: %+v", sym, &value)
	}
}

func SetInfo() {
	flag.Parse()
	pool = RedigoPool(*host, *password)

	conn := pool.Get()//从连接池获取连接
	defer conn.Close()//用完后放回连接池

	//redis操作
	value, err := conn.Do("SET", "name", "aaa")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)//返回ok

	value, err = redis.String(conn.Do("GET", "name"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)//返回name的值

}

func AddSession(sessionId,name,mobile string){
	flag.Parse()
	pool = RedigoPool(*host, *password)

	conn := pool.Get() //从连接池获取连接
	defer conn.Close() //用完后放回连接池

	myuser := map[string]*MyUser{
		sessionId: &MyUser{UserName: name, UserPhone: mobile},
	}

	//保存Map
	for sym, row := range myuser {
		if _, err := conn.Do("HMSET", redis.Args{sym}.AddFlat(row)...); err != nil {
			log.Fatal(err)
		}
	}

	value, err := conn.Do("EXPIRE", sessionId,1800)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)//返回ok
}
