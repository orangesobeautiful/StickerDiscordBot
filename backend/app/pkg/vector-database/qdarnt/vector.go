package qdarnt

import (
	"context"

	vectordatabase "backend/app/pkg/vector-database"

	pb "github.com/qdrant/go-client/qdrant"
	"golang.org/x/xerrors"
)

func newPointIDs(ids []uint) []*pb.PointId {
	pointIDs := make([]*pb.PointId, len(ids))
	for i, id := range ids {
		pointIDs[i] = newPointID(id)
	}

	return pointIDs
}

func newPointID(id uint) *pb.PointId {
	return &pb.PointId{
		PointIdOptions: &pb.PointId_Num{
			Num: uint64(id),
		},
	}
}

func (q *qdrant) Upsert(ctx context.Context, request *vectordatabase.UpsertRequest) (err error) {
	pointStructs := make([]*pb.PointStruct, 0, len(request.Vectors))

	for _, vectorData := range request.Vectors {
		pointStruct := &pb.PointStruct{
			Id: newPointID(vectorData.ID),
			Vectors: &pb.Vectors{
				VectorsOptions: &pb.Vectors_Vector{
					Vector: &pb.Vector{
						Data: vectorData.Data,
					},
				},
			},
		}
		if len(vectorData.Metadata) > 0 {
			var payload map[string]*pb.Value
			payload, err = NewPayload(vectorData.Metadata)
			if err != nil {
				return xerrors.Errorf("convert payload: %w", err)
			}

			pointStruct.Payload = payload
		}

		pointStructs = append(pointStructs, pointStruct)
	}

	wait := true
	_, err = q.pointsClient.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: q.collectionName,
		Wait:           &wait,
		Points:         pointStructs,
	})
	if err != nil {
		return xerrors.Errorf("insert data: %w", err)
	}

	return nil
}

func (q *qdrant) Search(
	ctx context.Context, request *vectordatabase.SearchRequest,
) (resp vectordatabase.SearchResponse, err error) {
	indexOnly := true

	searchReq := &pb.SearchPoints{
		CollectionName: q.collectionName,
		Vector:         request.Vector,
		Limit:          uint64(request.TopK),
		Params: &pb.SearchParams{
			IndexedOnly: &indexOnly,
		},
	}
	if request.Filter != nil {
		searchReq.Filter, err = convertSearchFilter(request.Filter)
		if err != nil {
			return vectordatabase.SearchResponse{}, xerrors.Errorf("convert search filter: %w", err)
		}
	}

	searchResp, err := q.pointsClient.Search(ctx, searchReq)
	if err != nil {
		return vectordatabase.SearchResponse{}, xerrors.Errorf("search: %w", err)
	}

	scorePoints := searchResp.GetResult()
	for _, scorePoint := range scorePoints {
		resp.Data = append(resp.Data, vectordatabase.SearchResponseData{
			ID:    uint(scorePoint.GetId().GetNum()),
			Score: scorePoint.GetScore(),
		})
	}

	return resp, nil
}

func (q *qdrant) Delete(ctx context.Context, request *vectordatabase.DeleteRequest) (err error) {
	wait := true
	_, err = q.pointsClient.Delete(ctx, &pb.DeletePoints{
		CollectionName: q.collectionName,
		Wait:           &wait,
		Points: &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Points{
				Points: &pb.PointsIdsList{
					Ids: newPointIDs(request.IDs),
				},
			},
		},
	})
	if err != nil {
		return xerrors.Errorf("delete: %w", err)
	}

	return nil
}
