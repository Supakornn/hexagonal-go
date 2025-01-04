package ordersHandlersgo

import (
	"strings"
	"time"

	"github.com/Supakornn/hexagonal-go/config"
	"github.com/Supakornn/hexagonal-go/modules/entities"
	"github.com/Supakornn/hexagonal-go/modules/orders"
	"github.com/Supakornn/hexagonal-go/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type orderHandlerErrcode string

const (
	FindOneOrderErr orderHandlerErrcode = "orders-001"
	FindOrdersErr   orderHandlerErrcode = "orders-002"
	InsertOrderErr  orderHandlerErrcode = "orders-003"
	UpdateOrderErr  orderHandlerErrcode = "orders-004"
)

type IOrdersHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindOrders(c *fiber.Ctx) error
	InsertOrder(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	cfg           config.IConfig
	ordersUsecase ordersUsecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, ordersUsecase ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		ordersUsecase: ordersUsecase,
	}
}

func (h *ordersHandler) FindOneOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.ordersUsecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(FindOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}

func (h *ordersHandler) FindOrders(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(FindOrdersErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(FindOrdersErr),
				"invalid start date",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(FindOrdersErr),
				"invalid end date",
			).Res()
		}
		req.EndDate = end.Format("2006-01-02")
	}

	orders := h.ordersUsecase.FindOrders(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, orders).Res()
}

func (h *ordersHandler) InsertOrder(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	req := &orders.Order{
		Products: make([]*orders.ProductsOrder, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertOrderErr),
			err.Error(),
		).Res()
	}

	if len(req.Products) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertOrderErr),
			"products is empty",
		).Res()
	}

	if c.Locals("userRoleId").(int) != 2 {
		req.UserId = userId
	}

	req.Status = "waiting"
	req.TotalPaid = 0

	order, err := h.ordersUsecase.InsertOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(InsertOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, order).Res()
}

func (h *ordersHandler) UpdateOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")
	req := new(orders.Order)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateOrderErr),
			err.Error(),
		).Res()
	}

	req.Id = orderId

	statusMap := map[string]string{
		"waiting":   "waiting",
		"shipping":  "shipping",
		"completed": "completed",
		"cancelled": "cancelled",
	}

	if c.Locals("userRoleId").(int) == 2 {
		req.Status = statusMap[strings.ToLower(req.Status)]
	} else if strings.ToLower(req.Status) == statusMap["cancelled"] {
		req.Status = statusMap["cancelled"]
	}

	if req.TransferSlip != nil {
		if req.TransferSlip.Id == "" {
			req.TransferSlip.Id = uuid.NewString()
		}

		if req.TransferSlip.CreatedAt == "" {
			loc, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return entities.NewResponse(c).Error(
					fiber.StatusInternalServerError,
					string(UpdateOrderErr),
					err.Error(),
				).Res()
			}

			now := time.Now().In(loc)
			req.TransferSlip.CreatedAt = now.Format("2006-01-02 15:04:05")
		}
	}

	order, err := h.ordersUsecase.UpdateOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(UpdateOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
