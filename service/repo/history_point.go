package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IHistoryPointRepo interface {
	CreateHistoryPoint(ctx *fiber.Ctx, data *model.HistoryPoint) error
	ListHistoryPoint(ctx *fiber.Ctx, req dto.ListHistoryPointRequest) ([]model.HistoryPoint, error)
}

func NewHistoryPointRepo(mgo *mongo.Client) IHistoryPointRepo {
	return &historyPointRepo{
		mgo: mgo,
	}
}

type historyPointRepo struct {
	mgo *mongo.Client
}

func (h *historyPointRepo) getCollection() *mongo.Collection {
	return h.mgo.Database(database.DatabaseMalo).Collection(database.CollectionHistoryPoint)
}

func (h *historyPointRepo) CreateHistoryPoint(ctx *fiber.Ctx, data *model.HistoryPoint) error {
	_, err := h.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (h *historyPointRepo) ListHistoryPoint(ctx *fiber.Ctx, req dto.ListHistoryPointRequest) ([]model.HistoryPoint, error) {
	matching := bson.M{}
	matching["customer_id"] = req.CustomerID
	cursor, err := h.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	histories := make([]model.HistoryPoint, 0)
	if err := cursor.All(ctx.Context(), &histories); err != nil {
		return nil, err
	}

	return histories, nil
}
