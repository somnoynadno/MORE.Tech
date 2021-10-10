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
	user.Sells = 0
	user.Balance = user.BaseBalance

	err := db.GetDB().Create(&user).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	db.GetDB().Create(&models.UserInstrument{InstrumentID: 21, UserID: user.ID, CurrentPrice: 7740})
	db.GetDB().Create(&models.UserInstrument{InstrumentID: 20, UserID: user.ID, CurrentPrice: 5320})
	db.GetDB().Create(&models.UserInstrument{InstrumentID: 19, UserID: user.ID, CurrentPrice: 6480})

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
				userInstrument.PriceChangedRate = float64(userInstrument.CurrentPrice) / float64(userInstrument.Instrument.BasePrice) - 1
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
		UserID:           uint(id),
		InstrumentID:     uint(instrumentID),
		CurrentPrice:     instrument.BasePrice,
		PriceChangedRate: 0.0,
		PriceChanged:     0,
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
	user.Sells += 1
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
		UserID:       uint(id),
	}

	err := db.GetDB().Create(&userTestAnswer).Error
	if err != nil {
		handleInternalError(c, err)
		return
	}

	handleOK(c, R{Message: "OK", Status: true})
}

func GetGameResult(c *gin.Context) {
	var user models.User
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		handleBadRequest(c, err)
		return
	}

	err = db.GetDB().Preload("InvestProfile").Preload("GameWeek").Preload("TestAnswers").
		Preload("UserInstruments").Preload("UserInstruments.Instrument").
		First(&user, id).Error
	if err != nil {
		handleBadRequest(c, err)
		return
	}

	correctAnswers := 0
	for _, i := range user.TestAnswers {
		if i.IsCorrect {
			correctAnswers += 1
		}
	}

	income := 0
	instrumentSum := 0
	for _, i := range user.UserInstruments {
		income += i.CurrentPrice
		instrumentSum += int(i.InstrumentID)
	}

	counter := 0
	instrumentsBalanceVerdict := "Вам не удалось сбалансировать портфель"
	if len(user.UserInstruments) > 4 {
		counter += 1
		instrumentsBalanceVerdict = "В целом, вы неплохо сбалансировали свой портфель"

		if len(user.UserInstruments) > 6 {
			instrumentsBalanceVerdict = "Ваш портфель вышел отлично сбалансированным"
		}
	}
	if income + user.Balance > user.BaseBalance {
		counter += 2
	}
	if instrumentSum > 9 {
		counter += 1
	}

	financialCushionVerdict := "Вы инвестировали почти все средства, не оставив финансовой подушки"
	if float64(income / user.BaseBalance) < 0.35 {
		financialCushionVerdict = "Вы очень аккуратно распределяли свои денежные средства"
		counter += 2
	} else if float64(income / user.BaseBalance) < 0.6 {
		financialCushionVerdict = "Вы создали некоторую финансовую подушку, как мы вам и советовали"
		counter += 1
	}

	tradingStrategyVerdict := "Вы придерживались краткосрочной стратегии продаж"
	if user.Sells <= 3 {
		tradingStrategyVerdict = "Вы приобретали ценные бумаги только на длительный срок"
		counter += 2
	} else if user.Sells <= 8 {
		tradingStrategyVerdict = "Вы предпочитали открывать позиции на длинный и средний срок"
		counter += 1
	}

	testAnswersVerdict := "К сожалению, вы плохо усвоили финансовую теорию"
	if correctAnswers >= 4 {
		testAnswersVerdict = "Вы прекрасно усвоили основы инвестирования"
		counter += 2
	} else if correctAnswers >= 2 {
		testAnswersVerdict = "Вы усвоили теорию, но во многих вопросах всё ещё плохо разбираетесь"
		counter += 1
	}

	totalVerdict := "Конечно, вы можете попробовать себя в инвестициях, но мы советуем вам лучше к этому подготовиться"
	if counter > 6 {
		totalVerdict = "Мы считаем, что вы превосходно справились с испытанием и можете спокойно начинать инвестировать"
	} else if counter > 3 {
		totalVerdict = "Мы считаем, что вы неплохо справились с испытанием и готовы начать инвестировать"
	}

	analytics := models.Analytics{
		TotalIncome: user.Balance + income,
		TotalIncomeRate: float64((user.Balance +income) / user.BaseBalance) - 1,

		TestAnswersVerdict: testAnswersVerdict,
		InstrumentsBalanceVerdict: instrumentsBalanceVerdict,
		TradingStrategyVerdict: tradingStrategyVerdict,
		FinancialCushionVerdict: financialCushionVerdict,
		TotalVerdict: totalVerdict,

		InvestProfile: *user.InvestProfile,
		InvestProfileID: user.InvestProfileID,
	}

	db.GetDB().Create(&analytics)
	user.AnalyticsID = analytics.ID
	db.GetDB().Model(&user).Updates(user)

	handleOK(c, analytics)
}
