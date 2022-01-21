package model

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	recIndPrefix  = "#recind:"
	pageIndPrefix = "#pageind:"
	tknPrefix     = "#tkn:"
)

type PageReqParams struct {
	afterCursor *string
	count       int32
}

func NewPageReqParams(afterCursor *string, count int32) PageReqParams {
	return PageReqParams{
		afterCursor: afterCursor,
		count:       count,
	}
}

func (r PageReqParams) Count() int32 {
	return r.count
}

func (r PageReqParams) HasCursor() bool {
	return r.afterCursor != nil && *r.afterCursor != ""
}

func (r PageReqParams) AfterOffset() int32 {
	if r.afterCursor == nil {
		return 0
	}
	afterCursor := *r.afterCursor
	if strings.HasPrefix(afterCursor, recIndPrefix) {
		offset, err := strconv.ParseInt((afterCursor)[len(recIndPrefix):], 10, 32)
		if err != nil {
			panic(fmt.Errorf("error parsing pagination cursor value: %s", *r.afterCursor))
		}
		return int32(offset) + 1
	} else if strings.HasPrefix(afterCursor, pageIndPrefix) {
		offset, err := strconv.ParseInt((afterCursor)[len(pageIndPrefix):], 10, 32)
		if err != nil {
			panic(fmt.Errorf("error parsing pagination cursor value: %s", *r.afterCursor))
		}
		return int32(offset)
	}
	panic(fmt.Errorf("invalid pagination cursor format: %s", *r.afterCursor))
}

func (r PageReqParams) AfterToken() *string {
	if r.afterCursor == nil {
		return nil
	}
	afterCursor := *r.afterCursor
	if strings.HasPrefix(afterCursor, tknPrefix) {
		token := (*r.afterCursor)[len(tknPrefix):]
		return &token
	}
	panic(fmt.Errorf("invalid pagination cursor format: %s", *r.afterCursor))
}

func NewPageInfoWithRecordIndexAndTotal(startRecordIndex int32, responseRecords int32, totalRecords int32) *PageInfo {
	hasNextPage := startRecordIndex+responseRecords < totalRecords
	var endCursor *string
	if hasNextPage {
		cursor := fmt.Sprintf("%s%d", recIndPrefix, startRecordIndex+responseRecords-1)
		endCursor = &cursor
	}
	return &PageInfo{hasNextPage, endCursor}
}
