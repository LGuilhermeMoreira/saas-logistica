package service

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"auth/internal/domain/contract"
	"auth/pkg/logger"
)

type OPASyncService struct {
	repo    contract.RoleRepositoryInterface
	storage contract.StorageServiceInterface
	log     *slog.Logger
}

func NewOPASyncService(
	repo contract.RoleRepositoryInterface,
	storage contract.StorageServiceInterface,
	log *slog.Logger,
) contract.OPASyncServiceInterface {
	return &OPASyncService{
		repo:    repo,
		storage: storage,
		log:     log,
	}
}

func (o *OPASyncService) SyncPolicies(ctx context.Context) error {
	reqID := logger.ExtractRequestID(ctx)
	o.log.Info("starting OPA policies sync", slog.String("request_id", reqID))

	roles, err := o.repo.FindAll(ctx)
	if err != nil {
		o.log.Error("failed to fetch roles for OPA sync", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao buscar roles do banco: %w", err)
	}

	payload := map[string]any{
		"roles": roles,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		o.log.Error("failed to marshal roles to json", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao encodar roles para json: %w", err)
	}

	// NOVO: manifest restringindo o bundle a só possuir "roles",
	// deixando o authz.rego local intocado.
	manifest := map[string]any{
		"roots": []string{"roles"},
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		o.log.Error("failed to marshal bundle manifest", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao encodar manifest do bundle: %w", err)
	}

	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// NOVO: escreve o .manifest no tar
	manifestHdr := &tar.Header{
		Name: ".manifest",
		Mode: 0600,
		Size: int64(len(manifestBytes)),
	}
	if err := tw.WriteHeader(manifestHdr); err != nil {
		o.log.Error("failed to write manifest tar header", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao escrever cabeçalho do manifest: %w", err)
	}
	if _, err := tw.Write(manifestBytes); err != nil {
		o.log.Error("failed to write manifest tar body", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao escrever corpo do manifest: %w", err)
	}

	hdr := &tar.Header{
		Name: "data.json",
		Mode: 0600,
		Size: int64(len(jsonBytes)),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		o.log.Error("failed to write tar header", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao escrever cabeçalho do tar: %w", err)
	}
	if _, err := tw.Write(jsonBytes); err != nil {
		o.log.Error("failed to write tar body", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao escrever corpo do tar: %w", err)
	}

	if err := tw.Close(); err != nil {
		o.log.Error("failed to close tar writer", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao fechar tar writer: %w", err)
	}
	if err := gw.Close(); err != nil {
		o.log.Error("failed to close gzip writer", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao fechar gzip writer: %w", err)
	}

	filename := "bundle.tar.gz"
	contentType := "application/gzip"
	reader := bytes.NewReader(buf.Bytes())
	size := int64(buf.Len())

	o.log.Debug("uploading bundle to storage", slog.String("request_id", reqID), slog.String("filename", filename), slog.Int64("size", size))

	_, _, err = o.storage.UploadFile(ctx, filename, contentType, reader, size)
	if err != nil {
		o.log.Error("failed to upload bundle to storage", slog.String("request_id", reqID), slog.Any("error", err))
		return fmt.Errorf("erro ao fazer upload do bundle para o storage: %w", err)
	}

	o.log.Info("OPA policies synced and uploaded successfully", slog.String("request_id", reqID), slog.String("filename", filename))

	return nil
}
