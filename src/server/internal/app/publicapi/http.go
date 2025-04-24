package publicapi

import (
	"clusterlizer/internal/handler/publicapi"
	documentsrvc "clusterlizer/internal/service/document"
	requestsrvc "clusterlizer/internal/service/request"
	s3srvc "clusterlizer/internal/service/s3"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

const (
	publicApiPrefix = "/api/v1/"
)

func registerHTPP(
	cfg *Config,
	log *zap.SugaredLogger,
	documentSrvc documentsrvc.Service,
	requestSrvc requestsrvc.Service,
	s3Srvc s3srvc.Service,
	//authMiddleware := keyauth.New(keyauth.Config{

) *fiber.App {
	app := newFiber(cfg, log)
	//
	publicApi := publicapi.New(
		log,
		documentSrvc,
		requestSrvc,
		s3Srvc,
	)
	//authMiddleware := keyauth.New(keyauth.Config{
	//	KeyLookup: "cookie:access_token",
	//	Validator: func(ctx *fiber.Ctx, token string) (bool, error) {
	//		res, err := authClient.VerifyToken(ctx.Context(), authclient.VerifyTokenParams{AccessToken: token})
	//		return res, err
	//	},
	//})
	////user
	publicApiGroup := app.Group(publicApiPrefix)

	publicApiGroup.Post("/uploadFiles", publicApi.UploadFiles)
	publicApiGroup.Get("/getClusterizations/:id", publicApi.GetClusterizations)
	//publicApiGroup.Get("/getCurrentQueue/:uuid", publicApi.GetCurrentQueue)

	//front
	app.Static("/", cfg.Front.Static)

	return app
}

func newFiber(cfg *Config, logger *zap.SugaredLogger) *fiber.App {
	app := fiber.New(
		fiber.Config{
			AppName:   cfg.App.Name,
			BodyLimit: 100 * 1024 * 1024, // this is the default limit of 4MB
		})
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Desugar(),
	}))
	// app.user(log)
	app.Use(compress.New())

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins: "*",
		//AllowCredentials: true,
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	return app
}
