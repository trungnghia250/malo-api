package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IProductRepo interface {
	GetProductByID(ctx *fiber.Ctx, productID string) (resp *model.Product, err error)
	ListProduct(ctx *fiber.Ctx, req dto.ListProductRequest) ([]model.Product, error)
	CreateProduct(ctx *fiber.Ctx, data *model.Product) error
	UpdateProductByID(ctx *fiber.Ctx, data *model.Product) error
	DeleteProductByID(ctx *fiber.Ctx, ids []string) error
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

func (p *productRepo) DeleteProductByID(ctx *fiber.Ctx, ids []string) error {
	_, err := p.getCollection().DeleteMany(ctx.Context(), bson.M{
		"product_id": bson.M{
			"$in": ids,
		},
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
	if len(req.ProductIDs) > 0 {
		matching["product_id"] = bson.M{
			"$in": req.ProductIDs,
		}
	}

	if len(req.CreatedAt) > 0 {
		matching["created_at"] = bson.M{
			"$gte": time.Unix(int64(req.CreatedAt[0]), 0),
			"$lte": time.Unix(int64(req.CreatedAt[1]), 0),
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
