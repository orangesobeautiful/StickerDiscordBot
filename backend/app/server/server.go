package server

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"backend/app/config"
	domainresponse "backend/app/domain-response"
	"backend/app/ent"
	discordcommandRepo "backend/app/model/discord-command/repository"
	discordmessage "backend/app/model/discord-message"
	discordcommand "backend/app/pkg/discord-command"
	objectstorage "backend/app/pkg/object-storage"
	vectordatabase "backend/app/pkg/vector-database"
	"backend/app/pkg/vector-database/qdarnt"

	"entgo.io/ent/dialect"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/xerrors"
)

type Server struct {
	dbClient         *ent.Client
	redisClient      *redis.Client
	vectorDB         vectordatabase.VectorDatabase
	bucketHandler    objectstorage.BucketObjectHandler
	openaiCli        *openai.Client
	hs               *http.Server
	dcCommandManager discordcommand.Manager
	dcSess           *discordgo.Session
}

func NewAndRun(ctx context.Context, cfg config.Config) error {
	var err error

	s := new(Server)

	uni := newValidateTranslator()
	validate, err := newValidate(uni)
	if err != nil {
		return xerrors.Errorf("new validate: %w", err)
	}

	eh := newErrHandler(uni)

	e := newGinEngine(cfg.GetServer().GetCORS(), validate, eh)

	err = s.initDBClient(cfg.GetDatabase())
	if err != nil {
		return xerrors.Errorf("new db client: %w", err)
	}
	err = s.initRedisClient(cfg.GetRedis())
	if err != nil {
		return xerrors.Errorf("new redis client: %w", err)
	}
	err = s.initVectorDatabase(ctx, cfg.GetVectorDatabase())
	if err != nil {
		return xerrors.Errorf("init vector database client: %w", err)
	}

	s.initOpenaiClient(cfg.GetOpenai())

	sessStore := s.newSessStore(
		cfg.GetServer().GetSessionKey().GetUserAuth(),
		cfg.GetServer().GetCookie(),
	)

	bucketHandler, err := objectstorage.NewBucketHandler(ctx, cfg.GetObjectStorage())
	if err != nil {
		return xerrors.Errorf("new bucket handler: %w", err)
	}
	s.bucketHandler = bucketHandler
	rd := domainresponse.New(bucketHandler)

	dcCommandRepo := discordcommandRepo.New(s.dbClient)

	s.initDCCommandManager(validate, eh, dcCommandRepo)

	dcMsgHandler, err := s.setModel(sessStore, e, s.dcCommandManager, rd)
	if err != nil {
		return xerrors.Errorf("set model: %w", err)
	}

	err = s.run(ctx, cfg.GetServer(), cfg.GetDiscord(), e, dcMsgHandler)
	if err != nil {
		return xerrors.Errorf("run: %w", err)
	}

	return nil
}

func newValidateTranslator() *ut.UniversalTranslator {
	english := en_US.New()
	uni := ut.New(english, english)

	return uni
}

func (s *Server) initDBClient(dbCfg config.Database) error {
	dbClient, err := ent.Open(dialect.Postgres, dbCfg.GetDSN())
	if err != nil {
		return xerrors.Errorf("open db connection: %w", err)
	}

	if dbCfg.GetAutoMigrate() {
		if err := dbClient.Schema.Create(context.Background()); err != nil {
			return xerrors.Errorf("auto migrate: %w", err)
		}
	}

	s.dbClient = dbClient
	return nil
}

func (s *Server) initRedisClient(redisCfg config.Redis) error {
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.GetAddr(),
		Username: redisCfg.GetUsername(),
		Password: redisCfg.GetPassword(),
		DB:       redisCfg.GetDB(),
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return xerrors.Errorf("redis ping: %w", err)
	}

	s.redisClient = client
	return nil
}

func (s *Server) initOpenaiClient(openaiCfg config.Openai) {
	client := openai.NewClient(openaiCfg.GetToken())
	s.openaiCli = client
}

func (s *Server) initVectorDatabase(ctx context.Context, vectorDBCfg config.VectorDatabase) (err error) {
	switch vectorDBCfg.GetType() {
	case config.VectorDatabaseTypeQdrant:
		s.vectorDB, err = qdarnt.New(vectorDBCfg)
		if err != nil {
			return xerrors.Errorf("new qdarnt: %w", err)
		}
	default:
		return xerrors.Errorf("unsupported vector database type: %s", vectorDBCfg.GetType())
	}

	if vectorDBCfg.GetToIntializeCollection() {
		const vectorDim = 1536
		err = s.vectorDB.CreateCollectionIfNotExist(ctx, vectorDim, vectordatabase.DistanceTypeCosine)
		if err != nil {
			return xerrors.Errorf("create collection: %w", err)
		}
	}

	return nil
}

func (s *Server) newSessStore(sessionKeyCfg config.SessionKey, cookieCfg config.Cookie) sessions.Store {
	gob.Register(uuid.UUID{})

	cookieStore := sessions.NewCookieStore(sessionKeyCfg.SessionKeyPair()...)
	if cookieCfg != nil {
		cookieStore.Options = &sessions.Options{
			MaxAge:   int(cookieCfg.GetMaxAge().Seconds()),
			Secure:   cookieCfg.GetSecure(),
			HttpOnly: cookieCfg.GetHTTPOnly(),
			SameSite: cookieCfg.GetSameSite(),
		}
	}

	return cookieStore
}

func (s *Server) run(
	ctx context.Context, serverCfg config.Server, dcCfg config.Discord, httpHandler http.Handler, dcMsgHandler discordmessage.HandlerInterface,
) (err error) {
	err = s.runDiscordBot(ctx, dcCfg, dcMsgHandler)
	if err != nil {
		return xerrors.Errorf("run discord bot: %w", err)
	}

	go func() {
		runHTTPServerErr := s.runHTTPServer(serverCfg, httpHandler)
		if runHTTPServerErr != nil {
			slog.Error("run http server", slog.Any("err", runHTTPServerErr))
		}
	}()

	<-ctx.Done()
	err = s.Close()
	if err != nil {
		return xerrors.Errorf("close: %w", err)
	}

	return nil
}

func (s *Server) runDiscordBot(
	ctx context.Context,
	dcCfg config.Discord,
	dcMsgHandler discordmessage.HandlerInterface,
) (err error) {
	discordSess, err := discordgo.New("Bot " + dcCfg.GetToken())
	if err != nil {
		return xerrors.Errorf("new discord session: %w", err)
	}
	discordSess.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		slog.Info(fmt.Sprintf("logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator))
	})
	discordSess.AddHandler(dcMsgHandler.GetHandler())

	err = discordSess.Open()
	if err != nil {
		return xerrors.Errorf("discord session open: %w", err)
	}

	if !dcCfg.GetDisableRegisterCommand() {
		slog.Info("migrate all command")

		err = s.dcCommandManager.MigrateAllCommand(ctx, discordSess, "")
		if err != nil {
			return xerrors.Errorf("migrate all command: %w", err)
		}

		slog.Info("migrate all command done")
	}

	discordSess.AddHandler(s.dcCommandManager.GetHandler())

	s.dcSess = discordSess
	return nil
}

func (s *Server) runHTTPServer(serverCfg config.Server, httpHandler http.Handler) (err error) {
	serverAddr := serverCfg.GetAddr()
	hs := http.Server{
		Addr:        serverAddr,
		Handler:     httpHandler,
		ReadTimeout: 60 * time.Second,
	}
	s.hs = &hs

	slog.Info(fmt.Sprintf("server start at %s", serverAddr))
	err = s.hs.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return xerrors.Errorf("server.ListenAndServe: %w", err)
	}

	return nil
}

func (s *Server) Close() (err error) {
	err = s.dcSess.Close()
	if err != nil {
		return xerrors.Errorf("discord session close: %w", err)
	}

	const closeSererTimeout = 1 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), closeSererTimeout)
	defer cancel()
	err = s.hs.Shutdown(ctx)
	if err != nil {
		return xerrors.Errorf("http server shutdown: %w", err)
	}

	err = s.redisClient.Close()
	if err != nil {
		return xerrors.Errorf("redis client close: %w", err)
	}

	err = s.dbClient.Close()
	if err != nil {
		return xerrors.Errorf("db client close: %w", err)
	}

	return nil
}

func newValidate(uni *ut.UniversalTranslator) (*validator.Validate, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.SetTagName("binding")

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		transTag := field.Tag.Get("trans")
		if transTag == "" {
			return field.Name
		}
		return transTag
	})

	en, _ := uni.GetTranslator("en")
	err := en_translations.RegisterDefaultTranslations(validate, en)
	if err != nil {
		return nil, err
	}

	return validate, nil
}
