package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"log"
	"time"
	"strconv"
	"gopkg.in/mgo.v2/bson"

	"models/user"

	"fmt"
)
//用户路由
func SetUserRouter(router *gin.Engine) *gin.Engine {
	userRoutert := router.Group("api/v1/")
	userRoutert.Use()
	{
		userRoutert.POST("user/register",UserRegister) //用户注册，如：http://0.0.0.0:8888/user/register 提交服务端参数在工具中创建
		userRoutert.GET("user/getuserinfo",GetUserInfo)   //根据uid获取用户信息，如：http://0.0.0.0:8888/user/getuserinfo?Uid=5a167c7265b39931c4c57861
		userRoutert.POST("user/login",UserLogin)                 //用户登录,如：http://0.0.0.0:8000/user/login?Phone=caimin&PassWord=123qwe
		userRoutert.GET("user/list",GetUsers)                   //获取用户列表,如：http://0.0.0.0:8000/api/v1/user/list
		userRoutert.PUT("user/userinfo",UserInfo)                 //修改用户信息,以表单的形式接受 如：http://0.0.0.0:8000/user/userinfo 提交服务端参数在工具中创建
		userRoutert.PUT("user/updatepwd",UpdateUserPassWord)	  //修改用户密码,以表单的形式接受 如：http://0.0.0.0:8000/user/updatepwd 提交服务端参数在工具中创建
		userRoutert.GET("user/fans",GetFans)		//获取所有粉丝  http://0.0.0.0:8000/user/fans?Uid=5a167c7265b39931c4c57861
		userRoutert.POST("user/addfollow",AddFollow)		//新增关注 	http://0.0.0.0:8000/user/addfollow?User_UID=5a2a35f2bfb1481f9cf54c7a&Following_UID=5a2a4b61bfb1481734be3ae1&User_name=test&Following_Name=ftest
		userRoutert.GET("user/follows",GetFollows)		//获取所有关注用户 http://0.0.0.0:8000/user/follows?Uid=5a2a35f2bfb1481f9cf54c7a
		userRoutert.DELETE("user/delfollow",DelFollow)		//http://127.0.0.1:8888/api/v1/user/delfollow?User_UID=5a2a35f2bfb1481f9cf54c7a&Following_UID=5a333679bfb1481ee4fe16a4
		userRoutert.POST("user/addfavorite",AddFavorite)		//新增收藏 	http://0.0.0.0:8000/user/addfollow?User_UID=5a2a35f2bfb1481f9cf54c7a&Following_UID=5a2a4b61bfb1481734be3ae1&User_name=test&Following_Name=ftest
		userRoutert.GET("user/favoritebyid",GetFavoriteByID)	//查看收藏文章 	http://0.0.0.0:8000/user/favoritebyid?Uid=5a2a35f2bfb1481f9cf54c7a
		userRoutert.DELETE("user/delfavorite",DelFavorite)	//// 删除 	http://127.0.0.1:8888/api/v1/user/delfavorite?User_UID=5a2a35f2bfb1481f9cf54c7a&Article_ID=5a2a35f2bfb1481f9cf54c7a
		userRoutert.POST("user/addbrowsehistory",AddBrowseHistory)	//http://127.0.0.1:8888/api/v1/user/addbrowsehistory?User_UID=5a3366aabfb1481940f4c672&Article_ID=5a3366aabfb1481940f4c672&Article_Title=test&Article_Author=zuozhe&Author_Picture=photo&Article_Time=2017-12-15 15:25:00&Created=2017-12-15 15:25:00
		userRoutert.GET("user/getbrowsehistory",GetBrowseHistory)	//查看浏览记录http://127.0.0.1:8888/api/v1/user/getbrowsehistory?User_UID=5a3366aabfb1481940f4c672
		userRoutert.DELETE("user/delbrowsehistory",DelBrowseHistory)

		}
	return router
}

//用户注册
func UserRegister(c *gin.Context) {
	_name := c.PostForm("Name")
	_phone := c.PostForm("Phone")
	_password := c.PostForm("PassWord")
	var result gin.H


	//根据输入值 判断用户信息是否存在
	ur,err:=user.UserRegister(_name,_phone,_password)

	if ur==nil{
		result=gin.H{"code":400,"msg":1,"start":0,"result":err}
	}else if ur!=nil{

		result=gin.H{"code":200,"msg":1,"start":1,"result":"注册成功"}

	}

	c.JSON(http.StatusOK,result)



	//c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("注册完成 %s\n", _name, " ", _phone, " ", _password)))

}

//uid获取用户信息
func GetUserInfo(c *gin.Context){

	var result gin.H

	uid:= c.Query("Uid")

	ur,r:=user.GetUserInfo(uid)

	//如果为nil 返回错误信息
	if ur==nil {
		result=gin.H{"code":400,"msg":1,"start":0,"result":r}
	}else {
		result=gin.H{"code":200,"msg":1,"start":1,"result":"success","context":ur}
	}

	c.JSON(http.StatusOK,result)


	return
}

//用户登录

func UserLogin(c *gin.Context){

	phone := c.PostForm("Phone")
	password := c.PostForm("PassWord")

	u,err:=user.UserLogin(phone,password)

	//用户不存在 或 账号密码错误 返回界面提示用户信息
	if u==nil {
		//user:=u[0]
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"retult":err})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"登录成功","User":u})
	}





}

//获取用户列表
func GetUsers(c *gin.Context){

	num:=c.Query("PageNum")
	count:=c.Query("PageCount")

	//如果没传 则赋默认值
	var pn,pc int
	var err error

	if num!="" {
		pn,err=strconv.Atoi(num)
	}else {
		pn=1
	}
	if count!="" {
		pc,err=strconv.Atoi(count)
	}else{
		pc=10
	}


	if err!=nil {
		log.Fatal(err)
	}
	fmt.Printf("pc,%s,pc,%s\n",pn,pc)
	ulist,PageNum,PageCount,PageSum,PageMax,_:=user.GetUsers(pn,pc)

	fmt.Println("页数",PageNum)
	fmt.Println("每页显示几行",PageCount)
	fmt.Println("总共几条数据",PageSum)
	fmt.Println("最大页数",PageMax)

	c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success","pageNum":PageNum,"pageCount":PageCount,"pageSum":PageSum,"pageMax":PageMax,"context":ulist})


}

//修改用户信息
func UserInfo (c *gin.Context){

	var u user.User
	var err error
/*
	Interest: 传输格式
{"Interest":[{"SubCategoryID":"id1","SubCategory":"111"},
{"SubCategoryID":"id2","SubCategory":"222"}]}

*/
	Interest:=c.PostForm("Interest")
	json.Unmarshal([]byte(Interest),&u)

	g:=c.PostForm("Gender")

	if g!=""{
		u.Gender,_=strconv.ParseInt(g,10,64)
	}else {
		u.Gender=3	//标记一下 表示字段取查询出来的数据
	}


	u.Uid=bson.ObjectIdHex(c.PostForm("Uid"))	//用户id

	u.Describe=c.PostForm("Describe")	//个人介绍
	u.Location=c.PostForm("Location")	//所在地
	u.Company=c.PostForm("Company")		//公司
	u.University=c.PostForm("University")		//学习
	u.WebSite=c.PostForm("WebSite")		//展示网站
	u.Profile_image_url=c.PostForm("Profile_image_url")		//头像地址

	formTime:=c.PostForm("LastLogin")		//time 格式：2014-06-15 10:10:10

	//获取时间 如果时间不为空 就转化为time类型 否则就取当前时间
	if formTime!=""{
		u.LastLogin,err= time.Parse("2006-01-02 15:04:05", formTime)
		if err!=nil {
			log.Fatal(err)
		}
	}else{
		u.LastLogin=time.Now()
	}


	result:=user.UpdateUserInfo(&u)
	//如果返回值不为“” 则错误 返回错误信息
	if result!="" {
		//user:=u[0]
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success"})
	}


}

//修改用户密码 表单形式 根据uid修改密码
func UpdateUserPassWord(c *gin.Context){

	var u user.User
	u.Uid=bson.ObjectIdHex(c.PostForm("Uid"))
	u.PassWord=c.PostForm("PassWord")

	result:=user.UpdateUserPassWord(&u)

	if result!="" {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success"})
	}



}

//uid查询所有粉丝
func GetFans(c *gin.Context){
	uid:=c.Query("Uid")

	u,result:=user.GetFans(uid)

	//查询失败 返回错误信息 成功 返回成功信息和粉丝详细信息
	if result!="" {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success","context":u})
	}



}

//新增关注
func AddFollow(c *gin.Context){

	var fo user.Follow

	fo.User_UID=bson.ObjectIdHex(c.PostForm("User_UID"))
	fo.Following_UID=bson.ObjectIdHex(c.PostForm("Following_UID"))
	fo.User_name=c.PostForm("User_name")
	fo.Following_Name=c.PostForm("Following_Name")
	fo.IsEnabled=1



	formTime:=c.PostForm("Created")

	//如果时间不为空 转化为time格式
	if	formTime!=""{
		var err error
		fo.Created,err=time.Parse("2006-01-02 15:04:05", formTime)
		fmt.Println(fo.Created)
		if err!=nil {
			log.Fatal(err)
		}
	}else{
		fo.Created=time.Now()
	}



	result:=user.AddFollow(&fo)

	if result!="" {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"text":result})

	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"text":"success"})
	}



}

//查看所有关注 根据用户id查看自己的所有关注
func GetFollows(c *gin.Context){
	var err error
	var pagenum,pagecount int



	uid:=c.Query("Uid")
	pnum:=c.Query("PageNum")
	pcount:=c.Query("PageCount")


	if pnum!="" {
		pagenum,err=strconv.Atoi(pnum)
	}else{
		pagenum=1
	}

	if pcount!="" {
		pagecount,err=strconv.Atoi(pcount)
	}else{
		pagecount=10
	}

	if	err!=nil{
		log.Fatal(err)
	}


	ulist,result,PageNum,PageCount,PageSum,PageMax:=user.GetFollows(uid,pagenum,pagecount)

	if result!="" {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success","pageNum":PageNum,"pageCount":PageCount,"pageSum":PageSum,"pageMax":PageMax,"context":ulist})
	}



}

//删除关注 传用户id和关注人id
func DelFollow(c *gin.Context){

	var fo user.Follow

	fo.User_UID=bson.ObjectIdHex(c.Query("User_UID"))
	fo.Following_UID=bson.ObjectIdHex(c.Query("Following_UID"))

	result:=user.DelFollow(fo)

	if result!="" {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success"})
	}




}

//新增收藏
func AddFavorite(c *gin.Context){

	var fr user.Favorite

	fr.User_UID=bson.ObjectIdHex(c.PostForm("User_UID"))
	fr.Article_ID=bson.ObjectIdHex(c.PostForm("Article_ID"))
	fr.Article_Title=c.PostForm("Article_Title")
	fr.Article_Author=c.PostForm("Article_Author")
	fr.Author_Picture=c.PostForm("Author_Picture")

	fr.IsEnabled=1

	formTime:=c.PostForm("Article_Time")
	//如果时间不为空 转化为time格式
	if	formTime!=""{
		var err error
		fr.Article_Time,err=time.Parse("2006-01-02 15:04:05", formTime)
		if err!=nil {
			log.Fatal(err)
		}
	}else{
		fr.Article_Time=time.Now()
	}

	formTime1:=c.PostForm("Created")
	//如果时间不为空 转化为time格式
	if	formTime!=""{
		var err error
		fr.Created,err=time.Parse("2006-01-02 15:04:05", formTime1)
		if err!=nil {
			log.Fatal(err)
		}
	}else{
		fr.Created=time.Now()
	}


	result:=user.AddFavorite(&fr)

	if result!="" {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})

	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success"})
	}


}

//查看收藏的所有文章 根据收藏人id  !!(现在只查询收藏表 还没有关联查询文章详细信息)
func GetFavoriteByID(c *gin.Context){

	uid:=c.Query("Uid")

	favoritelist:=user.GetFavoriteByID(uid)

	c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success","context":favoritelist})


}

//取消收藏  传收藏人id 和文章id
func DelFavorite(c *gin.Context){

	var fo user.Favorite

	fo.User_UID=bson.ObjectIdHex(c.Query("User_UID"))
	fo.Article_ID=bson.ObjectIdHex(c.Query("Article_ID"))

	result:=user.DelFavorite(&fo)
	if result!="" {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
	c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success"})
	}


}

//新增浏览记录		!!(现在直接新增 还没有校验文章id是否存在)
func AddBrowseHistory(c *gin.Context){

	var err error
 	var bh user.BrowseHistory

	bh.User_UID=bson.ObjectIdHex(c.PostForm("User_UID"))
	bh.Article_ID=bson.ObjectIdHex(c.PostForm("Article_ID"))
	bh.Article_Title=c.PostForm("Article_Title")
	bh.Article_Author=c.PostForm("Article_Author")
	bh.Author_Picture=c.PostForm("Author_Picture")
	article_Time:=c.PostForm("Article_Time")
	created:=c.PostForm("Created")
	bh.IsEnabled=1

	if article_Time!="" {
		bh.Article_Time,err=time.Parse("2006-01-02 15:04:05", article_Time)
		if err!=nil {
			log.Fatal(err)
		}
	}
	if created!="" {
		bh.Created,err=time.Parse("2006-01-02 15:04:05", created)
		if err!=nil {
			log.Fatal(err)
		}
	}

	result:=user.AddBrowseHistory(&bh)

	if	result!=""{
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success"})
	}



}

//根据用户id查看浏览记录
func GetBrowseHistory(c *gin.Context) {

	uid:=c.Query("User_UID")

	list:=user.GetBrowseHistory(uid)


	c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success","context":list})


}

//删除浏览记录 根据用户id 文章id
func DelBrowseHistory(c *gin.Context){

	uid:=c.Query("User_UID")
	article_ID:=c.Query("Article_ID")

	result:=user.DelBrowseHistory(uid,article_ID)

	if	result!=""{

		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":1,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":0,"result":"success"})

	}

}