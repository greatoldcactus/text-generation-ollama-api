package textgenerationapiollama

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SETTING_ANSWER_TOKENS = "ANSWER_TOKENS"
	SETTING_URL           = "URL"
	SETTING_MODEL         = "MODEL"
	SETTING_LIST_MODELS   = "LIST_MODELS"
)

func (a *TextGenerationAPIOllama) SetSetting(setting string, val string) error {
	switch setting {
	case SETTING_URL:
		a.Url = val
	case SETTING_ANSWER_TOKENS:
		tokens, err := strconv.Atoi(val)
		if err != nil {
			return ErrIntConvertFailed
		}
		a.AnswerTokens = tokens
	}

	return nil

}
func (a *TextGenerationAPIOllama) GetSetting(setting string) (result string, err error) {
	switch setting {
	case SETTING_URL:
		result = a.Url
	case SETTING_ANSWER_TOKENS:
		result = strconv.Itoa(a.AnswerTokens)
	case SETTING_MODEL:
		result = a.Model
		return
	case SETTING_LIST_MODELS:

		models, err_list := a.ListModels()
		if err_list != nil {
			err = fmt.Errorf("failed to list models: %w", err_list)
			return
		}
		result = strings.Join(models, ",")
		return
	}

	return
}
