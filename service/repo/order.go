package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IOrderRepo interface {
	GetOrderByID(ctx *fiber.Ctx, orderID string) (resp *model.Order, err error)
	CreateOrder(ctx *fiber.Ctx, data *dto.Order) error
	UpdateOrderByID(ctx *fiber.Ctx, data *dto.Order) error
	DeleteOrderByID(ctx *fiber.Ctx, id string) error
	ListOrder(ctx *fiber.Ctx, limit, offset int32) ([]model.Order, error)
	CountOrder(ctx *fiber.Ctx) (int32, error)
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

func (o *orderRepo) CreateOrder(ctx *fiber.Ctx, data *dto.Order) error {
	data.CreateAt = time.Now()
	_, err := o.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
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

func (o *orderRepo) DeleteOrderByID(ctx *fiber.Ctx, id string) error {
	_, err := o.getCollection().DeleteOne(ctx.Context(), bson.M{
		"order_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepo) ListOrder(ctx *fiber.Ctx, limit, offset int32) ([]model.Order, error) {
	matching := bson.M{}

	cursor, err := o.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$skip", offset}},
		bson.D{{"$limit", limit}},
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
