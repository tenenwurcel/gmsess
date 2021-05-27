package main

import (
	_entity "gmsess/api/entity"
	_handler "gmsess/api/handler"
	_repo "gmsess/api/repository"
	"gmsess/config"
	"gmsess/proto"
	"gmsess/utils"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	utils.SetupCypher()
	utils.SetupVerifier()
	config.SetupRedis()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	redisRepository := _repo.NewRedisRepository(config.GetRedisCli())

	sessionEntity := _entity.NewSesssionEntity(redisRepository)
	sessionHandler := _handler.NewSessionHandler(sessionEntity)

	authorizationEntity := _entity.NewAuthorizationEntity(redisRepository)
	authorizationHandler := _handler.NewAuthorizationHandler(authorizationEntity)

	grpcServer := grpc.NewServer()

	proto.RegisterAuthenticatorServer(grpcServer, sessionHandler)
	proto.RegisterAuthorizationServer(grpcServer, authorizationHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
