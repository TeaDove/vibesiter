package webpresentation

import (
	"fmt"
	"strings"
	"time"
	"vibesiter/pkg/dto"
	"vibesiter/pkg/managerservice"
	"vibesiter/pkg/validators"

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
		StructValidator: &fiber_utils.StructValidator{Validator: validators.Validator},
	})
	app.Use(recover2.New(recover2.Config{EnableStackTrace: true}))
	app.Use(fiber_utils.MiddlewareLogger())
	app.Use(fiber_utils.MiddlewareCtxTimeout(3 * time.Minute))
	app.Use(cors.New(cors.ConfigDefault))

	appGroup := app.Group("/app")
	appGroup.Post("/", func(c fiber.Ctx) error {
		body, err := fiber_utils.BindJSON[GenerateAppRequest](c)
		if err != nil {
			return err
		}

		resp, err := r.managerService.GenerateApp(c.Context(), body.UserPrompt)
		if err != nil {
			return errors.Wrap(err, "generate app")
		}

		return c.JSON(resp)
	})
	appGroup.Post("/:slug/raw", func(c fiber.Ctx) error {
		slug := c.Params("slug")

		req, err := fiber_utils.BindJSON[dto.HTTPRequest](c)
		if err != nil {
			return err
		}

		resp, err := r.managerService.HandleHTTP(c.Context(), slug, &req)
		if err != nil {
			return errors.Wrap(err, "handle http request")
		}

		c = c.Status(resp.Status).Type(resp.ContentType)

		if resp.ContentType == fiber.MIMEApplicationJSON {
			return c.JSON(resp)
		}

		switch body := resp.Body.(type) { // TODO вынести в сервис
		case string:
			return c.SendString(body)
		case []byte:
			return c.Send(body)
		default:
			return c.JSON(body)
		}
	})
	appGroup.Get("/:slug", func(c fiber.Ctx) error {
		return c.Redirect().To(fmt.Sprintf("/app/%s/index.html", c.Params("slug")))
	})
	appGroup.Get("/:slug/", func(c fiber.Ctx) error {
		return c.Redirect().To(fmt.Sprintf("/app/%s/index.html", c.Params("slug")))
	})
	appGroup.Get("/:slug/:path", func(c fiber.Ctx) error {
		slug := c.Params("slug")

		path := c.Params("path")
		if path == "" {
			return c.Redirect().To(fmt.Sprintf("/app/%s/index.html", slug))
		}

		resp, err := r.managerService.Serve(c.Context(), slug, path)
		if err != nil {
			return errors.Wrap(err, "handle http request")
		}

		return c.Type(strings.Split(path, ".")[len(strings.Split(path, "."))-1]).Send(resp)
	})

	return app
}

type GenerateAppRequest struct {
	UserPrompt string `json:"userPrompt" validate:"required"`
}
