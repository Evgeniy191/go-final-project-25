package nextdate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is required")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("incorrect rule format d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("incorrect number of days (1-400)")
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("incorrect rule format w")
		}

		// Парсим дни недели
		daysStr := strings.Split(parts[1], ",")
		var weekdays [8]bool // индексы 1-7 для дней недели

		for _, dayStr := range daysStr {
			day, err := strconv.Atoi(strings.TrimSpace(dayStr))
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("incorrect weekday (1-7)")
			}
			weekdays[day] = true
		}

		// Ищем следующий подходящий день
		for i := 0; i < 400; i++ {
			date = date.AddDate(0, 0, 1)
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7 // воскресенье = 7
			}
			if weekdays[wd] && afterNow(date, now) {
				break
			}
		}

	case "m":
		if len(parts) != 2 {
			return "", errors.New("incorrect rule format m")
		}

		// Парсим дни месяца
		daysStr := strings.Split(parts[1], ",")
		var monthdays [32]bool // индексы 1-31
		hasLastDay := false

		for _, dayStr := range daysStr {
			dayStr = strings.TrimSpace(dayStr)
			if dayStr == "-1" {
				hasLastDay = true
				continue
			}
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 31 {
				return "", errors.New("incorrect monthday (1-31 or -1)")
			}
			monthdays[day] = true
		}

		// Ищем следующий подходящий день
		for i := 0; i < 400; i++ {
			date = date.AddDate(0, 0, 1)
			day := date.Day()

			// Проверяем: это последний день месяца?
			nextDay := date.AddDate(0, 0, 1)
			isLastDay := nextDay.Month() != date.Month()

			if (monthdays[day] || (hasLastDay && isLastDay)) && afterNow(date, now) {
				break
			}
		}

	default:
		return "", errors.New("unsupported rule:" + rule)
	}

	return date.Format(DateFormat), nil
}
