package domain

type StickerListOptionFunc func(*listOption)

func StickerListWithSearchName(searchName string) StickerListOptionFunc {
	return func(o *listOption) {
		o.searchName = searchName
	}
}

func StickerListWithImages(limit uint) StickerListOptionFunc {
	return func(o *listOption) {
		o.withImages = true
		o.withImagesLimit = limit
	}
}

type listOption struct {
	searchName      string
	withImages      bool
	withImagesLimit uint
}

func NewStickerListOption(opts ...StickerListOptionFunc) *listOption {
	var option listOption
	for _, opt := range opts {
		opt(&option)
	}

	return &option
}

func (l *listOption) GetSearchName() string {
	return l.searchName
}

func (l *listOption) GetWithImages() bool {
	return l.withImages
}

func (l *listOption) GetWithImagesLimit() uint {
	return l.withImagesLimit
}
