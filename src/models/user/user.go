package user

import(
	"time"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	_ID			 	bson.ObjectId 	`json:"_ID"            bson:"_ID"`			//记录id
	Uid 		 	string		 	`json:"Uid"            bson:"Uid"`           //用户id
	Name 	 	 	string 			`json:"Name"    	   bson:"Name"`     		 //用户名、昵称
	Slug 		 	string 			`json:"Slug"    	   bson:"Slug"`     		 //昵称+唯一编号
	Phone 		 	string 			`json:"Phone"          bson:"Phone"`    		//手机
	PassWord     	string          `json:"PassWord" 	   bson:"PassWord"` 		//密码
	RegisterDate 	time.Time 		`json:"RegisterDate"   bson:"RegisterDate"`  //注册日期
	Location	 	string			`json:"Location"  	   bson:"Location"` 		//所在地
	University	 	string			`json:"University"     bson:"University"` 	//毕业院校
	Company		 	string			`json:"Company"   	   bson:"Company"`		//所在公司
	WebSite	        string		    `json:"WebSite"   	   bson:"WebSite"`		//展示网站
	Follower_count 	int64			`json:"Follower_count" bson:"Follower_count"`//关注我的人数（粉丝）
	Following_count int64			`json:"Following_count" bson:"Following_count"`	//我关注的人数
	Browse_count    int64			`json:"Browse_count" bson:"Browse_count"`//浏览数（文章 ）
	Article_count   int64			`json:"Article_count" bson:"Article_count"` //发表文章数  
	Describe	    string	 	    `json:"Describe"   	   bson:"Describe"`//个人介绍
	Profile_image_url	 string	    `json:"Profile_image_url"   	   bson:"Profile_image_url"`//头像地址
	LastLogin		time.Time 		`json:"LastLogin"   bson:"LastLogin"` //最后登录时间
	Interest		[]Interest		`json:"Interest" bson:"Interest"`//感兴趣（多选）
	IsEnabled	    int64			`json:"IsEnabled" bson:"IsEnabled"`//是否可用 （1可用 0不可用）
}

type Interest struct{
	SubCategoryID string   	`json:"SubCategoryID"          bson:"SubCategoryID"`		//子分类ID
	SubCategory string		`json:"SubCategory"            bson:"SubCategory"`			//子分类名称
}

























