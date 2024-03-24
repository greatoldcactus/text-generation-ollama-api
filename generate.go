package textgenerationapiollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/greatoldcactus/textgenerationapi"
)

func (sa *TextGenerationAPIOllama) GenerateText(history textgenerationapi.History) (message textgenerationapi.Message, err error) {
	if sa.Url == "" {
		err = ErrBaseUrlNotDefined
		return
	}

	err = sa.CheckModel()
	if err != nil {
		return
	}

	type Options struct {
		Num_predict int `json:"num_predict"`
		Num_thread  int `json:"num_thread"`
	}

	tokens := sa.AnswerTokens
	if tokens == 0 {
		tokens = 50
	}

	type history_message_node struct {
		Role    string `json:"role"`
		Message string `json:"content"`
	}

	final_history := []history_message_node{}
	for _, history_node := range history.Messages {
		final_name := history_node.AuthorName
		switch final_name {
		case "Assistant", "System", "User":
		default:
			final_name = "Assistant"
		}
		final_history = append(final_history, history_message_node{final_name, history_node.Message})
	}

	payload := struct {
		Model    string                 `json:"model"`
		Messages []history_message_node `json:"messages"`
		Stream   bool                   `json:"stream"`
		Options  Options                `json:"options"`
	}{
		Model:    sa.Model,
		Messages: final_history,
		Stream:   false,
		Options: Options{
			Num_predict: int(tokens),
			Num_thread:  2,
		},
	}

	var json_payload []byte
	json_payload, err = json.Marshal(payload)
	if err != nil {
		return
	}

	url := sa.Url + "/api/chat"

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_payload))

	if err != nil {
		err = fmt.Errorf("unable to create request %w", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	if resp.StatusCode != 200 {

		body, err_read := io.ReadAll(resp.Body)
		if err_read == nil {
			err = fmt.Errorf("incorrect response from server: %d, response: [====\n%s\n====]", resp.StatusCode, body)
			return
		}

		err = fmt.Errorf("incorrect response from server: %d", resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("unable to read body after request: %v, error: %w", resp.Body, err)
		return
	}

	var response struct {
		Response history_message_node `json:"message"`
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		err = fmt.Errorf("unable to unmarshall response body %w", err)
		return
	}

	text_answer := response.Response.Message

	message = textgenerationapi.Message{
		Message: text_answer,
	}

	return

}
