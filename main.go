package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CEP struct {
	Cep string `json:"cep"`
}

func main() {
	c1 := make(chan CEP)
	c2 := make(chan CEP)
	var cep string = "35660124"

	go func() {
		viacep, err := ViaCep(cep)
		if err != nil {
			fmt.Println(viacep)
		}
		c1 <- *viacep
	}()

	go func() {
		brasilApi, err := BrasilApi(cep)
		if err != nil {
			fmt.Println(err)
		}
		c2 <- *brasilApi
	}()

	for {
		select {
		case msg := <-c1:
			fmt.Printf("received from ViaCep - %+v\n", msg.Cep)
		case msg := <-c2:
			fmt.Printf("received from BrasilApi - %+v\n", msg.Cep)
		case <-time.After(8 * time.Second):
			println("timeout")
		}
	}
}

func ViaCep(cep string) (*CEP, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var c CEP
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func BrasilApi(cep string) (*CEP, error) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1" + cep)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var c CEP
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
