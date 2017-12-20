package article

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Category struct {
	ID	bson.ObjectId `json:"_ID" bson:"_id"`//记录id
	CategoryID int64 `json:"CategoryID" bson:"CategoryID"`//分类id
	CategoryName string `json:"CategoryName" bson:"CategoryName"`//分类名称
	ParentCategory_ID int64 `json:"ParentCategory_ID" bson:"ParentCategory_ID"`//父分类id
	ParentSort  int64 `json:"ParentSort" bson:"ParentSort"`//父分类排序顺序
	SubSort  int64 `json:"SubSort" bson:"SubSort"`//子分类排序顺序
	Define  string `json:"Define" bson:"Define"`//类别定义（语言、框架、工具。。。）
	Tags []Tag `json:"Tags" bson:"Tags"`//标签
	IsEnabled int64 `json:"IsEnabled" bson:"IsEnabled"`//是否可用 （1可用 0不可用）
}

type Tag struct {
	TagID string `json:"TagID" bson:"TagID"`//标签ID
	Tag string `json:"Tag" bson:"Tag"`//标签ID
}
type Article struct {
	ID	bson.ObjectId  `json:"_ID" bson:"_id"`//记录id
	ArticleID bson.ObjectId  `json:"ArticleID" bson:"ArticleID"`	//文章id
	Attribute string  `json:"Attribute" bson:"Attribute"`	//属性（入门教程、安装搭建、配置环境、解决问题、实战案例、疑惑吐槽、推荐介绍、心得分享、技术详解）
	Title string `json:"Title" bson:"Title"`	//标题
	Abstract string `json:"Abstract" bson:"Abstract"`	//摘要
	Author_uid bson.ObjectId `json:"Author_uid" bson:"Author_uid"`//记录id // (外键)	作者UID
	Author string `json:"Author" bson:"Author"`	//作者
	Created	time.Time `json:"Created" bson:"Created"`	//发表时间
	Original int64 `json:"Original" bson:"Original"`	//原创或非原创（转载） value：1 原创 2 转载 3 翻译
	Source_url string `json:"Source_url" bson:"Source_url"`	//源网址（2 转载 3 翻译）    如：www.cnblog.com/java/1234123123
	Source_domain string `json:"Source_domain" bson:"Source_domain"`	//源域名（2 转载 3 翻译）    如：www.cnblog.com
	Source_chinese string `json:"Source_chinese" bson:"Source_chinese"`	//源中文（2 转载 3 翻译）    如：博客园
	CategoryParent_ID string `json:"CategoryParent_ID" bson:"CategoryParent_ID"` //(外键)	父分类id
	CategoryParent string `json:"CategoryParent" bson:"CategoryParent"`	//父分类名称
	CategorySub_ID string `json:"CategorySub_ID" bson:"CategorySub_ID"` //(外键)	子分类id
	CategorySub string `json:"CategorySub" bson:"CategorySub"`	//子分类名
	Tags []string	`json:"Tags" bson:"Tags"`//标签 【数组】如：tags:["java","sprint","git"],
	Like_count	int64 `json:"Like_count" bson:"Like_count"`//喜欢数,默认0
	Dislike_count int64 `json:"Dislike_count" bson:"Dislike_count"`	//不喜欢,默认0
	Comment_count int64	`json:"Comment_count" bson:"Comment_count"`//评论数
	Collected_count	int64 `json:"Collected_count" bson:"Collected_count"`//收藏数
	Read_count int64 `json:"Read_count" bson:"Read_count"`	//阅读数
	IsEnabled int64 `json:"IsEnabled" bson:"IsEnabled"`	//是否可用 （1可用 0不可用）
	Text	string `json:"Text" bson:"Text"`//内容
}

type Comments struct{

	ID	bson.ObjectId `json:"_ID" bson:"_id"`//记录id
	Article_id  bson.ObjectId `json:"Article_id" bson:"Article_id"`// (外键)	文章id
	User_uid  bson.ObjectId `json:"User_uid" bson:"User_uid"`// (外键)	用户id
	Text string `json:"Text" bson:"Text"`//	内容
	Created time.Time `json:"Created" bson:"Created"`//	评论时间
	Like_count int64 `json:"Like_count" bson:"Like_count"`//	评论点赞数
	Replied_count int64 `Replied_count:"_ID" bson:"Replied_count"`//	回复评论总数
	Replieds []string `json:"_ID" bson:"_id"`//	"评论的回复数组
	IsEnabled int64 `json:"_ID" bson:"_id"`//	是否可用 （1可用 0不可用）
}

type Replied struct {
	User_name string	`json:"User_name" bson:"User_name"`//用户名
	UID bson.ObjectId   `json:"UID" bson:"UID"`//用户id
	Text string			`json:"Text" bson:"Text"`//评论内容
	Created time.Time	`json:"Created" bson:"Created"`//评论日期
	IsEnabled int64		`json:"IsEnabled" bson:"IsEnabled"`//是否可用
	Like_count int64	`json:"Like_count" bson:"Like_count"`//点赞数
}