package Controllers

import (
	"archivemicroservice/ApiHelpers"
	"archivemicroservice/Models"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// STORE A ncryptedFile

func StoreEncryptedFile(c *gin.Context) {

	user := c.MustGet("userInfo").(Models.User)
	
	var sentObj ApiHelpers.StoreFileRequest
	err := c.BindJSON(&sentObj)
	if err != nil {
		fmt.Println(err.Error())
		ApiHelpers.RespondJSON(c, 400, err.Error())
		return
	}

	if sentObj.FILENAME != "" {
		var search Models.EncryptedFile
		err = Models.GetFileByFilename(&search, sentObj.FILENAME)
		if err == nil {
			ApiHelpers.RespondJSON(c, 400, "File Name already used already exists")
			return
		} else {
			fmt.Println(err.Error())
		}
	}else {
		sentObj.FILENAME = "IMG_" + strconv.Itoa(int(time.Now().Unix()))
	}

	if !checkExt(sentObj.EXTENSION) {
		ApiHelpers.RespondJSON(c, 400, "Invalid extension")
		return
	}

	ef := Models.EncryptedFile{
		FILENAME:    sentObj.FILENAME,
		USER_ID:     user.ID,
		CREATION_DT: time.Now(),
		SALT:        sentObj.SALT,
		IV:          sentObj.IV,
		EXTENSION:   sentObj.EXTENSION,
	}

	err = Models.CreateFile(&ef)
	if err != nil {
		fmt.Println(err.Error())
		ApiHelpers.RespondJSON(c, 400, "Error creating file")
		return
	}

	if !createFile(ef.FILENAME, sentObj.ENC_DATA) {
		ApiHelpers.RespondJSON(c, 500, "Error creating or writing file")
		return
	}
	
	ApiHelpers.RespondJSON(c, 200, ef)
}

func createFile(filename string, enc_data string) bool {
	f, err := os.Create("/Storage/" + filename)
    if err != nil {
		fmt.Println(err)
		return false
	}

	defer f.Close()

	_, err = f.WriteString(enc_data)

	if err != nil { 
		fmt.Println(err)
		return false
	}
	return true
}

func checkExt(s string) bool {
	switch s {
	case "jpg", "jpeg", "png", "webp":
		return true
	}
	return false
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
