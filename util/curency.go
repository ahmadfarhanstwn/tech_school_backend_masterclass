package util

const (
	USD = "USD"
	IDR = "IDR"
	EUR = "EUR"
)

func IsValidCurency(currency string) bool {
	switch currency {
	case USD, IDR, EUR:
		return true
	}
	return false
}