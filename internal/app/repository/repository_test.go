package repository

import (
	"context"
	"quizon_bot/internal/generated/postgres/public/model"
	"quizon_bot/internal/pkg/testsupport"
	"quizon_bot/internal/pkg/timesupport"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

func Test_repository_Games(t *testing.T) {
	db := testsupport.ConnectToTestPostgres()
	ctx := context.Background()

	t.Cleanup(func() {
		testsupport.TruncateGames(t, db)
	})

	tests := []struct {
		name    string
		from    time.Time
		prep    func() []model.Games
		wantErr bool
	}{
		{
			name:    "1. Successful test, no rows.",
			from:    time.Date(2025, 3, 4, 5, 6, 7, 8, timesupport.LocMsk),
			prep:    func() []model.Games { return []model.Games{} },
			wantErr: false,
		},
		{
			name: "2. Successful test, multiple rows.",
			from: time.Date(2023, 3, 4, 5, 6, 7, 0, time.UTC),
			prep: func() []model.Games {
				in := []model.Games{
					// old
					{
						ID:          testsupport.RandInt64(),
						CreatedAt:   timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, time.UTC)),
						UdpatedAt:   timesupport.Pretty(time.Date(2023, 2, 3, 4, 5, 6, 7, time.UTC)),
						Location:    "kek",
						Date:        timesupport.Pretty(time.Date(2023, 3, 4, 5, 6, 6, 0, time.UTC)),
						Description: lo.ToPtr("lol"),
					},
					// ok
					{
						ID:          testsupport.RandInt64(),
						CreatedAt:   timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, time.UTC)),
						UdpatedAt:   timesupport.Pretty(time.Date(2023, 2, 3, 4, 5, 6, 7, time.UTC)),
						Location:    "kek",
						Date:        timesupport.Pretty(time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)),
						Description: lo.ToPtr("lol"),
					},
					{
						ID:          testsupport.RandInt64(),
						CreatedAt:   timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, time.UTC)),
						UdpatedAt:   timesupport.Pretty(time.Date(2023, 2, 3, 4, 5, 6, 7, time.UTC)),
						Location:    "kek",
						Date:        timesupport.Pretty(time.Date(2025, 3, 4, 5, 6, 8, 0, time.UTC)),
						Description: lo.ToPtr("lol"),
					},
					{
						ID:          testsupport.RandInt64(),
						CreatedAt:   timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, time.UTC)),
						UdpatedAt:   timesupport.Pretty(time.Date(2023, 2, 3, 4, 5, 6, 7, time.UTC)),
						Location:    "kek",
						Date:        timesupport.Pretty(time.Date(2025, 3, 4, 5, 6, 9, 0, time.UTC)),
						Description: lo.ToPtr("lol"),
					},
				}
				for i := 0; i < len(in); i++ {
					testsupport.InsertIntoGames(t, db, in[i])
				}
				return in[1:]
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repository{db: db}

			want := tt.prep()

			got, err := r.Games(ctx, tt.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.Games() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got = lo.Map(got, func(item model.Games, index int) model.Games {
				item.CreatedAt = timesupport.Pretty(item.CreatedAt)
				item.UdpatedAt = timesupport.Pretty(item.UdpatedAt)
				item.Date = timesupport.Pretty(item.Date)
				return item
			})

			slices.SortFunc(want, func(a, b model.Games) bool { return a.ID < b.ID })
			slices.SortFunc(got, func(a, b model.Games) bool { return a.ID < b.ID })

			if !cmp.Equal(got, want) {
				t.Errorf("repository.Games() got != want, diff:\n %v", cmp.Diff(want, got))
			}
		})
	}
}

func Test_repository_Create(t *testing.T) {
	db := testsupport.ConnectToTestPostgres()
	ctx := context.Background()

	t.Cleanup(func() {
		testsupport.TruncateGames(t, db)
	})

	tests := []struct {
		name    string
		want    model.Games
		wantErr bool
	}{
		{
			name: "1. Successful test.",
			want: model.Games{
				CreatedAt:   timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, timesupport.LocMsk)),
				UdpatedAt:   timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, timesupport.LocMsk)),
				Location:    "lol",
				Date:        timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, timesupport.LocMsk)),
				Description: lo.ToPtr("kek"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repository{db: db}
			var err error
			tt.want.ID, err = r.Create(ctx, tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := testsupport.SelectGames(t, db, tt.want.ID)
			if !cmp.Equal(tt.want, got) {
				t.Errorf("want != got, diff:\n%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func Test_repository_UpsertAdmin(t *testing.T) {
	db := testsupport.ConnectToTestPostgres()
	ctx := context.Background()

	tests := []struct {
		name    string
		want    model.Admins
		wantErr bool
	}{
		{
			name: "1. Successful test.",
			want: model.Admins{
				ID:        1,
				DateUntil: timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, timesupport.LocMsk)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repository{db: db}
			err := r.UpsertAdmin(ctx, tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.UpsertAdmin() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := testsupport.SelectAdmins(t, db, tt.want.ID)
			if !cmp.Equal(tt.want, got) {
				t.Errorf("got != want, diff:\n%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func Test_repository_CheckAuth(t *testing.T) {
	db := testsupport.ConnectToTestPostgres()
	ctx := context.Background()

	t.Cleanup(func() {
		testsupport.TruncateAdmins(t, db)
	})

	tests := []struct {
		name    string
		prep    func() model.Admins
		wantErr bool
	}{
		{
			name: "1. Successful test.",
			prep: func() model.Admins {
				in := model.Admins{
					ID:        testsupport.RandInt64(),
					DateUntil: timesupport.Pretty(time.Date(2023, 1, 2, 3, 4, 5, 6, timesupport.LocMsk)),
				}

				testsupport.InsertIntoAdmins(t, db, in)
				return in
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repository{db: db}
			want := tt.prep()
			got, err := r.CheckAuth(ctx, want.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.CheckAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, want) {
				t.Errorf("repository.CheckAuth() = %v, want %v", got, want)
			}
		})
	}
}
