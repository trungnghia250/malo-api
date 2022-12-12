package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ICustomerReportRepo interface {
	GetCustomerReport(ctx *fiber.Ctx, start, end time.Time) ([]dto.CustomerReport, error)
}

func NewCustomerReportRepo(mgo *mongo.Client) ICustomerReportRepo {
	return &customerReportRepo{
		mgo: mgo,
	}
}

type customerReportRepo struct {
	mgo *mongo.Client
}

func (c *customerReportRepo) getCollection() *mongo.Collection {
	return c.mgo.Database(database.DatabaseMalo).Collection(database.CollectionCustomerReport)
}

func (c *customerReportRepo) GetCustomerReport(ctx *fiber.Ctx, start, end time.Time) ([]dto.CustomerReport, error) {
	matching := bson.M{
		"date": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	cursor, err := c.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$group", bson.M{
			"_id":               "$phone",
			"name":              bson.M{"$first": "$name"},
			"email":             bson.M{"$first": "$email"},
			"total_orders":      bson.M{"$sum": "$total_orders"},
			"cancel_orders":     bson.M{"$sum": "$cancel_orders"},
			"success_orders":    bson.M{"$sum": "$success_orders"},
			"processing_orders": bson.M{"$sum": "$processing_orders"},
			"total_revenue":     bson.M{"$sum": "$revenue"},
		}}},
	})

	if err != nil {
		return nil, err
	}

	customers := make([]dto.CustomerReport, 0)
	if err := cursor.All(ctx.Context(), &customers); err != nil {
		return nil, err
	}

	return customers, nil
}
