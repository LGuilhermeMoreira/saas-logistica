package delivery

import (
	"log/slog"
	"net/http"

	"gateway/internal/domain/contract"
	deliveryv1 "proto/gen/delivery/v1"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/structpb"
)

type DeliveryController struct {
	deliveryService contract.DeliveryService
	log             *slog.Logger
}

func NewDeliveryController(deliveryService contract.DeliveryService, log *slog.Logger) *DeliveryController {
	return &DeliveryController{deliveryService: deliveryService, log: log}
}

func (c *DeliveryController) CreateDelivery(ctx *gin.Context) {
	requestID := ctx.Writer.Header().Get("X-Request-ID")
	var input CreateDeliveryInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metadata, err := structpb.NewStruct(input.Metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.deliveryService.CreateDelivery(ctx.Request.Context(), &deliveryv1.CreateDeliveryRequest{
		To:       addressToProto(input.To),
		From:     addressToProto(input.From),
		Weight:   input.Weight,
		ClientId: input.ClientID,
		Metadata: metadata,
	})
	if err != nil {
		c.log.Error("delivery service create failed", slog.String("request_id", requestID), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, DeliveryResponse{Delivery: deliveryFromProto(response.GetDelivery())})
}

func (c *DeliveryController) FindByID(ctx *gin.Context) {
	response, err := c.deliveryService.FindByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, DeliveryResponse{Delivery: deliveryFromProto(response.GetDelivery())})
}

func (c *DeliveryController) AssignToDriver(ctx *gin.Context) {
	var input AssignToDriverInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.deliveryService.AssignToDriver(ctx.Request.Context(), input.DeliveryID, input.DriverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, DeliveryResponse{Delivery: deliveryFromProto(response.GetDelivery())})
}

func (c *DeliveryController) FindUnassociated(ctx *gin.Context) {
	response, err := c.deliveryService.FindUnassociated(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deliveries := make([]DeliveryOutput, len(response.GetDeliveries()))
	for index, item := range response.GetDeliveries() {
		deliveries[index] = deliveryFromProto(item)
	}
	ctx.JSON(http.StatusOK, DeliveriesResponse{Deliveries: deliveries})
}

func addressToProto(input AddressInput) *deliveryv1.Address {
	return &deliveryv1.Address{
		Street: input.Street, Number: input.Number, Neighborhood: input.Neighborhood,
		City: input.City, ZipCode: input.ZipCode,
	}
}

func deliveryFromProto(item *deliveryv1.Delivery) DeliveryOutput {
	if item == nil {
		return DeliveryOutput{}
	}
	return DeliveryOutput{
		ID: item.GetId(), To: addressFromProto(item.GetTo()), From: addressFromProto(item.GetFrom()),
		Weight: item.GetWeight(), Metadata: item.GetMetadata().AsMap(), ClientID: item.GetClientId(),
		DriverID: item.GetDriverId(), Status: item.GetStatus().String(),
	}
}

func addressFromProto(item *deliveryv1.Address) AddressOutput {
	if item == nil {
		return AddressOutput{}
	}
	return AddressOutput{ID: item.GetId(), Street: item.GetStreet(), Number: item.GetNumber(), Neighborhood: item.GetNeighborhood(), City: item.GetCity(), ZipCode: item.GetZipCode()}
}
