package ginauth

import (
	"net/http"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/pkg/hserr"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/xerrors"
)

const (
	sessionName  = "token"
	sessionIDKey = "session_id"

	userCtxKey = "user"
)

type AuthInterface interface {
	GetRequiredAuthMiddleware() gin.HandlerFunc
	SetSession(ctx *gin.Context, sessionID uuid.UUID) error
	MustGetUserFromContext(ctx *gin.Context) *ent.DiscordUser
}

var _ AuthInterface = (*authController)(nil)

type authController struct {
	dcWebLoginUsecase domain.DiscordWebLoginVerificationUsecase

	sessionStore   sessions.Store
	errRespHandler func(*gin.Context, error)
}

func New(
	sessStore sessions.Store, dcWebLoginUsecase domain.DiscordWebLoginVerificationUsecase, opts ...Option,
) *authController {
	o := newOptions(opts...)

	return &authController{
		dcWebLoginUsecase: dcWebLoginUsecase,

		sessionStore:   sessStore,
		errRespHandler: o.errRespHandler,
	}
}

func (c *authController) getSession(ctx *gin.Context) (sess *sessions.Session, err error) {
	sess, err = c.sessionStore.Get(ctx.Request, sessionName)
	if err != nil {
		return nil, xerrors.Errorf("get session: %w", err)
	}

	return sess, nil
}

func (c *authController) SetSession(ctx *gin.Context, sessionID uuid.UUID) (err error) {
	sess, err := c.getSession(ctx)
	if err != nil {
		return xerrors.Errorf("get session: %w", err)
	}
	sess.Values[sessionIDKey] = sessionID
	err = sess.Save(ctx.Request, ctx.Writer)
	if err != nil {
		return hserr.NewInternalError(err, "save session")
	}

	return nil
}

func (c *authController) GetRequiredAuthMiddleware() gin.HandlerFunc {
	return c.requiredAuthMiddleware()
}

func (c *authController) requiredAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := c.getUserFromSession(ctx)
		if err != nil {
			c.errRespHandler(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Set(userCtxKey, user)
	}
}

func (c *authController) getUserFromSession(ctx *gin.Context) (user *ent.DiscordUser, err error) {
	sessionID, err := c.getUserSessionID(ctx)
	if err != nil {
		return nil, xerrors.Errorf("get user session id: %w", err)
	}

	user, err = c.dcWebLoginUsecase.GetDiscordUserBySessionID(ctx, sessionID)
	if err != nil {
		return nil, xerrors.Errorf("get dc user by session id: %w", err)
	}
	if user == nil {
		return nil, hserr.New(http.StatusUnauthorized, "unauthorized")
	}

	return user, nil
}

func (c *authController) getUserSessionID(ctx *gin.Context) (sessionID uuid.UUID, err error) {
	sess, err := c.getSession(ctx)
	if err != nil {
		return uuid.Nil, xerrors.Errorf("get session: %w", err)
	}

	sessionIDAny := sess.Values[sessionIDKey]
	if sessionIDAny == nil {
		return uuid.Nil, hserr.New(http.StatusUnauthorized, "unauthorized")
	}

	sessionID, ok := sessionIDAny.(uuid.UUID)
	if !ok {
		sess.Options.MaxAge = -1
		_ = sess.Save(ctx.Request, ctx.Writer)
		return uuid.Nil, hserr.New(http.StatusUnauthorized, "unauthorized")
	}

	return sessionID, nil
}

func (c *authController) MustGetUserFromContext(ctx *gin.Context) *ent.DiscordUser {
	user, exist := ctx.Get(userCtxKey)
	if !exist {
		panic("user not exist in context")
	}

	return user.(*ent.DiscordUser)
}
