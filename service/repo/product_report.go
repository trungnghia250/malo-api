package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IProductReportRepo interface {
	GetProductReport(ctx *fiber.Ctx, start, end time.Time) ([]dto.ProductReport, error)
}

func NewProductReportRepo(mgo *mongo.Client) IProductReportRepo {
	return &productReportRepo{
		mgo: mgo,
	}
}

type productReportRepo struct {
	mgo *mongo.Client
}

func (p *productReportRepo) getCollection() *mongo.Collection {
	return p.mgo.Database(database.DatabaseMalo).Collection(database.CollectionProductReport)
}

func (p *productReportRepo) GetProductReport(ctx *fiber.Ctx, start, end time.Time) ([]dto.ProductReport, error) {
	matching := bson.M{
		"date": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	cursor, err := p.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$group", bson.M{
			"_id":               "$sku",
			"name":              bson.M{"$first": "$name"},
			"total_orders":      bson.M{"$sum": "$total_orders"},
			"cancel_orders":     bson.M{"$sum": "$cancel_orders"},
			"success_orders":    bson.M{"$sum": "$success_orders"},
			"processing_orders": bson.M{"$sum": "$processing_orders"},
			"total_sales":       bson.M{"$sum": "$total_orders"},
			"total_revenue":     bson.M{"$sum": "$revenue"},
		}}},
	})

	if err != nil {
		return nil, err
	}

	products := make([]dto.ProductReport, 0)
	if err := cursor.All(ctx.Context(), &products); err != nil {
		return nil, err
	}

	return products, nil
}
