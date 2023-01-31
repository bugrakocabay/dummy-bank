package util

const (
	USD = "USD"
	CAD = "CAD"
	TRY = "TRY"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, CAD, TRY:
		return true
	default:
		return false
	}
}
