package repo

import (
	"context"
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/common/sqlhelp"
	"time"
)

type FilmEntity struct {
	SyntheticTotalCount int32     `db:"total_count" json:",omitempty"`
	ID                  int32     `db:"id"`
	Title               string    `db:"title"`
	Description         *string   `db:"description"`
	ReleaseYear         int32     `db:"release_year"`
	LanguageID          int32     `db:"language_id"`
	Length              int32     `db:"length"`
	Rating              string    `db:"rating"`
	LastUpdate          time.Time `db:"last_update"`
	SpecialFeatures     *string   `db:"special_features"`
}

func (r DBRepo) FetchFilms(
	ctx context.Context, page sqlhelp.DBPageReq,
) (films []*FilmEntity, total int32, error error) {
	/*language=SQL*/ query := `
SELECT count(*) OVER() as total_count, f.* FROM films.films f
ORDER BY title
LIMIT :limit OFFSET :offset`
	queryArgs := page.ToMap()
	nstmt, err := sqlhelp.PrepareNamedStmt(ctx, r.db, "FetchFilms", query, queryArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching films: %w", err)
	}
	defer nstmt.Close()
	err = nstmt.Select(&films, queryArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("error selecting films: %w", err)
	}
	var totalCount int32
	if len(films) > 0 {
		totalCount = films[0].SyntheticTotalCount
	}
	return films, totalCount, err
}
