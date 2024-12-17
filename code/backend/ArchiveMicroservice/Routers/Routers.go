package Routers

import (
	"archivemicroservice/ApiHelpers"
	"archivemicroservice/Controllers"
	"archivemicroservice/Models"
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	
		if err != nil || res.StatusCode != 200 {
			ApiHelpers.RespondJSON(c, res.StatusCode, "Error validating token")
			return
		}

		var u Models.User
		err = json.NewDecoder(res.Body).Decode(&u);

		if err != nil {
			ApiHelpers.RespondJSON(c, 500, "Error decoding user")
			return
		}

		c.Set("userInfo", u)

		c.Next()
    })
}


func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("efile", SecureRoute(), Controllers.StoreEncryptedFile)
	r.GET("efile", SecureRoute(), Controllers.GetEncryptedFile)

	return r
}