package usecase

import (
	"context"
	"log/slog"

	"delivery/internal/application/input"
	"delivery/internal/application/output"
	"delivery/internal/domain/contract"
	"delivery/internal/domain/entity"
	"delivery/pkg/logger"

	"github.com/google/uuid"
)

type DeliveryUsecaseInterface interface {
	FindByID(ctx context.Context, dto input.FindByIDDeliveryInput) (*output.DeliveryOutput, error)
	AssignToDriver(ctx context.Context, dto input.AssignDeliveryToDriverInput) (*output.DeliveryOutput, error)
	FindUnassociated(ctx context.Context) ([]output.DeliveryOutput, error)
	Create(ctx context.Context, dto input.CreateDeliveryInput) (*output.DeliveryOutput, error)
}

type DeliveryUsecase struct {
	repo contract.DeliveryRepositoryInterface
	log  *slog.Logger
}

func NewDeliveryUsecase(repo contract.DeliveryRepositoryInterface, log *slog.Logger) DeliveryUsecaseInterface {
	return &DeliveryUsecase{repo: repo, log: log}
}

func (d *DeliveryUsecase) FindByID(ctx context.Context, dto input.FindByIDDeliveryInput) (*output.DeliveryOutput, error) {
	d.log.Debug("FindByID called", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", dto.ID))

	id, err := uuid.Parse(dto.ID)
	if err != nil {
		d.log.Error("invalid delivery ID format", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", dto.ID))
		return nil, err
	}

	model, err := d.repo.FindByID(ctx, id)
	if err != nil {
		d.log.Error("failed to find delivery by ID", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", dto.ID))
		return nil, err
	}

	d.log.Info("delivery found successfully", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", dto.ID))
	result := deliveryOutput(model)
	return &result, nil
}

func (d *DeliveryUsecase) AssignToDriver(ctx context.Context, dto input.AssignDeliveryToDriverInput) (*output.DeliveryOutput, error) {
	d.log.Debug("AssignToDriver called", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", dto.DeliveryID), slog.String("driver_id", dto.DriverID))

	deliveryID, err := uuid.Parse(dto.DeliveryID)
	if err != nil {
		d.log.Error("invalid delivery ID format", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", dto.DeliveryID))
		return nil, err
	}

	driverID, err := uuid.Parse(dto.DriverID)
	if err != nil {
		d.log.Error("invalid driver ID format", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("driver_id", dto.DriverID))
		return nil, err
	}

	model, err := d.repo.FindByID(ctx, deliveryID)
	if err != nil {
		d.log.Error("failed to find delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", dto.DeliveryID))
		return nil, err
	}

	if err := model.AssingDriver(driverID); err != nil {
		d.log.Error("failed to assign driver to delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", dto.DeliveryID), slog.String("driver_id", dto.DriverID))
		return nil, err
	}

	if err := d.repo.AssingToDriver(ctx, model); err != nil {
		d.log.Error("failed to save driver assignment", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", dto.DeliveryID), slog.String("driver_id", dto.DriverID))
		return nil, err
	}

	d.log.Info("driver assigned to delivery successfully", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", dto.DeliveryID), slog.String("driver_id", dto.DriverID))
	result := deliveryOutput(model)
	return &result, nil
}

func (d *DeliveryUsecase) FindUnassociated(ctx context.Context) ([]output.DeliveryOutput, error) {
	d.log.Debug("FindUnassociated called")

	models, err := d.repo.FindUnassociatedDeliveries(ctx)
	if err != nil {
		d.log.Error("failed to find unassociated deliveries", slog.String("error", err.Error()))
		return nil, err
	}

	results := make([]output.DeliveryOutput, len(models))
	for i := range models {
		results[i] = deliveryOutput(&models[i])
	}

	d.log.Info("unassociated deliveries found", slog.Int("count", len(results)))
	return results, nil
}

func (d *DeliveryUsecase) Create(ctx context.Context, dto input.CreateDeliveryInput) (*output.DeliveryOutput, error) {
	d.log.Debug("Create delivery called", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("client_id", dto.ClientID))

	to, err := newAddress(dto.To)
	if err != nil {
		d.log.Error("invalid destination address", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("client_id", dto.ClientID))
		return nil, err
	}

	from, err := newAddress(dto.From)
	if err != nil {
		d.log.Error("invalid origin address", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("client_id", dto.ClientID))
		return nil, err
	}

	clientID, err := uuid.Parse(dto.ClientID)
	if err != nil {
		d.log.Error("invalid client ID format", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("client_id", dto.ClientID))
		return nil, err
	}

	model, err := entity.NewDelivery(*to, *from, dto.Weight, clientID, dto.Metadata)
	if err != nil {
		d.log.Error("failed to create delivery entity", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("client_id", dto.ClientID))
		return nil, err
	}

	if err := d.repo.Create(ctx, model); err != nil {
		d.log.Error("failed to save delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("client_id", dto.ClientID), slog.String("delivery_id", model.ID.String()))
		return nil, err
	}

	d.log.Info("delivery created successfully", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", model.ID.String()), slog.String("client_id", dto.ClientID))
	result := deliveryOutput(model)
	return &result, nil
}

func newAddress(dto input.AddressInput) (*entity.Address, error) {
	return entity.NewAddress(dto.Street, dto.Number, dto.Neighborhood, dto.City, dto.ZipCode)
}

func deliveryOutput(model *entity.Delivery) output.DeliveryOutput {
	return output.DeliveryOutput{
		ID:       model.ID,
		To:       addressOutput(model.To),
		From:     addressOutput(model.From),
		Weight:   model.Weight,
		Metadata: model.Metadata,
		ClientID: model.ClientID,
		DriverID: model.DriverID,
		Status:   string(model.Status),
	}
}

func addressOutput(address entity.Address) output.AddressOutput {
	return output.AddressOutput{
		ID:           address.ID,
		Street:       address.Street,
		Number:       address.Number,
		Neighborhood: address.Neighborhood,
		City:         address.City,
		ZipCode:      address.ZipCode,
	}
}
