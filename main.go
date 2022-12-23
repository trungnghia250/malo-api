package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/database"
	campaign_delivery "github.com/trungnghia250/malo-api/service/domain/campaign/delivery"
	campaign_uc "github.com/trungnghia250/malo-api/service/domain/campaign/usecase"
	customer_delivery "github.com/trungnghia250/malo-api/service/domain/customer/delivery"
	customer_uc "github.com/trungnghia250/malo-api/service/domain/customer/usecase"
	customer_group_delivery "github.com/trungnghia250/malo-api/service/domain/customer_group/delivery"
	customer_group_usecase "github.com/trungnghia250/malo-api/service/domain/customer_group/usecase"
	integrate_delivery "github.com/trungnghia250/malo-api/service/domain/integrate/delivery"
	integrate_uc "github.com/trungnghia250/malo-api/service/domain/integrate/usecase"
	loyalty_delivery "github.com/trungnghia250/malo-api/service/domain/loyalty/delivery"
	loyalty_uc "github.com/trungnghia250/malo-api/service/domain/loyalty/usecase"
	order_delivery "github.com/trungnghia250/malo-api/service/domain/order/delivery"
	order_uc "github.com/trungnghia250/malo-api/service/domain/order/usecase"
	product_delivery "github.com/trungnghia250/malo-api/service/domain/product/delivery"
	product_uc "github.com/trungnghia250/malo-api/service/domain/product/usecase"
	"github.com/trungnghia250/malo-api/service/domain/report/delivery"
	report_uc "github.com/trungnghia250/malo-api/service/domain/report/usecase"
	template_delivery "github.com/trungnghia250/malo-api/service/domain/template_message/delivery"
	template_uc "github.com/trungnghia250/malo-api/service/domain/template_message/usecase"
	user_delivery "github.com/trungnghia250/malo-api/service/domain/user/delivery"
	user_uc "github.com/trungnghia250/malo-api/service/domain/user/usecase"
	crm_repo "github.com/trungnghia250/malo-api/service/repo"
	"log"
	"os"
)

func main() {
	schema := config.NewSchema()
	// connect to mongodb
	mongoClient := database.ConnectMongo(schema.Mongodb.User, schema.Mongodb.Pass, schema.Mongodb.Host)
	// Check the connection
	err := mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to mongodb")
	defer func() {
		err = mongoClient.Disconnect(context.Background())
		if err != nil {
			fmt.Printf("failed to disconnect mongodb %v \n", err)
		}
	}()
	if err = database.ConnectAws(); err != nil {
		log.Fatal(err)
	}
	repo := crm_repo.NewRepo(mongoClient)
	//usecase
	customerUseCase := customer_uc.NewCustomerUseCase(repo)
	userUseCase := user_uc.NewUserUseCase(repo)
	productUseCase := product_uc.NewProductUseCase(repo)
	orderUseCase := order_uc.NewOrderUseCase(repo)
	integrateUseCase := integrate_uc.NewIntegrateUseCase(repo)
	campaignUseCase := campaign_uc.NewCampaignUseCase(repo)
	customerGroupUseCase := customer_group_usecase.NewCustomerGroupUseCase(repo)
	templateUseCase := template_uc.NewTemplateUseCase(repo)
	loyaltyUseCase := loyalty_uc.NewLoyaltyUseCase(repo)
	reportUseCase := report_uc.NewReportUseCase(repo)

	//handler
	customerHandler := customer_delivery.NewCustomerHandler(customerUseCase)
	userHandler := user_delivery.NewUserHandler(userUseCase)
	productHandler := product_delivery.NewProductHandler(productUseCase)
	orderHandler := order_delivery.NewOrderHandler(orderUseCase)
	integrateHandler := integrate_delivery.NewIntegrateHandler(integrateUseCase)
	campaignHandler := campaign_delivery.NewCampaignHandler(campaignUseCase)
	customerGroupHandler := customer_group_delivery.NewCustomerGroupHandler(customerGroupUseCase)
	templateHandler := template_delivery.NewTemplateHandler(templateUseCase)
	loyaltyHandler := loyalty_delivery.NewLoyaltyHandler(loyaltyUseCase)
	reportHandler := delivery.NewReportHandler(reportUseCase)

	//router
	router := fiber.New()
	router.Use(cors.New())
	customerHandler.InternalCustomerAPIRoute(router)
	userHandler.InternalUserAPIRoute(router)
	productHandler.InternalProductAPIRoute(router)
	orderHandler.InternalOrderAPIRoute(router)
	integrateHandler.InternalIntegrateAPIRoute(router)
	campaignHandler.InternalCampaignAPIRoute(router)
	customerGroupHandler.InternalCustomerGroupAPIRoute(router)
	templateHandler.InternalTemplateAPIRoute(router)
	loyaltyHandler.InternalLoyaltyAPIRoute(router)
	reportHandler.InternalReportAPIRoute(router)
	router.Use(cors.New())

	//mongoDB := database.NewMongoDB("order", "customer", "product",
	//	"counter", "history_point", "partner", "customer_report", "product_report")
	//loyaltyConfig, _ := mongoDB.GetLoyaltyConfig()
	//var rankConfig []int32
	//for _, rank := range loyaltyConfig.Ranks {
	//	rankConfig = append(rankConfig, rank.MinimumScore)
	//}
	//configRankPoint := model.RankPointConfig{
	//	Point: loyaltyConfig.Formula,
	//	Rank:  rankConfig,
	//}
	//
	//go amqp.OrderStream(mongoClient.Database("malo").Collection("order"))
	//go consumers.RunConsumer(config.Config.RabbitMq.AMQP,
	//	config.Config.RabbitMq.Tags,
	//	config.Config.RabbitMq.Exchange,
	//	"direct",
	//	config.Config.RabbitMq.Queue,
	//	config.Config.RabbitMq.RoutingKey,
	//	mongoDB, configRankPoint)

	port := os.Getenv("PORT")
	if err != nil {
		port = "3000"
	}

	_ = router.Listen(":" + port)
	//_ = router.Listen(":3000")

}
