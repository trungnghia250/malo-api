package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IPartnerRepo interface {
	GetPartnerConfig(ctx *fiber.Ctx, partner string) (resp *model.PartnerConfig, err error)
	CreatePartnerConfig(ctx *fiber.Ctx, data *model.PartnerConfig) error
	UpdatePartnerByID(ctx *fiber.Ctx, data *model.PartnerConfig) error
}

func NewPartnerRepo(mgo *mongo.Client) IPartnerRepo {
	return &partnerRepo{
		mgo: mgo,
	}
}

type partnerRepo struct {
	mgo *mongo.Client
}

func (p *partnerRepo) getCollection() *mongo.Collection {
	return p.mgo.Database(database.DatabaseMalo).Collection(database.CollectionPartner)
}

func (p *partnerRepo) GetPartnerConfig(ctx *fiber.Ctx, partner string) (resp *model.PartnerConfig, err error) {
	if err := p.getCollection().FindOne(ctx.Context(), bson.M{"_id": partner}).Decode(&resp); err != nil {
		return &model.PartnerConfig{}, err
	}
	return resp, nil
}

func (p *partnerRepo) CreatePartnerConfig(ctx *fiber.Ctx, data *model.PartnerConfig) error {
	_, err := p.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (p *partnerRepo) UpdatePartnerByID(ctx *fiber.Ctx, data *model.PartnerConfig) error {
	_, err := p.getCollection().UpdateOne(ctx.Context(), bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}
