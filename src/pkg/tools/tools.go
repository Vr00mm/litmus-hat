package tools

import (
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Request struct {
	METHOD   string
	ENDPOINT string
	REQUEST  string
}

func formatToken(url string, token string) string {
	tmp := strings.Split(url, "/")
	if tmp[3] == "auth" {
		return "Bearer " + token
	}
	return token
}

func SetLogLevel() {
	logLevel := log.WarnLevel
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatalf("Fatal error: failed to set logLevel\n%s\n", err)
	}
	log.SetLevel(logLevel)
}

func GraphqlRequest(request Request, token string) string {

	client := http.Client{}
	req, err := http.NewRequest(request.METHOD, request.ENDPOINT, strings.NewReader(request.REQUEST))
	if err != nil {
		log.Fatalf("Fatal error: initiate request\n%s\n", err)
	}
	log.Debugf("Request : %s\n", request.REQUEST)

	tokenData := formatToken(request.ENDPOINT, token)

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {tokenData},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Fatal error: cannot execute request\n%s\n", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Debugf("Response Body: %s\n", string(body))

	return string(body)
}
