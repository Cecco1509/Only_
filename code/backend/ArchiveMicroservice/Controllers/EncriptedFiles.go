package Controllers

import (
	"archivemicroservice/ApiHelpers"
	"archivemicroservice/Models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// STORE A ncryptedFile

func StoreEncryptedFile(c *gin.Context) {

	user := c.MustGet("userInfo").(Models.User)
	
	var sentObj Models.EncryptedFile
	c.BindJSON(&sentObj)

	if sentObj.FILENAME != "" {
		var search Models.EncryptedFile
		err := Models.GetFileByFilename(&search, sentObj.FILENAME)
		if err != nil {
			ApiHelpers.RespondJSON(c, 400, "File Name already used already exists")
			return
		}
	}

	ef := Models.EncryptedFile{
		FILENAME:    sentObj.FILENAME,
		USER_ID:     user.ID,
		CREATION_DT: time.Now(),
		SALT:        sentObj.SALT,
		IV:          sentObj.IV,
		EXTENSION:   sentObj.EXTENSION,
	}

	err := Models.CreateFile(&ef)
	if err != nil {
		fmt.Println(err.Error())
		ApiHelpers.RespondJSON(c, 400, "Error creating file")
		return
	}
	
	ApiHelpers.RespondJSON(c, 200, ef)
}

// GET A ncryptedFile

func GetEncryptedFile(c *gin.Context) {

	//var user Models.User
	//user = c.MustGet("userInfo").(Models.User)

	strId := c.Params.ByName("id")
	var ef Models.EncryptedFile
	id, err := strconv.Atoi(strId)
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "Invalid ID")
		return
	}

	err = Models.GetFile(&ef, id)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, ef)
		return
	}

	ef.LAST_RETREIVED_DT = time.Now()
	Models.UpdateFile(&ef, ef.ID)
	ApiHelpers.RespondJSON(c, 200, ef)
	
}

// DELETE A ncryptedFile

// UPDATE A ncryptedFile

// GET ALL ncryptedFileS (pagination?)
