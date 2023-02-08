package api

import (
	"fmt"
	db "github.com/bugrakocabay/dummy-bank/db/sqlc"
	"github.com/bugrakocabay/dummy-bank/token"
	"github.com/bugrakocabay/dummy-bank/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	authRoutes.GET("/users/:username", server.getUser)

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PATCH("/accounts", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	authRoutes.GET("/transfers/:id", server.getTransfer)
	authRoutes.POST("/transfers", server.createTransfer)
	authRoutes.GET("/transfers", server.listTransfers)

	server.router = router
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
