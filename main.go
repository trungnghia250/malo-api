package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/database"
	customer_delivery "github.com/trungnghia250/malo-api/service/domain/customer/delivery"
	customer_uc "github.com/trungnghia250/malo-api/service/domain/customer/usecase"
	order_delivery "github.com/trungnghia250/malo-api/service/domain/order/delivery"
	order_uc "github.com/trungnghia250/malo-api/service/domain/order/usecase"
	product_delivery "github.com/trungnghia250/malo-api/service/domain/product/delivery"
	product_uc "github.com/trungnghia250/malo-api/service/domain/product/usecase"
	user_delivery "github.com/trungnghia250/malo-api/service/domain/user/delivery"
	user_uc "github.com/trungnghia250/malo-api/service/domain/user/usecase"
	crm_repo "github.com/trungnghia250/malo-api/service/repo"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
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

	repo := crm_repo.NewRepo(mongoClient)
	//usecase
	customerUseCase := customer_uc.NewCustomerUseCase(repo)
	userUseCase := user_uc.NewUserUseCase(repo)
	productUseCase := product_uc.NewProductUseCase(repo)
	orderUseCase := order_uc.NewOrderUseCase(repo)

	//handler
	customerHandler := customer_delivery.NewCustomerHandler(customerUseCase)
	userHandler := user_delivery.NewUserHandler(userUseCase)
	productHandler := product_delivery.NewProductHandler(productUseCase)
	orderHandler := order_delivery.NewOrderHandler(orderUseCase)

	//router
	router := fiber.New()
	router.Use(cors.New())
	customerHandler.InternalCustomerAPIRoute(router)
	userHandler.InternalUserAPIRoute(router)
	productHandler.InternalProductAPIRoute(router)
	orderHandler.InternalOrderAPIRoute(router)

	_ = router.Listen(":3000")
}
