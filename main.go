package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/domain/customer/delivery"
	customer_uc "github.com/trungnghia250/malo-api/service/domain/customer/usecase"
	delivery2 "github.com/trungnghia250/malo-api/service/domain/user/delivery"
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

	//handler
	customerHandler := delivery.NewCustomerHandler(customerUseCase)
	userHandler := delivery2.NewUserHandler(userUseCase)

	//router
	router := fiber.New()
	router.Use(cors.New())
	customerHandler.InternalCustomerAPIRoute(router)
	userHandler.InternalUserAPIRoute(router)

	_ = router.Listen(":3000")
}
