package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	TimeSpend   string `json:"time_spend"`
}

type BrasilApi struct {
	Cep       string `json:"cep"`
	State     string `json:"state"`
	City      string `json:"city"`
	Street    string `json:"street"`
	Service   string `json:"service"`
	TimeSpend string `json:"time_spend"`
}

type CEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	APi         string `json:"api"`
	TimeSpend   string `json:"time_spend"`
}

func (response *ViaCep) convertToCEP() CEP {
	return CEP{
		Cep:         response.Cep,
		Logradouro:  response.Logradouro,
		Complemento: response.Complemento,
		Bairro:      response.Bairro,
		Localidade:  response.Localidade,
		Uf:          response.Uf,
		APi:         "ViaCep",
		TimeSpend:   response.TimeSpend,
	}
}

func (response *BrasilApi) convertToCEP() CEP {
	return CEP{
		Cep:         response.Cep,
		Logradouro:  response.Street,
		Complemento: "",
		Bairro:      "",
		Localidade:  response.City,
		Uf:          response.State,
		APi:         "BrasilApi",
		TimeSpend:   response.TimeSpend,
	}
}

func main() {
	c1 := make(chan CEP)
	var cep string = "35660124"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go findByViaCep(ctx, cep, c1)
	go findByBrasilApi(ctx, cep, c1)

	select {
	case cep := <-c1:
		fmt.Printf("received from api %s\n", cep.APi)
		fmt.Printf("Time spend: %s\n", cep.TimeSpend)
		fmt.Printf("CEP: %s\n", cep.Cep)
		fmt.Printf("Logradouro: %s\n", cep.Logradouro)
		fmt.Printf("Complemento: %s\n", cep.Complemento)
		fmt.Printf("Bairro: %s\n", cep.Bairro)
		fmt.Printf("Localidade: %s\n", cep.Localidade)
		fmt.Printf("UF: %s\n", cep.Uf)
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}

func findByViaCep(ctx context.Context, cep string, c1 chan<- CEP) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+cep+"/json", nil)
	duration := time.Since(start)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var c ViaCep
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		panic(err)
	}
	c.TimeSpend = duration.String()
	c1 <- c.convertToCEP()
}

func findByBrasilApi(ctx context.Context, cep string, c1 chan<- CEP) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
	duration := time.Since(start)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var c BrasilApi
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		panic(err)
	}
	c.TimeSpend = duration.String()
	c1 <- c.convertToCEP()
}
