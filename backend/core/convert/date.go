package convert

import (
	"fmt"
	"math"
	"strconv"
)

var ordinalDictionary = map[int]string{
	0: "th",
	1: "st",
	2: "nd",
	3: "rd",
	4: "th",
	5: "th",
	6: "th",
	7: "th",
	8: "th",
	9: "th",
}

var shortToLongDays = map[string]string{
	"mon": "Monday",
	"tue": "Tuesday",
	"wed": "Wednesday",
	"thu": "Thursday",
	"fri": "Friday",
	"sat": "Saturday",
	"sun": "Sunday",
}

func Ordinalize(n int) string {
	n = int(math.Abs(float64(n)))

	if ((n % 100) >= 11) && ((n % 100) <= 13) {
		return strconv.Itoa(n) + "th"
	}

	return strconv.Itoa(n) + ordinalDictionary[n%10]
}

func ShortToLongWeekdays(shortDay string) (string, error) {
	if val, ok := shortToLongDays[shortDay]; ok {
		return val, nil
	}
	return "", fmt.Errorf("invalid shortday string >%s<", shortDay)
}
