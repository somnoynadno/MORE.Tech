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
	var user models.User
	id := c.Param("id")

	err := db.GetDB().Preload("InvestProfile").Preload("GameWeek").
		Preload("UserInstruments").Preload("UserInstruments.Instrument").
		Preload("UserInstruments.Instrument.InstrumentType").
		Preload("GameWeek.InstrumentRateChanges").Preload("GameWeek.News").
		Preload("GameWeek.Advices").Preload("GameWeek.TestQuestions").
		Preload("GameWeek.TestQuestions.TestAnswers").
		Preload("GameWeek.Instruments").Preload("GameWeek.Instruments.InstrumentType").
		Preload("Analytics").Preload("Analytics.InvestProfile").Preload("TestAnswers").
		First(&user, id).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, user)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handleBadRequest(c, err)
		return
	}

	user.GameWeekID = 1
	user.Balance = user.BaseBalance

	err := db.GetDB().Create(&user).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, user)
}

func GetGameWeek(c *gin.Context) {
	var gameWeek models.GameWeek
	id := c.Param("id")

	err := db.GetDB().Preload("InstrumentRateChanges").
		Preload("News").Preload("Advices").Preload("TestQuestions").
		Preload("Instruments").Preload("Instruments.InstrumentType").
		First(&gameWeek, id).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, gameWeek)
}

func NextWeek(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	err := db.GetDB().First(&user, id).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	var instrumentRateChanges []models.InstrumentRateChange
	err = db.GetDB().Where("game_week_id = ?", user.GameWeekID).Find(&instrumentRateChanges).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	var userInstruments []models.UserInstrument
	err = db.GetDB().Preload("Instrument").Where("user_id = ?", user.ID).Find(&userInstruments).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	balanceDelta := 0
	for _, userInstrument := range userInstruments {
		for _, rateChange := range instrumentRateChanges {
			if rateChange.InstrumentID == userInstrument.InstrumentID {
				delta := int(rateChange.PriceChangeRate * float64(userInstrument.CurrentPrice) / 100)
				balanceDelta += delta

				userInstrument.CurrentPrice += delta
				userInstrument.PriceChangedRate = float64(userInstrument.CurrentPrice) / float64(userInstrument.Instrument.BasePrice)
				userInstrument.PriceChanged = userInstrument.CurrentPrice - userInstrument.Instrument.BasePrice

				db.GetDB().Model(&userInstrument).Updates(userInstrument)

				if rateChange.AdditionalPayment != nil {
					delta += *rateChange.AdditionalPayment
				}
			}
		}
	}

	user.GameWeekID += 1
	user.Balance += balanceDelta

	err = db.GetDB().Model(&user).Updates(user).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, user)
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
		CurrentPrice: instrument.BasePrice,
		PriceChangedRate: 0.0,
		PriceChanged: 0,
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