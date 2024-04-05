package qdarnt

import (
	vectordatabase "backend/app/pkg/vector-database"

	pb "github.com/qdrant/go-client/qdrant"
)

func distanceConvert(distance vectordatabase.DistanceType) pb.Distance {
	if distance == vectordatabase.DistanceTypeCosine {
		return pb.Distance_Cosine
	}

	return pb.Distance_UnknownDistance
}
