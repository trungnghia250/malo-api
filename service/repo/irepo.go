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
	NewProductRepo() IProductRepo
	NewOrderRepo() IOrderRepo
	NewCounterRepo() ICounterRepo
	NewPartnerRepo() IPartnerRepo
}

func (r repo) NewCustomerRepo() ICustomerRepo {
	return NewCustomerRepo(r.mgo)
}

func (r repo) NewUserRepo() IUserRepo {
	return NewUserRepo(r.mgo)
}

func (r repo) NewProductRepo() IProductRepo {
	return NewProductRepo(r.mgo)
}

func (r repo) NewOrderRepo() IOrderRepo {
	return NewOrderRepo(r.mgo)
}

func (r repo) NewCounterRepo() ICounterRepo {
	return NewCounterRepo(r.mgo)
}

func (r repo) NewPartnerRepo() IPartnerRepo {
	return NewPartnerRepo(r.mgo)
}
