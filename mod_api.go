package textgenerationapiollama

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrBaseUrlNotDefined = errors.New("base URL not defined for API")
var ErrIncorrectResponse = errors.New("incorrect response")
var ErrIncorrectModel = errors.New("incorrect model")
var ErrUnknownMode = errors.New("mode unknown")

func (sa *TextGenerationAPIOllama) CheckModel() error {
	models, err := sa.ListModels()
	if err != nil {
		return fmt.Errorf("unable to get model list from api %w", err)
	}

	model_ok := false

	for _, model := range models {
		if model == sa.Model {
			model_ok = true
			break
		}
	}

	if !model_ok {
		return fmt.Errorf("%w model: %s", ErrIncorrectModel, sa.Model)
	}

	return nil
}

func (sa *TextGenerationAPIOllama) ListModels() (models []string, err error) {
	url := sa.Url + "/api/tags"

	resp, err := http.Get(url)

	if err != nil {
		err = fmt.Errorf("unable to GET %w", err)
		return
	}

	if resp.StatusCode != 200 {

		err = fmt.Errorf("incorrect response from server: %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("unable to read from response body %w", err)
		return
	}

	var models_response struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	err = json.Unmarshal(body, &models_response)

	if err != nil {
		return
	}

	models = []string{}
	for _, model := range models_response.Models {
		models = append(models, model.Name)
	}

	return

}
