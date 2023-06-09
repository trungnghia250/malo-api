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
	GetDashboard(ctx *fiber.Ctx, start, end time.Time, phones []string) ([]dto.CustomerReport, error)
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

	if len(req.Phone) > 0 {
		matching["phone"] = bson.M{
			"$in": req.Phone,
		}
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
		bson.D{{"$match", filter}},
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

func (p *productReportRepo) GetDashboard(ctx *fiber.Ctx, start, end time.Time, phones []string) ([]dto.CustomerReport, error) {
	matching := bson.M{
		"date": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	if len(phones) > 0 {
		matching["phone"] = bson.M{
			"$in": phones,
		}
	}

	cursor, err := p.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$group", bson.M{
			"_id":          bson.M{"$dateToString": bson.M{"format": "%d-%m-%Y", "date": "$date"}},
			"total_orders": bson.M{"$sum": "$total_orders"},
			//"cancel_orders":     bson.M{"$sum": "$cancel_orders"},
			//"success_orders":    bson.M{"$sum": "$success_orders"},
			//"processing_orders": bson.M{"$sum": "$processing_orders"},
			"total_revenue": bson.M{"$sum": "$revenue"},
			"new":           bson.M{"$sum": "$new"},
		}}},
		bson.D{{"$sort", bson.M{"_id": 1}}},
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
