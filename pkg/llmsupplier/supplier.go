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
	"github.com/openai/openai-go/v3/shared"
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

func (r *Supplier) chat(
	ctx context.Context,
	systemPrompt string,
	userPrompt string,
	output any,
) ([]openai.ChatCompletionMessageParamUnion, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfSystem: new(openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{OfString: openai.String(systemPrompt)},
			}),
		},
		{
			OfUser: new(openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{OfString: openai.String(userPrompt)},
			}),
		},
	}

	return r.chatContinue(ctx, messages, output)
}

func (r *Supplier) chatContinue(
	ctx context.Context,
	messages []openai.ChatCompletionMessageParamUnion,
	output any,
) ([]openai.ChatCompletionMessageParamUnion, error) {
	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		return nil, errors.New("output should be a pointer")
	}

	if len(messages) == 0 {
		return nil, errors.New("no messages passed")
	}

	t0 := time.Now()

	resp, err := r.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Temperature:     openai.Float(0),
		Model:           r.model,
		Messages:        messages,
		ReasoningEffort: shared.ReasoningEffortLow,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &shared.ResponseFormatJSONObjectParam{Type: "json_object"},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "llm request")
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("llm response has no choices")
	}

	raw := resp.Choices[0].Message.Content
	messages = append(messages, resp.Choices[0].Message.ToParam())

	zerolog.Ctx(ctx).Debug().
		Str("raw", redact_utils.TrimSized(raw, 100)).
		Str("elapsed", time_utils.RoundDuration(time.Since(t0))).
		Msg("llm.called")

	err = json.Unmarshal([]byte(extractJSON(raw)), &output)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal: %s", raw)
	}

	err = validators.Validator.Struct(output)
	if err != nil {
		return nil, errors.Wrapf(err, "validate: %s", raw)
	}

	return messages, nil
}

func extractJSON(raw string) string {
	start := strings.Index(raw, "{")

	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return raw
	}

	return raw[start : end+1]
}
