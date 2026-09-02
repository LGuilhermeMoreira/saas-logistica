package authorization

import (
	"bytes"
	"delivery/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrForbidden = errors.New("acesso negado pelas políticas de autorização")

type OPAInterface interface {
	Validate(data OPAInput) error
}

type OPA struct {
	address string
}

func NewOPA(env *config.Env) OPAInterface {
	return &OPA{
		address: env.OPA_ADDRESS,
	}
}

func (o *OPA) Validate(data OPAInput) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	payload := map[string]any{
		"input": data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, o.address, bytes.NewReader(body))
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
		return fmt.Errorf("opa returned unexpected status: %d", resp.StatusCode)
	}

	var response struct {
		Result struct {
			Allow bool `json:"allow"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.Result.Allow {
		return ErrForbidden
	}

	return nil
}

type OPAInput struct {
	User struct {
		Role string `json:"role"`
	} `json:"user"`
	Action string `json:"method"`
	Path   string `json:"path"`
}

func (o *OPAInput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}
