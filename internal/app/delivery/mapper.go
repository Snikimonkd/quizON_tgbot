package delivery

import (
	"fmt"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

func mapGame(game model.Games) string {
	res := fmt.Sprintf("Игра №: %d\nМесто: %v\nДата и время: %v\n", game.ID, game.Location, utils.PrettyTime(game.Date))
	if game.Description != nil && *game.Description != "" {
		res += fmt.Sprintf("Комментарий: %v\n", *game.Description)
	}

	return res
}

func mapMonth(in string) (time.Month, error) {
	switch in {
	case "январь":
		return time.January, nil
	case "февраль":
		return time.February, nil
	case "март":
		return time.March, nil
	case "апрель":
		return time.April, nil
	case "май":
		return time.May, nil
	case "июнь":
		return time.June, nil
	case "июль":
		return time.July, nil
	case "август":
		return time.August, nil
	case "сентябрь":
		return time.September, nil
	case "октябрь":
		return time.October, nil
	case "ноябрь":
		return time.November, nil
	case "декабрь":
		return time.December, nil
	}

	return 0, fmt.Errorf("unknown month: %v", in)
}

// /create 26 июня 2023 18:00 локация коммент
func mapArgsIntoGame(args []string) (model.Games, error) {
	if len(args) < 4 {
		return model.Games{}, fmt.Errorf("len(args) < 4")
	}
	ret := model.Games{}

	day, err := strconv.Atoi(args[0])
	if err != nil {
		return model.Games{}, fmt.Errorf("can't atoi day: %w", err)
	}

	month, err := mapMonth(args[1])
	if err != nil {
		return model.Games{}, err
	}

	year, err := strconv.Atoi(args[2])
	if err != nil {
		return model.Games{}, fmt.Errorf("can't atoi year: %w", err)
	}

	hourMinute := strings.Split(args[3], ":")
	hour, err := strconv.Atoi(hourMinute[0])
	if err != nil {
		return model.Games{}, fmt.Errorf("can't atoi hours: %w", err)
	}

	minutes, err := strconv.Atoi(hourMinute[1])
	if err != nil {
		return model.Games{}, fmt.Errorf("can't atoi minutes: %w", err)
	}

	date := time.Date(year, month, day, hour, minutes, 0, 0, utils.LocMsk)
	if err != nil {
		return model.Games{}, fmt.Errorf("can't parse time: %w", err)
	}
	ret.Date = date

	if args[4] == "" {
		return model.Games{}, fmt.Errorf("empty location")
	}
	ret.Location = args[4]

	if len(args) > 4 {
		ret.Description = lo.Reduce(args[5:], func(r *string, t string, i int) *string {
			if i == len(args[5:])-1 {
				*r += t
			} else {
				*r += t + " "
			}
			return r
		}, lo.ToPtr(""))
	}

	return ret, nil
}
