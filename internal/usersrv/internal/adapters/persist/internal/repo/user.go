package repo

import (
	"context"
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/common/sqlhelp"
)

type UserEntity struct {
	SyntheticTotalCount int32  `db:"total_count" json:",omitempty"`
	ID                  int32  `db:"id"`
	Email               string `db:"email"`
	Password            string `db:"password"`
}

func (r DBRepo) FetchUsers(
	ctx context.Context, page sqlhelp.DBPageReq,
) (users []*UserEntity, total int32, error error) {
	/*language=SQL*/ query := `
SELECT count(*) OVER() as total_count, u.* FROM users.users u
ORDER BY email
LIMIT :limit OFFSET :offset`
	queryArgs := page.ToMap()
	nstmt, err := sqlhelp.PrepareNamedStmt(ctx, r.db, "FetchUsers", query, queryArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching users: %w", err)
	}
	defer nstmt.Close()
	err = nstmt.Select(&users, queryArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("error selecting users: %w", err)
	}
	var totalCount int32
	if len(users) > 0 {
		totalCount = users[0].SyntheticTotalCount
	}
	return users, totalCount, err
}
