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
	"sort"
	"fmt"
)

// 路由设置
func SetArtiCleRouter(router *gin.Engine) *gin.Engine {
	articleRoutert := router.Group("api/v1/article")
	articleRoutert.Use()
	{
		articleRoutert.POST("/addcategory",AddCategory)
		articleRoutert.POST("/addarticle",AddArticle)
		articleRoutert.POST("/addcomments",AddComments)
		articleRoutert.GET("/getarticlebyid",GetArticleByID)
		articleRoutert.GET("/getarticleall",GetArticleAll)


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


	pid:=c.PostForm("CategoryParent_ID")//父分类id
	if pid!="" {

		a.CategoryParent_ID,err=strconv.ParseInt(pid,10,64)

	}

	cid:=c.PostForm("CategorySub_ID")//子分类id

	if cid!="" {
		a.CategorySub_ID,err=strconv.ParseInt(cid,10,64)
	}


	a.CategoryParent=c.PostForm("CategoryParent")//父分类名称


	
	a.CategorySub=c.PostForm("CategorySub")//子分类名

	Tags:=c.PostForm("Tags")//子分类名	//{"Tags":["java","sprint","git"]}
	//子分类名转换格式
	json.Unmarshal([]byte(Tags),&a)

	//标签去重
	slice:=a.Tags
	sort.Strings(slice)
	i:= 0
	var j int
	for{
		if i >= len(slice)-1 {
			break
		}

		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {
		}
		slice= append(slice[:i+1], slice[j:]...)
		i++
	}
	a.Tags=slice
	fmt.Println("tag",a.Tags)
	a.Text=c.PostForm("Text")//内容

	result:=article.AddArticle(&a)

	c.JSON(http.StatusOK,result)



}
//新增评论
func AddComments(c *gin.Context){

	var ct article.Comments
	var err error

	ct.Article_id=bson.ObjectIdHex(c.PostForm("Article_id"))
	ct.User_uid=bson.ObjectIdHex(c.PostForm("User_uid"))
	ct.Text=c.PostForm("Text")






	//Replieds

	rs:=c.PostForm("Replieds")

	json.Unmarshal([]byte(rs),&ct)
	fmt.Println("Replieds:",ct.Replieds)

	ctime:=c.PostForm("Created")
	if ctime!="" {
		ct.Created,err=time.Parse("2006-01-02 15:04:05",ctime)
	}else{
		ct.Created=time.Now()
	}
	checkERR(err)

	article.AddComments(&ct)


}

//获取文章根据id
func GetArticleByID(c *gin.Context){

	aid:=c.Query("ArticleID")

	alist,result:=article.GetArticleByID(aid)

	if result!=""{
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":1,"start":0,"result":result})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success","context":alist})
	}
}
//查询所有文章(翻页)
func GetArticleAll(c *gin.Context){
	var pn,pc int
	var err error

	pnum:=c.Query("PageNum")
	pCount:=c.Query("PageCount")

	if pnum!="" {
		pn,err=strconv.Atoi(pnum)
	}else{
		pn=1
	}
	checkERR(err)
	if pCount!="" {
		pc,err=strconv.Atoi(pCount)
	}else{
		pc=10
	}
	checkERR(err)

	alist,PageNum,PageCount,PageSum,PageMax:=article.GetArticleAll(pn,pc)
	c.JSON(http.StatusOK,gin.H{"code":200,"msg":1,"start":1,"result":"success","pageNum":PageNum,"pageCount":PageCount,"pageSum":PageSum,"pageMax":PageMax,"context":alist})

}

func checkERR(err error){

	if err!=nil {
		log.Fatal(err)
	}
}
