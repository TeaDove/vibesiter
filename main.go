package main

import (
	"vibesiter/pkg/llmsupplier"
	"vibesiter/pkg/managerrepo"
	"vibesiter/pkg/managerservice"
	"vibesiter/pkg/settings"
	"vibesiter/pkg/webpresentation"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/teadove/teasutils/service_utils/db_utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func build() (*fiber.App, error) {
	db, err := gorm.Open(sqlite.Open(".data/db.sqlite"),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
			TranslateError: true,
			Logger:         db_utils.ZerologAdapter{},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "open gorm db")
	}

	err = db.AutoMigrate(new(managerrepo.Application), new(managerrepo.ApplicationKV))
	if err != nil {
		return nil, errors.Wrap(err, "auto migrate")
	}

	msgRepo := managerrepo.New(db)

	// option.WithBaseURL("http://localhost:11434/v1"),
	//		option.WithAPIKey("ollama"),
	llmSupplier := llmsupplier.NewSupplier(new(openai.NewClient(
		option.WithBaseURL(settings.Settings.AIURL),
		option.WithAPIKey(settings.Settings.AIAPIKEY),
	)))

	chatService := managerservice.NewService(msgRepo, llmSupplier)

	presentation := webpresentation.NewPresentation(chatService)

	return presentation.BuildApp(), nil
}

func main() {
	app, err := build()
	if err != nil {
		panic(err)
	}

	err = app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
