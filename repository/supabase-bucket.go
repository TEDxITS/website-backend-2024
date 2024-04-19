package repository

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/dto"
)

type (
	BucketRepository interface {
		UploadFile(string, *multipart.FileHeader) error
		DownloadFile(string, string) ([]byte, error)
	}

	bucketRepository struct {
		bucket config.SupabaseBucket
	}

	DownloadRequest struct {
		Path string `json:"path" form:"path"`
	}
)

func NewSupabaseBucketRepository(b *config.SupabaseBucket) BucketRepository {
	return &bucketRepository{
		bucket: *b,
	}
}

func (r *bucketRepository) UploadFile(folder string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multipartWriter.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	if err := multipartWriter.Close(); err != nil {
		return err
	}

	url := r.bucket.BucketURL + folder + "/" + fileHeader.Filename
	request, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	resp, err := r.bucket.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return dto.ErrFailedToStorePaymentFile
	}

	return nil
}

func (r *bucketRepository) DownloadFile(folder, filename string) ([]byte, error) {
	url := r.bucket.BucketURL + folder + "/" + filename
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.bucket.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, dto.ErrFileNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, dto.ErrFailedToDownloadFile
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
