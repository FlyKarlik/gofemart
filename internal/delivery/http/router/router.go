package router

import (
	"net/http/pprof"

	"github.com/FlyKarlik/gofemart/internal/delivery/http/handler"
	"github.com/FlyKarlik/gofemart/internal/delivery/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/FlyKarlik/gofemart/docs"
)

type HTTPRouter struct {
	middleware *middleware.Middleware
	handler    *handler.Handler
}

func New(middleware *middleware.Middleware, handler *handler.Handler) *HTTPRouter {
	return &HTTPRouter{
		middleware: middleware,
		handler:    handler,
	}
}

func (h *HTTPRouter) InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", h.handler.Ping)
	registerPprof(router)

	api := router.Group("api", h.middleware.JSONMiddleware())
	{
		h.registerUserRoutes(api)
	}

	return router
}

func (h *HTTPRouter) registerUserRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("user")
	{
		userGroup.POST("/register", h.handler.RegisterUser)
		userGroup.POST("/login", h.handler.LoginUser)

		ordersGroup := userGroup.Group("orders", h.middleware.Identity)
		{
			ordersGroup.POST("/", h.handler.CreateOrder)
			ordersGroup.GET("/", h.handler.GetUserOrders)
		}

		balanceGroup := userGroup.Group("balance", h.middleware.Identity)
		{
			balanceGroup.GET("/", h.handler.GetUserBalance)
			balanceGroup.POST("/withdraw", h.handler.WithdrawUserBalance)
		}

		withdrawalsGroup := userGroup.Group("withdrawals", h.middleware.Identity)
		{
			withdrawalsGroup.GET("/", h.handler.GetUserWithdrawals)
		}

	}
}

func registerPprof(router *gin.Engine) {
	pprofGroup := router.Group("/debug/pprof")
	{
		router.GET("/debug/vars", gin.WrapH(pprof.Handler("vars")))
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}
