package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golem/template/gogolem_test"
	"golem/template/roundtrip"
	"io/ioutil"

	"net/http"
)

type RequestBody struct {
	CurrentTotal uint64
}

type ResponseBody struct {
	Message string
}

func init() {
	a := GogolemTestImpl{}
	gogolem_test.SetExportsGolemTemplateApi(a)
}

// total State can be stored in global variables
var total uint64

type GogolemTestImpl struct {
	total uint64
}

// Implementation of the exported interface

func (e GogolemTestImpl) Add(value uint64) {
	total += value
}

func (e GogolemTestImpl) Get() uint64 {
	return total
}

func (e GogolemTestImpl) Publish() gogolem_test.Result[struct{}, string] {
	http.DefaultClient.Transport = roundtrip.WasiHttpTransport{}
	var result gogolem_test.Result[struct{}, string]

	postBody, _ := json.Marshal(RequestBody{
		CurrentTotal: total,
	})
	resp, err := http.Post("http://localhost:9999/post-example", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		result.SetErr(fmt.Sprintln(err))
		return result
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result.SetErr(fmt.Sprintln(err))
		return result
	}

	var response ResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		result.SetErr(fmt.Sprintln(err))
		return result
	}

	fmt.Println(response.Message)

	result.Set(struct{}{})
	return result
}

func (e GogolemTestImpl) Pause() {
	promise := gogolem_test.GolemApiHostGolemCreatePromise()
	gogolem_test.GolemApiHostGolemAwaitPromise(promise)
}

func main() {
}
