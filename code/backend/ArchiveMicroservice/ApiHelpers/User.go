package ApiHelpers

import "time"

type User struct {
	ID uint `json:"id"`
	REGISTRATION_DT time.Time `json:"registration_dt"`
	LAST_LOGIN_DT time.Time `json:"last_login_dt"` 
	USERNAME string `json:"username"`
	EMAIL string `json:"email"`
	FIRSTNAME string `json:"firstname"`
	LASTNAME string `json:"lastname"`
  }