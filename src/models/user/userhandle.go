package user

import (
	"com"
	"log"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func UserRegister(nickName,phone,password string) (uid string) {
	var user *User = new(User)
	user.Uid          = bson.NewObjectId()
	user.NickName     = nickName
	user.Phone    	  = phone
	user.PassWord     = password
	user.RegisterDate = time.Now()

	query := func(c *mgo.Collection) (error) {
		return c.Insert(user)
	}

	err := com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserRegister时报错: %s\n", err)
	}
	return user.Uid.Hex()
}
//

func UserLogin(phone,password string)(user []*User,err error){

	query := func(c *mgo.Collection) (error) {
		return c.Find(bson.M{"Phone":phone,"password":password}).All(&user)

	}

	err = com.GetCollection("User",query)
	if err != nil{
		log.Fatalf("User-UserLogin: %s\n", err)
	}

	if len(user)<1{
		return nil,errors.New("用户不存在")
	}else {
		return user,nil
	}

}