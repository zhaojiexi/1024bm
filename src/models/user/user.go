package user

import(
	"time"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Uid 		 bson.ObjectId 	`json:"Uid"           bson:"Uid"`           //用户id
	NickName 	 string 		`json:"NickName"      bson:"NickName"`      //昵称
	Phone 		 string 		`json:"Phone"         bson:"Phone"`    		//手机
	PassWord     string         `json:"PassWord" 	  bson:"PassWord"` 		//密码
	RegisterDate time.Time 		`json:"RegisterDate"  bson:"RegisterDate"`  //注册日期
}

























