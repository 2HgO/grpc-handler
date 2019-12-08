package main

import (
	"context"
	"log"
	"reflect"
	"strings"

	pb "api/config/api"
	"api/models"
	"api/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/grpc"
	"gopkg.in/go-playground/validator.v9"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			var name string
			if tag, ok := fld.Tag.Lookup("form"); ok {
				name = strings.SplitN(tag, ",", 2)[0]
			} else {
				name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			}
			return name
		})
	}

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
