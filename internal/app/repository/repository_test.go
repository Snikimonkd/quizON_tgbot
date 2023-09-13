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
		_, err := db.Exec(ctx, "TRUNCATE games CASCADE;")
		if err != nil {
			t.Errorf("can't clean db: %v", err)
		}
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

// func Test_repository_Create(t *testing.T) {
// 	db := testsupport.ConnectToTestPostgres()
// 	ctx := context.Background()
//
// 	t.Cleanup(func() {
// 		_, err := db.Exec(ctx, "TRUNCATE games CASCADE;")
// 		if err != nil {
// 			t.Errorf("can't clean db: %v", err)
// 		}
// 	})
//
// 	tests := []struct {
// 		name    string
// 		req     model.Games
// 		check   func(expect model.Games)
// 		wantErr bool
// 	}{
// 		{
// 			name: "1. Successful test.",
// 			req: model.Games{
// 				ID:          1,
// 				CreatedAt:   time.Now(),
// 				UdpatedAt:   time.Now(),
// 				Location:    "lol",
// 				Date:        time.Now(),
// 				Description: lo.ToPtr("kek"),
// 			},
// 			check: func(expect model.Games) {
// 				stmt := table.Games.Select(
// 					table.Games.AllColumns,
// 				).WHERE(
// 					table.Games.ID.EQ(expect.ID),
// 				)
//
//
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := repository{db: db}
// 			err := r.Create(ctx, tt.req)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("repository.Create() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			tt.check(tt.req)
// 		})
// 	}
// }
