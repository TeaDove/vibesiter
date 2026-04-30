package llmsupplier

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"time"
	"vibesiter/pkg/llmsupplier/prompts"
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

	prompts *prompts.Prompts
}

func NewSupplier(client *openai.Client) *Supplier {
	r := &Supplier{client: client, model: "deepseek-v4-flash"} // qwen3:8b

	promptsTemplates, err := prompts.NewPrompts()
	if err != nil {
		panic(errors.Wrap(err, "create prompts templates"))
	}

	r.prompts = promptsTemplates

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

	err = json.Unmarshal([]byte(extractJSON(raw)), &output)
	if err != nil {
		return errors.Wrapf(err, "unmarshal: %s", raw)
	}

	err = validators.Validator.Struct(output)
	if err != nil {
		return errors.Wrapf(err, "validate: %s", raw)
	}

	return nil
}

func extractJSON(raw string) string {
	start := strings.Index(raw, "{")

	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return raw
	}

	return raw[start : end+1]
}
