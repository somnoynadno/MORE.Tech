package server

import (
	"MORE.Tech/backend/db"
	"MORE.Tech/backend/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type R struct {
	Message string
	Status  bool
}

func GetUser(c *gin.Context) {
	var result models.User
	id := c.Param("id")

	err := db.GetDB().Preload("InvestProfile").Preload("GameWeek").
		Preload("UserInstruments").Preload("UserInstruments.Instrument").
		Preload("UserInstruments.Instrument.InstrumentType").
		Preload("GameWeek.InstrumentRateChanges").Preload("GameWeek.News").
		Preload("GameWeek.Advices").Preload("GameWeek.TestQuestions").
		Preload("GameWeek.TestQuestions.TestAnswers").
		Preload("GameWeek.Instruments").Preload("GameWeek.Instruments.InstrumentType").
		Preload("Analytics").Preload("Analytics.InvestProfile").
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

	payload.GameWeekID = 1
	payload.Balance = payload.BaseBalance

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
		Preload("News").Preload("Advices").Preload("TestQuestions").
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

func BuyInstrument(c *gin.Context) {

}

func SellInstrument(c *gin.Context) {

}

func AddTestAnswer(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	testAnswerID, _ := strconv.ParseUint(c.Param("test_answer_id"), 10, 64)


	userTestAnswer := models.UserTestAnswer{
		TestAnswerID: uint(testAnswerID),
		UserID: uint(id),
	}

	err := db.GetDB().Create(&userTestAnswer).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, R{Message: "OK", Status: true})
}

func GetGameResult(c *gin.Context) {

}