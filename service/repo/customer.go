package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICustomerRepo interface {
	GetCustomerByID(ctx *fiber.Ctx, customerID string) (resp *model.Customer, err error)
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
