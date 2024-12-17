package Controllers

import (
	"authmicroservice/ApiHelpers"
	"authmicroservice/Models"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	Userid   uint `json:"userid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var invalidTokens = []string{}

var secretKey = []byte("secret-key")

func createToken(username string, userid uint) (string, error) {

	claims := UserClaims{
		Userid: userid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

 	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is signed with the expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Optionally, check token claims
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		for _, invalidToken := range invalidTokens {
			if tokenString == invalidToken {
				return nil, fmt.Errorf("invalid token")
			}
		}

		return token, nil

	}else {
		fmt.Println(token)
	}

	return nil, fmt.Errorf("invalid token claims")
}


func Login(c *gin.Context) {
	var req ApiHelpers.AuthRequest
	err := c.BindJSON(&req)

	if err != nil {
		ApiHelpers.RespondJSON(c, 400, err.Error())
		return
	}

	var user Models.User
	err = Models.GetUserByUsername(&user, req.Username)

	if err != nil {
		ApiHelpers.RespondJSON(c, 404, err.Error())
		return
	}

	password := []byte(req.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.PASSHASH), password)
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "Wrong password")
		return
	}

	user.LAST_LOGIN_DT = time.Now()
	err = Models.UpdateUser(&user, user.ID)

	if err != nil {
		ApiHelpers.RespondJSON(c, 500, err.Error())
		return
	}

	token, err := createToken(user.USERNAME, user.ID)

	if err != nil {
		ApiHelpers.RespondJSON(c, 500, err.Error())
		return
	}

	var resp ApiHelpers.AuthResponse
	resp.Token = token

	ApiHelpers.RespondJSON(c, 200, resp)
}

func Register(c *gin.Context) {
	var req ApiHelpers.RegisterRequest
	err := c.BindJSON(&req)

	if err != nil {
		ApiHelpers.RespondJSON(c, 400, err.Error())
		return
	}

	var user Models.User
	err = Models.GetUserByUsername(&user, req.Username)

	if err == nil {
		ApiHelpers.RespondJSON(c, 400, "User already exists")
		return
	}

	password := []byte(req.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		ApiHelpers.RespondJSON(c, 500, err.Error())
		return
	}

	new_user := Models.User{
		USERNAME: req.Username,
		PASSHASH: string(hashedPassword),
		REGISTRATION_DT: time.Now(),
		EMAIL: req.Email,
		FIRSTNAME: req.Firstname,
		LASTNAME: req.Lastname,
	}

	err = Models.CreateUser(&new_user)

	if err != nil {
		ApiHelpers.RespondJSON(c, 500, err.Error())
		return
	}

	ApiHelpers.RespondJSON(c, 200, "User created")
}

func Logout(c *gin.Context) {
	token := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	fmt.Println(token)
	_, err := verifyToken(token)

	if err != nil {
		ApiHelpers.RespondJSON(c, 401, "Invalid token")
		return
	}

	invalidTokens = append(invalidTokens, token)

	ApiHelpers.RespondJSON(c, 200, "Succesfully logged out")
}

func Delete(c *gin.Context) {
	
}

func VerifyToken(c *gin.Context) {
	token := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	fmt.Println(token)
	decodedToken, err := verifyToken(token)

	if err != nil {
		ApiHelpers.RespondJSON(c, 401, err.Error())
		return
	}

	decodedClaims := map[string]interface{}{
		"userid": decodedToken.Claims.(jwt.MapClaims)["userid"],
		"username": decodedToken.Claims.(jwt.MapClaims)["username"],
	}

	var user Models.User
	userId := decodedClaims["userid"].(float64)
	err = Models.GetUserByUserId(&user, int(userId))

	if err != nil {
		ApiHelpers.RespondJSON(c, 404, "User not found")
	}

	ApiHelpers.RespondJSON(c, 200, user)
}
