package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"backend/app/config"
	domainresponse "backend/app/domain-response"
	"backend/app/ent"
	discordmessage "backend/app/model/discord-message"
	discordcommand "backend/app/pkg/discord-command"
	objectstorage "backend/app/pkg/object-storage"

	"entgo.io/ent/dialect"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/redis/go-redis/v9"
	"golang.org/x/xerrors"
)

type Server struct {
	dbClient         *ent.Client
	redisClient      *redis.Client
	bucketHandler    objectstorage.BucketObjectHandler
	hs               *http.Server
	dcCommandManager discordcommand.Manager
	dcSess           *discordgo.Session
}

func NewAndRun(ctx context.Context, cfg *config.CfgInfo) error {
	var err error

	s := new(Server)

	uni := newValidateTranslator()
	validate, err := newValidate(uni)
	if err != nil {
		return xerrors.Errorf("new validate: %w", err)
	}

	eh := newErrHandler(uni)

	e := newGinEngine(cfg, validate, eh)
	s.initDCCommandManager(validate, eh)

	err = s.initDBClient(cfg)
	if err != nil {
		return xerrors.Errorf("new db client: %w", err)
	}
	err = s.initRedisClient(cfg)
	if err != nil {
		return xerrors.Errorf("new redis client: %w", err)
	}
	bucketHandler, err := objectstorage.NewBucketHandler(ctx, cfg)
	if err != nil {
		return xerrors.Errorf("new bucket handler: %w", err)
	}
	s.bucketHandler = bucketHandler
	rd := domainresponse.New(bucketHandler)

	sessStore := newSessStore(cfg)
	_ = sessStore
	dcMsgHandler, err := s.setModel(e, s.dcCommandManager, rd)
	if err != nil {
		return xerrors.Errorf("set model: %w", err)
	}

	err = s.run(ctx, cfg, e, dcMsgHandler)
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

func (s *Server) initDBClient(cfg *config.CfgInfo) error {
	dbClient, err := ent.Open(dialect.Postgres, cfg.Database.DSN)
	if err != nil {
		return xerrors.Errorf("open db connection: %w", err)
	}

	if cfg.Database.AutoMigrate {
		if err := dbClient.Schema.Create(context.Background()); err != nil {
			return xerrors.Errorf("auto migrate: %w", err)
		}
	}

	s.dbClient = dbClient
	return nil
}

func (s *Server) initRedisClient(cfg *config.CfgInfo) error {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return xerrors.Errorf("redis ping: %w", err)
	}

	s.redisClient = client
	return nil
}

func (s *Server) run(
	ctx context.Context, cfg *config.CfgInfo, httpHandler http.Handler, dcMsgHandler discordmessage.HandlerInterface,
) (err error) {
	err = s.runDiscordBot(cfg, dcMsgHandler)
	if err != nil {
		return xerrors.Errorf("run discord bot: %w", err)
	}

	go func() {
		runHTTPServerErr := s.runHTTPServer(cfg, httpHandler)
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

func (s *Server) runDiscordBot(cfg *config.CfgInfo, dcMsgHandler discordmessage.HandlerInterface) (err error) {
	discordSess, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		return xerrors.Errorf("new discord session: %w", err)
	}
	discordSess.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		slog.Info(fmt.Sprintf("logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator))
	})
	discordSess.AddHandler(dcMsgHandler.GetHandler())

	err = discordSess.Open()
	if err != nil {
		return xerrors.Errorf("discord session open: %w", err)
	}
	err = s.dcCommandManager.RegisterAllCommand(discordSess, "")
	if err != nil {
		return xerrors.Errorf("register all command: %w", err)
	}
	discordSess.AddHandler(s.dcCommandManager.GetHandler())

	s.dcSess = discordSess
	return nil
}

func (s *Server) runHTTPServer(cfg *config.CfgInfo, httpHandler http.Handler) (err error) {
	hs := http.Server{
		Addr:        cfg.Server.Addr,
		Handler:     httpHandler,
		ReadTimeout: 60 * time.Second,
	}
	s.hs = &hs

	slog.Info(fmt.Sprintf("server start at %s", cfg.Server.Addr))
	err = s.hs.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return xerrors.Errorf("server.ListenAndServe: %w", err)
	}

	return nil
}

func (s *Server) Close() (err error) {
	err = s.dcCommandManager.DeleteAllCommand(s.dcSess, "")
	if err != nil {
		return xerrors.Errorf("delete all command: %w", err)
	}
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
