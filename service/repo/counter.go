package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICounterRepo interface {
	GetSequenceNextValue(ctx *fiber.Ctx, seqName string) (int32, error)
}

func NewCounterRepo(mgo *mongo.Client) ICounterRepo {
	return &counterRepo{
		mgo: mgo,
	}
}

type counterRepo struct {
	mgo *mongo.Client
}

func (c *counterRepo) getCollection() *mongo.Collection {
	return c.mgo.Database(database.DatabaseMalo).Collection(database.CollectionCustomer)
}

func (c *counterRepo) GetSequenceNextValue(ctx *fiber.Ctx, seqName string) (int32, error) {
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var result bson.M
	if err := c.getCollection().FindOneAndUpdate(ctx.Context(), bson.M{
		"_id": seqName,
	}, bson.M{
		"$inc": bson.M{
			"seq": 1,
		},
	}, &opt).Decode(&result); err != nil {
		return -1, err
	}

	seq := result["seq"].(int32)

	return seq, nil
}
