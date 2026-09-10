package transport

import (
	"context"
	"log/slog"

	"delivery/internal/application/input"
	"delivery/internal/application/output"
	"delivery/internal/application/usecase"
	"delivery/pkg/authorization"
	"delivery/pkg/logger"

	deliveryv1 "proto/gen/delivery/v1"

	"google.golang.org/protobuf/types/known/structpb"
)

type DeliveryTransport struct {
	usecase usecase.DeliveryUsecaseInterface
	opa     authorization.OPAInterface
	log     *slog.Logger
	deliveryv1.UnimplementedDeliveryServiceServer
}

func NewDeliveryTransport(
	uc usecase.DeliveryUsecaseInterface,
	opa authorization.OPAInterface,
	log *slog.Logger,
) *DeliveryTransport {
	return &DeliveryTransport{
		usecase: uc,
		opa:     opa,
		log:     log,
	}
}

func (d *DeliveryTransport) CreateDelivery(ctx context.Context, req *deliveryv1.CreateDeliveryRequest) (*deliveryv1.CreateDeliveryResponse, error) {
	d.log.Debug("CreateDelivery called", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("client_id", req.ClientId))

	dto := input.CreateDeliveryInput{
		To: input.AddressInput{
			Street:       req.To.Street,
			Number:       req.To.Number,
			Neighborhood: req.To.Neighborhood,
			City:         req.To.City,
			ZipCode:      req.To.ZipCode,
		},
		From: input.AddressInput{
			Street:       req.From.Street,
			Number:       req.From.Number,
			Neighborhood: req.From.Neighborhood,
			City:         req.From.City,
			ZipCode:      req.From.ZipCode,
		},
		Weight:   req.Weight,
		ClientID: req.ClientId,
		Metadata: req.Metadata.AsMap(),
	}

	result, err := d.usecase.Create(ctx, dto)
	if err != nil {
		d.log.Error("failed to create delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("client_id", req.ClientId))
		return nil, err
	}

	d.log.Info("delivery created via gRPC", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", result.ID.String()), slog.String("client_id", req.ClientId))

	return &deliveryv1.CreateDeliveryResponse{
		Delivery: outputToProto(result),
	}, nil
}

func (d *DeliveryTransport) FindByID(ctx context.Context, req *deliveryv1.FindByIDRequest) (*deliveryv1.FindByIDResponse, error) {
	d.log.Debug("FindByID called", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", req.Id))

	dto := input.FindByIDDeliveryInput{
		ID: req.Id,
	}

	result, err := d.usecase.FindByID(ctx, dto)
	if err != nil {
		d.log.Error("failed to find delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", req.Id))
		return nil, err
	}

	d.log.Info("delivery found via gRPC", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", req.Id))

	return &deliveryv1.FindByIDResponse{
		Delivery: outputToProto(result),
	}, nil
}

func (d *DeliveryTransport) AssignToDriver(ctx context.Context, req *deliveryv1.AssignToDriverRequest) (*deliveryv1.AssignToDriverResponse, error) {
	d.log.Debug("AssignToDriver called", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", req.DeliveryId), slog.String("driver_id", req.DriverId))

	dto := input.AssignDeliveryToDriverInput{
		DeliveryID: req.DeliveryId,
		DriverID:   req.DriverId,
	}

	result, err := d.usecase.AssignToDriver(ctx, dto)
	if err != nil {
		d.log.Error("failed to assign driver", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", req.DeliveryId), slog.String("driver_id", req.DriverId))
		return nil, err
	}

	d.log.Info("driver assigned via gRPC", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", req.DeliveryId), slog.String("driver_id", req.DriverId))

	return &deliveryv1.AssignToDriverResponse{
		Delivery: outputToProto(result),
	}, nil
}

func (d *DeliveryTransport) FindUnassociated(ctx context.Context, req *deliveryv1.FindUnassociatedRequest) (*deliveryv1.FindUnassociatedResponse, error) {
	d.log.Debug("FindUnassociated called", slog.String("request-id", logger.ExtractRequestID(ctx)))

	results, err := d.usecase.FindUnassociated(ctx)
	if err != nil {
		d.log.Error("failed to find unassociated deliveries", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()))
		return nil, err
	}

	deliveries := make([]*deliveryv1.Delivery, len(results))
	for i := range results {
		deliveries[i] = outputToProto(&results[i])
	}

	d.log.Info("unassociated deliveries found via gRPC", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.Int("count", len(deliveries)))

	return &deliveryv1.FindUnassociatedResponse{
		Deliveries: deliveries,
	}, nil
}

func outputToProto(output *output.DeliveryOutput) *deliveryv1.Delivery {
	metadata, _ := structpb.NewValue(output.Metadata)
	status := deliveryv1.DeliveryStatus_DELIVERY_STATUS_UNSPECIFIED

	switch output.Status {
	case "created":
		status = deliveryv1.DeliveryStatus_DELIVERY_STATUS_CREATED
	case "in_transit":
		status = deliveryv1.DeliveryStatus_DELIVERY_STATUS_IN_TRANSIT
	case "delivered":
		status = deliveryv1.DeliveryStatus_DELIVERY_STATUS_DELIVERED
	case "cancelled":
		status = deliveryv1.DeliveryStatus_DELIVERY_STATUS_CANCELLED
	}

	return &deliveryv1.Delivery{
		Id: output.ID.String(),
		To: &deliveryv1.Address{
			Id:           output.To.ID.String(),
			Street:       output.To.Street,
			Number:       output.To.Number,
			Neighborhood: output.To.Neighborhood,
			City:         output.To.City,
			ZipCode:      output.To.ZipCode,
		},
		From: &deliveryv1.Address{
			Id:           output.From.ID.String(),
			Street:       output.From.Street,
			Number:       output.From.Number,
			Neighborhood: output.From.Neighborhood,
			City:         output.From.City,
			ZipCode:      output.From.ZipCode,
		},
		Weight:   output.Weight,
		Metadata: metadata.GetStructValue(),
		ClientId: output.ClientID.String(),
		DriverId: output.DriverID.String(),
		Status:   status,
	}
}
