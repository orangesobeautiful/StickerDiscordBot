package qdarnt

import (
	"context"

	"backend/app/config"
	vectordatabase "backend/app/pkg/vector-database"

	pb "github.com/qdrant/go-client/qdrant"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var _ vectordatabase.VectorDatabase = (*qdrant)(nil)

type qdrant struct {
	collectionName string

	collectionsClient pb.CollectionsClient
	pointsClient      pb.PointsClient
}

func New(cfg config.VectorDatabase) (db vectordatabase.VectorDatabase, err error) {
	qdrantCfg := cfg.GetQdrant()

	prePRCCred := newGrpcAPIKeyAuth(qdrantCfg.GetAPICredentials())

	conn, err := grpc.Dial(
		qdrantCfg.GetAddr(),
		grpc.WithPerRPCCredentials(prePRCCred),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, xerrors.Errorf("did not connect: %w", err)
	}

	return &qdrant{
		collectionName: cfg.GetCollectionName(),

		collectionsClient: pb.NewCollectionsClient(conn),
		pointsClient:      pb.NewPointsClient(conn),
	}, nil
}

var _ credentials.PerRPCCredentials = grpcAPIKeyAuth{}

type grpcAPIKeyAuth struct {
	requestMetadata map[string]string

	requireTransportSecurity bool
}

func newGrpcAPIKeyAuth(cfg config.QdrantAPICredentials) grpcAPIKeyAuth {
	apiKey := cfg.GetAPIKey()
	requestMetadata := map[string]string{
		"api-key": apiKey,
	}

	requireTransportSecurity := cfg.GetAPIKeyRequireTransportSecurity()

	return grpcAPIKeyAuth{
		requestMetadata:          requestMetadata,
		requireTransportSecurity: requireTransportSecurity,
	}
}

func (auth grpcAPIKeyAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return auth.requestMetadata, nil
}

func (auth grpcAPIKeyAuth) RequireTransportSecurity() bool {
	return auth.requireTransportSecurity
}
