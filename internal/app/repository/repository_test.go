package repository

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/generated/postgres/public/table"
	"quizon_bot/internal/pkg/testsupport"
	"testing"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/samber/lo"
)

func Test_repository_RegisterAvailable(t *testing.T) {
	db := testsupport.ConnectToTestPostgres()
	ctx := context.Background()

	tests := []struct {
		name    string
		prep    func()
		want    bool
		wantErr bool
	}{
		{
			name: "1. Successful test, returt true.",
			prep: func() {
				testsupport.TruncateRegistrations(t, db)
				testsupport.TruncateGames(t, db)
				stmt := table.Games.INSERT(table.Games.Reserve).VALUES(postgres.Int64(10))
				query, args := stmt.Sql()
				_, err := db.Exec(ctx, query, args...)
				if err != nil {
					t.Errorf("prep error: %v", err.Error())
				}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "2. Successful test, return false.",
			prep: func() {
				testsupport.TruncateRegistrations(t, db)
				testsupport.TruncateGames(t, db)
				testsupport.InsertRegistration(t, db, model.Registrations{
					TgContact:   "123",
					TeamID:      lo.ToPtr("1"),
					TeamName:    "123",
					CaptainName: "123",
					Phone:       "123",
					GroupName:   "123",
					Amount:      "123",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				stmt := table.Games.INSERT(table.Games.Reserve).VALUES(postgres.Int64(1))
				query, args := stmt.Sql()
				_, err := db.Exec(ctx, query, args...)
				if err != nil {
					t.Errorf("prep error: %v", err.Error())
				}
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repository{db: db}
			tt.prep()
			got, err := r.RegisterAvailable(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.RegisterAvailable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repository.RegisterAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
