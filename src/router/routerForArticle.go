package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"models/article"
	"encoding/json"
	"strconv"
	"log"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func SetArtiCleRouter(router *gin.Engine) *gin.Engine {
	articleRoutert := router.Group("api/v1/article")
	articleRoutert.Use()
	{
		articleRoutert.POST("/addcategory",AddCategory)
		articleRoutert.POST("/addarticle",AddArticle)



	}
	return router
}
//新增分类
func AddCategory(c *gin.Context){
	var err error
	var a article.Category
	a.CategoryName=c.PostForm("CategoryName")
	pid:=c.PostForm("ParentCategory_ID")
	if pid!="" {
		a.ParentCategory_ID,err=strconv.ParseInt(pid,10,64)
	}

	checkERR(err)

	a.ID=bson.NewObjectId()

	a.Define=c.PostForm("Define")
	ts:=c.PostForm("Tags")	//{"Tags":[{"TagID":"id1","Tag":"111"}, {"TagID":"id2","Tag":"222"}]}

	json.Unmarshal([]byte(ts),&a)

	result:=article.AddCategory(&a)

	if result!=""{
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success"})
	}

}

//新增文章
func AddArticle(c *gin.Context){

	var a article.Article
	var err error

	a.Attribute=c.PostForm("Attribute")//属性（入门教程、安装搭建、配置环境、解决问题、实战案例、疑惑吐槽、推荐介绍、心得分享、技术详解）
	a.Title=c.PostForm("Title")//
	a.Abstract=c.PostForm("Abstract")//摘要
	a.Author_uid=bson.ObjectIdHex(c.PostForm("Author_uid"))//作者id
	a.Author=c.PostForm("Author")//

	//
	t:=c.PostForm("Created")
	//时间类型转换
	if t!=""{
		a.Created,err=time.Parse("2006-01-02 15:04:05",t)
	}else{
		a.Created=time.Now()
	}
	checkERR(err)
	orn:=c.PostForm("Original")
	if orn!="" {
		a.Original,err=strconv.ParseInt(orn,10,64)
	}
	checkERR(err)
	a.Source_url=c.PostForm("Source_url")//
	a.Source_domain=c.PostForm("Source_domain")//
	a.Source_chinese=c.PostForm("Source_chinese")//
	a.CategoryParent_ID=c.PostForm("CategoryParent_ID")//父分类id

	a.CategoryParent=c.PostForm("CategoryParent")//父分类名称

	a.CategorySub_ID=c.PostForm("CategorySub_ID")//子分类id
	a.CategorySub=c.PostForm("CategorySub")//子分类名

	Tags:=c.PostForm("Tags")//子分类名	//{"Tags":["java","sprint","git"]}
	//子分类名转换格式
	json.Unmarshal([]byte(Tags),&a)

	a.Text=c.PostForm("Text")//内容
	lc:=c.PostForm("Like_count")
	if lc!="" {
		a.Like_count,err=strconv.ParseInt(lc,10,64)
	}
	checkERR(err)
	dc:=c.PostForm("Dislike_count")
	if lc!="" {
		a.Dislike_count,err=strconv.ParseInt(dc,10,64)
	}
	checkERR(err)
	cc:=c.PostForm("Comment_count")
	if lc!="" {
		a.Comment_count,err=strconv.ParseInt(cc,10,64)
	}
	checkERR(err)
	cdc:=c.PostForm("Collected_count")
	if lc!="" {
		a.Collected_count,err=strconv.ParseInt(cdc,10,64)
	}
	checkERR(err)
	rc:=c.PostForm("Read_count")
	if lc!="" {
		a.Read_count,err=strconv.ParseInt(rc,10,64)
	}
	checkERR(err)



	c.JSON(http.StatusOK,article.AddArticle(&a))





}

func checkERR(err error ){

	if err!=nil {
		log.Fatal(err)
	}
}
