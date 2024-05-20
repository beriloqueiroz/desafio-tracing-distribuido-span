package entity

import (
	"errors"
	"regexp"
)

type ZipCode struct {
	Value string
}

func NewZipCode(code string) (*ZipCode, error) {

	if len(code) != 8 {
		return nil, errors.New("invalid zipcode")
	}

	var re = regexp.MustCompile(`^[0-9]+$`)

	if !re.MatchString(code) {
		return nil, errors.New("invalid zipcode")
	}

	return &ZipCode{
		Value: code,
	}, nil
}

type TemperatureLocation struct {
	City    string
	ZipCode *ZipCode
	TempC   float64
	TempF   float64
	TempK   float64
}

func NewTemperatureLocation(zipCode *ZipCode, tempC float64, city string) (*TemperatureLocation, error) {
	return &TemperatureLocation{
		City:    city,
		ZipCode: zipCode,
		TempC:   tempC,
		TempF:   tempC*1.8 + 32,
		TempK:   tempC + 273,
	}, nil

}
