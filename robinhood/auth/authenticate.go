package auth

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const authURL = "https://api.robinhood.com/api-token-auth/"

// Credentials are the login credentials required to start a robinhood session
type Credentials struct {
	username, password string
	mfaCode            *string
}

// NewCredentials returns a Credentials holding the given username and password
func NewCredentials(username, password string) *Credentials {
	return &Credentials{
		username: username,
		password: password,
	}
}

func (a *Credentials) addMultiFactorCode(code string) {
	a.mfaCode = &code
}

func (a *Credentials) generatePostForm() io.Reader {
	d := make(url.Values)
	d.Add("username", a.username)
	d.Add("password", a.password)

	if a.mfaCode != nil {
		d.Add("mfa_code", *a.mfaCode)
	}

	return strings.NewReader(d.Encode())
}

type authResp struct {
	T       *string  `json:"token"`
	MFA     *bool    `json:"mfa_required"`
	MFAType string   `json:"mfa_type"`
	Errors  []string `json:"non_field_errors"`
}

// RequestToken uses a set of account Credentials to request an authorization token from
// the Robinhood API. In the event that the Robinhood account has multi-factor authentication
// activated, RequestToken will prompt for an MFA code.
func RequestToken(creds *Credentials) (string, error) {
	resp, err := requestAuthenticate(creds)
	if err != nil {
		return "", err
	}

	if resp.MFA != nil && *resp.MFA == true {
		code, err := getMFACode(resp.MFAType)
		if err != nil {
			return "", err
		}
		creds.addMultiFactorCode(code)
		resp, err = requestAuthenticate(creds)
		if err != nil {
			return "", err
		}
	}

	if resp.T != nil {
		return *resp.T, nil
	}

	return "", errors.New("failed to retrieve auth token")
}

func requestAuthenticate(creds *Credentials) (*authResp, error) {
	req, err := http.NewRequest("POST", authURL, creds.generatePostForm())
	if err != nil {
		return nil, err
	}

	req.Header.Add("Host", "api.robinhood.com")
	req.Header.Add("User-Agent", "go-robinhood")
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

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	a := new(authResp)
	if err := json.Unmarshal(b, a); err != nil {
		return nil, err
	}

	return a, nil
}

func getMFACode(mfaType string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("You will recieve a two factor authentication code via %s\nEnter your code: ", mfaType)
	code, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return code, nil
}
