package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepo interface {
	GetUserByEmail(ctx *fiber.Ctx, customerID string) (resp *model.User, err error)
	UpdateUser(ctx *fiber.Ctx, user *model.User) error
	RemoveUserData(ctx *fiber.Ctx, userID string, remove bson.M) error
	ListUser(ctx *fiber.Ctx, req dto.ListUserRequest) ([]model.User, error)
	CreateUser(ctx *fiber.Ctx, data *model.User) error
	DeleteUsers(ctx *fiber.Ctx, IDs []string) error
	GetUserByID(ctx *fiber.Ctx, id string) (resp *model.User, err error)
}

func NewUserRepo(mgo *mongo.Client) IUserRepo {
	return &userRepo{
		mgo: mgo,
	}
}

type userRepo struct {
	mgo *mongo.Client
}

func (u *userRepo) getCollection() *mongo.Collection {
	return u.mgo.Database(database.DatabaseMalo).Collection(database.CollectionUser)
}

func (u *userRepo) GetUserByEmail(ctx *fiber.Ctx, email string) (resp *model.User, err error) {
	if err := u.getCollection().FindOne(ctx.Context(), bson.M{"email": email}).Decode(&resp); err != nil {
		return &model.User{}, err
	}

	return resp, nil
}

func (u *userRepo) GetUserByID(ctx *fiber.Ctx, id string) (resp *model.User, err error) {
	if err := u.getCollection().FindOne(ctx.Context(), bson.M{"user_id": id}).Decode(&resp); err != nil {
		return &model.User{}, err
	}

	return resp, nil
}

func (u *userRepo) UpdateUser(ctx *fiber.Ctx, user *model.User) error {
	updateInfo := bson.M{
		"$set": user,
	}
	query := bson.M{
		"user_id": user.UserID,
	}
	_, err := u.getCollection().UpdateOne(ctx.Context(), query, &updateInfo)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) RemoveUserData(ctx *fiber.Ctx, userID string, remove bson.M) error {
	updateInfo := bson.M{
		"$unset": remove,
	}
	query := bson.M{
		"user_id": userID,
	}
	_, err := u.getCollection().UpdateOne(ctx.Context(), query, &updateInfo)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) ListUser(ctx *fiber.Ctx, req dto.ListUserRequest) ([]model.User, error) {
	matching := bson.M{}

	cursor, err := u.getCollection().Aggregate(ctx.Context(), mongo.Pipeline{
		bson.D{{"$match", matching}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$setWindowFields", bson.D{{"output", bson.D{{"totalCount", bson.D{{"$count", bson.M{}}}}}}}}},
		bson.D{{"$skip", req.Offset}},
		bson.D{{"$limit", req.Limit}},
	})

	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)
	if err := cursor.All(ctx.Context(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userRepo) CreateUser(ctx *fiber.Ctx, data *model.User) error {
	_, err := u.getCollection().InsertOne(ctx.Context(), data)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) DeleteUsers(ctx *fiber.Ctx, IDs []string) error {
	_, err := u.getCollection().DeleteMany(ctx.Context(), bson.M{"user_id": bson.M{"$in": IDs}})
	if err != nil {
		return err
	}

	return nil
}
