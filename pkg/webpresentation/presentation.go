package webpresentation

import (
	"time"
	"ws-lan-chat/pkg/managerservice"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	recover2 "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/teadove/teasutils/fiber_utils"
)

type Presentation struct {
	managerService *managerservice.Service
}

func NewPresentation(managerService *managerservice.Service) *Presentation {
	return &Presentation{managerService}
}

func (r *Presentation) BuildApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Immutable:       true,
		ErrorHandler:    fiber_utils.ErrHandler(),
		StructValidator: fiber_utils.NewDefaultStructValidator(),
	})
	app.Use(recover2.New(recover2.Config{EnableStackTrace: true}))
	app.Use(fiber_utils.MiddlewareLogger())
	app.Use(fiber_utils.MiddlewareCtxTimeout(29 * time.Second))
	app.Use(cors.New(cors.ConfigDefault))

	appGroup := app.Group("/app")
	appGroup.Post("/", func(c fiber.Ctx) error {
		body, err := fiber_utils.BindJSON[GenerateAppRequest](c)
		if err != nil {
			return err
		}

		resp, err := r.managerService.GenerateApp(c, body.UserPrompt)
		if err != nil {
			return errors.Wrap(err, "generate app")
		}

		return c.JSON(resp)
	})

	return app
}

type GenerateAppRequest struct {
	UserPrompt string `json:"userPrompt" validate:"required"`
}
