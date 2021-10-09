package models

import "time"

type BaseModelCompact struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type GameWeek struct {
	BaseModelCompact
	InstrumentRateChanges []InstrumentRateChange
	News                  []News
	Instruments           []Instrument
	Advices               []Advice
	TestQuestions         []TestQuestion
}

type Instrument struct {
	BaseModelCompact
	Name             string
	ImageURL         string
	InstrumentTypeID uint
	InstrumentType   *InstrumentType `json:",omitempty"`
	AdditionalInfo   string
	Description      string
	Sector           string
	Legend           string
	Rating           int
	BasePrice        int
	BaseAmount       int
	GameWeekID       uint
	GameWeek         *GameWeek `json:",omitempty"`
}

type InstrumentType struct {
	BaseModelCompact
	Name        string
	Description string
}

type InstrumentRateChange struct {
	BaseModelCompact
	InstrumentID      uint
	Instrument        *Instrument `json:",omitempty"`
	GameWeekID        uint
	GameWeek          *GameWeek `json:",omitempty"`
	PriceChangeRate   float64
	AdditionalPayment *int
	PaymentName       *string
}

type News struct {
	BaseModelCompact
	Text       string
	GameWeekID uint
	GameWeek   *GameWeek `json:",omitempty"`
}

type Advice struct {
	BaseModelCompact
	Title      string
	Text       string
	GameWeekID uint
	GameWeek   *GameWeek `json:",omitempty"`
}

type UserInstrument struct {
	BaseModelCompact
	CurrentPrice     int
	PriceChanged     int
	PriceChangedRate float64
	InstrumentID     uint
	Instrument       *Instrument `json:",omitempty"`
	UserID           uint
	User             *User `json:",omitempty"`
}

type User struct {
	BaseModelCompact
	Name            string
	Balance         int
	BaseBalance     int
	Sells           int
	GameWeekID      uint
	GameWeek        *GameWeek `json:",omitempty"`
	InvestProfileID uint
	InvestProfile   *InvestProfile `json:",omitempty"`
	AnalyticsID     uint
	Analytics       *Analytics        `json:",omitempty"`
	UserInstruments []*UserInstrument `json:",omitempty"`
	TestAnswers     []*TestAnswer     `json:",omitempty" gorm:"many2many:user_test_answers;"`
}

type TestQuestion struct {
	BaseModelCompact
	Name        string
	Text        string
	TestAnswers []TestAnswer
	GameWeekID  uint
	GameWeek    *GameWeek `json:",omitempty"`
}

type UserTestAnswer struct {
	BaseModelCompact
	UserID       uint
	User         *User `json:",omitempty"`
	TestAnswerID uint
	TestAnswer   TestAnswer `json:",omitempty"`
}

type TestAnswer struct {
	BaseModelCompact
	Name           string
	Feedback       string
	IsCorrect      bool
	TestQuestionID uint
	TestQuestion   *TestQuestion `json:",omitempty"`
	Users          []*User       `json:",omitempty" gorm:"many2many:user_test_answers;"`
}

type InvestProfile struct {
	BaseModelCompact
	Name        string
	Description string
}

type Analytics struct {
	BaseModelCompact
	TotalIncome     int
	TotalIncomeRate float64

	InvestProfileID uint
	InvestProfile   InvestProfile

	InstrumentsBalanceVerdict string
	FinancialCushionVerdict   string
	TestAnswersVerdict        string
	TradingStrategyVerdict    string
	TotalVerdict              string
}
