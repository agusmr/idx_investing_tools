package tools

import (
	"fmt"
	"strconv"
	"time"
)

func YearToDate(year int) (time.Time, error) {
	yearString := strconv.Itoa(year)
	dateString := fmt.Sprintf("30 Sep %s 00:00 GMT+7", yearString[2:])
	t, err := time.Parse(time.RFC822, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
