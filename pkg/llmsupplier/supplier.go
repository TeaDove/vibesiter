package llmsupplier

import (
	"context"
	_ "embed"
	"encoding/json"
	"reflect"
	"text/template"
	"time"
	"vibesiter/pkg/validators"

	"github.com/cockroachdb/errors"
	"github.com/openai/openai-go/v3"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/redact_utils"
	"github.com/teadove/teasutils/utils/time_utils"
)

type Supplier struct {
	client *openai.Client
	model  string

	systemPromptTemplateDesignSite *template.Template
	systemPromptTemplateGenerate   *template.Template
	userPromptTemplateGenerate     *template.Template
	systemPromptTemplateHTTP       *template.Template
	userPromptTemplateHTTP         *template.Template
}

var (
	//go:embed design.gohtml
	systemPromptDesignSite string
	//go:embed generate_system.gohtml
	systemPromptGenerate string
	//go:embed generate_user.gohtml
	userPromptGenerate string
	//go:embed http_system.gohtml
	systemPromptHTTP string
	//go:embed http_user.gohtml
	userPromptHTTP string
)

func NewSupplier(client *openai.Client) *Supplier {
	r := &Supplier{client: client, model: "deepseek-v4-flash"} // qwen3:8b

	var err error
	// TODO move to separate system
	r.systemPromptTemplateDesignSite, err = template.New("example").Parse(systemPromptDesignSite)
	if err != nil {
		panic(errors.Wrap(err, "parse system prompt template"))
	}

	r.systemPromptTemplateGenerate, err = template.New("example").Parse(systemPromptGenerate)
	if err != nil {
		panic(errors.Wrap(err, "parse system prompt template"))
	}

	r.userPromptTemplateGenerate, err = template.New("example").Parse(userPromptGenerate)
	if err != nil {
		panic(errors.Wrap(err, "parse system prompt template"))
	}

	r.systemPromptTemplateHTTP, err = template.New("example").Parse(systemPromptHTTP)
	if err != nil {
		panic(errors.Wrap(err, "parse system prompt template"))
	}

	r.userPromptTemplateHTTP, err = template.New("example").Parse(userPromptHTTP)
	if err != nil {
		panic(errors.Wrap(err, "parse system prompt template"))
	}

	return r
}

func (r *Supplier) chat(ctx context.Context, systemPrompt string, userPrompt string, output any) error {
	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		return errors.New("output should be a pointer")
	}

	t0 := time.Now()

	resp, err := r.client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Temperature: openai.Float(0),
			Model:       r.model,
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfSystem: new(openai.ChatCompletionSystemMessageParam{
						Content: openai.ChatCompletionSystemMessageParamContentUnion{
							OfString: openai.String(systemPrompt),
						},
					}),
				},
				{
					OfUser: new(openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{OfString: openai.String(userPrompt)},
					}),
				},
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "llm request")
	}

	if len(resp.Choices) == 0 {
		return errors.New("llm response has no choices")
	}

	raw := resp.Choices[0].Message.Content

	zerolog.Ctx(ctx).Debug().
		Str("user_prompt", redact_utils.TrimSized(userPrompt, 100)).
		Str("raw", redact_utils.TrimSized(raw, 100)).
		Str("elapsed", time_utils.RoundDuration(time.Since(t0))).
		Msg("llm.called")

	err = json.Unmarshal([]byte(raw), &output)
	if err != nil {
		return errors.Wrapf(err, "unmarshal: %s", raw)
	}

	err = validators.Validator.Struct(output)
	if err != nil {
		return errors.Wrapf(err, "validate: %s", raw)
	}

	return nil
}
