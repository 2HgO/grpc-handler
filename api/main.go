package main

import (
	"context"
	"log"

	pb "api/config/api"
	"api/models"
	"api/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	conn, err := grpc.Dial("localhost:55045", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserClient(conn)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(utils.PanicHandler())
	
	r.POST("/", func(c *gin.Context) {
		user := models.User{}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(utils.ResponsePayload(nil, nil, err))
			return
		}

		req := pb.UserInfo{}
		utils.Unmarshal(user, &req)

		res, err := client.CreateUser(context.Background(), &req)
		c.JSON(utils.ResponsePayload(models.User{}, res, err))
	})

	r.GET("/", func(c *gin.Context) {
		id := struct {
			ID bson.ObjectId `json:"id" binding:"required"`
		}{}

		if err := c.ShouldBindJSON(&id); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(utils.ResponsePayload(nil, nil, err))
			return
		}

		req := pb.UserID{}
		utils.Unmarshal(id, &req)

		var user models.User
		res, err := client.GetUser(context.Background(), &req)
		c.JSON(utils.ResponsePayload(&user, res, err))
	})

	if err := r.Run(":55055"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
