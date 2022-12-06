package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICustomerGroupRepo interface {
	GetCustomerGroupByID(ctx *fiber.Ctx, ID string) (resp *model.CustomerGroup, err error)
	CreateCustomerGroup(ctx *fiber.Ctx, data *model.CustomerGroup) error
	UpdateCustomerGroupByID(ctx *fiber.Ctx, data *model.CustomerGroup) error
	DeleteCustomerGroupByID(ctx *fiber.Ctx, ids []string) error
	ListCustomerGroup(ctx *fiber.Ctx, req dto.ListCustomerGroupRequest) ([]model.CustomerGroup, error)
	ListCustomerGroupByCustomerID(ctx *fiber.Ctx, customerID string) (resp []model.CustomerGroup, err error)
}

func NewCustomerGroupRepo(mgo *mongo.Client) ICustomerGroupRepo {
	return &customerGroupRepo{
		mgo: mgo,
	}
}

type customerGroupRepo struct {
	mgo *mongo.Client
}

func (c *customerGroupRepo) getCollection() *mongo.Collection {
	return c.mgo.Database(database.DatabaseMalo).Collection(database.CollectionCustomerGroup)
}

func (c *customerGroupRepo) GetCustomerGroupByID(ctx *fiber.Ctx, ID string) (resp *model.CustomerGroup, err error) {
	if err := c.getCollection().FindOne(ctx.Context(), bson.M{"_id": ID}).Decode(&resp); err != nil {
		return &model.CustomerGroup{}, err
	}

	return resp, nil
}

func (c *customerGroupRepo) ListCustomerGroupByCustomerID(ctx *fiber.Ctx, customerID string) (resp []model.CustomerGroup, err error) {
	results, err := c.getCollection().Find(ctx.Context(), bson.M{"customer_ids": customerID})
	if err != nil {
		return resp, err
	}
	if err = results.All(ctx.Context(), &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *customerGroupRepo) CreateCustomerGroup(ctx *fiber.Ctx, data *model.CustomerGroup) error {
	_, err := c.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (c *customerGroupRepo) UpdateCustomerGroupByID(ctx *fiber.Ctx, data *model.CustomerGroup) error {
	_, err := c.getCollection().UpdateOne(ctx.Context(), bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *customerGroupRepo) DeleteCustomerGroupByID(ctx *fiber.Ctx, ids []string) error {
	_, err := c.getCollection().DeleteMany(ctx.Context(), bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *customerGroupRepo) ListCustomerGroup(ctx *fiber.Ctx, req dto.ListCustomerGroupRequest) ([]model.CustomerGroup, error) {
	matching := bson.M{}
	if len(req.IDs) > 0 {
		matching["_id"] = bson.M{
			"$in": req.IDs,
		}
	}

	cursor, err := c.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	groups := make([]model.CustomerGroup, 0)
	if err := cursor.All(ctx.Context(), &groups); err != nil {
		return nil, err
	}

	return groups, nil
}
