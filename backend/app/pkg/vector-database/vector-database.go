package vectordatabase

import (
	"context"
	"errors"

	searchfilter "backend/app/pkg/vector-database/search-filter"
)

type DistanceType int

const (
	DistanceTypeUnknow DistanceType = iota
	DistanceTypeCosine
)

func (d DistanceType) String() string {
	switch d {
	case DistanceTypeCosine:
		return "cosine"
	default:
		return "unknow"
	}
}

type VectorDatabase interface {
	CreateCollectionIfNotExist(ctx context.Context, dim uint, distance DistanceType) (err error)

	Upsert(ctx context.Context, request *UpsertRequest) (err error)

	Search(ctx context.Context, request *SearchRequest) (resp SearchResponse, err error)

	Delete(ctx context.Context, request *DeleteRequest) (err error)
}

type UpsertRequestVector struct {
	ID uint `validate:"required"`

	Data []float32 `validate:"required"`

	Metadata map[string]any
}

type UpsertRequest struct {
	Vectors []UpsertRequestVector `validate:"required,dive,min=1"`
}

type SearchRequest struct {
	Vector []float32 `validate:"required"`
	TopK   uint      `validate:"required"`

	Filter *searchfilter.FilterInstance
}

var _ error = (*SearchFilterConvertError)(nil)

type SearchFilterConvertError struct {
	message string
}

func (e *SearchFilterConvertError) Error() string {
	return e.message
}

func NewSearchFilterConvertError(message string) error {
	return &SearchFilterConvertError{
		message: message,
	}
}

func IsSearchFilterConvertError(err error) bool {
	if err == nil {
		return false
	}

	var e *SearchFilterConvertError
	return errors.As(err, &e)
}

type SearchResponseData struct {
	ID    uint64
	Score float32
}

type SearchResponse struct {
	Data []SearchResponseData
}

type DeleteRequest struct {
	IDs []uint `validate:"required,min=1"`
}
