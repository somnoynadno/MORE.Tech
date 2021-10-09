package server

import (
	"MORE.Tech/backend/db"
	"MORE.Tech/backend/models"
	"github.com/gin-gonic/gin"
)

func GetTestQuestions(c *gin.Context) {
	var result []models.TestQuestion

	err := db.GetDB().Preload("TestAnswers").Find(&result).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, result)
}

func GetUser(c *gin.Context) {
	var result models.User
	id := c.Param("id")

	err := db.GetDB().Preload("InvestProfile").Preload("GameWeek").
		Preload("Instruments").Preload("GameWeek.InstrumentRateChanges").
		Preload("GameWeek.News").Preload("GameWeek.Advices").
		Preload("GameWeek.Instruments").Preload("GameWeek.Instruments.InstrumentType").
		First(&result, id).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, result)
}

func CreateUser(c *gin.Context) {
	var payload models.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		handleBadRequest(c, err)
		return
	}

	payload.InvestProfileID = 1
	payload.GameWeekID = 1
	payload.Balance = 100000

	err := db.GetDB().Create(&payload).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, payload)
}

func GetGameWeek(c *gin.Context) {
	var result models.GameWeek
	id := c.Param("id")

	err := db.GetDB().Preload("InstrumentRateChanges").
		Preload("News").Preload("Advices").
		Preload("Instruments").Preload("Instruments.InstrumentType").
		First(&result, id).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, result)
}

func NextWeek(c *gin.Context) {
	var result models.User
	id := c.Param("id")

	err := db.GetDB().First(&result, id).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	result.GameWeekID += 1
	// TODO: make balance changes

	err = db.GetDB().Model(&result).Updates(result).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, result)
}