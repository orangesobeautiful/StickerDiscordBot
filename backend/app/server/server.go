package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"backend/app/config"
	"backend/app/ent"
	"backend/app/pkg/log"

	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/xerrors"
)

type Server struct {
	dbClient    *ent.Client
	redisClient *redis.Client
	hs          *http.Server
}

func NewAndRun(ctx context.Context, cfg *config.CfgInfo) error {
	var err error

	e, err := newGinEngine(cfg)
	if err != nil {
		log.Errorf("newEngine failed, err=%s", err)
		return err
	}

	s := new(Server)
	err = s.initDBClient(cfg)
	if err != nil {
		return xerrors.Errorf("new db client: %w", err)
	}
	err = s.initRedisClient(cfg)
	if err != nil {
		return xerrors.Errorf("new redis client: %w", err)
	}

	sessStore := newSessStore(cfg)
	_ = sessStore
	err = s.setModel(ctx, e, cfg)
	if err != nil {
		return xerrors.Errorf("set model: %w", err)
	}

	err = s.run(e, cfg)
	if err != nil {
		return xerrors.Errorf("run: %w", err)
	}

	return nil
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

func (s *Server) run(e *gin.Engine, cfg *config.CfgInfo) error {
	hs := http.Server{
		Addr:        cfg.Server.Addr,
		Handler:     e,
		ReadTimeout: 60 * time.Second,
	}
	s.hs = &hs

	slog.Info(fmt.Sprintf("server start at %s", cfg.Server.Addr))
	err := s.hs.ListenAndServe()
	if err != nil {
		return xerrors.Errorf("server.ListenAndServe: %w", err)
	}

	return nil
}
