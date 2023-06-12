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

type ICustomerReportRepo interface {
	GetCustomerReport(ctx *fiber.Ctx, start, end time.Time, req dto.GetReportRequest) ([]dto.CustomerReport, error)
	GetDashboard(ctx *fiber.Ctx, start, end time.Time, phones []string) ([]dto.CustomerReport, error)
	GetReturn(ctx *fiber.Ctx, start, end time.Time) (res int32, err error)
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

func (c *customerReportRepo) GetCustomerReport(ctx *fiber.Ctx, start, end time.Time, req dto.GetReportRequest) ([]dto.CustomerReport, error) {
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
	if len(req.Email) > 0 {
		filter["email"] = bson.D{
			{"$in", req.Email},
		}
	}
	if len(req.Phone) > 0 {
		filter["_id"] = bson.D{
			{"$in", req.Phone},
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
	if len(req.TotalRevenue) > 0 {
		filter["total_revenue"] = bson.M{
			"$gte": req.TotalRevenue[0],
			"$lte": req.TotalRevenue[1],
		}
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
			"new":               bson.M{"$sum": "$new"},
		}}},
		bson.D{{"$match", filter}},
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

func (c *customerReportRepo) GetDashboard(ctx *fiber.Ctx, start, end time.Time, phones []string) ([]dto.CustomerReport, error) {
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

	cursor, err := c.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$group", bson.M{
			"_id":          bson.M{"$dateToString": bson.M{"format": "%d-%m-%Y", "date": "$date"}},
			"total_orders": bson.M{"$sum": "$total_orders"},
			//"cancel_orders":     bson.M{"$sum": "$cancel_orders"},
			//"success_orders":    bson.M{"$sum": "$success_orders"},
			//"processing_orders": bson.M{"$sum": "$processing_orders"},
			"total_revenue": bson.M{"$sum": "$revenue"},
			"new":           bson.M{"$sum": "$new"},
			"date":          bson.M{"$first": "$date"},
		}}},
		bson.D{{"$sort", bson.M{"date": 1}}},
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

func (c *customerReportRepo) GetReturn(ctx *fiber.Ctx, start, end time.Time) (res int32, err error) {
	matching := bson.M{
		"date": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := c.getCollection().Distinct(ctx.Context(), "phone", matching)
	if err != nil {
		return 0, err
	}

	return int32(len(cursor)), nil
}
