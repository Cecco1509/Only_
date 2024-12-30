package Controllers

import (
	"archivemicroservice/ApiHelpers"
	"archivemicroservice/Models"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// STORE A ncryptedFile

func StoreEncryptedFile(c *gin.Context) {

	user := c.MustGet("userInfo").(ApiHelpers.User)
	
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

	res, ext := checkExt(sentObj.EXTENSION)

	if !res {
		ApiHelpers.RespondJSON(c, 400, "Invalid extension")
		return
	}

	fmt.Println("user: ", user)

	ef := Models.EncryptedFile{
		FILENAME:    sentObj.FILENAME,
		USER_ID:     user.ID,
		EXTENSION:   ext,
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
	
	fmt.Println("ef: ", ef)
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

func checkExt(s string) (bool, string) {
	ext := strings.Split(s, "/")[1]
	switch ext {
		case "jpg", "jpeg", "png", "webp":
			return true, ext
	}
	switch s {
	case "jpg", "jpeg", "png", "webp":
		return true, s
	}
	return false, ""
}

// GET A ncryptedFile

func GetEncryptedFile(c *gin.Context) {

	user := c.MustGet("userInfo").(ApiHelpers.User)

	strId := c.Params.ByName("id")
	var ef Models.EncryptedFile
	id, err := strconv.Atoi(strId)
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "Invalid ID")
		return
	}

	err = Models.GetFile(&ef, uint(id))
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, "File not found")
		return
	}

	if ef.USER_ID != user.ID {
		ApiHelpers.RespondJSON(c, 403, "Unauthorized")
		return
	}

	ef.LAST_RETREIVED_DT = time.Now()
	Models.UpdateFile(&ef, ef.ID)
	ApiHelpers.RespondJSON(c, 200, ef)
	
}

func GetUsersEncryptedFiles(c *gin.Context) {

	user := c.MustGet("userInfo").(ApiHelpers.User)

	var ef []Models.EncryptedFile
	err := Models.GetUsersFiles(&ef, user.ID)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, ef)
		return
	}

	fmt.Println(ef)

	ApiHelpers.RespondJSON(c, 200, ef)
}

func GetStoredFile(c *gin.Context) {
	
	filename := c.Params.ByName("filename")
	file, err := os.Open("/Storage/" + filename)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, "File not found")
		return
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	_, err = file.Read(bytes)
	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error reading file")
		return
	}

	c.Writer.WriteString(string(bytes))
}

func ShareEncryptedFile(c *gin.Context) {
	user := c.MustGet("userInfo").(ApiHelpers.User)

	var sharedFileReq ApiHelpers.ShareFileRequest
	err := c.BindJSON(&sharedFileReq)
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, err.Error())
		return
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	userInfoReq, err := http.NewRequest("GET", "https://authservice:5000/userinfo/"+sharedFileReq.USERNAME, nil)
	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error creating validate request")
		return
	}
	userInfoReq.Header.Set("Authorization", c.Request.Header.Get("Authorization"))

	res, err := http.DefaultClient.Do(userInfoReq)

	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error validating token")
		return
	}

	if res.StatusCode != 200 {
		ApiHelpers.RespondJSON(c, res.StatusCode, "Error validating token")
		return
	}

	var body ApiHelpers.VerifyResponseData
	err = json.NewDecoder(res.Body).Decode(&body);

	fmt.Println("body: ", body)

	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error decoding user")
		return
	}

	shareWith := body.Data

	var ef Models.EncryptedFile
	err = Models.GetFile(&ef, sharedFileReq.FILE_ID)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, "File not found")
		return
	}

	if ef.USER_ID != user.ID {
		ApiHelpers.RespondJSON(c, 401, "Unauthorized")
		return
	}

	var checkShared Models.Shared
	err = Models.GetShared(&checkShared, ef.ID, user.ID, shareWith.ID)

	if err == nil {
		ApiHelpers.RespondJSON(c, 400, "Cannot share twice")
		return
	}

	var shared Models.Shared
	shared.FILE_ID = ef.ID
	shared.SHARED_BY_USER_ID = user.ID
	shared.SHARED_WITH_USERNAME = shareWith.USERNAME
	shared.SHARED_WITH_USER_ID = shareWith.ID
	shared.SHARED_BY_USERNAME = user.USERNAME

	err = Models.CreateShared(&shared)
	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error sharing file")
		return
	}

	ApiHelpers.RespondJSON(c, 200, shared)
}

func DeleteStoredFile(c *gin.Context) {
	filename := c.Params.ByName("filename")

	var ef Models.EncryptedFile
	err := Models.DeleteFileByFilename(&ef, filename)

	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error deleting file")
		return
	}

	var shared []Models.Shared
	err = Models.DeleteShares(&shared, ef.ID)

	err = os.Remove("/Storage/" + filename)
	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error deleting file")
		return
	}

	ApiHelpers.RespondJSON(c, 200, "File deleted")
}

// DELETE A ncryptedFile

// UPDATE A ncryptedFile

// GET ALL ncryptedFileS (pagination?)
