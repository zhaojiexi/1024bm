package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"models/user"
	"encoding/json"
	"log"
	"time"
	"fmt"
	"strconv"
)
//用户路由
func SetUserRouter(router *gin.Engine) *gin.Engine {
	userRoutert := router.Group("api/v1/")
	userRoutert.Use()
	{
		userRoutert.POST("user/register",UserRegister) //用户注册，如：http://0.0.0.0:8888/user/register 提交服务端参数在工具中创建
		userRoutert.GET("user/uid/:id",GetUserInfo)   //根据uid获取用户信息，如：http://0.0.0.0:8888/user/uid/5a167c7265b39931c4c57861
		userRoutert.POST("user/login",UserLogin)                 //用户登录,如：http://0.0.0.0:8000/user/login?name=caimin&password=123qwe
		userRoutert.GET("user/list",GetUsers)                   //获取用户列表,如：http://0.0.0.0:8000/api/v1/user/list
		userRoutert.POST("user/userinfo",UserInfo)                 //修改用户信息,以表单的形式接受 如：http://0.0.0.0:8000/user/userinfo 提交服务端参数在工具中创建
		userRoutert.POST("user/updatepwd",UpdateUserPassWord)	  //修改用户密码,以表单的形式接受 如：http://0.0.0.0:8000/user/updatepwd 提交服务端参数在工具中创建
		}
	return router
}

//用户注册
func UserRegister(c *gin.Context) {
	_name := c.PostForm("name")
	_phone := c.PostForm("phone")
	_password := c.PostForm("password")
	var result gin.H


	//根据输入值 判断用户信息是否存在
	ur,err:=user.UserRegister(_name,_phone,_password)

	if ur==nil{
		result=gin.H{"code":400,"msg":1,"start":0,"text":err}
	}else if ur!=nil{

		result=gin.H{"code":200,"msg":1,"start":1,"text":"注册成功"}

	}

	c.JSON(http.StatusOK,result)



	//c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("注册完成 %s\n", _name, " ", _phone, " ", _password)))

}

//获取用户信息
func GetUserInfo(c *gin.Context){

	var result gin.H

	uid:= c.Param("id")

	ur,r:=user.GetUserInfo(uid)

	//如果为nil 返回错误信息
	if ur==nil {
		result=gin.H{"code":400,"msg":1,"start":0,"text":r}
	}else {
		result=gin.H{"code":400,"msg":1,"start":0,"text":"success","UserInfo":ur}
	}

	c.JSON(http.StatusOK,result)


	return
}

//用户登录

func UserLogin(c *gin.Context){

	phone := c.Query("phone")
	password := c.Query("password")

	u,err:=user.UserLogin(phone,password)

	//用户不存在 或 账号密码错误 返回界面提示用户信息
	if u==nil {
		//user:=u[0]
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"text":err})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"text":"登录成功"})
	}





}

//获取用户列表
func GetUsers(c *gin.Context){


	ulist,_:=user.GetUsers()


	c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"text":" ","user":ulist})


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


	fmt.Println("gender",g)
	u.Uid=c.PostForm("Uid")	//用户id
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
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"text":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"text":"修改成功"})
	}


}

//修改用户密码 表单形式 根据uid修改密码
func UpdateUserPassWord(c *gin.Context){

	var u user.User
	u.Uid=c.PostForm("Uid")
	u.PassWord=c.PostForm("PassWord")

	result:=user.UpdateUserPassWord(&u)

	if result!="" {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"text":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"text":"修改成功"})
	}



}





