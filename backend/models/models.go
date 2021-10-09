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
	InstrumentType   *InstrumentType `json:",omitempty"`
	AdditionalInfo   string
	Description      string
	Sector           string
	Legend           string
	BasePrice        int
	BaseAmount       int
	GameWeekID       uint
	GameWeek         *GameWeek `json:",omitempty"`
	Users            []*User   `gorm:"many2many:user_instruments;"`
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
	PriceChange       int
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
	InstrumentID uint
	Instrument   *Instrument `json:",omitempty"`
	UserID       uint
	User         *User `json:",omitempty"`
}

type User struct {
	BaseModelCompact
	Name            string
	Balance         int
	Instruments     []*Instrument `gorm:"many2many:user_instruments;"`
	GameWeekID      uint
	GameWeek        *GameWeek `json:",omitempty"`
	InvestProfileID uint
	InvestProfile   *InvestProfile `json:",omitempty"`
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
	TestQuestion   *TestQuestion `json:",omitempty"`
}

type InvestProfile struct {
	BaseModelCompact
	Name        string
	Description string
	MinScore    int
}
