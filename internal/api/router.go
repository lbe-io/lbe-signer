package api

import (
	"io/fs"
	"lbe_crypto_signer/internal/api/wallet"
	"lbe_crypto_signer/internal/service"
	"lbe_crypto_signer/static"
	"net/http"

	_ "lbe_crypto_signer/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swag init -g .\cmd\active\main.go -o ./docs --parseDependency --parseInternal

const (
	PrefixRouter = "/lbe-signer-api"
)

func NewGinRouter(config *service.ActiveConfig) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// web服务静态资源目录
	distFS, err := fs.Sub(static.Static, "web")
	if err != nil {
		panic(err)
	}
	// 将 /static 路径映射到 distFS
	r.StaticFS("/static", http.FS(distFS))

	r.StaticFileFS("/favicon.ico", "/web/favicon.ico", http.FS(distFS))
	// 将根路径重定向到 /static/index.html
	// r.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	// })

	r.Use(gin.Recovery(), CorsMiddleware())

	service.Init(config)

	api := r.Group(PrefixRouter)
	{
		if config.ApiConf.Api.Swagger {
			swagHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
			if swagHandler != nil {
				api.GET("/doc/swagger/*any", swagHandler)
			}
		}

		f := NewMessageHandler()
		api.HEAD("/ping", f.Ping)
		api.GET("/version", f.Version)

		keysApi := wallet.NewKeysHandler()
		accountRouter := api.Group("/key")
		{
			accountRouter.POST("/addMnemonic", keysApi.KeysAddMnemonic)
			accountRouter.POST("/list", keysApi.KeysList)
			accountRouter.POST("/addAccounts", keysApi.AddAccounts)
			accountRouter.POST("/addDepositAddr", keysApi.AddDepositAddr)
			accountRouter.POST("/listDepositAddr", keysApi.ListDepositAddr)
		}
		txApi := wallet.NewTransactionHandler()
		txRouter := api.Group("/tx")
		{
			txRouter.POST("/sign", txApi.Sign)
		}
	}

	return r
}
