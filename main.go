package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/rabitmq/handler"
	customer_delivery "github.com/trungnghia250/malo-api/service/domain/customer/delivery"
	customer_uc "github.com/trungnghia250/malo-api/service/domain/customer/usecase"
	integrate_delivery "github.com/trungnghia250/malo-api/service/domain/integrate/delivery"
	integrate_uc "github.com/trungnghia250/malo-api/service/domain/integrate/usecase"
	order_delivery "github.com/trungnghia250/malo-api/service/domain/order/delivery"
	order_uc "github.com/trungnghia250/malo-api/service/domain/order/usecase"
	product_delivery "github.com/trungnghia250/malo-api/service/domain/product/delivery"
	product_uc "github.com/trungnghia250/malo-api/service/domain/product/usecase"
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

	configData, err := json.MarshalIndent(schema, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal indent schema, %v", err)
	}
	fmt.Printf("config: %s \n", configData)

	// connect to queue
	queue := handler.NewSession(&schema.RabbitMq)

	useCase := handler.NewUseCase(mongoClient.Database("malo").Collection("campaign"), queue)
	go useCase.Process()

	repo := crm_repo.NewRepo(mongoClient)
	//usecase
	customerUseCase := customer_uc.NewCustomerUseCase(repo)
	userUseCase := user_uc.NewUserUseCase(repo)
	productUseCase := product_uc.NewProductUseCase(repo)
	orderUseCase := order_uc.NewOrderUseCase(repo)
	integrateUseCase := integrate_uc.NewIntegrateUseCase(repo)

	//handler
	customerHandler := customer_delivery.NewCustomerHandler(customerUseCase)
	userHandler := user_delivery.NewUserHandler(userUseCase)
	productHandler := product_delivery.NewProductHandler(productUseCase)
	orderHandler := order_delivery.NewOrderHandler(orderUseCase)
	integrateHandler := integrate_delivery.NewIntegrateHandler(integrateUseCase)

	//router
	router := fiber.New()
	router.Use(cors.New())
	customerHandler.InternalCustomerAPIRoute(router)
	userHandler.InternalUserAPIRoute(router)
	productHandler.InternalProductAPIRoute(router)
	orderHandler.InternalOrderAPIRoute(router)
	integrateHandler.InternalIntegrateAPIRoute(router)

	router.Use(cors.New())

	port := os.Getenv("PORT")
	if err != nil {
		port = "3000"
	}

	_ = router.Listen(":" + port)
	//_ = router.Listen(":3000")
}
