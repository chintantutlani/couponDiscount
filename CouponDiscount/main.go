package main

import (
	"context"
	"fmt"
	"log"
	"monk_commerce/configs"
	controller "monk_commerce/controllers"
	services "monk_commerce/services"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx         context.Context
	mongoclient *mongo.Client
	err         error
	cc          *controller.CouponController
	cs          *services.CouponServices
)

func init() {

	cfg := configs.LoadConfig()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	uri := os.Getenv("SERVER_URI")

	mongoconn := options.Client().ApplyURI(uri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	db := mongoclient.Database(cfg.Database.DB)

	coupon := db.Collection(cfg.Database.Couponcollection)
	cs = services.NewCouponService(coupon, ctx)
	cc = controller.NewCouponController(cs)

}

func main() {
	server := gin.Default()
	defer mongoclient.Disconnect(ctx)

	// config := config.LoadConfig()

	if err != nil {
		log.Println("unable to load config", err)
		return
	}
	basepath := server.Group("/v1")
	cc.RegisterCouponRoutes(basepath)

	err := server.Run("localhost:8080")
	if err != nil {
		log.Fatal("not able to start server ,", err)
	}

}
