package user

import (
	"com"
	"log"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
	"github.com/garyburd/redigo/redis"
	"flag"
	"fmt"
)

//用户注册
func UserRegister(nickName,phone,password string) (user1 *User,result string) {

	//验证用户输入
	if len(nickName)<1{
		return nil,"用户名长度小于1"
	}
	if len(phone)<1{
		return nil,"手机长度不能小于1"
	}
	if len(password)<1{
		return nil,"密码用长度不能小于1"
	}

	var ulist []User

	//验证用户是否存在
	query := func(c *mgo.Collection) (error) {

		return c.Find(bson.M{"Phone":phone}).All(&ulist)

	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserRegister时报错: %s\n", err)
	}
	//如果存在 直接返回
	if len(ulist)>1 {
		return nil,"该手机已被注册"
	}


	var user *User = new(User)
	user.Uid          = bson.NewObjectId().Hex()
	user.Name    	  = nickName
	user.Phone    	  = phone
	user.PassWord     = password
	user.RegisterDate = time.Now()
	user.Slug		  =user.Uid+nickName
	user.Location		=""
	user.University		=""
	user.Company		=""
	user.WebSite		=""
	user.Follower_count	=0
	user.Following_count=0
	user.Browse_count=0
	user.Article_count=0
	user.Describe =""
	user.Profile_image_url=""
	user.LastLogin	=time.Now()
	user.Interest  =nil
	user.IsEnabled	=1	//默认可用

	query = func(c *mgo.Collection) (error) {
		return c.Insert(user)
	}

	err = com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserRegister时报错: %s\n", err)
	}
	return user,""
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


		user.Uid:&User{Uid:user.Uid,Name:user.Name,Phone:user.Phone,PassWord:user.PassWord,RegisterDate:user.RegisterDate},
	}



	//保存Map
	for sym, row := range myuser {
		if _, err := conn.Do("HMSET", redis.Args{sym}.AddFlat(row)...); err != nil {
			log.Fatal(err)
		}
	}
	//20分钟缓存时间
	value, err := conn.Do("EXPIRE", user.Uid,1800)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)//返回ok
}

//根据uid获取用户信息
func GetUserInfo(uid string)(user *User,result string){

	var users []User

	//
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Uid":uid}).All(&users)
	}
	
	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-GetUserInfo: %s\n", err)
	}
	if len(users)<1 {
		return nil,"该用户不存在"
	}

	return &users[0],""


}

//查找所有用户信息

func GetUsers()([]User,error){

	var users []User

	//
	query := func(c *mgo.Collection) (error) {
		return c.Find(nil).All(&users)
	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("getUsers: %s\n", err)
	}


	return users,err


}
