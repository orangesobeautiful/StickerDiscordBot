package repository

import (
	"context"
	"encoding/json"
	"time"

	"backend/app/pkg/hserr"

	"github.com/meilisearch/meilisearch-go"
)

type stickerMeilisearchEntity struct {
	ID int `json:"id"`

	Name string `json:"name"`

	GuildID string `json:"discord_guild_sticker"`

	CreatedAt int64 `json:"created_at"`
}

func newStickerMeilisearchEntity(
	id int,
	name string,
	guildID string,
	createdAt time.Time,
) *stickerMeilisearchEntity {
	return &stickerMeilisearchEntity{
		ID:        id,
		Name:      name,
		GuildID:   guildID,
		CreatedAt: createdAt.Unix(),
	}
}

type meilisearchStickerResponse struct {
	Hits []*stickerMeilisearchEntity `json:"hits"`

	EstimatedTotalHits int `json:"estimatedTotalHits"`

	Limit int `json:"limit"`

	ProcessingTimeMs int `json:"processingTimeMs"`

	Query string `json:"query"`
}

func (r *stickerRepository) searchWithMeilisearch(
	ctx context.Context,
	guildID string,
	query string,
	offset,
	limit int,
) (*meilisearchStickerResponse, error) {
	searchRespRaw, err := r.meilisearchSticker.SearchRawWithContext(
		ctx,
		query,
		&meilisearch.SearchRequest{
			Offset: int64(offset),
			Limit:  int64(limit),
			Filter: "discord_guild_sticker = " + guildID,
		},
	)
	if err != nil {
		return nil, hserr.NewInternalError(err, "search sticker")
	}

	var searchResp meilisearchStickerResponse

	err = json.Unmarshal(*searchRespRaw, &searchResp)
	if err != nil {
		return nil, hserr.NewInternalError(err, "unmarshal search response")
	}

	return &searchResp, nil
}
