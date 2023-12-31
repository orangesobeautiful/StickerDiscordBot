package usecase

import (
	"context"
	"math/rand"

	"backend/app/domain"
	"backend/app/ent"

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

func (s *stickerUsecase) AddImageByURL(ctx context.Context, name, imageURL string) (err error) {
	stickerID, err := s.stickerRepository.CreateIfNotExist(ctx, name)
	if err != nil {
		return xerrors.Errorf("sticker creat if not exist: %w", err)
	}

	_, err = s.imageRepository.CreateWithURL(ctx, stickerID, imageURL)
	if err != nil {
		return xerrors.Errorf("image create: %w", err)
	}

	return nil
}

func (s *stickerUsecase) RandSelectImage(ctx context.Context, stickerName string) (result *ent.Image, err error) {
	images, err := s.GetStickerAllImages(ctx, stickerName)
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
	ctx context.Context, offset, limit int, opts ...domain.StickerListOptionFunc,
) (stickers domain.ListStickerResult, err error) {
	stickers, err = s.stickerRepository.List(ctx, offset, limit, opts...)
	if err != nil {
		return stickers, xerrors.Errorf(": %w", err)
	}

	return stickers, nil
}

func (s *stickerUsecase) GetStickerAllImages(ctx context.Context, stickerName string) (result []*ent.Image, err error) {
	sticker, err := s.stickerRepository.FindByName(ctx, stickerName)
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

		err = s.imageRepository.DeleteByImageEnt(ctx, images...)
		if err != nil {
			return xerrors.Errorf("image delete: %w", err)
		}

		err = s.stickerRepository.Delete(ctx, ids...)
		if err != nil {
			return xerrors.Errorf("sticker delete: %w", err)
		}
	}

	return nil
}
