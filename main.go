package textgenerationapiollama

import (
	"errors"
	"fmt"

	textgenerationapi "github.com/greatoldcactus/textgenerationapi"
)

type TextGenerationAPIOllama struct {
	Url          string
	AnswerTokens int
	Model        string
}

// func (a *TextGenerationAPIOllama) ClearHistory() {
// 	a.History.Clear()
// }

// func (a *TextGenerationAPIOllama) AddMessage(msg textgenerationapi.Message) {
// 	a.History.Add(msg)
// }

// func (a *TextGenerationAPIOllama) GetHistory() (history textgenerationapi.History, err error) {
// 	history = a.History

// 	return
// }

func (a *TextGenerationAPIOllama) Continue(msg textgenerationapi.Message) (result_msg textgenerationapi.Message, err error) {
	history := textgenerationapi.History{}
	msg.AuthorName = "Assistant"
	history.Add(msg)
	result_msg, err = a.GenerateText(history)
	return
}

func (a *TextGenerationAPIOllama) Generate(history textgenerationapi.History) (result_msg textgenerationapi.Message, err error) {
	result_msg, err = a.GenerateText(history)
	return
}

func (a *TextGenerationAPIOllama) Answer(msg textgenerationapi.Message) (result_msg textgenerationapi.Message, err error) {
	history := textgenerationapi.History{}
	history.Add(msg)
	result_msg, err = a.GenerateText(history)
	return
}

func (a *TextGenerationAPIOllama) SetUrl(url string) (err error) {
	a.Url = url
	return
}

func (a *TextGenerationAPIOllama) GetUrl() (url string, err error) {
	url = a.Url
	return
}

var ErrFailedToListModels = errors.New("failed to list models")
var ErrUnknownModel = errors.New("error unknown model")

func (a *TextGenerationAPIOllama) SetModel(model string) error {
	models, err := a.ListModels()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToListModels, err)
	}

	model_ok := false

	for _, model_check := range models {
		if model_check == model {
			model_ok = true
			break
		}
	}

	if !model_ok {
		return fmt.Errorf("%w: %s", ErrUnknownModel, model)
	}

	a.Model = model

	return nil

}

func (a *TextGenerationAPIOllama) GetModel() string {
	return a.Model
}
