package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"main/services"
	"net/http"
	"os"
)

type RouteReportRequest struct {
	Station []Stop  `json:"Stops"`
	Options Options `json:"Options"`
	Report  string  `json:"Report"`
}

type Stop struct {
	Format   string `json:"Format"`
	Name     string `json:"Name"`
	Railroad string `json:"Railroad"`
}

type Options struct {
	RoutingPreference      string `json:"RoutingPreference"`
	TerminalSwitching      bool   `json:"TerminalSwitching"`
	AmtrakRoutes           bool   `json:"AmtrakRoutes"`
	IntermodalOnlyStations bool   `json:"IntermodalOnlyStations"`
	DistUnit               string `json:"DistUnit"`
}

type Response struct {
	Report       Report `json:"Report"`
	ErrorDetails string `json:"ErrorDetails"`
}

type Report struct {
	Lines []Line `json:"Lines"`
}

type Line struct {
	RailRoad string `json:"Railroad"`
	Distance string `json:"Distance"`
}

const (
	railApi = "https://pcmrail.alk.com/REST/v24.1//Service.svc/route/Report"
)

func handleRequest(ctx context.Context, routeReportRequest RouteReportRequest) (string, error) {
	reportBytes, err := json.Marshal(routeReportRequest)
	if err != nil {
		return "error", nil
	}
	butler := services.Butler{
		Url:           railApi,
		HttpMethod:    http.MethodPost,
		Authorization: os.Getenv("RAIL_KEY"),
		Body:          &reportBytes,
		DistanceFunc: func(data []byte) (string, error) {
			var response Response
			if err := json.Unmarshal(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")), &response); err != nil {
				return "", err
			}
			fmt.Println(response)
			if response.ErrorDetails != "" {
				return "", errors.New(response.ErrorDetails)
			}
			total := ""
			for _, v := range response.Report.Lines {
				total = v.Distance
			}
			return total, nil
		},
	}
	br, err := butler.Do()
	brBytes, _ := json.Marshal(br)
	return string(brBytes), err

}

func main() {
	lambda.Start(handleRequest)
}
