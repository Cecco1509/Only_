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

type StoreFileRequest struct {
	FILENAME  string `json:"filename"`
	SALT      string `json:"salt"`
	IV        string `json:"iv"`
	EXTENSION string `json:"extension"`
	ENC_DATA  string `json:"enc_data"`
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