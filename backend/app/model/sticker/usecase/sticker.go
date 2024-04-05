package usecase

import (
	"context"
	"math/rand"
	"net/http"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/pkg/hserr"

	"golang.org/x/xerrors"
)

var _ domain.StickerUsecase = (*stickerUsecase)(nil)

type stickerUsecase struct {
	stickerRepository domain.StickerRepository
	imageRepository   domain.ImageRepository
}

func New(stickerRepo domain.StickerRepository, imageRepo domain.ImageRepository) domain.StickerUsecase {
	return &stickerUsecase{
		imageRepository:   imageRepo,
		stickerRepository: stickerRepo,
	}
}

func (s *stickerUsecase) AddImageByURL(ctx context.Context, guildID, name, imageURL string) (err error) {
	err = s.stickerRepository.WithTx(ctx, func(ctx context.Context) error {
		stickerID, txErr := s.stickerRepository.CreateIfNotExist(ctx, guildID, name)
		if txErr != nil {
			return xerrors.Errorf("sticker creat if not exist: %w", txErr)
		}

		_, txErr = s.imageRepository.CreateWithURL(ctx, stickerID, imageURL)
		if txErr != nil {
			return xerrors.Errorf("image create: %w", txErr)
		}

		return nil
	})
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (s *stickerUsecase) RandSelectImage(ctx context.Context, guildID, stickerName string) (result *ent.Image, err error) {
	images, err := s.GetStickerAllImages(ctx, guildID, stickerName)
	if err != nil {
		return nil, xerrors.Errorf("get sticker all image: %w", err)
	}
	if len(images) == 0 {
		return nil, nil
	}

	//nolint:gosec // not for security purpose
	randSelectIndex := rand.Intn(len(images))
	return images[randSelectIndex], nil
}

func (s *stickerUsecase) ListStickers(
	ctx context.Context, guildID string, offset, limit int, opts ...domain.StickerListOptionFunc,
) (stickers domain.ListStickerResult, err error) {
	stickers, err = s.stickerRepository.List(ctx, guildID, offset, limit, opts...)
	if err != nil {
		return stickers, xerrors.Errorf(": %w", err)
	}

	return stickers, nil
}

func (s *stickerUsecase) GetStickerAllImages(ctx context.Context, guildID, stickerName string) (result []*ent.Image, err error) {
	sticker, err := s.stickerRepository.FindByName(ctx, guildID, stickerName)
	if err != nil {
		return result, xerrors.Errorf("sticker find by id: %w", err)
	}
	if sticker == nil {
		return nil, nil
	}

	result, err = s.imageRepository.ListAllByStickerID(ctx, sticker.ID)
	if err != nil {
		return result, xerrors.Errorf("image list all: %w", err)
	}

	return result, nil
}

func (s *stickerUsecase) Delete(ctx context.Context, ids ...int) (err error) {
	for _, stickerID := range ids {
		images, err := s.imageRepository.ListAllByStickerID(ctx, stickerID)
		if err != nil {
			return xerrors.Errorf("image list all by sticker id: %w", err)
		}
		if len(images) > 0 {
			err = s.imageRepository.DeleteByImageEnt(ctx, images...)
			if err != nil {
				return xerrors.Errorf("image delete: %w", err)
			}
		}

		err = s.stickerRepository.Delete(ctx, ids...)
		if err != nil {
			return xerrors.Errorf("sticker delete: %w", err)
		}
	}

	return nil
}

func (s *stickerUsecase) DeleteByName(ctx context.Context, guildID, name string) (err error) {
	sticker, err := s.stickerRepository.FindByName(ctx, guildID, name)
	if err != nil {
		return xerrors.Errorf("sticker find by name: %w", err)
	}
	if sticker == nil {
		return hserr.New(http.StatusBadRequest, "sticker not exist")
	}

	err = s.Delete(ctx, sticker.ID)
	if err != nil {
		return xerrors.Errorf("delete: %w", err)
	}

	return nil
}
