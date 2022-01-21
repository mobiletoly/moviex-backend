package usecase

import (
	"fmt"
	"strconv"
)

func idToInt32(id string) (int32, error) {
	idval, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid id format: id=%s", id)
	}
	return int32(idval), nil
}
