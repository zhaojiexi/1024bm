package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"router"
	"time"
)

func main() {
	fmt.Println("hello world0000")

	_router := InitRouters()
	_server := &http.Server{
		Addr:           ":8888",
		Handler:        _router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("==========  Web API 接受调用[ok]...  ==========")
	_server.ListenAndServe()
	//fmt.Println(models.UserRegister("ss","ss","ss"))
}

//路由初始化
func InitRouters() *gin.Engine {
	_initrouter := gin.Default()
	_initrouter = router.SetUserRouter(_initrouter)
	//router.GET("/user/:name", handle.GetUserName) //http://0.0.0.0:8000/user/caimin
	//router.GET("/register", api.Userregister) //http://0.0.0.0:8000/reg?name=caimin
	//router.GET("/login", api.UserLogin) //http://0.0.0.0:8000/user?name=caimin
	return _initrouter
}
