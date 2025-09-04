package controller

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/knadh/koanf/v2"
)

type ControllerConfig struct {
	appConfig *koanf.Koanf
	awsConfig *aws.Config
}

type ProblemDetails struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func InitControllerConfig(k *koanf.Koanf, a *aws.Config) *ControllerConfig {
	return &ControllerConfig{appConfig: k, awsConfig: a}
}

func writeProblemDetails(w http.ResponseWriter, status int, title, detail, problemType, instance string) {
	problem := ProblemDetails{
		Type:     problemType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem)
}
