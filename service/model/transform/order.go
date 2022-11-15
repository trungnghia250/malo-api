package transform

import (
	"fmt"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

func ConvertSapoOrder(sapoOrders []model.SapoOrder) []dto.Order {
	var resp []dto.Order
	for _, order := range sapoOrders {
		voucherCode := ""
		if len(order.DiscountCodes) > 0 {
			voucherCode = order.DiscountCodes[0].Code
		}
		shippingPrice := int32(0)
		if len(order.ShippingLines) > 0 {
			shippingPrice = int32(order.ShippingLines[0].Price)
		}
		var orderItems []dto.Item
		for _, item := range order.LineItems {
			orderItem := dto.Item{
				ProductName:    item.Name,
				SKU:            item.SKU,
				Quantity:       item.Quantity,
				UnitPrice:      int32(item.Price),
				TotalDiscount:  int32(item.TotalDiscount),
				TotalTaxAmount: 0,
				Subtotal:       int32(item.Price)*item.Quantity - int32(item.TotalDiscount),
			}
			orderItems = append(orderItems, orderItem)
		}
		status := ""
		if order.Status == "open" {
			status = "processing"
		}
		if order.Status == "closed" {
			status = "success"
		}
		if order.Status == "cancelled" {
			status = "cancel"
		}
		resp = append(resp, dto.Order{
			OrderID:              fmt.Sprintf("%d", order.ID),
			CustomerName:         fmt.Sprintf("%s %s", order.ShippingAddress.FirstName, order.ShippingAddress.LastName),
			PhoneNumber:          order.ShippingAddress.Phone,
			Email:                order.Email,
			Address:              fmt.Sprintf("%s, %s, %s", order.ShippingAddress.Address1, order.ShippingAddress.City, order.ShippingAddress.Province),
			Source:               "SAPO",
			Status:               status,
			Items:                orderItems,
			VoucherCode:          voucherCode,
			TotalLineItemsAmount: int32(order.TotalLineItemsPrice),
			ShippingPrice:        shippingPrice,
			TotalDiscount:        int32(order.TotalDiscounts),
			TotalTaxAmount:       0,
			TotalOrderAmount:     int32(order.TotalPrice),
			Note:                 order.Note,
			CreateAt:             order.CreatedOn,
		})
	}

	return resp
}
