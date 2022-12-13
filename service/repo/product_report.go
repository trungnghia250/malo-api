package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IProductReportRepo interface {
	GetProductReport(ctx *fiber.Ctx, start, end time.Time, req dto.GetReportRequest) ([]dto.ProductReport, error)
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

func (p *productReportRepo) GetProductReport(ctx *fiber.Ctx, start, end time.Time, req dto.GetReportRequest) ([]dto.ProductReport, error) {
	matching := bson.M{
		"date": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	filter := bson.M{}
	customerNameQuery := bson.A{}
	for _, name := range req.Name {
		customerNameQuery = append(customerNameQuery, primitive.Regex{Pattern: name})
	}
	if len(req.Name) > 0 {
		filter["name"] = bson.D{
			{"$in", customerNameQuery},
		}
	}
	if len(req.SKU) > 0 {
		filter["_id"] = bson.D{
			{"$in", req.SKU},
		}
	}
	if len(req.TotalOrders) > 0 {
		filter["total_orders"] = bson.M{
			"$gte": req.TotalOrders[0],
			"$lte": req.TotalOrders[1],
		}
	}
	if len(req.TotalSuccess) > 0 {
		filter["success_orders"] = bson.M{
			"$gte": req.TotalSuccess[0],
			"$lte": req.TotalSuccess[1],
		}
	}
	if len(req.TotalProcessing) > 0 {
		filter["processing_orders"] = bson.M{
			"$gte": req.TotalProcessing[0],
			"$lte": req.TotalProcessing[1],
		}
	}
	if len(req.TotalCancel) > 0 {
		filter["cancel_orders"] = bson.M{
			"$gte": req.TotalCancel[0],
			"$lte": req.TotalCancel[1],
		}
	}
	if len(req.TotalSales) > 0 {
		filter["total_sales"] = bson.M{
			"$gte": req.TotalSales[0],
			"$lte": req.TotalSales[1],
		}
	}
	if len(req.TotalRevenue) > 0 {
		filter["total_revenue"] = bson.M{
			"$gte": req.TotalRevenue[0],
			"$lte": req.TotalRevenue[1],
		}
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
