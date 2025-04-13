package publicapi

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

const (
	userBaseUrl = "/api/v1/user"
)

func newHTTP(cfg *Config, logger *zap.SugaredLogger) *fiber.App {
	app := fiber.New(fiber.Config{AppName: cfg.App.Name})
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

func registerHTPP(
	app *fiber.App,
	cfg *Config,

) {
	//
	//userHandler := userhandler.New(userSrvc, authClient)
	//authMiddleware := keyauth.New(keyauth.Config{
	//	KeyLookup: "cookie:access_token",
	//	Validator: func(ctx *fiber.Ctx, token string) (bool, error) {
	//		res, err := authClient.VerifyToken(ctx.Context(), authclient.VerifyTokenParams{AccessToken: token})
	//		return res, err
	//	},
	//})
	////user
	//userRouter := app.Group(userBaseUrl)
	//userRouter.Get("/auth/:hash", userHandler.Auth)
	//userRouter.Post("/saveFilter", authMiddleware, userHandler.SaveFilter)
	//userRouter.Get("/filter", authMiddleware, userHandler.GetFilter)
	//userRouter.Post("/updateNotify", authMiddleware, userHandler.UpdateNotify)

	//front
	app.Static("/", cfg.Front.Static)

	////vacancy
	//app.All("/graphql", vacancy.Handler(vacancySrvc, filterSrvc))

}
