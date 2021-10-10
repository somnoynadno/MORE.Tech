package server

import (
	"MORE.Tech/backend/db"
	"MORE.Tech/backend/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func RenderResult(c *gin.Context) {
	var user models.User
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	log.Error(id, errors.New("bad parameters"))

	if err != nil {
		handleBadRequest(c, errors.New("bad path parameters"))
		return
	}

	err = db.GetDB().Preload("InvestProfile").Preload("GameWeek").Preload("TestAnswers").
		Preload("UserInstruments").Preload("UserInstruments.Instrument").
		First(&user, id).Error
	if err != nil {
		handleBadRequest(c, errors.New("bad path parameters"))
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

	if user.InvestProfile == nil {
		user.InvestProfile = &models.InvestProfile{
			Name: "Сбалансированный",
			Description: "У вас получается найти баланс между доходностью и риском: ваш портфель сможет принести вам хорошую прибыль и будет защищен от негативных обстоятельств.",
		}
	}

	c.HTML(http.StatusOK, "view.html", gin.H{
		"income": income,
		"percent":  fmt.Sprintf("%.2f", float64(income / user.BaseBalance - 1) * 100),
		"l1": instrumentsBalanceVerdict,
		"l2": financialCushionVerdict,
		"l3": tradingStrategyVerdict,
		"l4": testAnswersVerdict,
		"verdict": totalVerdict,
		"profile": user.InvestProfile.Name,
		"description": user.InvestProfile.Description,
	})
}