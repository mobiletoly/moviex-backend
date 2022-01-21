package sqlhelp

type DBPageReq struct {
	Limit  int32
	Offset int32
}

func NewDBPageReq(numRecords, startRecord int32) DBPageReq {
	if numRecords == 0 {
		numRecords = 999999999
	}
	return DBPageReq{
		Limit:  numRecords,
		Offset: startRecord,
	}
}

func (r *DBPageReq) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"limit":  r.Limit,
		"offset": r.Offset,
	}
}
