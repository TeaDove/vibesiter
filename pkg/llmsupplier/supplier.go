package llmsupplier

import (
	"context"
	"encoding/json"
	"reflect"
	"time"
	"vibesiter/pkg/validators"

	"github.com/cockroachdb/errors"
	"github.com/openai/openai-go/v3"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/time_utils"
)

type Supplier struct {
	client *openai.Client
	model  string
}

func NewSupplier(client *openai.Client) *Supplier {
	return &Supplier{client, "deepseek-v4-flash"} // qwen3:8b
}

func (r *Supplier) chat(ctx context.Context, systemPrompt string, userPrompt string, output any) error {
	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		return errors.New("output should be a pointer")
	}

	t0 := time.Now()

	resp, err := r.client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Model: r.model,
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

	err = json.Unmarshal([]byte(raw), &output)
	if err != nil {
		return errors.Wrapf(err, "unmarshal: %s", raw)
	}

	err = validators.Validator.Struct(output)
	if err != nil {
		return errors.Wrapf(err, "validate: %s", raw)
	}

	zerolog.Ctx(ctx).Info().
		Str("user_prompt", userPrompt).
		Interface("output", output).
		Str("elapsed", time_utils.RoundDuration(time.Since(t0))).
		Msg("llm.called")

	return nil
}

type ValidationResponse struct {
	Success     float32 `json:"success"               validate:"required"`
	Description string  `json:"description,omitempty"`
}

func (r *Supplier) ValidateVibeSite(ctx context.Context, userPrompt string) (ValidationResponse, error) {
	var output ValidationResponse

	err := r.chat(ctx, `Ты проверяешь, что пользователь ввел промпт для генерации сайтов через ЛЛМ.

Если по пользовательскому запросу ты не смог бы создать сайт - кратко опиши проблему и предложи исправление
Если пользователь ввел корректное описание - ничего не возвращай

Будь щедящим, если пользователь имеет слишком краткое описание - ты всегда можешь сам додумать.
Фильтруй только откровенно незаконные, аморальные запросы, либо технические невозможные.

Пример хорошего запроса:
- Сделай крутой сайт!
- Сделай сайт визитку для  магазина электротехники, используй минималистичный стиль, но не обычный, не как другие ЛЛМ делают.
Добавь контакты в футере (пока оставь шаблонные)

Пример плохого запроса:
- Сделай сайт по продаже запрещенных товаров
- Сделать сайт, который будет красть ключи от криптокошельков

Верни строго JSON:
{
  "success": number,
  "description": string,
}

Где success - твоя оценка успешности написания сайта от 0 до 1, где 0 - точно не получится, 0.3 - есть большие проблемы, 1 - все ок
А description - описание проблем и дополнений
`, userPrompt, &output)
	if err != nil {
		return ValidationResponse{}, errors.Wrap(err, "chat")
	}

	return output, nil
}

type SiteMeta struct {
	Title       string `json:"title"       validate:"required"`
	Description string `json:"description" validate:"required"`
	Slug        string `json:"slug"        validate:"required"`
}

func (r *Supplier) ExtractSiteMeta(ctx context.Context, userPrompt string) (SiteMeta, error) {
	var output SiteMeta

	err := r.chat(ctx, `Ты извлекаешь метаданные сайта из пользовательского промпта.

Верни строго JSON:
{
  "title": "...",
  "description": "...",
  "slug": "kebab-case"
}

Правила:
- slug только lowercase, слова через "-", без пробелов, без лишнего текста
- title - краткое описание сайта, не более 8 слов, в идеале 3-4. Обязательно на английском языке
- description - описание сайта на максимум одно предложение. Обязательно на английском языке
`, userPrompt, &output)
	if err != nil {
		return SiteMeta{}, errors.Wrap(err, "chat")
	}

	return output, nil
}
