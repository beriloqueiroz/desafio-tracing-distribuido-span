package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateATemperatureLocation(t *testing.T) {
	zip, _ := NewZipCode("60541706")
	tl, err := NewTemperatureLocation(zip, 10.5, "Fortaleza")
	assert.Nil(t, err)
	assert.Equal(t, "60541706", tl.ZipCode.Value)
	assert.Equal(t, "Fortaleza", tl.City)
	assert.InDelta(t, 10.5, tl.TempC, 0.000001)
	assert.InDelta(t, 50.9, tl.TempF, 0.000001)
	assert.InDelta(t, 283.5, tl.TempK, 0.000001)
}

func TestCreateATemperatureLocationWhenInvalidZipCode(t *testing.T) {
	zip, err := NewZipCode("6054176")
	assert.NotNil(t, err)
	assert.Nil(t, zip)
	assert.Equal(t, "invalid zipcode", err.Error())

	zip, err = NewZipCode("1236547a")
	assert.NotNil(t, err)
	assert.Nil(t, zip)
	assert.Equal(t, "invalid zipcode", err.Error())
}
