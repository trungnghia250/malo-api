package repo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepo interface {
	GetUserByEmail(ctx *fiber.Ctx, customerID string) (resp *model.User, err error)
	UpdateUser(ctx *fiber.Ctx, user *model.User) error
	RemoveUserData(ctx *fiber.Ctx, userID string, remove bson.M) error
}

func NewUserRepo(mgo *mongo.Client) IUserRepo {
	return &userRepo{
		mgo: mgo,
	}
}

type userRepo struct {
	mgo *mongo.Client
}

func (c *userRepo) getCollection() *mongo.Collection {
	return c.mgo.Database(database.DatabaseMalo).Collection(database.CollectionUser)
}

func (c *userRepo) GetUserByEmail(ctx *fiber.Ctx, email string) (resp *model.User, err error) {
	if err := c.getCollection().FindOne(ctx.Context(), bson.M{"email": email}).Decode(&resp); err != nil {
		return &model.User{}, err
	}

	return resp, nil
}

func (c *userRepo) UpdateUser(ctx *fiber.Ctx, user *model.User) error {
	updateInfo := bson.M{
		"$set": user,
	}
	query := bson.M{
		"user_id": user.UserID,
	}
	_, err := c.getCollection().UpdateOne(ctx.Context(), query, &updateInfo)
	if err != nil {
		return err
	}

	return nil
}

func (c *userRepo) RemoveUserData(ctx *fiber.Ctx, userID string, remove bson.M) error {
	updateInfo := bson.M{
		"$unset": remove,
	}
	query := bson.M{
		"user_id": userID,
	}
	_, err := c.getCollection().UpdateOne(ctx.Context(), query, &updateInfo)
	if err != nil {
		return err
	}

	return nil
}
