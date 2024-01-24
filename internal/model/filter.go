package model

type Filter struct {
	AgeMin  uint
	AgeMax  uint
	Gender  string
	Country string
	Offset  uint64
	Limit   uint64
}
