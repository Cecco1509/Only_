package main

import (
	"archivemicroservice/Config"
	"archivemicroservice/Models"
	"archivemicroservice/Routers"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

const (
	host     = "archivedb"
	port     = 5432
	user     = "archiveuser"
	password = "password123" // Da segretare
	dbname   = "archivedb"
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
	Config.DB.AutoMigrate(&Models.EncryptedFile{})

	r := Routers.SetupRouter()
	// running
	r.Run()
}