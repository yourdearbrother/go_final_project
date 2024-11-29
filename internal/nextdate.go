package nextdate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	nowDateStr := now.Format("20060102")
	now, err := time.Parse("20060102", nowDateStr)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("правило повторения отсутствует")
	}

	rep := strings.Split(repeat, " ")

	if len(rep) < 1 {
		return "", errors.New("неправильное правило повторения")
	}

	timBase, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if rep[0] == "y" {

		origDay := timBase.Day()
		origMonth := timBase.Month()

		for {

			timBase = timBase.AddDate(1, 0, 0)

			if timBase.Day() == origDay && timBase.Month() == origMonth {

				if timBase.After(now) {
					break
				}
			} else {

				timBase = time.Date(timBase.Year(), time.March, 1, 0, 0, 0, 0, timBase.Location())
				if timBase.After(now) {
					break
				}
			}
		}
		return timBase.Format("20060102"), nil
	}

	if rep[0] == "d" {
		if len(rep) < 2 {
			return "", errors.New("неправильно указан режим повторения")
		}

		days, err := strconv.Atoi(rep[1])
		if err != nil {
			return "", err
		}

		if days > 400 {
			return "", errors.New("перенос события более чем на 400 дней недопустим")
		}

		for {
			timBase = timBase.AddDate(0, 0, days)
			if timBase.After(now) {
				break
			}
		}
		return timBase.Format("20060102"), nil
	}

	return "", errors.New("неправильное правило повторения")
}
