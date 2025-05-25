package usecase

import (
	"math"

	"github.com/FlyKarlik/gofemart/pkg/generics"
)

func convertMoneyValueToFloat64(v *int64) *float64 {
	if v != nil {
		value := float64(*v) / 100.0
		return &value
	}
	return generics.Pointer(0.0)
}

func convertMoneyValueToInt64(v *float64) *int64 {
	if v != nil {
		value := int64(math.Round(*v * 100.0))
		return &value
	}
	return generics.Pointer[int64](0)
}
