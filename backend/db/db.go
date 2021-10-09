package db

import (
	"MORE.Tech/backend/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"os"
)

var db *gorm.DB

func init() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	conn, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			dbHost, dbPort, username, dbName, password))
	if err != nil {
		panic(err)
	} else {
		db = conn
		log.Info("[DB] connected on " + dbHost)
	}

	err = migrateSchema()
	if err != nil {
		log.Fatal("[DB] error migrating schema: " + err.Error())
	} else {
		log.Info("[DB] schema migrated successfully")
	}
}

func GetDB() *gorm.DB {
	return db
}

func migrateSchema() error {
	err := db.AutoMigrate(
		models.UserInstrument{},
		models.UserTestAnswer{},

		models.Advice{},
		models.GameWeek{},
		models.Instrument{},
		models.InstrumentType{},
		models.InstrumentRateChange{},
		models.News{},
		models.User{},

		models.InvestProfile{},
		models.TestAnswer{},
		models.TestQuestion{},

		models.Analytics{},
	).Error

	return err
}
