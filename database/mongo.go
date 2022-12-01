package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var ctx = context.Background()

const (
	DatabaseMalo            = "malo"
	CollectionCustomer      = "customer"
	CollectionUser          = "user"
	CollectionProduct       = "product"
	CollectionOrder         = "order"
	CollectionCounter       = "counter"
	CollectionPartner       = "partner"
	CollectionCampaign      = "campaign"
	CollectionCustomerGroup = "customer_group"
)

func ConnectMongo(user, pass, host string) *mongo.Client {
	connectString := fmt.Sprintf(`mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority`, user, pass, host)
	opt := options.Client()
	opt.ApplyURI(connectString)
	opt.SetTLSConfig(&tls.Config{})
	err := opt.Validate()
	if err != nil {
		log.Fatal(err)
	}
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatalf("err connect to mongodb: %v", err)
	}
	return client
}
