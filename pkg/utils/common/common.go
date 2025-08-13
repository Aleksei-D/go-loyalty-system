package common

import "strconv"

type LoginKey string

func (c LoginKey) String() string {
	return string(c)
}

func CheckLuhnAlgorithm(orderNumber string) bool {
	sum := 0
	double := false

	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil {
			return false
		}

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	return sum%10 == 0
}
