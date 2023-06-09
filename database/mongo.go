package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var ctx = context.Background()

const (
	DatabaseMalo             = "malo"
	CollectionCustomer       = "customer"
	CollectionUser           = "user"
	CollectionProduct        = "product"
	CollectionOrder          = "order"
	CollectionCounter        = "counter"
	CollectionPartner        = "partner"
	CollectionCampaign       = "campaign"
	CollectionCustomerGroup  = "customer_group"
	CollectionTemplate       = "template"
	CollectionGift           = "gift"
	CollectionVoucher        = "voucher"
	CollectionVoucherUsage   = "voucher_usage"
	CollectionRewardRedeem   = "reward_redeem"
	CollectionHistoryPoint   = "history_point"
	CollectionCustomerReport = "customer_report"
	CollectionProductReport  = "product_report"
)

var (
	MgoClient   *mongo.Client
	MgoDatabase *mongo.Database
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
	MgoClient, err = mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatalf("err connect to mongodb: %v", err)
	}
	MgoDatabase = MgoClient.Database(config.Config.Mongodb.Database)

	return MgoClient
}

// NewMongoDB returns a new layer database instance
func NewMongoDB(order, customer, product, count, history, partner, customerReport, productReport string) *MogoDB {
	return &MogoDB{
		Order:          MgoDatabase.Collection(order),
		Customer:       MgoDatabase.Collection(customer),
		Product:        MgoDatabase.Collection(product),
		Count:          MgoDatabase.Collection(count),
		HistoryPoint:   MgoDatabase.Collection(history),
		Partner:        MgoDatabase.Collection(partner),
		CustomerReport: MgoDatabase.Collection(customerReport),
		ProductReport:  MgoDatabase.Collection(productReport),
	}
}

// MogoDB hold layer struct
type MogoDB struct {
	Order          *mongo.Collection
	Customer       *mongo.Collection
	Product        *mongo.Collection
	Count          *mongo.Collection
	HistoryPoint   *mongo.Collection
	Partner        *mongo.Collection
	CustomerReport *mongo.Collection
	ProductReport  *mongo.Collection
}

// ChangeStream implement
func ChangeStream(collection *mongo.Collection, pipeline interface{}, s func(*mongo.ChangeStream)) {
	ctx := context.Background()
	cso := options.ChangeStream()
	cso.SetFullDocumentBeforeChange(options.WhenAvailable)
	cur, _ := collection.Watch(ctx, pipeline, cso)

	defer cur.Close(ctx)
	//Handling change stream in a cycle
	for {
		select {
		case <-ctx.Done():
			err := cur.Close(ctx)
			if err != nil {
				fmt.Printf("change stream err:", err)
				break
			}
		default:
			s(cur)
		}
	}
}

func (mg *MogoDB) GetOrderByObjectID(key string) (resp *model.Order, err error) {
	id, err := primitive.ObjectIDFromHex(key)
	if err := mg.Order.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&resp); err != nil {
		return &model.Order{}, err
	}

	return resp, nil
}

func (mg *MogoDB) CheckCustomerExist(phone string) bool {
	count, _ := mg.Customer.CountDocuments(context.TODO(), bson.M{"phone_number": phone})
	if count > 0 {
		return true
	}
	return false
}

func (mg *MogoDB) GetSequenceNextValue(seqName string) (int32, error) {
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var result bson.M
	if err := mg.Count.FindOneAndUpdate(context.TODO(), bson.M{
		"_id": seqName,
	}, bson.M{
		"$inc": bson.M{
			"seq": 1,
		},
	}, &opt).Decode(&result); err != nil {
		return -1, err
	}

	seq := result["seq"].(int32)

	return seq, nil
}

func (mg *MogoDB) CreateCustomer(data *model.Customer) error {
	if len(data.Tags) == 0 {
		data.Tags = []string{}
	}
	_, err := mg.Customer.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}

func (mg *MogoDB) CreateHistoryPoint(data *model.HistoryPoint) error {
	_, err := mg.HistoryPoint.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}

func (mg *MogoDB) GetCustomerByPhone(phone string) (resp *model.Customer, err error) {
	if err := mg.Customer.FindOne(context.TODO(), bson.M{"phone_number": phone}).Decode(&resp); err != nil {
		return &model.Customer{}, err
	}

	return resp, nil
}

func (mg *MogoDB) UpdateCustomerByID(data *dto.Customer) error {
	_, err := mg.Customer.UpdateOne(context.TODO(), bson.M{"customer_id": data.CustomerID}, bson.M{
		"$set": data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (mg *MogoDB) GetLoyaltyConfig() (resp *model.LoyaltyConfig, err error) {
	if err := mg.Partner.FindOne(context.TODO(), bson.M{"_id": "LOYALTY_CONFIG"}).Decode(&resp); err != nil {
		return &model.LoyaltyConfig{}, err
	}
	return resp, nil
}

func (mg *MogoDB) CreateProduct(data *model.Product) error {
	_, err := mg.Product.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}

func (mg *MogoDB) UpdateProductByID(data *model.Product) error {
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	_ = mg.Product.FindOneAndUpdate(context.TODO(), bson.M{"sku": data.SKU}, bson.M{
		"$set": data,
	}, &opt)

	return nil
}

func (mg *MogoDB) UpsertCustomerReport(data *model.CustomerReport) error {
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	_ = mg.CustomerReport.FindOneAndUpdate(context.TODO(), bson.M{"date": data.Date, "phone": data.Phone}, bson.M{
		"$set": bson.M{
			"phone": data.Phone,
			"date":  data.Date,
			"name":  data.Name,
			"email": data.Email,
		},
		"$inc": bson.M{
			"total_orders":      data.TotalOrders,
			"success_orders":    data.SuccessOrders,
			"processing_orders": data.ProcessOrders,
			"cancel_orders":     data.CancelOrders,
			"revenue":           data.Revenue,
			"new":               data.New,
		},
	}, &opt)

	return nil
}

func (mg *MogoDB) UpsertProductReport(data *model.ProductReport) error {
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	_ = mg.ProductReport.FindOneAndUpdate(context.TODO(), bson.M{"date": data.Date, "sku": data.SKU, "phone": data.Phone}, bson.M{
		"$set": bson.M{
			"sku":   data.SKU,
			"date":  data.Date,
			"name":  data.Name,
			"phone": data.Phone,
		},
		"$inc": bson.M{
			"total_sales":       data.TotalSales,
			"total_orders":      data.TotalOrders,
			"success_orders":    data.SuccessOrders,
			"processing_orders": data.ProcessOrders,
			"cancel_orders":     data.CancelOrders,
			"revenue":           data.Revenue,
		},
	}, &opt)

	return nil
}
