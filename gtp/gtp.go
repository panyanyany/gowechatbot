package gtp

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/869413421/wechatbot/config"
    "github.com/869413421/wechatbot/pkg/logger"
    "github.com/parnurzeal/gorequest"
    "log"
    "time"
)

const BASEURL = "https://api.openai.com/v1/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
    ID      string                 `json:"id"`
    Object  string                 `json:"object"`
    Created int                    `json:"created"`
    Model   string                 `json:"model"`
    Choices []ChoiceItem           `json:"choices"`
    Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
    Text         string `json:"text"`
    Index        int    `json:"index"`
    Logprobs     int    `json:"logprobs"`
    FinishReason string `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
    Model            string  `json:"model"`
    Prompt           string  `json:"prompt"`
    MaxTokens        uint    `json:"max_tokens"`
    Temperature      float64 `json:"temperature"`
    TopP             int     `json:"top_p"`
    FrequencyPenalty int     `json:"frequency_penalty"`
    PresencePenalty  int     `json:"presence_penalty"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
    cfg := config.LoadConfig()
    requestBody := ChatGPTRequestBody{
        Model:            cfg.Model,
        Prompt:           msg,
        MaxTokens:        cfg.MaxTokens,
        Temperature:      cfg.Temperature,
        TopP:             1,
        FrequencyPenalty: 0,
        PresencePenalty:  0,
    }
    requestData, err := json.Marshal(requestBody)

    if err != nil {
        return "", err
    }

    apiKey := config.LoadConfig().ApiKey

    request := gorequest.New()
    request.Post(BASEURL + "completions").
        Timeout(60 * time.Second).
        Proxy("http://127.0.0.1:7890").
        Send(string(requestData))

    request.AppendHeader("Content-Type", "application/json")
    request.AppendHeader("Authorization", "Bearer "+apiKey)

    resp, body, errs := request.
        EndBytes()

    if len(errs) > 0 {
        return "", errors.New(fmt.Sprintf("gtp api has errors: %v", errs))
    }
    if resp.StatusCode != 200 {
        return "", errors.New(fmt.Sprintf("gtp api status code not equals 200,code is %d, %v", resp.StatusCode, body))
    }

    gptResponseBody := &ChatGPTResponseBody{}
    log.Println(string(body))
    err = json.Unmarshal(body, gptResponseBody)
    if err != nil {
        return "", err
    }

    var reply string
    if len(gptResponseBody.Choices) > 0 {
        reply = gptResponseBody.Choices[0].Text
    }
    logger.Info(fmt.Sprintf("gpt response text: %s ", reply))
    return reply, nil
}
