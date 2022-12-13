package consumers

import (
	"encoding/json"
	"fmt"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"time"
)

func processing(data []byte, mgo *database.MogoDB, config model.RankPointConfig) (err error) {
	var obj model.Mgostream
	if err = json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if obj.OPType == "insert" {
		order, _ := mgo.GetOrderByObjectID(obj.Key.ID)
		customer, products := extractCustomerAndProduct(order, config, mgo)
		var cancel, success, process, revenue, isNew int32
		if order.Status == "cancel" {
			cancel = 1
		}
		if order.Status == "success" {
			success = 1
			revenue = order.TotalOrderAmount
		}
		if order.Status == "processing" {
			process = 1
		}
		if customer.IsNew {
			isNew = 1
		}
		customerReport := model.CustomerReport{
			Phone:         customer.PhoneNumber,
			Date:          time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()),
			Name:          customer.CustomerName,
			Email:         customer.Email,
			TotalOrders:   1,
			SuccessOrders: success,
			ProcessOrders: process,
			CancelOrders:  cancel,
			Revenue:       revenue,
			New:           isNew,
		}
		mgo.UpsertCustomerReport(&customerReport)

		for _, product := range products {
			var revenueProduct, sales int32
			if order.Status == "success" {
				revenueProduct = product.Revenue
				sales = product.Sale
			}
			productReport := model.ProductReport{
				Date:          time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()),
				SKU:           product.SKU,
				Name:          product.ProductName,
				TotalSales:    sales,
				TotalOrders:   1,
				SuccessOrders: success,
				ProcessOrders: process,
				CancelOrders:  cancel,
				Revenue:       revenueProduct,
			}
			mgo.UpsertProductReport(&productReport)
		}
	}

	return nil
}

func extractCustomerAndProduct(order *model.Order, config model.RankPointConfig, mgo *database.MogoDB) (customerRes model.Customer, products []model.Product) {
	history := model.HistoryPoint{
		CustomerName:  order.CustomerName,
		CustomerPhone: order.PhoneNumber,
		RewardPoint:   order.TotalOrderAmount / config.Point[0] * config.Point[1],
		Type:          "+",
		OrderID:       order.OrderID,
		Content:       fmt.Sprintf("Cộng %d điểm từ đơn hàng #%s", order.TotalOrderAmount/config.Point[0]*config.Point[1], order.OrderID),
		CreatedAt:     time.Now(),
	}
	if order.Status == "success" {
		mgo.CreateHistoryPoint(&history)
	}
	for _, item := range order.Items {
		products = append(products, model.Product{
			ProductName: item.ProductName,
			SKU:         item.SKU,
			Sale:        item.Quantity,
			Revenue:     item.Subtotal,
		})
		mgo.UpdateProductByID(&model.Product{
			ProductName: item.ProductName,
			SKU:         item.SKU,
			ModifiedAt:  time.Now(),
		})
	}

	if !mgo.CheckCustomerExist(order.PhoneNumber) {
		customerID, _ := mgo.GetSequenceNextValue("customer_id")
		point := int32(0)
		customerType := "member"
		if order.Status == "success" {
			point = order.TotalOrderAmount / config.Point[0] * config.Point[1]
			for index, rank := range config.Rank {
				if point < rank {
					if index == 2 {
						customerType = "silver"
					}
					if index == 3 {
						customerType = "gold"
					}
					break
				}
				customerType = "platinum"
			}
		}

		customerRes = model.Customer{
			CustomerID:     fmt.Sprintf("C%d", customerID),
			CustomerName:   order.CustomerName,
			Gender:         order.Gender,
			PhoneNumber:    order.PhoneNumber,
			Email:          order.Email,
			Address:        order.Address,
			Province:       order.Province,
			CustomerType:   customerType,
			CustomerSource: order.Source,
			CreatedAt:      time.Now(),
			ModifiedAt:     time.Now(),
			RewardPoint:    point,
			RankPoint:      point,
			IsNew:          true,
		}
		mgo.CreateCustomer(&customerRes)
	} else {
		customer, _ := mgo.GetCustomerByPhone(order.PhoneNumber)
		customerRes = model.Customer{
			CustomerName: order.CustomerName,
			PhoneNumber:  order.PhoneNumber,
			Email:        order.Email,
			IsNew:        false,
		}
		if order.Status == "success" {
			customerType := customer.CustomerType
			for index, rank := range config.Rank {
				if order.TotalOrderAmount/config.Point[0]*config.Point[1]+customer.RankPoint < rank {
					if index == 2 {
						customerType = "silver"
					}
					if index == 3 {
						customerType = "gold"
					}
					break
				}
				customerType = "platinum"
			}
			updateCustomer := dto.Customer{
				CustomerID:   customer.CustomerID,
				CustomerType: customerType,
				ModifiedAt:   time.Now(),
				RewardPoint:  order.TotalOrderAmount/config.Point[0]*config.Point[1] + customer.RewardPoint,
				RankPoint:    order.TotalOrderAmount/config.Point[0]*config.Point[1] + customer.RankPoint,
			}
			mgo.UpdateCustomerByID(&updateCustomer)
		}
	}

	return customerRes, products
}
