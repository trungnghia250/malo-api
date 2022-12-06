package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IVoucherRepo interface {
	GetVoucherByID(ctx *fiber.Ctx, ID string) (resp *model.Voucher, err error)
	CreateVoucher(ctx *fiber.Ctx, data *model.Voucher) error
	UpdateVoucherByID(ctx *fiber.Ctx, data *model.Voucher) error
	DeleteVouchersByID(ctx *fiber.Ctx, ids []string) error
	ListVoucher(ctx *fiber.Ctx, req dto.ListVoucherRequest) ([]model.Voucher, error)
	ListValidateVoucherByGroupIDs(ctx *fiber.Ctx, groupIDs []string) (resp []model.Voucher, err error)
}

func NewVoucherRepo(mgo *mongo.Client) IVoucherRepo {
	return &voucherRepo{
		mgo: mgo,
	}
}

type voucherRepo struct {
	mgo *mongo.Client
}

func (v *voucherRepo) getCollection() *mongo.Collection {
	return v.mgo.Database(database.DatabaseMalo).Collection(database.CollectionVoucher)
}

func (v *voucherRepo) GetVoucherByID(ctx *fiber.Ctx, ID string) (resp *model.Voucher, err error) {
	if err := v.getCollection().FindOne(ctx.Context(), bson.M{"_id": ID}).Decode(&resp); err != nil {
		return &model.Voucher{}, err
	}

	return resp, nil
}

func (v *voucherRepo) CreateVoucher(ctx *fiber.Ctx, data *model.Voucher) error {
	_, err := v.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (v *voucherRepo) UpdateVoucherByID(ctx *fiber.Ctx, data *model.Voucher) error {
	_, err := v.getCollection().UpdateOne(ctx.Context(), bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (v *voucherRepo) DeleteVouchersByID(ctx *fiber.Ctx, ids []string) error {
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

func (v *voucherRepo) ListVoucher(ctx *fiber.Ctx, req dto.ListVoucherRequest) ([]model.Voucher, error) {
	matching := bson.M{}

	if len(req.Code) > 0 {
		matching["_id"] = bson.M{
			"$in": req.Code,
		}
	}

	if len(req.CreatedAt) > 0 {
		matching["created_at"] = bson.M{
			"$gte": time.Unix(int64(req.CreatedAt[0]), 0),
			"$lte": time.Unix(int64(req.CreatedAt[1]), 0),
		}
	}

	if len(req.Status) > 0 {
		matching["status"] = req.Status
	}

	if len(req.DiscountAmount) > 0 {
		matching["discount_amount"] = bson.M{
			"$gte": req.DiscountAmount[0],
			"$lte": req.DiscountAmount[1],
		}
	}

	if len(req.MinOrderAmount) > 0 {
		matching["min_order_amount"] = bson.M{
			"$gte": req.MinOrderAmount[0],
			"$lte": req.MinOrderAmount[1],
		}
	}

	if len(req.StartAt) > 0 {
		matching["start_at"] = bson.M{
			"$gte": req.StartAt[0],
			"$lte": req.StartAt[1],
		}
	}

	if len(req.ExpireAt) > 0 {
		matching["expire_at"] = bson.M{
			"$gte": req.ExpireAt[0],
			"$lte": req.ExpireAt[1],
		}
	}

	if len(req.LimitUsage) > 0 {
		matching["limit_usage"] = bson.M{
			"$gte": req.LimitUsage[0],
			"$lte": req.LimitUsage[1],
		}
	}

	if len(req.LimitPerCustomer) > 0 {
		matching["limit_per_customer"] = bson.M{
			"$gte": req.LimitPerCustomer[0],
			"$lte": req.LimitPerCustomer[1],
		}
	}

	cursor, err := v.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	vouchers := make([]model.Voucher, 0)
	if err := cursor.All(ctx.Context(), &vouchers); err != nil {
		return nil, err
	}

	return vouchers, nil
}

func (v *voucherRepo) ListValidateVoucherByGroupIDs(ctx *fiber.Ctx, groupIDs []string) (resp []model.Voucher, err error) {
	query := bson.M{
		"group_ids": bson.M{
			"$in": groupIDs,
		},
		"status": "ACTIVE",
		"start_at": bson.M{
			"$lte": time.Now().Unix(),
		},
		"expire_at": bson.M{
			"$gte": time.Now().Unix(),
		},
		"remain_amount": bson.M{
			"$gt": 0,
		},
	}

	results, err := v.getCollection().Find(ctx.Context(), query)
	if err != nil {
		return resp, err
	}
	if err = results.All(ctx.Context(), &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
