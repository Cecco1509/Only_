package ApiHelpers

import (
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Status int
	Meta   interface{}
	Data   interface{}
}


type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	ExpiresAt int64 `json:"expiresAt"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func RespondJSON(w *gin.Context, status int, payload interface{}) {
	//fmt.Println("status ", status)
	var res ResponseData

	res.Status = status
	//res.Meta = utils.ResponseMessage(status)
	res.Data = payload

	w.JSON(status, res)
}