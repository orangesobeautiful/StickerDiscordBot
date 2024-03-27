package repository

import (
	"context"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/pkg/hserr"
)

func BaseList[resultT any](
	ctx context.Context, queryFilter ent.QueryPaging[resultT], limit, offset int,
) (result domain.ListResult[*resultT], err error) {
	total, err := queryFilter.PagingClone().PagingCount(ctx)
	if err != nil {
		return result, hserr.NewInternalError(err, "query count")
	}

	ragReferencePools, err := queryFilter.PagingClone().
		PagingLimit(limit).
		PagingOffset(offset).
		PagingAll(ctx)
	if err != nil {
		return domain.ListResult[*resultT]{}, hserr.NewInternalError(err, "paging list all")
	}

	result = domain.NewListResult(total, ragReferencePools)
	return result, nil
}
