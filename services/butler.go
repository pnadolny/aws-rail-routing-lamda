package services

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type Butler struct {
	Authorization string
	Url           string
	Body          *[]byte
	HttpMethod    string
	DistanceFunc  func(data []byte) (string, error)
}

type ButlerResponse struct {
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
	Success bool        `json:"success"`
}

func (butler *Butler) Do() (ButlerResponse, error) {

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(butler.HttpMethod, butler.Url, bytes.NewBuffer(*butler.Body))
	if err != nil {
		return ButlerResponse{Errors: err.Error()}, err
	}
	req.Header.Set("Authorization", butler.Authorization)
	req.Header.Set("Content-Type", "application/json")
	res, err := netClient.Do(req)
	if err != nil {
		return ButlerResponse{Errors: err.Error()}, nil
	}
	defer res.Body.Close()
	httpResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ButlerResponse{Errors: err.Error()}, nil
	}
	distance, err := butler.DistanceFunc(httpResponse)
	if err != nil {
		return ButlerResponse{Errors: err.Error()}, nil
	}
	return ButlerResponse{Data: distance, Success: true}, nil

}
