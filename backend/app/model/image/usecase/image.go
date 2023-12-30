package usecase

import (
	"context"

	"backend/app/domain"

	"golang.org/x/xerrors"
)

var _ domain.ImageUsecase = (*imageUsecase)(nil)

type imageUsecase struct {
	imageRepository domain.ImageRepository
}

func New(imgRepo domain.ImageRepository) domain.ImageUsecase {
	return &imageUsecase{
		imageRepository: imgRepo,
	}
}

func (s *imageUsecase) Delete(ctx context.Context, ids ...int) (err error) {
	images, err := s.imageRepository.GetBatch(ctx, ids...)
	if err != nil {
		return xerrors.Errorf("get batch image: %w", err)
	}

	err = s.imageRepository.DeleteByImageEnt(ctx, images...)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}
