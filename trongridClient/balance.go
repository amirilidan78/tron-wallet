package trongridClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"tronWallet/enums"
	"tronWallet/httpClient"
)

type GetAccountRequestBody struct {
	Address string `json:"address"`
	Visible bool   `json:"visible"`
}

type GetAccountResponseBody struct {
	Address               string `json:"address"`
	Balance               uint64 `json:"balance"`
	CreateTime            uint64 `json:"create_time"`
	LatestOprationTime    uint64 `json:"latest_opration_time"`
	LatestConsumeFreeTime uint64 `json:"latest_consume_free_time"`
}

func GetAddressBalance(network enums.Network, address string) (GetAccountResponseBody, error) {

	url := string(network) + "/wallet/getaccount"

	requestBody := GetAccountRequestBody{
		Address: address,
		Visible: true,
	}

	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	responseBody := GetAccountResponseBody{}

	httpResponse, _, statusCode, err := httpClient.NewHttpClient().HttpPost(url, requestBody, header)

	if statusCode != http.StatusOK {
		return responseBody, errors.New(fmt.Sprintf("http status code is not 200 it is %d", statusCode))
	}

	if err != nil {
		return responseBody, err
	}

	err = json.Unmarshal(httpResponse, &responseBody)
	if err != nil {
		return responseBody, err
	}

	return responseBody, nil
}
