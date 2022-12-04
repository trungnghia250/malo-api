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

type IGiftRepo interface {
	GetGiftByID(ctx *fiber.Ctx, ID string) (resp *model.Gift, err error)
	CreateGift(ctx *fiber.Ctx, data *model.Gift) error
	UpdateGiftByID(ctx *fiber.Ctx, data *model.Gift) error
	DeleteGiftsByID(ctx *fiber.Ctx, ids []string) error
	ListGift(ctx *fiber.Ctx, req dto.ListGiftRequest) ([]model.Gift, error)
}

func NewGiftRepo(mgo *mongo.Client) IGiftRepo {
	return &giftRepo{
		mgo: mgo,
	}
}

type giftRepo struct {
	mgo *mongo.Client
}

func (g *giftRepo) getCollection() *mongo.Collection {
	return g.mgo.Database(database.DatabaseMalo).Collection(database.CollectionGift)
}

func (g *giftRepo) GetGiftByID(ctx *fiber.Ctx, ID string) (resp *model.Gift, err error) {
	if err := g.getCollection().FindOne(ctx.Context(), bson.M{"_id": ID}).Decode(&resp); err != nil {
		return &model.Gift{}, err
	}

	return resp, nil
}

func (g *giftRepo) CreateGift(ctx *fiber.Ctx, data *model.Gift) error {
	_, err := g.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (g *giftRepo) UpdateGiftByID(ctx *fiber.Ctx, data *model.Gift) error {
	_, err := g.getCollection().UpdateOne(ctx.Context(), bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *giftRepo) DeleteGiftsByID(ctx *fiber.Ctx, ids []string) error {
	_, err := g.getCollection().DeleteMany(ctx.Context(), bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *giftRepo) ListGift(ctx *fiber.Ctx, req dto.ListGiftRequest) ([]model.Gift, error) {
	matching := bson.M{}
	nameQuery := bson.A{}
	for _, name := range req.Name {
		nameQuery = append(nameQuery, primitive.Regex{Pattern: name})
	}
	if len(req.SKU) > 0 {
		matching["name"] = bson.D{
			{"$in", nameQuery},
		}
	}

	if len(req.CreatedAt) > 0 {
		matching["created_at"] = bson.M{
			"$gte": time.Unix(int64(req.CreatedAt[0]), 0),
			"$lte": time.Unix(int64(req.CreatedAt[1]), 0),
		}
	}

	if len(req.Status) > 0 {
		matching["status"] = req.Status
	}

	if len(req.Category) > 0 {
		matching["category"] = bson.M{
			"$in": req.Category,
		}
	}

	if len(req.GiftIDs) > 0 {
		matching["_id"] = bson.M{
			"$in": req.GiftIDs,
		}
	}

	if len(req.ReleaseAmount) > 0 {
		matching["release_amount"] = bson.M{
			"$gte": req.ReleaseAmount[0],
			"$lte": req.ReleaseAmount[1],
		}
	}

	if len(req.StockAmount) > 0 {
		matching["stock_amount"] = bson.M{
			"$gte": req.StockAmount[0],
			"$lte": req.StockAmount[1],
		}
	}

	if len(req.UsedAmount) > 0 {
		matching["used_amount"] = bson.M{
			"$gte": req.UsedAmount[0],
			"$lte": req.UsedAmount[1],
		}
	}

	if len(req.Price) > 0 {
		matching["price"] = bson.M{
			"$gte": req.Price[0],
			"$lte": req.Price[1],
		}
	}

	cursor, err := g.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	gifts := make([]model.Gift, 0)
	if err := cursor.All(ctx.Context(), &gifts); err != nil {
		return nil, err
	}

	return gifts, nil
}
