package repository

import (
	"context"
	"encoding/json"
	"time"

	"backend/app/domain"
	"backend/app/pkg/hserr"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

const verifcationExpDuration = 5 * time.Minute

var _ domain.DiscordWebLoginVerificationRepository = (*dcWebLoginVerifyRepo)(nil)

type (
	dcWebLoginVerifyRepo                  = discordWebLoginVerificationRepository
	discordWebLoginVerificationRepository struct {
		cache *cache.Cache[string]
	}
)

func NewMemoryDCWebLoginVerification() domain.DiscordWebLoginVerificationRepository {
	const cleanupInterval = 10 * time.Minute

	gocacheClient := gocache.New(verifcationExpDuration, cleanupInterval)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)

	return newWithCacheStore(gocacheStore)
}

func NewRedisDCWebLoginVerification(redisClient *redis.Client) domain.DiscordWebLoginVerificationRepository {
	redisStore := redis_store.NewRedis(redisClient,
		store.WithExpiration(verifcationExpDuration),
	)

	return newWithCacheStore(redisStore)
}

func newWithCacheStore(s store.StoreInterface) *dcWebLoginVerifyRepo {
	return &dcWebLoginVerifyRepo{
		cache: cache.New[string](s),
	}
}

func getDCWebLoginVerifyCacheKey(verifyCode string) string {
	return "StickerDiscordBot:discord-web-login-verify:" + verifyCode
}

type dcWebLoginVerifyInfo struct {
	DCUserID    string `json:"dc_user_id"`
	DCChannelID string `json:"dc_channel_id"`
}

func (r *dcWebLoginVerifyRepo) Create(ctx context.Context, verifyCode, dcUserID, dcChannelID string) error {
	bs, err := json.Marshal(dcWebLoginVerifyInfo{
		DCUserID:    dcUserID,
		DCChannelID: dcChannelID,
	})
	if err != nil {
		return hserr.NewInternalError(err, "marshal dbWebLoginVerifyInfo")
	}

	err = r.cache.Set(ctx,
		getDCWebLoginVerifyCacheKey(verifyCode),
		string(bs),
	)
	if err != nil {
		return hserr.NewInternalError(err, "set discord web login verify cache")
	}

	return nil
}

func (r *dcWebLoginVerifyRepo) FindByCode(ctx context.Context, code string) (*domain.DiscordWebLoginVerification, error) {
	jsonStr, err := r.cache.Get(ctx, code)
	if err != nil {
		return nil, hserr.NewInternalError(err, "get discord web login verify cache")
	}

	var saveResult dcWebLoginVerifyInfo
	err = json.Unmarshal([]byte(jsonStr), &saveResult)
	if err != nil {
		_ = r.cache.Delete(ctx, code)
		return nil, hserr.NewInternalError(err, "unmarshal discord web login verify json")
	}

	return convertDCWebLoginVerifyInfo(code, saveResult), nil
}
