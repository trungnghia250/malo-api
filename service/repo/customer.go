package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICustomerRepo interface {
	GetCustomerByID(ctx *fiber.Ctx, customerID string) (resp *model.Customer, err error)
	CreateCustomer(ctx *fiber.Ctx, data *dto.Customer) error
	UpdateCustomerByID(ctx *fiber.Ctx, data *dto.Customer) error
	DeleteCustomerByID(ctx *fiber.Ctx, id string) error
	ListCustomer(ctx *fiber.Ctx, limit, offset int32) ([]model.Customer, error)
	CountCustomer(ctx *fiber.Ctx) (int32, error)
}

func NewCustomerRepo(mgo *mongo.Client) ICustomerRepo {
	return &customerRepo{
		mgo: mgo,
	}
}

type customerRepo struct {
	mgo *mongo.Client
}

func (c *customerRepo) getCollection() *mongo.Collection {
	return c.mgo.Database(database.DatabaseMalo).Collection(database.CollectionCustomer)
}

func (c *customerRepo) GetCustomerByID(ctx *fiber.Ctx, customerID string) (resp *model.Customer, err error) {
	if err := c.getCollection().FindOne(ctx.Context(), bson.M{"customer_id": customerID}).Decode(&resp); err != nil {
		return &model.Customer{}, err
	}

	return resp, nil
}

func (c *customerRepo) CreateCustomer(ctx *fiber.Ctx, data *dto.Customer) error {
	_, err := c.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepo) UpdateCustomerByID(ctx *fiber.Ctx, data *dto.Customer) error {
	_, err := c.getCollection().UpdateOne(ctx.Context(), bson.M{"customer_id": data.CustomerID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepo) DeleteCustomerByID(ctx *fiber.Ctx, id string) error {
	_, err := c.getCollection().DeleteOne(ctx.Context(), bson.M{
		"customer_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepo) ListCustomer(ctx *fiber.Ctx, limit, offset int32) ([]model.Customer, error) {
	matching := bson.M{}

	cursor, err := c.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$skip", offset}},
		bson.D{{"$limit", limit}},
	})

	if err != nil {
		return nil, err
	}

	customers := make([]model.Customer, 0)
	if err := cursor.All(ctx.Context(), &customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (c *customerRepo) CountCustomer(ctx *fiber.Ctx) (int32, error) {
	value, err := c.getCollection().CountDocuments(ctx.Context(), bson.M{})
	if err != nil {
		return 0, err
	}

	return int32(value), nil
}
