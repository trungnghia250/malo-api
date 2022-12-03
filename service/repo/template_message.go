package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ITemplateRepo interface {
	GetTemplateByID(ctx *fiber.Ctx, ID string) (resp *model.Template, err error)
	CreateTemplate(ctx *fiber.Ctx, data *model.Template) error
	UpdateTemplateByID(ctx *fiber.Ctx, data *model.Template) error
	DeleteTemplateByID(ctx *fiber.Ctx, ids []string) error
	ListTemplate(ctx *fiber.Ctx, req dto.ListTemplateRequest) ([]model.Template, error)
}

func NewTemplateRepo(mgo *mongo.Client) ITemplateRepo {
	return &templateRepo{
		mgo: mgo,
	}
}

type templateRepo struct {
	mgo *mongo.Client
}

func (t *templateRepo) getCollection() *mongo.Collection {
	return t.mgo.Database(database.DatabaseMalo).Collection(database.CollectionTemplate)
}

func (t *templateRepo) GetTemplateByID(ctx *fiber.Ctx, ID string) (resp *model.Template, err error) {
	if err := t.getCollection().FindOne(ctx.Context(), bson.M{"_id": ID}).Decode(&resp); err != nil {
		return &model.Template{}, err
	}

	return resp, nil
}

func (t *templateRepo) CreateTemplate(ctx *fiber.Ctx, data *model.Template) error {
	_, err := t.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (t *templateRepo) UpdateTemplateByID(ctx *fiber.Ctx, data *model.Template) error {
	_, err := t.getCollection().UpdateOne(ctx.Context(), bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *templateRepo) DeleteTemplateByID(ctx *fiber.Ctx, ids []string) error {
	_, err := t.getCollection().DeleteMany(ctx.Context(), bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *templateRepo) ListTemplate(ctx *fiber.Ctx, req dto.ListTemplateRequest) ([]model.Template, error) {
	matching := bson.M{}

	cursor, err := t.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	templates := make([]model.Template, 0)
	if err := cursor.All(ctx.Context(), &templates); err != nil {
		return nil, err
	}

	return templates, nil
}
