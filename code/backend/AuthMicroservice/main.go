package main

import (
	"authmicroservice/Config"
	"authmicroservice/Models"
	"authmicroservice/Routers"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

const (
	host     = "authdb"
	port     = 5432
	user     = "authuser"
	password = "password123" // Da segretare
	dbname   = "authdb"
  )

var err error

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    						"password=%s dbname=%s sslmode=disable",
    						host, port, user, password, dbname)

	Config.DB, err = gorm.Open(postgres.Open(psqlInfo))

	if err != nil {
		fmt.Println("statuse: ", err)
	}
	Config.DB.AutoMigrate(&Models.User{})

	r := Routers.SetupRouter()
	// running
	r.RunTLS(":5000", "/run/secrets/auth_cert", "/run/secrets/auth_key")
}