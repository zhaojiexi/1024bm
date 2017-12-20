package article

import (
	"com"
	"log"
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
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


	a.ID=bson.NewObjectId()
	a.ArticleID=bson.NewObjectId()
	a.IsEnabled=1
	a.Like_count=0
	a.Dislike_count=0
	a.Comment_count=0
	a.Collected_count=0
	a.Read_count=0

	query := func(c *mgo.Collection) (error) {
		return c.Insert(a)
	}
	fmt.Println(a)
	err := com.GetCollection("Article",query)
	if err != nil{
		log.Fatalf("Article-AddCategory时报错: %s\n", err)
	}


	return ""
}