package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type ProxyAuthInput struct {
	Host     string
	Protocol string
	Path     string
}

type ProxyAuthDataBody struct {
	ProxyId       int    `json:"proxyId"`
	ProxyPassword string `json:"proxyPassword"`
}

type ProxyAuthOutput struct {
	Token string `json:"token"`
}

func AuthProxy(input ProxyAuthInput) (*ProxyAuthOutput, error) {
	godotenv.Load("./../../.env")
	fmt.Println(os.Getenv("PROXY_ID"))
	client := http.Client{
		Transport: http.DefaultTransport,
	}
	u := url.URL{
		Scheme: input.Protocol,
		Host:   input.Host,
		Path:   input.Path,
	}
	proxyId, err := strconv.Atoi(os.Getenv("PROXY_ID"))
	proxyPassword := os.Getenv("PROXY_PASSWORD")
	if err != nil {
		log.Printf("error on proxyId int parsing: %s", err.Error())
	}
	body, err := json.Marshal(&ProxyAuthDataBody{ProxyId: proxyId, ProxyPassword: proxyPassword})
	fmt.Println(string(body))
	if err != nil {
		log.Printf("error on body marshal: %s", err.Error())
	}
	buff := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", u.String(), buff)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("error on auth proxy request: %s", err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error on auth proxy response: %s", err.Error())
	}
	defer res.Body.Close()

	var output *ProxyAuthOutput
	parsedResponse, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error on auth proxy body decode: %s", err.Error())
	}
	err = json.Unmarshal(parsedResponse, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}