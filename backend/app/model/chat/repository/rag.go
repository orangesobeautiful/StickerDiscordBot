package repository

import (
	"context"
	"encoding/binary"
	"math"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/discordguild"
	"backend/app/ent/ragreferencepool"
	"backend/app/ent/ragreferencetext"
	"backend/app/ent/schema"
	commonRepo "backend/app/pkg/common/repository"
	"backend/app/pkg/hserr"
	vectordatabase "backend/app/pkg/vector-database"
	searchfilter "backend/app/pkg/vector-database/search-filter"

	"golang.org/x/xerrors"
)

const embedMetadataKeyRAGReferencePoolID = "rag_reference_pool_id"

func (r *chatRepository) CreateRAGReferencePool(
	ctx context.Context, guildID, name, description string,
) (id int, err error) {
	result, err := r.GetEntClient(ctx).RAGReferencePool.
		Create().
		SetGuildID(guildID).
		SetName(name).
		SetDescription(description).
		Save(ctx)
	if err != nil {
		return 0, hserr.NewInternalError(err, "create rag reference pool")
	}

	return result.ID, nil
}

func (r *chatRepository) ListRAGReferencePools(
	ctx context.Context, guildID string, limit, offset int,
) (result domain.ListRAGReferencePoolsResult, err error) {
	queryFilter := r.GetEntClient(ctx).RAGReferencePool.
		Query().
		Where(
			ragreferencepool.HasGuildWith(discordguild.ID(guildID)),
		)

	result, err = commonRepo.BaseList[ent.RAGReferencePool](ctx, queryFilter, limit, offset)
	if err != nil {
		return domain.ListRAGReferencePoolsResult{}, xerrors.Errorf(": %w", err)
	}

	return result, nil
}

func (r *chatRepository) GetRAGReferencePoolWithGuildByID(
	ctx context.Context, ragReferencePoolID int,
) (ragReferencePool *ent.RAGReferencePool, err error) {
	ragReferencePool, err = r.GetEntClient(ctx).RAGReferencePool.
		Query().
		WithGuild().
		Where(ragreferencepool.ID(ragReferencePoolID)).
		Only(ctx)
	if err != nil {
		return nil, hserr.NewInternalError(err, "get rag reference pool")
	}

	return ragReferencePool, nil
}

func (r *chatRepository) SearchRAGReferencePoolText(
	ctx context.Context, ragReferencePoolIDs []int, vector []float32, topK uint,
) (result []string, err error) {
	resp, err := r.vectorDB.Search(ctx, &vectordatabase.SearchRequest{
		Vector: vector,
		TopK:   topK,
		Filter: searchfilter.AndFilter(
			searchfilter.ConditionField(
				embedMetadataKeyRAGReferencePoolID,
				searchfilter.MatchInIntegers(ragReferencePoolIDs),
			),
		),
	})
	if err != nil {
		return nil, hserr.NewInternalError(err, "search vector")
	}

	result = make([]string, len(resp.Data))
	for i, data := range resp.Data {
		result[i], err = r.GetRAGReferenceTextContent(ctx, int(data.ID))
		if err != nil {
			return nil, xerrors.Errorf("get rag reference text content: %w", err)
		}
	}

	return result, nil
}

func (r *chatRepository) CreateRAGReferenceText(
	ctx context.Context, ragReferencePoolID int, text string, embedContent []float32,
) (id int, err error) {
	var newRAGText *ent.RAGReferenceText

	embedMetadata := schema.EmbedMetadata{
		embedMetadataKeyRAGReferencePoolID: ragReferencePoolID,
	}

	err = r.WithTx(ctx, func(ctx context.Context) error {
		newRAGText, err = r.GetEntClient(ctx).RAGReferenceText.
			Create().
			SetRefID(ragReferencePoolID).
			SetText(text).
			SetEmbedContent(float32EncodeToBinary(embedContent)).
			SetEmbedMetadata(
				embedMetadata,
			).
			Save(ctx)
		if err != nil {
			return hserr.NewInternalError(err, "create rag reference text")
		}

		err = r.vectorDB.Upsert(ctx, &vectordatabase.UpsertRequest{
			Vectors: []vectordatabase.UpsertRequestVector{
				{
					ID:       uint(newRAGText.ID),
					Data:     embedContent,
					Metadata: embedMetadata,
				},
			},
		})
		if err != nil {
			return hserr.NewInternalError(err, "upsert vector")
		}

		return nil
	})
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	return newRAGText.ID, nil
}

func float32EncodeToBinary(fs []float32) []byte {
	const sizeOfFloat32 = 4

	bf := make([]byte, 0, len(fs)*sizeOfFloat32)
	for _, f := range fs {
		bits := math.Float32bits(f)
		binary.LittleEndian.AppendUint32(bf, bits)
	}

	return bf
}

func (r *chatRepository) GetRAGReferenceTextContent(ctx context.Context, id int) (content string, err error) {
	ragText, err := r.GetEntClient(ctx).RAGReferenceText.
		Query().
		Where(ragreferencetext.ID(id)).
		Only(ctx)
	if err != nil {
		return "", hserr.NewInternalError(err, "get rag reference text")
	}

	return ragText.Text, nil
}

func (r *chatRepository) ListRAGReferenceTexts(
	ctx context.Context, ragReferencePoolID int, limit, offset int,
) (result domain.ListRAGReferenceTextsResult, err error) {
	queryFilter := r.GetEntClient(ctx).RAGReferenceText.
		Query().
		Where(
			ragreferencetext.HasRefWith(
				ragreferencepool.ID(ragReferencePoolID),
			),
		)

	result, err = commonRepo.BaseList(ctx, queryFilter, limit, offset)
	if err != nil {
		return domain.ListRAGReferenceTextsResult{}, xerrors.Errorf(": %w", err)
	}

	return result, nil
}

func (r *chatRepository) GetRAGReferenceTextWithGuildByID(
	ctx context.Context, ragReferenceTextID int,
) (ragReferenceText *ent.RAGReferenceText, err error) {
	ragReferenceText, err = r.GetEntClient(ctx).RAGReferenceText.
		Query().
		WithRef(func(query *ent.RAGReferencePoolQuery) {
			query.WithGuild()
		}).
		Where(ragreferencetext.ID(ragReferenceTextID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewHsNotFoundError("rag reference text")
		}

		return nil, hserr.NewInternalError(err, "get rag reference text")
	}

	return ragReferenceText, nil
}

func (r *chatRepository) DeleteRAGReferenceText(
	ctx context.Context, ragReferenceTextID int,
) (err error) {
	err = r.WithTx(ctx, func(ctx context.Context) error {
		err = r.GetEntClient(ctx).RAGReferenceText.
			DeleteOneID(ragReferenceTextID).
			Exec(ctx)
		if err != nil {
			return hserr.NewInternalError(err, "delete rag reference text")
		}

		err = r.vectorDB.Delete(ctx, &vectordatabase.DeleteRequest{
			IDs: []uint{uint(ragReferenceTextID)},
		})
		if err != nil {
			return hserr.NewInternalError(err, "delete vector")
		}

		return nil
	})
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}
