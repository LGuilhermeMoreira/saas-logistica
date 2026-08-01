package service

import (
	"archive/tar"
	"auth/internal/domain/contract"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
)

type OPASyncService struct {
	repo    contract.RoleRepositoryInterface
	storage contract.StorageServiceInterface
}

func NewOPASyncService(
	repo contract.RoleRepositoryInterface,
	storage contract.StorageServiceInterface,
) contract.OPASyncServiceInterface {
	return &OPASyncService{
		repo:    repo,
		storage: storage,
	}
}

func (o *OPASyncService) SyncPolicies(ctx context.Context) error {
	roles, err := o.repo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("erro ao buscar roles do banco: %w", err)
	}

	payload := map[string]any{
		"roles": roles,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao encodar roles para json: %w", err)
	}

	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	hdr := &tar.Header{
		Name: "data.json",
		Mode: 0600,
		Size: int64(len(jsonBytes)),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("erro ao escrever cabeçalho do tar: %w", err)
	}
	if _, err := tw.Write(jsonBytes); err != nil {
		return fmt.Errorf("erro ao escrever corpo do tar: %w", err)
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("erro ao fechar tar writer: %w", err)
	}
	if err := gw.Close(); err != nil {
		return fmt.Errorf("erro ao fechar gzip writer: %w", err)
	}

	filename := "bundle.tar.gz"
	contentType := "application/gzip"
	reader := bytes.NewReader(buf.Bytes())
	size := int64(buf.Len())

	_, _, err = o.storage.UploadFile(ctx, filename, contentType, reader, size)
	if err != nil {
		return fmt.Errorf("erro ao fazer upload do bundle para o storage: %w", err)
	}

	return nil
}
