package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IProductRepo interface {
	GetProductByID(ctx *fiber.Ctx, productID string) (resp *model.Product, err error)
	ListProduct(ctx *fiber.Ctx, req dto.ListProductRequest) ([]model.Product, error)
	CreateProduct(ctx *fiber.Ctx, data *model.Product) error
	UpdateProductByID(ctx *fiber.Ctx, data *model.Product) error
	DeleteProductByID(ctx *fiber.Ctx, id string) error
	CountProduct(ctx *fiber.Ctx) (int32, error)
}

func NewProductRepo(mgo *mongo.Client) IProductRepo {
	return &productRepo{
		mgo: mgo,
	}
}

type productRepo struct {
	mgo *mongo.Client
}

func (p *productRepo) getCollection() *mongo.Collection {
	return p.mgo.Database(database.DatabaseMalo).Collection(database.CollectionProduct)
}

func (p *productRepo) GetProductByID(ctx *fiber.Ctx, productID string) (resp *model.Product, err error) {
	if err := p.getCollection().FindOne(ctx.Context(), bson.M{"product_id": productID}).Decode(&resp); err != nil {
		return &model.Product{}, err
	}

	return resp, nil
}

func (p *productRepo) CreateProduct(ctx *fiber.Ctx, data *model.Product) error {
	_, err := p.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepo) UpdateProductByID(ctx *fiber.Ctx, data *model.Product) error {
	_, err := p.getCollection().UpdateOne(ctx.Context(), bson.M{"product_id": data.ProductID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepo) DeleteProductByID(ctx *fiber.Ctx, id string) error {
	_, err := p.getCollection().DeleteOne(ctx.Context(), bson.M{
		"product_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepo) ListProduct(ctx *fiber.Ctx, req dto.ListProductRequest) ([]model.Product, error) {
	matching := bson.M{}
	if len(req.SKU) > 0 {
		matching["sku"] = req.SKU
	}
	if len(req.Category) > 0 {
		matching["category"] = bson.M{
			"$in": req.Category,
		}
	}

	productNameQuery := bson.A{}
	for _, name := range req.Name {
		productNameQuery = append(productNameQuery, primitive.Regex{Pattern: name})
	}
	if len(req.Name) > 0 {
		matching["product_name"] = bson.D{
			{"$in", productNameQuery},
		}
	}
	cursor, err := p.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	products := make([]model.Product, 0)
	if err := cursor.All(ctx.Context(), &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productRepo) CountProduct(ctx *fiber.Ctx) (int32, error) {
	value, err := p.getCollection().CountDocuments(ctx.Context(), bson.M{})
	if err != nil {
		return 0, err
	}

	return int32(value), nil
}
