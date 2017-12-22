package article

import (
	"com"
	"log"
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"models/user"
	"time"
)
//新增分类
func AddCategory(category *Category)string{

	var Categorys []Category

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"CategoryName":category.CategoryName,"IsEnabled":1}).All(&Categorys)
	}

	err := com.GetCollection("Category",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
	if len(Categorys)>0 {
		return "该分类已存在"
	}




	//取分类id最大的值 +1 作为新的分类id的值
	query = func(c *mgo.Collection) (error) {
		return c.Find(nil).Sort("-CategoryID").Limit(1).All(&Categorys)
	}

	err = com.GetCollection("Category",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
	//如果没有分类记录 则赋默认值1
	if len(Categorys)==0 {
		category.CategoryID=1
	}else{
		category.CategoryID=Categorys[0].CategoryID+1
	}

	category.SubSort=1
	category.ParentSort=1
	category.IsEnabled=1

	query = func(c *mgo.Collection) (error) {
		return c.Insert(category)
	}
	fmt.Println(category)
	err = com.GetCollection("Category",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
	return ""
}

//新增文章
func AddArticle(a *Article)string{


	var cgs []Category
	var us []user.User

	//校验子分类是否存在
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"CategoryID":a.CategorySub_ID}).All(&cgs)
	}


	err := com.GetCollection("Category",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
	if len(cgs)<1 {
		return "子分类不存在"
	}

	query = func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"CategoryID":a.CategorySub_ID,"ParentCategory_ID":a.CategoryParent_ID}).All(&cgs)
	}

	err = com.GetCollection("Category",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
	if len(cgs)<1 {
		return "主分类下没有该子分类"
	}

	query = func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Uid":a.Author_uid,"IsEnabled":1}).All(&us)
	}

	err = com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("Article-AddCategory2时报错: %s\n", err)
	}
	if len(us)<1 {
		return "用户不存在或为不可用状态"
	}


	a.ID=bson.NewObjectId()
	a.ArticleID=bson.NewObjectId()
	a.IsEnabled=1
	a.Like_count=0
	a.Dislike_count=0
	a.Comment_count=0
	a.Collected_count=0
	a.Read_count=0

	query = func(c *mgo.Collection) (error) {
		return c.Insert(a)
	}

	err = com.GetCollection("Article",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}


	return ""
}

//根据文章id查找文章
func GetArticleByID(aid string)(article *Article,result string){

	abjectid:=bson.ObjectIdHex(aid)

	var arts []Article

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"ArticleID":abjectid,"IsEnabled":1}).All(&arts)
	}
	fmt.Println("cgs",arts)


	err := com.GetCollection("Article",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
	if len(arts)<1 {
		return nil,"文章不存在或状态不可用"
	}

	return &arts[0],""
}

func GetArticleAll(pagenum,pagecount int)([]Article,int,int,int,int){

	var arts []Article
	//查询总条数
	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"IsEnabled":1}).All(&arts)
	}

	err := com.GetCollection("Article",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
	pageSum:=len(arts)
	fmt.Println("pageSum",pageSum)

	PageMax:=pageSum/pagecount
	//取模 如果不能整除 最大页数+1
	if pageSum%pagecount!=0 {
		PageMax++
	}
	//如果输入的页数大于最大页数 则=最大页数
	if pagenum>PageMax&&PageMax!=0 {
		pagenum=PageMax
	}
	//如果输入的小于最大页数 则=1 第一页
	if pagenum<1 {
		pagenum=1
	}
	query = func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"IsEnabled":1}).Skip((pagenum-1)*pagecount).Limit(pagecount).All(&arts)
	}

	err = com.GetCollection("Article",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
	return arts,pagenum,pagecount,pageSum,PageMax
}
//新增评论
func AddComments(cm *Comments){

	cm.Like_count=0
	cm.Replied_count=0
	cm.IsEnabled=1
	cm.ID=bson.NewObjectId()
	vcd:=time.Now()

	fmt.Println(vcd)
	for k,_:=range cm.Replieds{

		cm.Replieds[k].Like_count=0
		cm.Replieds[k].Created=vcd
		cm.Replieds[k].IsEnabled=1
	}


	query := func(c *mgo.Collection) (error) {
		return c.Insert(cm)
	}

	err := com.GetCollection("Comments",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}
}
