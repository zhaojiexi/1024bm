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
func UserRegister(name,phone,password string) (user1 *User,result string) {

	//验证用户输入
	if len(name)<1{
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

		return c.Find(bson.M{"Name":name}).All(&ulist)

	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserRegister时报错: %s\n", err)
	}
	//如果存在 直接返回
	if len(ulist)>0 {
		return nil,"该用户名已存在"
	}


	//验证手机是否存在
	query = func(c *mgo.Collection) (error) {

		return c.Find(bson.M{"Phone":phone}).All(&ulist)

	}

	err = com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserRegister时报错: %s\n", err)
	}
	//如果存在 直接返回
	if len(ulist)>0 {
		return nil,"该手机已被注册"
	}



	var user *User = new(User)
	user.Uid          = bson.NewObjectId().Hex()
	user.Name    	  = name
	user.Phone    	  = phone
	user.PassWord     = password
	user.RegisterDate = time.Now()
	user.Slug		  =user.Uid+name
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
	//注册成功 放入缓存
	AddSession(user)

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
		return c.Find(bson.M{"Phone":phone,"PassWord":password}).All(&users)

	}

	err = com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserLogin2: %s\n", err)
	}


	if len(users)<1{
		return nil,"账号或密码错误"
	}else{
		//注册成功 放入缓存
		AddSession(users[0])
		return users[0],""
	}

}

//redis缓存整个用户信息
func AddSession(user *User){

	flag.Parse()
	pool := com.RedigoPool(*com.Host, *com.Password)

	conn := pool.Get() //从连接池获取连接
	defer conn.Close() //用完后放回连接池


	myuser := map[string]*User{


		user.Uid:&User{user._ID,user.Uid,user.Name,user.Slug,user.Phone,
		user.PassWord,user.RegisterDate,user.Location,user.University,
		user.Company,user.WebSite,user.Follower_count,user.Following_count,
		user.Browse_count,user.Article_count,user.Describe,user.Profile_image_url,
		user.LastLogin,user.Interest,user.IsEnabled,user.Gender},
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

//修改用户信息
func UserInfo(user *User)(string){


	//根据uid去查找用户，然后修改用户信息

	var ulist []User

	query := func(c *mgo.Collection) (error) {

		return c.Find(bson.M{"Uid":user.Uid}).All(&ulist)

	}
	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserInfo: %s\n", err)
	}
	//
	if len(ulist)<1 {
		return "用户不存在"
	}

// 判断用户输入值是否为空 否 则修改数据
	if user.Gender!=""{
		ulist[0].Gender=user.Gender//性别
	}
	if user.Describe!=""{
		ulist[0].Describe=user.Describe//个人介绍
	}
	if user.Location!=""{
		ulist[0].Location=user.Location//所在地
	}
	if user.Company!=""{
		ulist[0].Company=user.Company//公司
	}
	if user.University!=""{
		ulist[0].University=user.University//学习
	}
	if user.Interest!=nil{
		ulist[0].Interest=user.Interest//兴趣
	}
	if user.WebSite!=""{
		ulist[0].WebSite=user.WebSite//展示网站
	}
	if user.Profile_image_url!=""{
		ulist[0].Profile_image_url=user.Profile_image_url//展示网站
	}
	ulist[0].LastLogin=user.LastLogin

/*	ulist[0].Gender=user.Gender//性别
	ulist[0].Location=user.Location	//所在地
	ulist[0].Company=user.Company	//公司
	ulist[0].University=user.University//学习
	ulist[0].Interest=user.Interest//兴趣
	ulist[0].WebSite=user.WebSite	//展示网站
	ulist[0].Profile_image_url=user.Profile_image_url	//头像地址
	*/


	query = func(c *mgo.Collection) (error) {
		return c.Update(bson.M{"Uid":ulist[0].Uid},bson.M{"$set":bson.M{
			"Gender":ulist[0].Gender,
			"Describe":ulist[0].Describe,
			"Location":ulist[0].Location,
			"Company":ulist[0].Company,
			"University":ulist[0].University,
			"Interest":ulist[0].Interest,
			"WebSite":ulist[0].WebSite,
			"Profile_image_url":ulist[0].Profile_image_url,
			"LastLogin":ulist[0].LastLogin,

		}})
	}

	err = com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserInfo时报错: %s\n", err)
	}

	AddSession(&ulist[0])//缓存用户信息

	return ""
}
