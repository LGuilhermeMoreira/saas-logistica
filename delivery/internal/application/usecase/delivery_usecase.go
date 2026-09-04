package usecase

import (
	"context"

	"delivery/internal/application/input"
	"delivery/internal/application/output"
	"delivery/internal/domain/contract"
	"delivery/internal/domain/entity"

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
}

func NewDeliveryUsecase(repo contract.DeliveryRepositoryInterface) DeliveryUsecaseInterface {
	return &DeliveryUsecase{repo: repo}
}

func (d *DeliveryUsecase) FindByID(ctx context.Context, dto input.FindByIDDeliveryInput) (*output.DeliveryOutput, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return nil, err
	}

	model, err := d.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := deliveryOutput(model)
	return &result, nil
}

func (d *DeliveryUsecase) AssignToDriver(ctx context.Context, dto input.AssignDeliveryToDriverInput) (*output.DeliveryOutput, error) {
	deliveryID, err := uuid.Parse(dto.DeliveryID)
	if err != nil {
		return nil, err
	}

	driverID, err := uuid.Parse(dto.DriverID)
	if err != nil {
		return nil, err
	}

	model, err := d.repo.FindByID(ctx, deliveryID)
	if err != nil {
		return nil, err
	}

	if err := model.AssingDriver(driverID); err != nil {
		return nil, err
	}

	if err := d.repo.AssingToDriver(ctx, model); err != nil {
		return nil, err
	}

	result := deliveryOutput(model)
	return &result, nil
}

func (d *DeliveryUsecase) FindUnassociated(ctx context.Context) ([]output.DeliveryOutput, error) {
	models, err := d.repo.FindUnassociatedDeliveries(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]output.DeliveryOutput, len(models))
	for i := range models {
		results[i] = deliveryOutput(&models[i])
	}

	return results, nil
}

func (d *DeliveryUsecase) Create(ctx context.Context, dto input.CreateDeliveryInput) (*output.DeliveryOutput, error) {
	to, err := newAddress(dto.To)
	if err != nil {
		return nil, err
	}

	from, err := newAddress(dto.From)
	if err != nil {
		return nil, err
	}

	clientID, err := uuid.Parse(dto.ClientID)
	if err != nil {
		return nil, err
	}

	model, err := entity.NewDelivery(*to, *from, dto.Weight, clientID, dto.Metadata)
	if err != nil {
		return nil, err
	}

	if err := d.repo.Create(ctx, model); err != nil {
		return nil, err
	}

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
