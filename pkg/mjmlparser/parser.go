package mjmlparser

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type MJMLResponse struct {
	*MJMLReplyResponse
	*MJMLErrorResponse
}

type MJMLReplyResponse struct {
	Errors      []interface{} `json:"errors"`
	HTML        string        `json:"html"`
	MJML        string        `json:"mjml"`
	MJMLVersion string        `json:"mjml_version"`
}

type MJMLErrorResponse struct {
	Message   string     `json:"message"`
	RequestID string     `json:"request_id"`
	StartedAt *time.Time `json:"started_at"`
}

// ParserMJMLtoHTML parses mjml strings
func ParserMJMLtoHTML(mjmlStr string) *MJMLResponse {
	const MJMLApiURL = "https://api.mjml.io/v1/render"

	data := map[string]string{
		"mjml": mjmlStr,
	}
	j, _ := json.Marshal(data)

	statusCode, _, body, err := HTTPClient(http.MethodPost, MJMLApiURL, []byte(j), nil)
	if err != nil {
		log.Fatal(err)
	}
	response := &MJMLResponse{}

	if statusCode != http.StatusOK {
		errRes := &MJMLErrorResponse{}
		json.Unmarshal(body, errRes)

		response.MJMLErrorResponse = errRes

		return response
	}

	okRes := &MJMLReplyResponse{}

	fmt.Println(body)

	json.Unmarshal(body, okRes)

	response.MJMLReplyResponse = okRes

	return response
}
