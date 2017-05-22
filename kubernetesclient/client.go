package kubernetesclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const byNamePath string = "/api/v1/namespaces/%s/%s/%s"

func NewClient(apiURL string, debug bool) *Client {
	client := &Client{
		baseClient: baseClient{
			BaseURL: apiURL,
			debug:   debug,
		},
	}

	client.Pod = newPodClient(client)
	client.Namespace = newNamespaceClient(client)
	client.ReplicationController = newReplicationControllerClient(client)
	client.Service = newServiceClient(client)
	client.Node = newNodeClient(client)

	return client
}

type Client struct {
	baseClient
	Pod                   PodOperations
	Namespace             NamespaceOperations
	ReplicationController ReplicationControllerOperations
	Service               ServiceOperations
	Node                  NodeOperations
}

type baseClient struct {
	BaseURL string
	debug   bool
}

func (c *baseClient) doByName(resourceType string, namespace string, name string, responseObject interface{}) error {
	path := fmt.Sprintf(byNamePath, namespace, resourceType, name)
	err := c.doGet(path, responseObject)
	return err
}

func (c *baseClient) doGet(path string, respObject interface{}) error {
	url := c.BaseURL + path
	return c.doNoBodyRequest("GET", url, respObject)
}

func (c *baseClient) doDelete(path string, respObject interface{}) error {
	url := c.BaseURL + path
	return c.doNoBodyRequest("DELETE", url, respObject)
}

func (c *baseClient) doNoBodyRequest(method string, url string, respObject interface{}) error {
	client := c.newHttpClient()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", GetAuthorizationHeader())

	if c.debug {
		fmt.Println("Request => " + method + " " + url)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return newApiError(resp, url)
	}

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if c.debug {
		fmt.Println("Response <= " + string(byteContent))
	}

	return json.Unmarshal(byteContent, respObject)
}

func (c *baseClient) doPost(path string, inputObject interface{}, respObject interface{}) error {
	return c.doModify(path, "POST", inputObject, respObject)
}

func (c *baseClient) doPut(path string, inputObject interface{}, respObject interface{}) error {
	return c.doModify(path, "PUT", inputObject, respObject)
}

func (c *baseClient) doModify(path string, method string, inputObject interface{}, respObject interface{}) error {
	url := c.BaseURL + path

	var input io.Reader
	bodyContent, err := json.Marshal(inputObject)
	if err != nil {
		return err
	}
	input = bytes.NewBuffer(bodyContent)

	if c.debug {
		fmt.Println(method + " " + url)
		fmt.Println("Request => " + string(bodyContent))
	}

	client := c.newHttpClient()
	req, err := http.NewRequest(method, url, input)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", GetAuthorizationHeader())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return newApiError(resp, url)
	}

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if c.debug {
		fmt.Println("Response <= " + string(byteContent))
	}

	return json.Unmarshal(byteContent, respObject)
}

func (c *baseClient) newHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: GetTLSClientConfig(),
		},
	}
}

type ApiError struct {
	StatusCode int
	Url        string
	Msg        string
	Status     string
	Body       string
}

func (e ApiError) Error() string {
	return e.Msg
}

func newApiError(resp *http.Response, url string) *ApiError {
	contents, err := ioutil.ReadAll(resp.Body)
	var body string
	if err != nil {
		body = "Unreadable body."
	} else {
		body = string(contents)
	}
	formattedMsg := fmt.Sprintf("Bad response from [%s], statusCode [%d]. Status [%s]. Body: [%s]",
		url, resp.StatusCode, resp.Status, body)
	return &ApiError{
		Url:        url,
		Msg:        formattedMsg,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
	}
}
