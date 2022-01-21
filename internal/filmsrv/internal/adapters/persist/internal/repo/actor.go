package repo

import (
	"context"
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/common/sqlhelp"
	"time"
)

type ActorEntity struct {
	SyntheticTotalCount int32     `db:"total_count" json:",omitempty"`
	ID                  int32     `db:"id"`
	FirstName           string    `db:"first_name"`
	LastName            string    `db:"last_name"`
	LastUpdate          time.Time `db:"last_update"`
}

func (r DBRepo) FetchActorByID(ctx context.Context, id int32) (*ActorEntity, error) {
	/*language=SQL*/ query := "SELECT * FROM films.actors WHERE id = $1 LIMIT 1"
	var actor ActorEntity
	err := r.db.GetContext(ctx, &actor, query, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching actor by id=%d: %w", id, err)
	}
	return &actor, nil
}

func (r DBRepo) FetchActorsByFilmID(ctx context.Context, id int32) (actors []*ActorEntity, err error) {
	/*language=SQL*/ query := `
SELECT a.* FROM films.actors a
LEFT JOIN films.film_actor_links fal on a.id = fal.actor_id
WHERE fal.film_id = :filmId
`
	queryArgs := map[string]interface{}{"filmId": id}
	nstmt, err := sqlhelp.PrepareNamedStmt(ctx, r.db, "FetchFilms", query, queryArgs)
	if err != nil {
		return nil, fmt.Errorf("error fetching actors by film id=%d: %w", id, err)
	}
	defer nstmt.Close()
	err = nstmt.Select(&actors, queryArgs)
	if err != nil {
		return nil, fmt.Errorf("error selecting actors by film id=%d: %w", id, err)
	}
	return actors, nil
}
