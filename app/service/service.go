package service

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const layout string = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", errors.New("repeat has empty value")
	}

	dayMatch, _ := regexp.MatchString(`d \d{1,3}`, repeat)
	yearMatch, _ := regexp.MatchString(`y`, repeat)

	if dayMatch {
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d ")) // отсекаем "d ", чтобы осталось только число
		if err != nil {
			return "", err
		}

		if days > 400 {
			return "", errors.New("максимальное количество дней должно быть меньше 400")
		}

		parsedDate, err := time.Parse(layout, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(0, 0, days)

		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, days)
		}

		return newDate.Format(layout), nil

	} else if yearMatch {
		parsedDate, err := time.Parse(layout, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(1, 0, 0)

		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}

		return newDate.Format(layout), nil

	}

	return "", errors.New("repeat имеет неправильный формат")

}
