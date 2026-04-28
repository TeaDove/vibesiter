package llmsupplier

import (
	"context"

	"github.com/cockroachdb/errors"
)

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
