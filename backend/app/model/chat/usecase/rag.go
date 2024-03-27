package usecase

import (
	"context"

	"backend/app/domain"
	"backend/app/pkg/hserr"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/xerrors"
)

func (u *chatUsecase) CreateRAGReferencePool(
	ctx context.Context, guildID, name, description string,
) (id int, err error) {
	id, err = u.chatRepo.CreateRAGReferencePool(ctx, guildID, name, description)
	if err != nil {
		return 0, xerrors.Errorf("create rag reference pool: %w", err)
	}

	return id, nil
}

func (u *chatUsecase) ListRAGReferencePools(
	ctx context.Context, guildID string, limit, offset int,
) (result domain.ListRAGReferencePoolsResult, err error) {
	result, err = u.chatRepo.ListRAGReferencePools(ctx, guildID, limit, offset)
	if err != nil {
		return domain.ListRAGReferencePoolsResult{}, xerrors.Errorf("list rag reference pools: %w", err)
	}

	return result, nil
}

func (u *chatUsecase) CreateRAGReferenceText(
	ctx context.Context, ragReferencePoolID int, text string,
) (id int, err error) {
	resp, err := u.openaiClient.CreateEmbeddings(ctx,
		openai.EmbeddingRequestStrings{
			Input:          []string{text},
			Model:          openai.SmallEmbedding3,
			EncodingFormat: openai.EmbeddingEncodingFormatBase64,
		})
	if err != nil {
		return 0, hserr.NewInternalError(err, "create embedding")
	}

	id, err = u.chatRepo.CreateRAGReferenceText(
		ctx, ragReferencePoolID, text, resp.Data[0].Embedding)
	if err != nil {
		return 0, xerrors.Errorf("create rag reference text: %w", err)
	}

	return id, nil
}

func (u *chatUsecase) ListRAGReferenceTexts(
	ctx context.Context, ragReferencePoolID int, limit, offset int,
) (result domain.ListRAGReferenceTextsResult, err error) {
	result, err = u.chatRepo.ListRAGReferenceTexts(ctx, ragReferencePoolID, limit, offset)
	if err != nil {
		return domain.ListRAGReferenceTextsResult{}, xerrors.Errorf("list rag reference texts: %w", err)
	}

	return result, nil
}

func (u *chatUsecase) DeleteRAGReferenceText(
	ctx context.Context, ragReferenceTextID int,
) (err error) {
	err = u.chatRepo.DeleteRAGReferenceText(ctx, ragReferenceTextID)
	if err != nil {
		return xerrors.Errorf("delete rag reference text: %w", err)
	}

	return nil
}
