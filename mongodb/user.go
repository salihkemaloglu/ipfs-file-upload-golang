package mongodb

import (
	"gopkg.in/mgo.v2/bson"
)

// Mongodb User database structs
type User struct {
	Id          	 bson.ObjectId `bson:"_id" json:"id" `
	Name 	 		 string        `bson:"name" json:"name"`
	Surname 		 string        `bson:"surname" json:"surname"`
	Email 	    	 string        `bson:"email" json:"email"`
	Username 		 string        `bson:"username" json:"username"`
	Password 		 string        `bson:"password" json:"password"`
	Description 	 string        `bson:"description" json:"description"`
	CreatedDate  	 string        `bson:"createddate" json:"createddate"`
	UpdatedDate  	 string        `bson:"updateddate" json:"updateddate"`
	ImagePath 		 string        `bson:"imagepath" json:"imagepath"`
	TotalSpace  	 string        `bson:"totalspace" json:"totalspace"`
	LanguageType 	 string        `bson:"languagetype" json:"languagetype"`
}
