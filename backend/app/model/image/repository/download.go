package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"backend/app/domain"
	"backend/app/pkg/hserr"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	"golang.org/x/xerrors"
)

type downloadCheckResult struct {
	contentLen int64
	kind       types.Type
	headBuf    []byte
	body       io.ReadCloser
}

type downloadAndUploadResult struct {
	saveType       domain.ImgSaveType
	uploadKey      string
	sha256Checksum []byte
}

func (r *imageRepository) downloadAndUploadToObjectStorage(
	ctx context.Context, imageID int, imageURL string,
) (result *downloadAndUploadResult, err error) {
	const timeoutDuration = 1 * time.Minute
	var cancle context.CancelFunc
	ctx, cancle = context.WithTimeout(ctx, timeoutDuration)
	defer cancle()

	downCheckResult, err := downloadAndCheckImage(ctx, imageURL)
	if err != nil {
		return nil, xerrors.Errorf("download and check image: %w", err)
	}
	downloadBody := downCheckResult.body
	headBuf := downCheckResult.headBuf
	defer downloadBody.Close()

	sha256Hasher := sha256.New()

	pr, pw := io.Pipe()
	w := io.MultiWriter(sha256Hasher, pw)
	go func() {
		defer pw.Close()

		_, err = pw.Write(headBuf)
		if err != nil {
			pw.CloseWithError(err)
		}

		_, err = io.Copy(w, downloadBody)
		if err != nil {
			pw.CloseWithError(err)
		}
	}()

	uploadKey := path.Join(
		"sticker",
		"image",
		fmt.Sprintf("%d_%d.%s",
			imageID,
			time.Now().Unix(),
			downCheckResult.kind.Extension),
	)
	err = r.bucketBasics.Upload(ctx, uploadKey, pr,
		PutObjectWithContentLength(downCheckResult.contentLen),
		PutObjectWithContentType(downCheckResult.kind.MIME.Value),
	)
	if err != nil {
		return nil, xerrors.Errorf("put object: %w", err)
	}

	result = &downloadAndUploadResult{
		saveType:       domain.ImgSaveTypeCloudfare,
		uploadKey:      uploadKey,
		sha256Checksum: sha256Hasher.Sum(nil),
	}
	return result, nil
}

func downloadAndCheckImage(ctx context.Context, url string) (result *downloadCheckResult, err error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, xerrors.Errorf("new request: %w", err)
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, xerrors.Errorf("do request: %w", err)
	}
	contentLen := httpResp.ContentLength

	const headerSize = 16
	fileHeader := make([]byte, headerSize)

	_, err = httpResp.Body.Read(fileHeader)
	if err != nil {
		return nil, xerrors.Errorf("read body header: %w", err)
	}
	defer func() {
		if err != nil {
			httpResp.Body.Close()
		}
	}()

	imgKind, isSupported := checkImageType(fileHeader)
	if !isSupported {
		return nil, hserr.New(http.StatusBadRequest,
			"not supported image type",
			hserr.WithDetails("current only support jpeg, jpeg2000, png, gif, webp"),
		)
	}

	result = &downloadCheckResult{
		contentLen: contentLen,
		kind:       imgKind,
		headBuf:    fileHeader,
		body:       httpResp.Body,
	}
	return result, nil
}

func checkImageType(buf []byte) (kind types.Type, isSupported bool) {
	supportedTypeMatcher := matchers.Map{
		matchers.TypeJpeg:     matchers.Jpeg,
		matchers.TypeJpeg2000: matchers.Jpeg2000,
		matchers.TypePng:      matchers.Png,
		matchers.TypeGif:      matchers.Gif,
		matchers.TypeWebp:     matchers.Webp,
	}

	kind = filetype.MatchMap(buf, supportedTypeMatcher)
	return kind, kind != filetype.Unknown
}
