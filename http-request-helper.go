package echor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPError interface
type HTTPError interface {
	error
	Status() int
}

// HTTPStatusError struct
type HTTPStatusError struct {
	Code int
	Err  error
}

// Error func
func (se HTTPStatusError) Error() string {
	return se.Err.Error()
}

// Status func
func (se HTTPStatusError) Status() int {
	return se.Code
}

// DoJSONReq func
func DoJSONReq(mode string, url string, reqInput interface{}, header map[string]string, output interface{}) error {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	var postData io.Reader = nil
	var method = http.MethodGet
	if mode == "POST" || mode == "PUT" || mode == "PATCH" || mode == "OPTIONS" {
		if reqInput != nil {
			jsonStr, err := json.Marshal(reqInput)
			if err != nil {
				return err
			}
			postData = bytes.NewBuffer(jsonStr)
		}
		if mode == "POST" {
			method = http.MethodPost
		} else if mode == "PUT" {
			method = http.MethodPut
		} else if mode == "PATCH" {
			method = http.MethodPatch
		} else if mode == "OPTIONS" {
			method = http.MethodOptions
		}
	}
	req, err := http.NewRequest(method, url, postData)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	res, reqErr := client.Do(req)
	defer res.Body.Close()

	if reqErr != nil {
		return reqErr
	}
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		message := struct {
			Message string `json:"message"`
		}{
			Message: "",
		}
		jsonErr := json.Unmarshal(body, &message)
		if jsonErr != nil {
			return HTTPStatusError{res.StatusCode, errors.New(string(body))}
		}
		return HTTPStatusError{res.StatusCode, errors.New(message.Message)}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	jsonErr := json.Unmarshal(body, &output)
	if jsonErr != nil {
		return jsonErr
	}
	return nil
}

// DoJSONPost func
func DoJSONPost(url string, reqInput interface{}, header map[string]string, output interface{}) error {
	return DoJSONReq("POST", url, reqInput, header, output)
}

// DoJSONGet func
func DoJSONGet(url string, header map[string]string, output interface{}) error {
	return DoJSONReq("GET", url, nil, header, output)
}
