package common

import (
	"time"
)

type PriceType struct {
	PriceTypeId  int
	TaxRateToAdd float32
	Description  string
}

type Price struct {
	CurrencyId int
	Amount     float32
	StartDate  time.Time
	TypeId     int
}

type Measure struct {
	Value float32
	Unit  string
}

type Packing struct {
	Length      *Measure
	Width       *Measure
	Height      *Measure
	GrossWeight *Measure
	NetWeight   *Measure
}

type MoQ struct {
	Quantity  int
	StartDate time.Time
}
