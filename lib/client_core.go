package lib

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// RequestClient -
type RequestClient struct {
	HostURL     string
	HTTPClient  *http.Client
	Token       string
	ContentType string
}

// NewRequestClient .
func NewRequestClient(ctURL string, ctAPIKey string, skipSSLValidation bool) *RequestClient {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: skipSSLValidation}

	c := RequestClient{
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	c.HostURL = ctURL
	c.Token = ctAPIKey

	return &c
}

// GET - Returns an element.
func (c *RequestClient) GET(urlPath string, returnData interface{}) error {
	if returnData != nil {
		// Ensure the correct returnData was passed in.
		v := reflect.ValueOf(returnData)
		if v.Kind() != reflect.Ptr {
			return errors.New("data must pass a pointer, not a value")
		}
	}

	//out := make([]interface{}, 0)

	pageURL := fmt.Sprintf("%s%s", c.HostURL, urlPath)

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return err
	}

	body, _, _, err := c.doRequest(req)
	if err != nil {
		return err
	}

	//fmt.Println("Body:", string(body))

	//page := make([]interface{}, 0)

	err = json.Unmarshal(body, returnData)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %v", string(body))
	}

	//out = append(out, page...)

	// if returnData != nil {
	// 	// Convert to bytes.
	// 	b, err := json.Marshal(out)
	// 	if err != nil {
	// 		return fmt.Errorf("could not marshal response body: %v", err.Error())
	// 	}

	// 	// Unmarshal back to struct.
	// 	err = json.Unmarshal(b, returnData)
	// 	if err != nil {
	// 		return fmt.Errorf("could not unmarshal response body: %v", string(b))
	// 	}
	// }

	return nil
}

// POST - Create an element.
func (c *RequestClient) POST(urlPath string, sendData interface{}, returnData interface{}) error {
	if returnData != nil {
		// Ensure the correct returnData was passed in.
		v := reflect.ValueOf(returnData)
		if v.Kind() != reflect.Ptr {
			return errors.New("data must pass a pointer, not a value")
		}
	}

	pageURL := fmt.Sprintf("%s%s", c.HostURL, urlPath)

	sendJSON, err := json.Marshal(sendData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", pageURL, bytes.NewReader(sendJSON))
	if err != nil {
		return err
	}

	body, _, _, err := c.doRequest(req)
	if err != nil {
		return err
	}

	//fmt.Println("Body:", string(body))

	err = json.Unmarshal(body, returnData)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %v", string(body))
	}

	return nil
}

// PUT - Update an element.
func (c *RequestClient) PUT(urlPath string, returnData interface{}) error {
	if returnData != nil {
		// Ensure the correct returnData was passed in.
		v := reflect.ValueOf(returnData)
		if v.Kind() != reflect.Ptr {
			return errors.New("data must pass a pointer, not a value")
		}
	}

	pageURL := fmt.Sprintf("%s%s", c.HostURL, urlPath)

	req, err := http.NewRequest("PUT", pageURL, nil)
	if err != nil {
		return err
	}

	body, _, _, err := c.doRequest(req)
	if err != nil {
		return err
	}

	//fmt.Println("Body:", string(body))

	err = json.Unmarshal(body, returnData)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %v", string(body))
	}

	return nil
}

// PATCH - Update an element.
func (c *RequestClient) PATCH(urlPath string, sendData interface{}, returnData interface{}) error {
	if returnData != nil {
		// Ensure the correct returnData was passed in.
		v := reflect.ValueOf(returnData)
		if v.Kind() != reflect.Ptr {
			return errors.New("data must pass a pointer, not a value")
		}
	}

	pageURL := fmt.Sprintf("%s%s", c.HostURL, urlPath)
	if strings.HasPrefix(urlPath, "http") {
		// Just use the URL if it's a full page.
		pageURL = urlPath
	}

	sendJSON, err := json.Marshal(sendData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", pageURL, bytes.NewReader(sendJSON))
	if err != nil {
		return err
	}

	body, _, _, err := c.doRequest(req)
	if err != nil {
		return err
	}

	//fmt.Println("Body:", string(body))

	err = json.Unmarshal(body, returnData)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %v", string(body))
	}

	return nil
}

// Pagination -
type Pagination struct {
	TotalItems int
	TotalPages int
	NextPage   int
}

// NewPagination -
func NewPagination(h http.Header) Pagination {
	totalItems, _ := strconv.Atoi(h.Get("X-Total"))
	totalPages, _ := strconv.Atoi(h.Get("X-Total-Pages"))
	nextPage, _ := strconv.Atoi(h.Get("X-Next-Page"))

	return Pagination{
		TotalItems: totalItems,
		TotalPages: totalPages,
		NextPage:   nextPage,
	}
}

func (c *RequestClient) doRequest(req *http.Request) ([]byte, int, Pagination, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	if len(c.ContentType) > 0 {
		req.Header.Add("Content-Type", c.ContentType)
	} else {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, Pagination{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, Pagination{}, err
	}

	// fmt.Println("HEADERS:", res.Header)
	// fmt.Println("Total:", res.Header.Get("X-Total"))
	// fmt.Println("Pages:", res.Header.Get("X-Total-Pages"))
	// fmt.Println("Next:", res.Header.Get("X-Next-Page"))

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, res.StatusCode, Pagination{}, fmt.Errorf("url: %s, method: %s, status: %d, body: %s", req.URL.String(), req.Method, res.StatusCode, body)
	}

	return body, res.StatusCode, NewPagination(res.Header), err
}
