package soap

import (
	"github.com/gin-gonic/gin"
	soaphandler "github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/soap/handler"
)

func SetupSOAPRoutes(router *gin.Engine, movieHandler *soaphandler.MovieSOAPHandler) {
	router.GET("/soap/movie.wsdl", movieHandler.ServeWSDL)
	router.POST("/soap/movie", movieHandler.Handle)
}
