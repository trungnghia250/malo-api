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

type ICampaignRepo interface {
	GetCampaignByID(ctx *fiber.Ctx, campaignID string) (resp *model.Campaign, err error)
	ListCampaign(ctx *fiber.Ctx, req dto.ListCampaignRequest) ([]model.Campaign, error)
	CreateCampaign(ctx *fiber.Ctx, data *model.Campaign) error
	UpdateCampaignByID(ctx *fiber.Ctx, data *model.Campaign) error
	DeleteCampaignByID(ctx *fiber.Ctx, ids []string) error
}

func NewCampaignRepo(mgo *mongo.Client) ICampaignRepo {
	return &campaignRepo{
		mgo: mgo,
	}
}

type campaignRepo struct {
	mgo *mongo.Client
}

func (c *campaignRepo) getCollection() *mongo.Collection {
	return c.mgo.Database(database.DatabaseMalo).Collection(database.CollectionCampaign)
}

func (c *campaignRepo) GetCampaignByID(ctx *fiber.Ctx, campaignID string) (resp *model.Campaign, err error) {
	if err := c.getCollection().FindOne(ctx.Context(), bson.M{"_id": campaignID}).Decode(&resp); err != nil {
		return &model.Campaign{}, err
	}

	return resp, nil
}

func (c *campaignRepo) CreateCampaign(ctx *fiber.Ctx, data *model.Campaign) error {
	_, err := c.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (c *campaignRepo) UpdateCampaignByID(ctx *fiber.Ctx, data *model.Campaign) error {
	_, err := c.getCollection().UpdateOne(ctx.Context(), bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *campaignRepo) DeleteCampaignByID(ctx *fiber.Ctx, ids []string) error {
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

func (c *campaignRepo) ListCampaign(ctx *fiber.Ctx, req dto.ListCampaignRequest) ([]model.Campaign, error) {
	matching := bson.M{}
	if len(req.CreatedAt) > 0 {
		matching["created_at"] = bson.M{
			"$gte": time.Unix(int64(req.CreatedAt[0]), 0),
			"$lte": time.Unix(int64(req.CreatedAt[1]), 0),
		}
	}

	if len(req.Status) > 0 {
		matching["status"] = req.Status
	}

	if len(req.SendAt) > 0 {
		matching["send_at"] = bson.M{
			"$gte": req.SendAt[0],
			"$lte": req.SendAt[1],
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

	products := make([]model.Campaign, 0)
	if err := cursor.All(ctx.Context(), &products); err != nil {
		return nil, err
	}

	return products, nil
}
