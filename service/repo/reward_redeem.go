package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRewardRedeemRepo interface {
	GetRedeemByID(ctx *fiber.Ctx, ID string) (resp *model.RewardRedeem, err error)
	CreateRedeem(ctx *fiber.Ctx, data *model.RewardRedeem) error
	UpdateRedeemByID(ctx *fiber.Ctx, data *model.RewardRedeem) error
	DeleteRedeemsByID(ctx *fiber.Ctx, ids []string) error
	ListRedeems(ctx *fiber.Ctx, req dto.ListRewardRedeemRequest) ([]model.RewardRedeem, error)
}

func NewRewardRedeemRepo(mgo *mongo.Client) IRewardRedeemRepo {
	return &rewardRedeemRepo{
		mgo: mgo,
	}
}

type rewardRedeemRepo struct {
	mgo *mongo.Client
}

func (r *rewardRedeemRepo) getCollection() *mongo.Collection {
	return r.mgo.Database(database.DatabaseMalo).Collection(database.CollectionRewardRedeem)
}

func (r *rewardRedeemRepo) GetRedeemByID(ctx *fiber.Ctx, ID string) (resp *model.RewardRedeem, err error) {
	if err := r.getCollection().FindOne(ctx.Context(), bson.M{"_id": ID}).Decode(&resp); err != nil {
		return &model.RewardRedeem{}, err
	}

	return resp, nil
}

func (r *rewardRedeemRepo) CreateRedeem(ctx *fiber.Ctx, data *model.RewardRedeem) error {
	_, err := r.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (r *rewardRedeemRepo) UpdateRedeemByID(ctx *fiber.Ctx, data *model.RewardRedeem) error {
	_, err := r.getCollection().UpdateOne(ctx.Context(), bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *rewardRedeemRepo) DeleteRedeemsByID(ctx *fiber.Ctx, ids []string) error {
	_, err := r.getCollection().DeleteMany(ctx.Context(), bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *rewardRedeemRepo) ListRedeems(ctx *fiber.Ctx, req dto.ListRewardRedeemRequest) ([]model.RewardRedeem, error) {
	matching := bson.M{}
	if len(req.IDs) > 0 {
		matching["_id"] = bson.M{
			"$in": req.IDs,
		}
	}
	cursor, err := r.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	redeems := make([]model.RewardRedeem, 0)
	if err := cursor.All(ctx.Context(), &redeems); err != nil {
		return nil, err
	}

	return redeems, nil
}
