package router

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	//"github.com/golang/net/html/atom"
	//"gopkg.in/mgo.v2/bson"
	"models/user"


)

//用户路由
func SetUserRouter(router *gin.Engine) *gin.Engine {
	userRoutert := router.Group("api/v1/")
	userRoutert.Use()
	{
		userRoutert.POST("user/register",UserRegister) //用户注册，如：http://0.0.0.0:8888/user/register 提交服务端参数在工具中创建
		userRoutert.GET("user/uid/:id",GetUserInfo)   //根据uid获取用户信息，如：http://0.0.0.0:8888/user/uid/5a167c7265b39931c4c57861
		userRoutert.POST("user/login",UserLogin)                 //用户登录,如：http://0.0.0.0:8000/user/login?name=caimin&password=123qwe
		//userRoutert.GET("list",GetUsers)                   //获取用户列表,如：
	}
	return router
}

//用户注册
func UserRegister(c *gin.Context) {
	_name := c.PostForm("name")
	_phone := c.PostForm("phone")
	_password := c.PostForm("password")

	//c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("注册完成 %s\n", _name, " ", _phone, " ", _password)))
	c.JSON(http.StatusOK, user.UserRegister(_name,_phone,_password))
}

//获取用户信息
func GetUserInfo(c *gin.Context){
	value, exist := c.GetQuery("id")
	if !exist {
		value = "CaiMin"
	}

	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("ok! %s\n", value)))
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

	}
	//验证成功放入缓存

	user.AddSession(u)




	c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"text":"登录成功"})

}

func checkERR(err error){
	if err!=nil {
		log.Fatal(err)
	}
}
