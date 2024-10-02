package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/santhosh3/ECOM/Config"
	"github.com/santhosh3/ECOM/models"
	pb "github.com/santhosh3/ECOM/proto" // Replace with the correct module path
	"github.com/santhosh3/ECOM/services/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// server is used to implement proto.UserServiceServer and proto.ProductServiceServer
type server struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedProductServiceServer
	DB *gorm.DB
}

// GetUser implements proto.UserServiceServer
func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	token, err := auth.ValidateJWT(req.Token, []byte(config.Envs.AccessJWTSecret))
	if err != nil || !token.Valid {
		return nil, status.Error(codes.PermissionDenied, "Invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	str, ok := claims["UserId"].(string)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Invalid claims in token")
	}

	userId, err := strconv.Atoi(str)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	var user models.User
	if err := s.DB.First(&user, userId).Error; err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	if !user.Status {
		return nil, status.Error(codes.PermissionDenied, "User is inactive")
	}

	return &pb.UserResponse{
		Id:           int32(user.ID),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
		PhoneNumber:  user.PhoneNumber,
	}, nil
}

// GetProduct implements proto.ProductServiceServer
func (s *server) GetProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	productId := req.Id

	var product models.Product
	if err := s.DB.First(&product, productId).Error; err != nil {
		return nil, status.Error(codes.NotFound, "Product not found")
	}

	return &pb.ProductResponse{
		Quantity: int32(product.Quantity),
		Price: int32(product.Price),
	}, nil
}

// StartGRPCServer starts the gRPC server and registers services
func StartGRPCServer(db *gorm.DB) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Envs.GrpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{DB: db})
	pb.RegisterProductServiceServer(s, &server{DB: db}) // Register the Product service as well

	msg := fmt.Sprintf("gRPC server is running on port %s", config.Envs.GrpcPort)
	log.Println(msg)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
