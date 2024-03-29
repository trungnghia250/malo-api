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

type ICustomerRepo interface {
	GetCustomerByID(ctx *fiber.Ctx, customerID string) (resp *model.Customer, err error)
	CreateCustomer(ctx *fiber.Ctx, data *dto.Customer) error
	UpdateCustomerByID(ctx *fiber.Ctx, data *dto.Customer) error
	DeleteCustomersByID(ctx *fiber.Ctx, ids []string) error
	ListCustomer(ctx *fiber.Ctx, req dto.ListCustomerRequest) ([]model.Customer, error)
	CountCustomer(ctx *fiber.Ctx) (int32, error)
	UpdateListCustomers(ctx *fiber.Ctx, req dto.UpdateListCustomerRequest) error
	GetCustomerByPhone(ctx *fiber.Ctx, Phone string) (resp *model.Customer, err error)
	UpsertCustomer(ctx *fiber.Ctx, query bson.M, order dto.Customer) (int32, int32, error)
	CreateCustomerImport(ctx *fiber.Ctx, data *dto.CustomerImport) error
	UpdateCustomerImport(ctx *fiber.Ctx, data *dto.CustomerImport) error
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

func (c *customerRepo) GetCustomerByPhone(ctx *fiber.Ctx, Phone string) (resp *model.Customer, err error) {
	if err := c.getCollection().FindOne(ctx.Context(), bson.M{"phone_number": Phone}).Decode(&resp); err != nil {
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

func (c *customerRepo) CreateCustomerImport(ctx *fiber.Ctx, data *dto.CustomerImport) error {
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

func (c *customerRepo) UpdateCustomerImport(ctx *fiber.Ctx, data *dto.CustomerImport) error {
	_, err := c.getCollection().UpdateOne(ctx.Context(), bson.M{"customer_id": data.CustomerID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepo) DeleteCustomersByID(ctx *fiber.Ctx, ids []string) error {
	_, err := c.getCollection().DeleteMany(ctx.Context(), bson.M{
		"customer_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepo) ListCustomer(ctx *fiber.Ctx, req dto.ListCustomerRequest) ([]model.Customer, error) {
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

	if len(req.Phone) > 0 {
		matching["phone_number"] = req.Phone
	}

	if len(req.Address) > 0 {
		matching["address"] = primitive.Regex{Pattern: req.Address}
	}

	if len(req.CustomerType) > 0 {
		matching["customer_type"] = bson.M{
			"$in": req.CustomerType,
		}
	}

	if len(req.Tags) > 0 {
		matching["tags"] = bson.M{
			"$in": req.Tags,
		}
	}

	if len(req.Gender) > 0 {
		matching["gender"] = req.Gender
	}

	if len(req.Email) > 0 {
		matching["email"] = req.Email
	}

	if len(req.Source) > 0 {
		matching["customer_source"] = bson.M{
			"$in": req.Source,
		}
	}

	if len(req.CustomerIDs) > 0 {
		matching["customer_id"] = bson.M{
			"$in": req.CustomerIDs,
		}
	}
	if len(req.ExceptIDs) > 0 {
		matching["customer_id"] = bson.M{
			"$nin": req.ExceptIDs,
		}
	}

	if req.Limit == 0 {
		req.Limit = 1000
	}

	if len(req.CreatedAt) > 0 {
		matching["created_at"] = bson.M{
			"$gte": time.Unix(int64(req.CreatedAt[0]), 0),
			"$lte": time.Unix(int64(req.CreatedAt[1]), 0),
		}
	}

	if len(req.Revenue) > 0 {
		matching["rank_point"] = bson.M{
			"$gte": req.Revenue[0] / 1000,
			"$lte": req.Revenue[1] / 1000,
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

func (c *customerRepo) UpdateListCustomers(ctx *fiber.Ctx, req dto.UpdateListCustomerRequest) error {
	_, err := c.getCollection().UpdateMany(ctx.Context(),
		bson.M{
			"customer_id": bson.M{"$in": req.CustomerIDs},
		},
		bson.M{"$push": bson.M{
			"tags": bson.M{
				"$each": req.Tags,
			},
		}})

	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepo) UpsertCustomer(ctx *fiber.Ctx, query bson.M, order dto.Customer) (int32, int32, error) {
	opts := options.Update().SetUpsert(true)
	result, err := c.getCollection().UpdateOne(ctx.Context(), query, order, opts)
	if err != nil {
		return 0, 0, err
	}

	return int32(result.UpsertedCount), int32(result.ModifiedCount), nil
}
