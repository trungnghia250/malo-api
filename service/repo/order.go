package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type IOrderRepo interface {
	GetOrderByID(ctx *fiber.Ctx, orderID string) (resp *model.Order, err error)
	CreateOrder(ctx *fiber.Ctx, data *dto.Order) (int32, error)
	UpdateOrderByID(ctx *fiber.Ctx, data *dto.Order) error
	DeleteOrderByID(ctx *fiber.Ctx, ids []string) error
	ListOrder(ctx *fiber.Ctx, req dto.ListOrderRequest) ([]model.Order, error)
	CountOrder(ctx *fiber.Ctx) (int32, error)
	CheckOrderExist(ctx *fiber.Ctx, query bson.M) (bool, error)
	UpsertOrder(ctx *fiber.Ctx, query bson.M, order dto.Order) (int32, int32, error)
}

func NewOrderRepo(mgo *mongo.Client) IOrderRepo {
	return &orderRepo{
		mgo: mgo,
	}
}

type orderRepo struct {
	mgo *mongo.Client
}

func (o *orderRepo) getCollection() *mongo.Collection {
	return o.mgo.Database(database.DatabaseMalo).Collection(database.CollectionOrder)
}

func (o *orderRepo) GetOrderByID(ctx *fiber.Ctx, orderID string) (resp *model.Order, err error) {
	if err := o.getCollection().FindOne(ctx.Context(), bson.M{"order_id": orderID}).Decode(&resp); err != nil {
		return &model.Order{}, err
	}

	return resp, nil
}

func (o *orderRepo) CreateOrder(ctx *fiber.Ctx, data *dto.Order) (int32, error) {
	data.CreatedAt = time.Now()
	res, err := o.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return 0, err
	}
	if len(res.InsertedID.(primitive.ObjectID).Hex()) > 0 {
		return 1, nil
	}

	return 0, nil
}

func (o *orderRepo) UpdateOrderByID(ctx *fiber.Ctx, data *dto.Order) error {
	data.ModifiedAt = time.Now()
	_, err := o.getCollection().UpdateOne(ctx.Context(), bson.M{"order_id": data.OrderID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepo) DeleteOrderByID(ctx *fiber.Ctx, ids []string) error {
	_, err := o.getCollection().DeleteMany(ctx.Context(), bson.M{
		"order_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepo) ListOrder(ctx *fiber.Ctx, req dto.ListOrderRequest) ([]model.Order, error) {
	matching := bson.M{}
	customerNameQuery := bson.A{}
	for _, name := range req.CustomerName {
		customerNameQuery = append(customerNameQuery, primitive.Regex{Pattern: name})
	}
	if len(req.CustomerName) > 0 {
		matching["customer_name"] = bson.D{
			{"$in", customerNameQuery},
		}
	}
	if len(req.Email) > 0 {
		matching["email"] = req.Email
	}
	if len(req.Source) > 0 {
		matching["source"] = bson.M{
			"$in": req.Source,
		}
	}
	if len(req.TotalLineItemsAmount) > 0 {
		matching["total_line_items_amount"] = bson.M{
			"$gte": req.TotalLineItemsAmount[0],
			"$lte": req.TotalLineItemsAmount[1],
		}
	}
	if len(req.TotalDiscount) > 0 {
		matching["total_discount"] = bson.M{
			"$gte": req.TotalDiscount[0],
			"$lte": req.TotalDiscount[1],
		}
	}
	if len(req.TotalOrderAmount) > 0 {
		matching["total_order_amount"] = bson.M{
			"$gte": req.TotalOrderAmount[0],
			"$lte": req.TotalOrderAmount[1],
		}
	}
	if len(req.Phone) > 0 {
		matching["phone_number"] = req.Phone
	}
	if len(req.Address) > 0 {
		matching["address"] = primitive.Regex{Pattern: req.Address}
	}
	if len(req.VoucherCode) > 0 {
		matching["voucher_code"] = req.VoucherCode
	}
	if len(req.ShippingPrice) > 0 {
		matching["shipping_price"] = bson.M{
			"$gte": req.ShippingPrice[0],
			"$lte": req.ShippingPrice[1],
		}
	}
	if len(req.TotalTax) > 0 {
		matching["total_tax_amount"] = bson.M{
			"$gte": req.TotalTax[0],
			"$lte": req.TotalTax[1],
		}
	}
	if len(req.Status) > 0 {
		matching["status"] = bson.M{
			"$in": req.Status,
		}
	}
	if len(req.OrderIDs) > 0 {
		matching["order_id"] = bson.M{
			"$in": req.OrderIDs,
		}
	}
	cursor, err := o.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	orders := make([]model.Order, 0)
	if err := cursor.All(ctx.Context(), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderRepo) CountOrder(ctx *fiber.Ctx) (int32, error) {
	value, err := o.getCollection().CountDocuments(ctx.Context(), bson.M{})
	if err != nil {
		return 0, err
	}

	return int32(value), nil
}

func (o *orderRepo) CheckOrderExist(ctx *fiber.Ctx, query bson.M) (bool, error) {
	option := options.CountOptions{}
	option.SetLimit(1)
	value, err := o.getCollection().CountDocuments(ctx.Context(), query, &option)
	if err != nil {
		return false, err
	}

	if value > 0 {
		return true, nil
	}

	return false, nil
}

func (o *orderRepo) UpsertOrder(ctx *fiber.Ctx, query bson.M, order dto.Order) (int32, int32, error) {
	opts := options.Update().SetUpsert(true)
	result, err := o.getCollection().UpdateOne(ctx.Context(), query, order, opts)
	if err != nil {
		return 0, 0, err
	}

	return int32(result.UpsertedCount), int32(result.ModifiedCount), nil
}
