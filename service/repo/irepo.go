package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type repo struct {
	mgo *mongo.Client
}

func NewRepo(mgo *mongo.Client) IRepo {
	return &repo{
		mgo: mgo,
	}
}

type IRepo interface {
	NewCustomerRepo() ICustomerRepo
	NewUserRepo() IUserRepo
}

func (r repo) NewCustomerRepo() ICustomerRepo {
	return NewCustomerRepo(r.mgo)
}

func (r repo) NewUserRepo() IUserRepo {
	return NewUserRepo(r.mgo)
}
