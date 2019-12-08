package definitions

import (
	"context"

	"gopkg.in/mgo.v2/bson"

	"service/config/db"
	"service/utils"
	"service/models"

	pb "service/config/api"
)

// GetUser ...
func (s *Server) GetUser(ctx context.Context, in *pb.UserID) (*pb.Res, error) {
	user := models.User{}
	if err := db.Connection.FindById(bson.ObjectIdHex(in.GetId()), &user); err != nil {
		return nil, utils.ErrorHandler(err)
	}

	data, err := utils.Marshal(user); 
	if err != nil {
		return nil, err
	}

	return &pb.Res{
		Success: true,
		Message: "User gotten successfully",
		Data: data,
		Code: 200,
	}, nil	
}

// CreateUser ...
func (s *Server) CreateUser(ctx context.Context, in *pb.UserInfo) (*pb.Res, error) {
	user := models.User{}
	if err := utils.Unmarshal(in, &user); err != nil {
		return nil, err
	}

	if err := db.Connection.Save(&user); err != nil {
		return nil, utils.ErrorHandler(err)
	}
	
	data, err := utils.Marshal(user); 
	if err != nil {
		return nil, err
	}

	return &pb.Res{
		Success: true,
		Message: "User created successfully",
		Data: data,
		Code: 201,
	}, nil
}