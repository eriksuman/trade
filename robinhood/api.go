package robinhood

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const APIURL = "https://api.robinhood.com/"

func getWithAuthorization(creds *Credentials, url string) ([]byte, error) {
	if creds.token == nil {
		if err := RequestLogIn(creds); err != nil {
			return nil, err
		}
		defer RequestLogOut(creds)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Host", "api.robinhood.com")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", string(*creds.token)))
	req.Header.Add("Accept", "application/json")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("http error: %s", resp.Status)
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(msg)
		}
		return nil, errors.New(msg + fmt.Sprintf(", Body: %s", string(b)))
	}

	return ioutil.ReadAll(resp.Body)
}

func postWithAuthorization(creds *Credentials, url, form string) ([]byte, error) {
	if creds.token == nil {
		if err := RequestLogIn(creds); err != nil {
			return nil, err
		}
		defer RequestLogOut(creds)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(form))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Host", "api.robinhood.com")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", string(*creds.token)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("http error: %s", resp.Status)
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New(msg)
		}
		return nil, errors.New(msg + fmt.Sprintf(", Body: %s", string(b)))
	}

	return ioutil.ReadAll(resp.Body)
}
