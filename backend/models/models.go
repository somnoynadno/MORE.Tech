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
	InstrumentRateChanges []*InstrumentRateChange
	News                  []*News
	Instruments           []*Instrument
	Advices               []*Advice
}

type Instrument struct {
	BaseModelCompact
	Name             string
	ImageURL         string
	InstrumentTypeID uint
	InstrumentType   InstrumentType
	AdditionalInfo   string
	Description      string
	Sector           string
	Legend           string
	BasePrice        int
	BaseAmount       int
	GameWeekID       uint
	GameWeek         GameWeek
	Users            []*User `gorm:"many2many:user_instruments;"`
}

type InstrumentType struct {
	BaseModelCompact
	Name        string
	Description string
}

type InstrumentRateChange struct {
	BaseModelCompact
	InstrumentID      uint
	Instrument        Instrument
	GameWeekID        uint
	GameWeek          GameWeek
	PriceChange       int
	AdditionalPayment *int
	PaymentName       *string
}

type News struct {
	BaseModelCompact
	Text       string
	GameWeekID uint
	GameWeek   GameWeek
}

type Advice struct {
	BaseModelCompact
	Title      string
	Text       string
	GameWeekID uint
	GameWeek   GameWeek
}

type UserInstrument struct {
	BaseModelCompact
	InstrumentID uint
	Instrument   Instrument
	UserID       uint
	User         User
}

type User struct {
	BaseModelCompact
	Name            string
	Balance         int
	Instruments     []*Instrument `gorm:"many2many:user_instruments;"`
	GameWeekID      uint
	GameWeek        GameWeek
	InvestProfileID uint
	InvestProfile   InvestProfile
}

type TestQuestion struct {
	BaseModelCompact
	ImageURL    string
	Name        string
	Text        string
	TestAnswers []*TestAnswer
}

type TestAnswer struct {
	BaseModelCompact
	Name           string
	Score          int
	TestQuestionID uint
	TestQuestion   TestQuestion
}

type InvestProfile struct {
	BaseModelCompact
	Name        string
	Description string
	MinScore    int
}
