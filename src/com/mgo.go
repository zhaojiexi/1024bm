package com

import (
	"log"
	"time"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

const (
	//MongodbHosts = "localhost:27017;106.75.13.29:27017"
	AuthDatabase = "1024bm"

	MongodbHosts = "106.75.13.30:27017"
	//AuthUserName = "guest"
	//AuthPassword = "welcome"
	//TestDatabase = "goinggo"
)

//获取session
func getSession() *mgo.Session {
	if session == nil {
		var err error
		//如果不存在，就新创建session
		mgoDialInfo := &mgo.DialInfo {
			Addrs: 		[]string{MongodbHosts},
			Direct: 	false,
			Timeout: 	60 * time.Second,
			Database: 	AuthDatabase,
			PoolLimit: 	4096,
			//Username:"cm",
			//Password:"123456"
		}
		/*
		DialWithInfo方法链接服务器或者服务器群.不同的是DialWithInfo方法可以提供额外的值给服务器.
		DialWithInfo和服务器(群)建立一个新的session. DialWithInfo方法也可以自定义值,当链接服务器的时候.
		当使用Dial方法建立链接,默认的超时时间为10秒,使用DialWithInfo可以自己设置超时时间.
		*/
		session, err = mgo.DialWithInfo(mgoDialInfo)
		if err != nil {
			log.Fatalf("getSession-mgo.DislwithInfo()时报错: %s\n",err)
		}

		session.SetMode(mgo.Monotonic,true)
	}

	return session.Clone()
}

//有了session就可以获取collection, f 函数类型的参数
func GetCollection(c string, f func(*mgo.Collection)error) error {
	session := getSession()
	defer session.Close()
	//defer func(){
		//session.Close()
		//if err := recover(); err != nil{
		//	log.Fatalf("GetCollection-session.close时报错: %s\n", err)
		//	panic(err)
		//}
	//}()

	//获取表对象
	collection := session.DB(AuthDatabase).C(c)

	return f(collection)
}