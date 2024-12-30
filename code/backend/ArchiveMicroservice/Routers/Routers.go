package Routers

import (
	"archivemicroservice/ApiHelpers"
	"archivemicroservice/Controllers"
	"archivemicroservice/Models"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func SecureRoute() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		req, err := http.NewRequest("GET", "https://authservice:5000/verify", nil)
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

		var body ApiHelpers.VerifyResponseData
		err = json.NewDecoder(res.Body).Decode(&body);

		if err != nil {
			ApiHelpers.RespondJSON(c, 500, "Error decoding user")
			return
		}

		c.Set("userInfo", body.Data)

		c.Next()
    })
}

func SuperSecureRoute() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		user := c.MustGet("userInfo").(ApiHelpers.User)

		filename := strings.Split(c.Request.URL.String(), "/")[2]
		fmt.Println("filename : ", filename)
		var ef Models.EncryptedFile
		err := Models.GetFileByFilename(&ef, filename)
		if err != nil {
			ApiHelpers.RespondJSON(c, 404, "File not found")
			c.Abort()
			return
		}

		//fmt.Println("file userID : ", ef.USER_ID, "\nuser id and name : ", user.ID, user.USERNAME)

		if ef.USER_ID != user.ID && c.Request.Method != "GET" { 
			ApiHelpers.RespondJSON(c, 401, "Unauthorized")
			c.Abort()
			return
		}

		if user.ID != ef.USER_ID {

			var sf Models.Shared
			err = Models.GetShared(&sf, ef.ID ,ef.USER_ID, user.ID)

			if err != nil {
				ApiHelpers.RespondJSON(c, 401, "Unauthorized")
				c.Abort()
				return
			}
		}

		fmt.Println("SuperSecureRoute passed")
		c.Next()
	})
}


func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.OPTIONS("efile", CORSMiddleware())
	r.POST("efile", CORSMiddleware(), SecureRoute(), Controllers.StoreEncryptedFile)
	r.GET("efile", CORSMiddleware(), SecureRoute(), Controllers.GetUsersEncryptedFiles)


	// Files shared with a user
	r.OPTIONS("shared/efile/:username", CORSMiddleware())
	r.GET("shared/efile/:username", CORSMiddleware(), SecureRoute(), Controllers.GetUsersSharedEncryptedFiles)

	// Get all the users that are sharing files with the user
	r.OPTIONS("shared/users", CORSMiddleware())
	r.GET("shared/users", CORSMiddleware(), SecureRoute(), Controllers.GetUsersSharingFiles)

	// Get users a file is shared with
	r.OPTIONS("efile/:id/shared", CORSMiddleware())
	r.GET("efile/:id/shared", CORSMiddleware(), SecureRoute(), Controllers.GetSharedWithUsers)

	// Share a file with a user
	r.OPTIONS("share/efile", CORSMiddleware())
	r.OPTIONS("share/efile/:fileid", CORSMiddleware())
	r.POST("share/efile", CORSMiddleware(), SecureRoute(), Controllers.ShareEncryptedFile)
	r.DELETE("share/efile/:fileid", CORSMiddleware(), SecureRoute(), Controllers.RevokeAccess)

	r.GET("storage/:filename", CORSMiddleware(), SecureRoute(), SuperSecureRoute(), Controllers.GetStoredFile)
	r.DELETE("storage/:filename", CORSMiddleware(), SecureRoute(), SuperSecureRoute(), Controllers.DeleteStoredFile)
	r.OPTIONS("storage/:filename", CORSMiddleware())

	return r
}