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

//redis缓存整个用户信息
func AddSession(user *User){

	flag.Parse()
	pool := com.RedigoPool(*com.Host, *com.Password)

	conn := pool.Get() //从连接池获取连接
	defer conn.Close() //用完后放回连接池


	myuser := map[string]*User{


		user.Uid.Hex():&User{user.ID,user.Uid,user.Name,user.Slug,user.Phone,
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
		fmt.Println(redis.Args{sym}.AddFlat(row))
	}
	//20分钟缓存时间	//根据上面存入的 sym：uid 设置缓存时间
	value, err := conn.Do("EXPIRE", user.Uid,1800)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)//返回ok
}

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
	user.Uid          = bson.NewObjectId()
	user.Name    	  = name
	user.Phone    	  = phone
	user.PassWord     = password
	user.RegisterDate = time.Now()
	user.Slug		  =user.Uid.Hex()+name
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
		return c.Find(bson.M{"Phone":phone,"PassWord":password,"IsEnabled":1}).All(&users)

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

<<<<<<< HEAD
//redis缓存整个用户信息
func AddSession(user *User){

	flag.Parse()
	pool := com.RedigoPool(*com.Host, *com.Password)

	conn := pool.Get() //从连接池获取连接
	defer conn.Close() //用完后放回连接池


	myuser := map[string]*User{


		user.Uid.Hex():&User{user._ID,user.Uid,user.Name,user.Slug,user.Phone,
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
		fmt.Println(redis.Args{sym}.AddFlat(row))
	}
	//20分钟缓存时间	//根据上面存入的 sym：uid 设置缓存时间
	value, err := conn.Do("EXPIRE", user.Uid.Hex(),1800)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)//返回ok
}
=======
>>>>>>> dev

//根据uid获取用户信息
func GetUserInfo(uid string)(user *User,result string){

	var users []User
	ubjectid:=bson.ObjectIdHex(uid)
	//
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Uid":ubjectid,"IsEnabled":1}).All(&users)
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
		return c.Find(bson.M{"IsEnabled":1}).All(&users)
	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("getUsers: %s\n", err)
	}


	return users,err


}

//修改用户信息
func UpdateUserInfo(user *User)(string){


	//根据uid去查找用户，然后修改用户信息

	var ulist []User

	query := func(c *mgo.Collection) (error) {

		return c.Find(bson.M{"Uid":user.Uid,"IsEnabled":"1"}).All(&ulist)

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
	if user.Gender!=3 {
		ulist[0].Gender=user.Gender//性别
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
		return c.Update(bson.M{"Uid":ulist[0].Uid,"IsEnabled":1},bson.M{"$set":bson.M{
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

//修改密码

func UpdateUserPassWord(user *User)string{

	var ulist []User

	//先判断用户是否存在
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Uid":user.Uid}).All(&ulist)
	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("getUsers: %s\n", err)
	}
	if len(ulist)<1 {
		return "用户不存在"
		
	}


	query = func(c *mgo.Collection) (error) {
		return c.Update(bson.M{"Uid":user.Uid,"IsEnabled":1},bson.M{"$set":bson.M{
			"PassWord":user.PassWord,
		}})
	}

	err = com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("getUsers: %s\n", err)
	}


	return ""
}

//根据关注人的id Following_UID 获取粉丝
func GetFans(uid string)([]User,string){

	//判断传入的 Following_UID 是否存在
	ubjectid:=bson.ObjectIdHex(uid)

	var follows []Follow

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Following_UID":ubjectid,"IsEnabled":1}).All(&follows)
	}

	err := com.GetCollection("Follow",query)
	if err != nil{
		log.Fatalf("getUsers: %s\n", err)
	}
	if len(follows)<1 {
		return nil,"用户不存在"
	}



	// 根据关注表 获取用户信息 再用一个总的用户数组保存 并返回所有粉丝
	var users []User
	var usersAll []User

	fmt.Println("len",len(follows))
	for   i:=0;i<len(follows) ;i++  {
		fmt.Println("uid",follows[i].User_UID)

		query = func(c *mgo.Collection) (error) {
			return c.Find(bson.M{"Uid":follows[i].User_UID,"IsEnabled":1}).All(&users)
		}
		err= com.GetCollection("User",query)
		if err != nil{
			log.Fatalf("GetFollows: %s\n", err)
		}
		if len(users)>0 {
			usersAll=append(usersAll,users[0] )
		}



	}
	fmt.Println(usersAll)


	return usersAll,""
}

//新增关注
func AddFollow (fo *Follow)string{

	var follows []Follow
	var users []User
	//校验关注人是否存在
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Uid":fo.Following_UID}).All(&users)
	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-AddFollow: %s\n", err)
	}
	if len(users)<1{
		return "关注用户id不存在"
	}

	//校验是否关注
	query = func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":fo.User_UID,"Following_UID":fo.Following_UID}).All(&follows)
	}

	err = com.GetCollection("Follow",query)
	if err != nil{
		log.Fatalf("User-AddFollow: %s\n", err)
	}
	if len(follows)>0{
		return "不能重复关注"
	}



	query = func(c *mgo.Collection) (error) {
		return c.Insert(&fo)
	}

	err = com.GetCollection("Follow",query)
	if err != nil{
		log.Fatalf("User-AddFollow: %s\n", err)
	}
	return ""

}

//根据uid获取所有关注的用户信息
func GetFollows(uid string)([]User,string){

	var follows []Follow
	var users []User
	bjectid:=bson.ObjectIdHex(uid)
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Uid":bjectid,"IsEnabled":1}).All(&users)
	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("Follows: %s\n", err)
	}
	if len(users)<1 {
		return nil,"用户不存在"
	}

	query = func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":bjectid,"IsEnabled":1}).All(&follows)
	}

	err = com.GetCollection("Follow",query)
	if err != nil{
		log.Fatalf("Follows: %s\n", err)
	}

	var user2 []User

	//查询关注的用户详细信息
	for  i:=0;i<len(follows) ;i++  {

		query = func(c *mgo.Collection) (error) {
			return c.Find(bson.M{"Uid":follows[i].Following_UID}).All(&users)
		}

		err = com.GetCollection("User",query)
		if err != nil{
			log.Fatalf("Follows: %s\n", err)
		}
		if len(users)>0 {
			user2=append(user2, users[0])
		}

	}


	return user2,""


}
//逻辑删除关注 把IsEnabled设置0
func DelFollow(fo Follow)string{
	var f []Follow

	//校验关注人是否存在
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":fo.User_UID,"Following_UID":fo.Following_UID,"IsEnabled":1}).All(&f)
	}
	err := com.GetCollection("Follow",query)

	if err != nil{
		log.Fatalf("User-DelFollow: %s\n", err)
	}
	if len(f)<1{
		return "用户不存在或状态为不可用"
	}

	query = func(c *mgo.Collection) (error) {
		return c.Update(bson.M{"_id":f[0].ID},bson.M{"$set":bson.M{
			"IsEnabled":0,
		}})
	}

	err = com.GetCollection("Follow",query)
	if err != nil{
		log.Fatalf("DelFollow: %s\n", err)
	}

	return ""


}

//新增收藏
func AddFavorite(fr *Favorite)string{

	var frlist []Favorite

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":fr.User_UID,"Article_ID":fr.Article_ID}).All(&frlist)
	}

	err := com.GetCollection("Favorite",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}
	if len(frlist)>0 {
		return "已关注此文章，请勿重复关注"
	}



	query = func(c *mgo.Collection) (error) {
		return c.Insert(fr)
	}

	err = com.GetCollection("Favorite",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}
	return ""
}

//查询所有收藏的文章
func GetFavoriteByID(uid string)([]Favorite){

	var frlist []Favorite

	ubjiectid:=bson.ObjectIdHex(uid)

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":ubjiectid,"IsEnabled":1}).All(&frlist)
	}

	err := com.GetCollection("Favorite",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}

	return frlist

}

//逻辑删除收藏  把IsEnabled设置0
func DelFavorite(fo *Favorite)(string){

	var folist []Favorite

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":fo.User_UID,"Article_ID":fo.Article_ID,"IsEnabled":1}).All(&folist)
	}

	err := com.GetCollection("Favorite",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}
	if len(folist)<1 {
		return "找不到对应的文章"
	}


	query = func(c *mgo.Collection) (error) {
		return c.Update(bson.M{"User_UID":fo.User_UID,"Article_ID":fo.Article_ID},bson.M{"$set":bson.M{
			"IsEnabled":0,
		}})
	}

	err = com.GetCollection("Favorite",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}


	return ""

}
//新增浏览记录
func AddBrowseHistory(bh *BrowseHistory)string{

	var bhlist []BrowseHistory

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":bh.User_UID,"Article_ID":bh.Article_ID,"IsEnabled":1}).All(&bhlist)
	}

	err := com.GetCollection("BrowseHistory",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}
	if len(bhlist)>1 {
		return "此浏览记录已存在"
	}


	query = func(c *mgo.Collection) (error) {
		return c.Insert(&bh)
	}

	err = com.GetCollection("BrowseHistory",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}
	return ""

}
//查看浏览记录
func GetBrowseHistory(uid string)[]BrowseHistory{


	var bhlist []BrowseHistory
	ubjectid:=bson.ObjectIdHex(uid)

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":ubjectid,"IsEnabled":1}).All(&bhlist)
	}

	err := com.GetCollection("BrowseHistory",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}

	return bhlist

}

//删除记录
func DelBrowseHistory(uid,article_ID string)string{

	var bhlist []BrowseHistory
	ubjectid:=bson.ObjectIdHex(uid)
	abject:=bson.ObjectIdHex(article_ID)

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"User_UID":ubjectid,"Article_ID":abject,"IsEnabled":1}).All(&bhlist)
	}

	err := com.GetCollection("BrowseHistory",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}
	if len(bhlist)<1{
		return "该用户没有该浏览记录"
	}
	query = func(c *mgo.Collection) (error) {
		return c.Update(bson.M{"_id":bhlist[0].ID},bson.M{"$set":bson.M{
			"IsEnabled":0,
		}})

	}
	err = com.GetCollection("BrowseHistory",query)
	if err != nil{
		log.Fatalf("addFavorite: %s\n", err)
	}

	return ""

}