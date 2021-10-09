package server

import (
	"MORE.Tech/backend/db"
	"MORE.Tech/backend/models"
	"errors"
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
	id, err1 := strconv.ParseUint(c.Param("id"), 10, 64)
	instrumentID, err2 := strconv.ParseUint(c.Param("instrument_id"), 10, 64)

	if err1 != nil || err2 != nil {
		handleBadRequest(c, errors.New("bad path parameters"))
		return
	}

	var user models.User
	err := db.GetDB().First(&user, id).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	var instrument models.Instrument
	err = db.GetDB().First(&instrument, instrumentID).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	if instrument.BasePrice > user.Balance {
		handleBadRequest(c, errors.New("not enough money"))
		return
	}

	user.Balance -= instrument.BasePrice
	db.GetDB().Model(&user).Updates(user)

	userInstrument := models.UserInstrument{
		UserID: uint(id),
		InstrumentID: uint(instrumentID),
	}

	err = db.GetDB().Create(&userInstrument).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, R{Message: "OK", Status: true})
}

func SellInstrument(c *gin.Context) {
	id, err1 := strconv.ParseUint(c.Param("id"), 10, 64)
	instrumentID, err2 := strconv.ParseUint(c.Param("instrument_id"), 10, 64)

	if err1 != nil || err2 != nil {
		handleBadRequest(c, errors.New("bad path parameters"))
		return
	}

	var user models.User
	err := db.GetDB().First(&user, id).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	var userInstrument models.UserInstrument
	err = db.GetDB().Where("user_id = ? and instrument_id = ?", id, instrumentID).First(&userInstrument).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	user.Balance += userInstrument.CurrentPrice
	db.GetDB().Model(&user).Updates(user)
	db.GetDB().Delete(&userInstrument)

	handleOK(c, R{Message: "OK", Status: true})
}

func AddTestAnswer(c *gin.Context) {
	id, err1 := strconv.ParseUint(c.Param("id"), 10, 64)
	testAnswerID, err2 := strconv.ParseUint(c.Param("test_answer_id"), 10, 64)

	if err1 != nil || err2 != nil {
		handleBadRequest(c, errors.New("bad path parameters"))
		return
	}

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