package ApiHelpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Status int
	Meta   interface{}
	Data   interface{}
}

type VerifyResponseData struct {
	Status int
	Meta   interface{}
	Data   User
}

type StoreFileRequest struct {
	FILENAME  string `json:"filename"`
	EXTENSION string `json:"extension"`
	ENC_DATA  string `json:"enc_data"`
}

type ShareFileRequest struct {
	FILE_ID uint `json:"file_id"`
	USERNAME string `json:"username"`
}

func RespondJSON(w *gin.Context, status int, payload interface{}) {
	//fmt.Println("status ", status)
	var res ResponseData

	res.Status = status
	//res.Meta = utils.ResponseMessage(status)
	res.Data = payload

	fmt.Printf("set status %d", status)
	w.JSON(status, res)
}