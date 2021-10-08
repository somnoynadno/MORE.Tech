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
	InstrumentTypeID uint
	InstrumentType   InstrumentType
	AdditionalInfo   string
	Description      string
	Legend           string
	BasePrice        int
	BaseAmount       int
	GameWeekID       uint
	GameWeek         GameWeek
}

type InstrumentType struct {
	BaseModelCompact
	Name        string
	Description string
}

type InstrumentRateChange struct {
	BaseModelCompact
	InstrumentID uint
	Instrument   Instrument
	GameWeekID   uint
	GameWeek     GameWeek
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
	Name        string
	Balance     int
	Instruments []*Instrument
	GameWeekID  uint
	GameWeek    GameWeek
}
