package discordmessage

import (
	"context"
	"log/slog"

	"backend/app/domain"
	domainresponse "backend/app/domain-response"

	"github.com/bwmarrin/discordgo"
)

type HandlerInterface interface {
	GetHandler() func(s *discordgo.Session, m *discordgo.MessageCreate)
}

var _ HandlerInterface = (*handler)(nil)

type handler struct {
	stickerUsecase domain.StickerUsecase
	rd             *domainresponse.DomainResponse
}

func New(stickerUsecase domain.StickerUsecase, rd *domainresponse.DomainResponse) HandlerInterface {
	return &handler{
		stickerUsecase: stickerUsecase,
		rd:             rd,
	}
}

func (d *handler) GetHandler() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return d.Handle
}

func (d *handler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "" {
		return
	}

	ctx := context.Background()
	img, err := d.stickerUsecase.RandSelectImage(ctx, m.Content)
	if err != nil {
		slog.Error("stickerUsecase.RandSelectImage failed", slog.Any("err", err))
		return
	}
	if img == nil {
		return
	}

	respImg := d.rd.NewImageFromEnt(img)
	_, err = s.ChannelMessageSend(m.ChannelID, respImg.URL)
	if err != nil {
		slog.Error("discord message send failed", slog.Any("err", err))
		return
	}
}
