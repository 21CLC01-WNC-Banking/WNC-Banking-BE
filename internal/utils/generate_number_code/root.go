package generate_number_code

import (
	"fmt"
	"math"
	"math/rand"
)

func GenerateRandomNumber(digits int) string {
	if digits <= 0 {
		return ""
	}

	maxValue := int64(math.Pow10(digits))

	randomNumber := rand.Int63n(maxValue)

	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, randomNumber)
}
