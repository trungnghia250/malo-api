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

			customer := model.Customer{
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
			}
			mgo.CreateCustomer(&customer)
		} else {
			if order.Status == "success" {
				customer, _ := mgo.GetCustomerByPhone(order.PhoneNumber)
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

	}

	return nil
}
