package user

import (
	"com"
	"log"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
	"github.com/garyburd/redigo/redis"
	//"github.com/syndtr/goleveldb/leveldb/errors"
	"flag"
	"fmt"
)


func UserRegister(nickName,phone,password string) (uid string) {
	var user *	User = new(User)
	user.Uid          = bson.NewObjectId()
	user.Name     = nickName
	user.Phone    	  = phone
	user.PassWord     = password
	user.RegisterDate = time.Now()

	query := func(c *mgo.Collection) (error) {
		return c.Insert(user)
	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserRegister时报错: %s\n", err)
	}
	return user.Uid.Hex()
}

//用户登录
func UserLogin(phone,password string)(user *User,s string){
	var users []*User

	//根据手机号判断是否注册
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Phone":phone}).All(&users)
	}
	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserLogin1: %s\n", err)
	}

	if len(users)<1{
		return nil,"该手机尚未注册"
	}


    //校验密码准确性
	query = func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Phone":phone,"password":password}).All(&users)

	}

	err = com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserLogin2: %s\n", err)
	}


	if len(users)<1{
		return nil,"账号或密码错误"
	}else{
		return users[0],""
	}

}

//存在redis
func AddSession(user *User){

	flag.Parse()
	pool := com.RedigoPool(*com.Host, *com.Password)

	conn := pool.Get() //从连接池获取连接
	defer conn.Close() //用完后放回连接池

	myuser := map[string]*User{
		user.Uid.Hex():&User{Uid:user.Uid,Name:user.Name,Phone:user.Phone,PassWord:user.PassWord,RegisterDate:user.RegisterDate},
	}



	//保存Map
	for sym, row := range myuser {
		if _, err := conn.Do("HMSET", redis.Args{sym}.AddFlat(row)...); err != nil {
			log.Fatal(err)
		}
	}
	//20分钟缓存时间
	value, err := conn.Do("EXPIRE", user.Uid.Hex(),1800)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)//返回ok
}