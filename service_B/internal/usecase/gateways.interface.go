package usecase

import "context"

type TemperatureGateway interface {
	GetTemperatureByZipCode(ctx context.Context, location string) (*float64, *string, error)
}
