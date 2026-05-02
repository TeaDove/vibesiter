package llmsupplier

import (
	"context"
	"vibesiter/pkg/dto"
	"vibesiter/pkg/kvrepo"
	"vibesiter/pkg/managerrepo"

	"github.com/cockroachdb/errors"
	"github.com/openai/openai-go/v3"
	"github.com/teadove/teasutils/utils/test_utils"
)

func (r *Supplier) HTTP(
	ctx context.Context,
	app *managerrepo.Application,
	req *dto.HTTPRequest,
) (dto.LLMHTTPResponse, []openai.ChatCompletionMessageParamUnion, error) {
	systemPrompt, err := r.prompts.Render("http_system", map[string]any{"Slug": app.Slug})
	if err != nil {
		return dto.LLMHTTPResponse{}, nil, errors.Wrap(err, "execute design site")
	}

	userPrompt, err := r.prompts.Render("http_user",
		map[string]any{
			"Slug":       app.Slug,
			"UserPrompt": app.UserPrompt,
			"MetaJson":   map[string]any{"slug": app.Slug, "title": app.Title, "description": app.Description},
			"Design":     app.Design,
			"Request":    req,
		},
	)
	if err != nil {
		return dto.LLMHTTPResponse{}, nil, errors.Wrap(err, "execute design site")
	}

	var output dto.LLMHTTPResponse

	messages, err := r.chat(ctx, systemPrompt, userPrompt, &output)
	if err != nil {
		return dto.LLMHTTPResponse{}, nil, errors.Wrap(err, "chat")
	}

	return output, messages, nil
}

func (r *Supplier) HTTPContinue(
	ctx context.Context,
	messages []openai.ChatCompletionMessageParamUnion,
	kvs []kvrepo.ApplicationKV,
) (dto.LLMHTTPResponse, []openai.ChatCompletionMessageParamUnion, error) {
	userPrompt, err := r.prompts.Render("http_actions_user", map[string]any{"KVs": kvs})
	if err != nil {
		return dto.LLMHTTPResponse{}, nil, errors.Wrap(err, "render")
	}

	messages = append(messages, openai.ChatCompletionMessageParamUnion{
		OfUser: new(openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{OfString: openai.String(userPrompt)},
		}),
	})
	test_utils.Pprint(messages)

	var output dto.LLMHTTPResponse

	messages, err = r.chatContinue(ctx, messages, &output)
	if err != nil {
		return dto.LLMHTTPResponse{}, nil, errors.Wrap(err, "chat continue")
	}

	return output, messages, nil
}
