package usecase

import (
	"bytes"
	"context"
	"fmt"
	"quizon_bot/internal/generated/postgres/public/model"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
	"github.com/olekukonko/tablewriter"
)

type TableRepostiory interface {
	Registrations(ctx context.Context) ([]model.Registrations, error)
}

func (u usecase) Table(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	arr, err := u.tableRepostiory.Registrations(ctx)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.From.ID, DefaultErrorMessage), err
	}

	table := [][]string{}

	for i, v := range arr {
		row := []string{}

		row = append(row, fmt.Sprintf("%d", i+1))
		row = append(row, v.TeamName)
		row = append(row, fmt.Sprintf("%d", v.TeamID))
		row = append(row, v.CaptainName)
		row = append(row, v.TgContact)
		row = append(row, v.Phone)
		row = append(row, v.GroupName)
		row = append(row, v.Amount)

		table = append(table, row)
	}

	buffer := new(bytes.Buffer)
	writer := tablewriter.NewWriter(buffer)
	writer.SetHeader([]string{"№", "Название", "ID", "Имя капитана", "Телега", "Телефон", "Группа", "Количество человек"})

	for i := 0; i < len(table); i++ {
		writer.Append(table[i])
	}

	writer.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	writer.SetCenterSeparator("|")
	writer.Render()

	msg := tgbotapi.NewMessage(update.Message.From.ID, "<pre>"+buffer.String()+"</pre>")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ParseMode = "HTML"
	return msg, nil
}
