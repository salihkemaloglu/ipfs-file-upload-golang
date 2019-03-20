
package mongodb

import (
	"gopkg.in/mgo.v2/bson"
)

type File struct {
	Id          	 bson.ObjectId `bson:"_id" json:"id" `
	UserId 	    	 string        `bson:"userid" json:"userid"`
	FolderId 		 string        `bson:"folderid" json:"folderid"`
	Name 	 		 string        `bson:"name" json:"name"`
	Description 	 string        `bson:"description" json:"description"`
	CreatedDate 	 string        `bson:"createddate" json:"createddate"`
	UpdatedDate 	 string        `bson:"updateddate" json:"updateddate"`
	FileHash    	 string        `bson:"filehash" json:"filehash"`
	IsBuried    	 bool          `bson:"isburied" json:"isburied"`
	IsFolderFile	 bool          `bson:"isfolderfile" json:"isfolderfile"`
	IsStarred   	 bool          `bson:"isstarred" json:"isstarred"`
	IsTrash     	 bool          `bson:"istrash" json:"istrash"`
	IsDeleted   	 bool          `bson:"isdeleted" json:"isdeleted"`
}
