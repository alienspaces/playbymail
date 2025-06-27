package server

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

type MapperFn[Rec any, Data any] func(l logger.Logger, rec *Rec) (Data, error)

func Paginate[Rec any, Data any](l logger.Logger, recs []*Rec, mapper MapperFn[Rec, Data], pageSize int) ([]*Data, error) {
	res := make([]*Data, 0) // empty slice is needed for json.Marshal to return [] instead of null

	for _, rec := range recs {
		if pageSize == 0 {
			break
		}

		responseData, err := mapper(l, rec)
		if err != nil {
			l.Warn("(core) failed to paginate record set to response data")
			return nil, err
		}
		res = append(res, &responseData)

		pageSize--
	}

	return res, nil
}
