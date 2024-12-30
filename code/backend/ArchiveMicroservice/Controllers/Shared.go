package Controllers

import (
	"archivemicroservice/ApiHelpers"
	"archivemicroservice/Models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func GetUsersSharedEncryptedFiles(c *gin.Context) {

	user := c.MustGet("userInfo").(ApiHelpers.User)

	username := c.Param("username")

	req, err := http.NewRequest("GET", "https://authservice:5000/userinfo/"+username, nil)
	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error creating validate request")
		return
	}
	req.Header.Set("Authorization", c.Request.Header.Get("Authorization"))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error validating token")
		return
	}

	if res.StatusCode != 200 {
		ApiHelpers.RespondJSON(c, res.StatusCode, "Error validating token")
		return
	}

	var resData ApiHelpers.VerifyResponseData
	err = json.NewDecoder(res.Body).Decode(&resData);
	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error decoding user")
		return
	}

	u := resData.Data

	var ef []Models.EncryptedFile
	err = Models.GetUserSharedWithFiles(&ef, u.ID, user.ID)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, ef)
		return
	}

	ApiHelpers.RespondJSON(c, 200, ef)
}

func GetSharedWithUsers(c *gin.Context) {

	user := c.MustGet("userInfo").(ApiHelpers.User)

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "Invalid id")
		return
	}

	var ef Models.EncryptedFile
	err = Models.GetFile(&ef, uint(id))
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, "File not found")
		return
	}

	if ef.USER_ID != user.ID {
		ApiHelpers.RespondJSON(c, 401, "Unauthorized")
		return
	}

	var sharedWith []Models.Shared
	err = Models.GetSharedWithUsers(&sharedWith, ef.ID)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, sharedWith)
		return
	}

	ApiHelpers.RespondJSON(c, 200, sharedWith)
}

func RevokeAccess(c *gin.Context) {
	
	user := c.MustGet("userInfo").(ApiHelpers.User)

	id, err := strconv.Atoi(c.Param("fileid"))

	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "Invalid id")
		return
	}

	userId, err := strconv.Atoi(c.Query("userid"))
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "Invalid id")
		return
	}

	var shared Models.Shared
	err = Models.GetShared(&shared, uint(id), user.ID ,uint(userId))
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, "Shared file not found")
		return
	}

	var ef Models.EncryptedFile
	err = Models.GetFile(&ef, shared.FILE_ID)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, "File not found")
		return
	}

	if ef.USER_ID != user.ID {
		ApiHelpers.RespondJSON(c, 401, "Unauthorized")
		return
	}

	err = Models.DeleteShared(&shared)
	if err != nil {
		ApiHelpers.RespondJSON(c, 500, "Error revoking access")
		return
	}

	ApiHelpers.RespondJSON(c, 200, "Access revoked")
}

func GetUsersSharingFiles(c *gin.Context) {

	user := c.MustGet("userInfo").(ApiHelpers.User)

	var shared []Models.Shared
	err := Models.GetShareUsers(&shared, user.ID)
	if err != nil {
		fmt.Println(err)
		ApiHelpers.RespondJSON(c, 404, shared)
		return
	}

	fmt.Println(shared)

	ApiHelpers.RespondJSON(c, 200, shared)
}