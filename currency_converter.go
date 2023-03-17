package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CurrencyConverter struct {
	addr   string
	apiKey string
}

func NewCurrencyConverter(addr string, apiKey string) *CurrencyConverter {
	return &CurrencyConverter{
		addr:   addr,
		apiKey: apiKey,
	}
}

func (c *CurrencyConverter) Convert(amount int, baseCurrency string, resultCurrency string) (float64, error) {
	url := fmt.Sprintf(
		"%s/convert?from=%s&to=%s&amount=%d",
		c.addr, baseCurrency, resultCurrency, amount,
	)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("apikey", c.apiKey)

	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}

	type Response struct {
		Result float64 `json:"result"`
	}

	var result Response
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Result, nil
}

func (c *CurrencyConverter) GetCurrencies() (map[string]string, error) {
	url := c.addr + "/symbols"

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	type Response struct {
		Symbols map[string]string `json:"symbols"`
	}

	result := Response{Symbols: make(map[string]string)}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Symbols, nil
}
