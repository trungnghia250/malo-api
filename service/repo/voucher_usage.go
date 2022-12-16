package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IVoucherUsageRepo interface {
	GetVoucherUsageByID(ctx *fiber.Ctx, ID string) (resp *model.VoucherUsage, err error)
	CreateVoucherUsage(ctx *fiber.Ctx, data *model.VoucherUsage) error
	UpdateVoucherUsageByID(ctx *fiber.Ctx, data *model.VoucherUsage) error
	DeleteVoucherUsagesByID(ctx *fiber.Ctx, ids []string) error
	ListVoucherUsages(ctx *fiber.Ctx, req dto.ListVoucherUsageRequest) ([]model.VoucherUsage, error)
	CountCustomerUseVoucher(ctx *fiber.Ctx, phone, code string) (int32, error)
}

func NewVoucherUsageRepo(mgo *mongo.Client) IVoucherUsageRepo {
	return &voucherUsageRepo{
		mgo: mgo,
	}
}

type voucherUsageRepo struct {
	mgo *mongo.Client
}

func (v *voucherUsageRepo) getCollection() *mongo.Collection {
	return v.mgo.Database(database.DatabaseMalo).Collection(database.CollectionVoucherUsage)
}

func (v *voucherUsageRepo) GetVoucherUsageByID(ctx *fiber.Ctx, ID string) (resp *model.VoucherUsage, err error) {
	if err := v.getCollection().FindOne(ctx.Context(), bson.M{"_id": ID}).Decode(&resp); err != nil {
		return &model.VoucherUsage{}, err
	}

	return resp, nil
}

func (v *voucherUsageRepo) CreateVoucherUsage(ctx *fiber.Ctx, data *model.VoucherUsage) error {
	_, err := v.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (v *voucherUsageRepo) UpdateVoucherUsageByID(ctx *fiber.Ctx, data *model.VoucherUsage) error {
	_, err := v.getCollection().UpdateOne(ctx.Context(), bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (v *voucherUsageRepo) DeleteVoucherUsagesByID(ctx *fiber.Ctx, ids []string) error {
	_, err := v.getCollection().DeleteMany(ctx.Context(), bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (v *voucherUsageRepo) ListVoucherUsages(ctx *fiber.Ctx, req dto.ListVoucherUsageRequest) ([]model.VoucherUsage, error) {
	matching := bson.M{}
	customerNameQuery := bson.A{}
	for _, name := range req.CustomerName {
		customerNameQuery = append(customerNameQuery, primitive.Regex{Pattern: name})
	}
	if len(req.CustomerName) > 0 {
		matching["customer_name"] = bson.D{
			{"$in", customerNameQuery},
		}
	}
	if len(req.Code) > 0 {
		matching["code"] = bson.M{
			"$in": req.Code,
		}
	}
	if len(req.Phone) > 0 {
		matching["phone"] = bson.M{
			"$in": req.Phone,
		}
	}
	if len(req.OrderID) > 0 {
		matching["order_id"] = bson.M{
			"$in": req.OrderID,
		}
	}
	if len(req.CreatedAt) > 0 {
		matching["created_at"] = bson.M{
			"$gte": time.Unix(int64(req.CreatedAt[0]), 0),
			"$lte": time.Unix(int64(req.CreatedAt[1]), 0),
		}
	}
	if len(req.DiscountAmount) > 0 {
		matching["discount_amount"] = bson.M{
			"$gte": req.DiscountAmount[0],
			"$lte": req.DiscountAmount[1],
		}
	}
	if len(req.UsageDate) > 0 {
		matching["used_date"] = bson.M{
			"$gte": req.UsageDate[0],
			"$lte": req.UsageDate[1],
		}
	}
	cursor, err := v.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output",
			bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}},
				{"total_discount", bson.M{"$sum": "$discount_amount"}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	usages := make([]model.VoucherUsage, 0)
	if err := cursor.All(ctx.Context(), &usages); err != nil {
		return nil, err
	}

	return usages, nil
}

func (v *voucherUsageRepo) CountCustomerUseVoucher(ctx *fiber.Ctx, phone, code string) (int32, error) {
	query := bson.M{
		"phone": phone,
		"code":  code,
	}
	count, err := v.getCollection().CountDocuments(ctx.Context(), query)
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}
