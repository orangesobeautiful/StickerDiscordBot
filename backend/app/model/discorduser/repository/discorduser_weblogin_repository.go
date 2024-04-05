package repository

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/discorduser"
	"backend/app/ent/webloginsession"
	"backend/app/pkg/hserr"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"golang.org/x/xerrors"
)

const verifcationExpDuration = 5 * time.Minute

var _ domain.DiscordWebLoginVerificationRepository = (*dcWebLoginVerifyRepo)(nil)

type (
	dcWebLoginVerifyRepo                  = discordWebLoginVerificationRepository
	discordWebLoginVerificationRepository struct {
		*domain.BaseEntRepo
		cache *cache.Cache[string]
	}
)

func NewMemoryDCWebLoginVerification(dbClient *ent.Client) domain.DiscordWebLoginVerificationRepository {
	const cleanupInterval = 10 * time.Minute

	gocacheClient := gocache.New(verifcationExpDuration, cleanupInterval)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)

	return newWithCacheStore(dbClient, gocacheStore)
}

func NewRedisDCWebLoginVerification(dbClient *ent.Client, redisClient *redis.Client) domain.DiscordWebLoginVerificationRepository {
	redisStore := redis_store.NewRedis(redisClient,
		store.WithExpiration(verifcationExpDuration),
	)

	return newWithCacheStore(dbClient, redisStore)
}

func newWithCacheStore(dbClient *ent.Client, s store.StoreInterface) *dcWebLoginVerifyRepo {
	bRepo := domain.NewBaseEntRepo(dbClient)
	return &dcWebLoginVerifyRepo{
		BaseEntRepo: bRepo,
		cache:       cache.New[string](s),
	}
}

func getDCWebLoginVerifyCacheKey(verifyCode string) string {
	return "StickerDiscordBot:discord-web-login-verify:" + verifyCode
}

type dcWebLoginVerifyInfo struct {
	DCUserID    string `json:"dc_user_id,omitempty"`
	DCGuildlID  string `json:"dc_channel_id,omitempty"`
	DCName      string `json:"dc_name,omitempty"`
	DCAvatarURL string `json:"dc_avatar_url,omitempty"`
}

func (r *dcWebLoginVerifyRepo) CreateLoginCode(ctx context.Context, verifyCode string) error {
	bs, err := json.Marshal(dcWebLoginVerifyInfo{})
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

func (r *dcWebLoginVerifyRepo) FindLoginCodeByCode(ctx context.Context, code string) (*domain.DiscordWebLoginVerification, error) {
	jsonStr, err := r.cache.Get(ctx, getDCWebLoginVerifyCacheKey(code))
	if err != nil {
		if errors.Is(err, store.NotFound{}) {
			return nil, nil
		}

		return nil, hserr.NewInternalError(err, "get discord web login verify cache")
	}

	var saveResult dcWebLoginVerifyInfo
	err = json.Unmarshal([]byte(jsonStr), &saveResult)
	if err != nil {
		_ = r.cache.Delete(ctx, getDCWebLoginVerifyCacheKey(code))
		return nil, hserr.NewInternalError(err, "unmarshal discord web login verify json")
	}

	return convertDCWebLoginVerifyInfo(code, saveResult), nil
}

func (r *dcWebLoginVerifyRepo) UpdateDiscordUserInfoByCode(
	ctx context.Context, verifyCode, userDiscordID, userGuildlID, name, avatarURL string,
) error {
	jsonStr, err := r.cache.Get(ctx, getDCWebLoginVerifyCacheKey(verifyCode))
	if err != nil {
		return hserr.NewInternalError(err, "get discord web login verify cache")
	}
	if jsonStr == "" {
		return hserr.NewInternalError(xerrors.New("verify code not found"), "get discord web login verify cache")
	}

	var saveResult dcWebLoginVerifyInfo
	saveResult.DCUserID = userDiscordID
	saveResult.DCGuildlID = userGuildlID
	saveResult.DCName = name
	saveResult.DCAvatarURL = avatarURL

	bs, err := json.Marshal(saveResult)
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

func (r *dcWebLoginVerifyRepo) DeleteLoginCode(ctx context.Context, verifyCode string) error {
	err := r.cache.Delete(ctx, getDCWebLoginVerifyCacheKey(verifyCode))
	if err != nil {
		return hserr.NewInternalError(err, "delete discord web login verify cache")
	}

	return nil
}

func (r *dcWebLoginVerifyRepo) CreateLoginSession(ctx context.Context, dcUserID int) (sessionID uuid.UUID, err error) {
	newSessionID := uuid.New()

	err = r.GetEntClient(ctx).WebLoginSession.
		Create().
		SetID(newSessionID).
		SetDiscordUserID(dcUserID).
		Exec(ctx)
	if err != nil {
		return uuid.Nil, hserr.NewInternalError(err, "create login session")
	}

	return newSessionID, nil
}

func (r *dcWebLoginVerifyRepo) GetDiscordUserBySessionID(ctx context.Context, sessionID uuid.UUID) (*ent.DiscordUser, error) {
	result, err := r.GetEntClient(ctx).DiscordUser.
		Query().
		Where(
			discorduser.HasWebLoginSessionWith(
				webloginsession.ID(sessionID),
			),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, hserr.New(http.StatusNotFound, "discord user not found")
		}

		return nil, hserr.NewInternalError(err, "get discord user")
	}

	return result, nil
}
