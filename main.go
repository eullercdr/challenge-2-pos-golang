package main

import (
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
}

type BrasilApi struct {
	Cep     string `json:"cep"`
	State   string `json:"state"`
	City    string `json:"city"`
	Street  string `json:"street"`
	Service string `json:"service"`
}

type CEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	APi         string `json:"api"`
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
	}
}

func main() {
	c1 := make(chan CEP)
	var cep string = "35660124"

	go findByViaCep(cep, c1)
	go findByBrasilApi(cep, c1)

	select {
	case cep := <-c1:
		fmt.Printf("received from api %s\n", cep.APi)
		fmt.Printf("CEP: %s\n", cep.Cep)
		fmt.Printf("Logradouro: %s\n", cep.Logradouro)
		fmt.Printf("Complemento: %s\n", cep.Complemento)
		fmt.Printf("Bairro: %s\n", cep.Bairro)
		fmt.Printf("Localidade: %s\n", cep.Localidade)
		fmt.Printf("UF: %s\n", cep.Uf)
	case <-time.After(time.Second):
		println("timeout")
	}
}

func findByViaCep(cep string, c1 chan<- CEP) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var c ViaCep
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		panic(err)
	}
	c1 <- c.convertToCEP()
}

func findByBrasilApi(cep string, c1 chan<- CEP) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var c BrasilApi
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		panic(err)
	}
	c1 <- c.convertToCEP()
}
