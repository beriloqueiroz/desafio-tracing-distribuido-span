package usecase

import "context"

type LocationGateway interface {
	GetLocationByZipCode(ctx context.Context, zipCode string) (string, error)
}

type TemperatureGateway interface {
	GetTemperatureByLocation(ctx context.Context, location string) (float64, error)
}
