package authorization

import (
	"auth/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrForbidden = errors.New("acesso negado pelas políticas de autorização")

type OPAInterface interface {
	Validate(data any) error
}

type OPA struct {
	address string
}

func NewOPA(env *config.Env) *OPA {
	return &OPA{
		address: env.OPA_ADDRESS,
	}
}

func (o *OPA) Validate(data any) error {
	client := &http.Client{
		Timeout: time.Second * 2,
	}

	payload := map[string]any{
		"input": data,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, o.address, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("opa retornou status inesperado: %d", resp.StatusCode)
	}

	var opaResponse struct {
		Result bool `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&opaResponse); err != nil {
		return err
	}

	if !opaResponse.Result {
		return ErrForbidden
	}

	return nil
}
