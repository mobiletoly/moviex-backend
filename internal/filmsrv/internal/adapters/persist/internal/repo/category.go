package repo

import (
	"context"
	"fmt"
	"time"
)

type CategoryEntity struct {
	SyntheticTotalCount int32     `db:"total_count" json:",omitempty"`
	ID                  int32     `db:"id"`
	Name                string    `db:"name"`
	LastUpdate          time.Time `db:"last_update"`
}

func (r DBRepo) FetchCategoryByID(ctx context.Context, id int32) (*CategoryEntity, error) {
	/*language=SQL*/ query := "SELECT * FROM films.categories WHERE id = $1 LIMIT 1"
	var cat CategoryEntity
	err := r.db.GetContext(ctx, &cat, query, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching category: %w", err)
	}
	return &cat, nil
}

func (r DBRepo) FetchCategoryByFilmID(ctx context.Context, id int32) (*CategoryEntity, error) {
	/*language=SQL*/ query := `
SELECT c.* FROM films.categories c
LEFT JOIN films.film_category_links fcl on c.id = fcl.category_id
WHERE fcl.film_id = $1 LIMIT 1
`
	var cat CategoryEntity
	err := r.db.GetContext(ctx, &cat, query, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching category: %w", err)
	}
	return &cat, nil
}
