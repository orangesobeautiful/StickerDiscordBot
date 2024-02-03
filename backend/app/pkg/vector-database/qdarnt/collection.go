package qdarnt

import (
	"context"

	vectordatabase "backend/app/pkg/vector-database"

	pb "github.com/qdrant/go-client/qdrant"
	"golang.org/x/xerrors"
)

func (q *qdrant) CreateCollectionIfNotExist(
	ctx context.Context, dim uint, distance vectordatabase.DistanceType,
) (err error) {
	exist, err := q.isCollectionExist(ctx, dim, distance)
	if err != nil {
		return xerrors.Errorf("check collection exist: %w", err)
	}
	if exist {
		return nil
	}

	err = q.createCollection(ctx, dim, distance)
	if err != nil {
		return xerrors.Errorf("create collection: %w", err)
	}

	return nil
}

func (q *qdrant) isCollectionExist(
	ctx context.Context, dim uint, distance vectordatabase.DistanceType,
) (exist bool, err error) {
	nameExist, err := q.isCollectionNameExist(ctx)
	if err != nil {
		return false, xerrors.Errorf("check collection name exist: %w", err)
	}
	if !nameExist {
		return false, nil
	}

	err = q.checkExistedCollection(ctx, dim, distance)
	if err != nil {
		return false, xerrors.Errorf("check existed collection: %w", err)
	}

	return true, nil
}

func (q *qdrant) isCollectionNameExist(ctx context.Context) (exist bool, err error) {
	collectionListResp, err := q.collectionsClient.List(ctx, &pb.ListCollectionsRequest{})
	if err != nil {
		return false, xerrors.Errorf("list collections: %w", err)
	}

	for _, collection := range collectionListResp.GetCollections() {
		if collection.GetName() == q.collectionName {
			return true, nil
		}
	}
	return false, nil
}

func (q *qdrant) checkExistedCollection(
	ctx context.Context, dim uint, distance vectordatabase.DistanceType,
) (err error) {
	collectionInfoResp, err := q.collectionsClient.Get(ctx,
		&pb.GetCollectionInfoRequest{
			CollectionName: q.collectionName,
		},
	)
	if err != nil {
		return xerrors.Errorf("get collection info: %w", err)
	}

	vectorConfigParams := collectionInfoResp.
		GetResult().
		GetConfig().
		GetParams().
		GetVectorsConfig().
		GetParams()
	if vectorConfigParams.Size != uint64(dim) {
		return xerrors.Errorf("collection dim not match, expect %d, got %d", dim, vectorConfigParams.Size)
	}
	if distanceConvert(distance) != vectorConfigParams.Distance {
		return xerrors.Errorf("collection distance not match, expect %s, got %s",
			distance, vectorConfigParams.Distance)
	}

	return nil
}

func (q *qdrant) createCollection(ctx context.Context, dim uint, distance vectordatabase.DistanceType) (err error) {
	_, err = q.collectionsClient.Create(ctx, &pb.CreateCollection{
		CollectionName: q.collectionName,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     uint64(dim),
					Distance: distanceConvert(distance),
				},
			},
		},
	})
	if err != nil {
		return xerrors.Errorf("create collection: %w", err)
	}

	return nil
}
