package robinhood

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

var authURL = APIURL + "api-token-auth/"

type Token string

// Credentials are the login credentials required to start a robinhood session
type Credentials struct {
	username, password string
	mfaCode            *string
	token              *Token
}

// NewCredentials returns a Credentials holding the given username and password
func NewCredentials(username, password string) *Credentials {
	return &Credentials{
		username: username,
		password: password,
	}
}

func (c *Credentials) addToken(t Token) {
	if c.token != nil {
		if err := RequestLogOut(c); err != nil {
			fmt.Printf("failed to log out token for %s\n", c.username)
		}
	}

	c.token = &t
}

func (c *Credentials) addMultiFactorCode(code string) {
	c.mfaCode = &code
}

func (c *Credentials) generateAuthPostForm() io.Reader {
	d := make(url.Values)
	d.Add("username", c.username)
	d.Add("password", c.password)

	if c.mfaCode != nil {
		d.Add("mfa_code", *c.mfaCode)
	}

	return strings.NewReader(d.Encode())
}

type authResp struct {
	T       *Token   `json:"token"`
	MFA     *bool    `json:"mfa_required"`
	MFAType string   `json:"mfa_type"`
	Errors  []string `json:"non_field_errors"`
}

// RequestLogIn uses the fields of a Credentials object to request log in
// to a Robinhood account. Upon successful login, the creds will have a
// token for a log in session on the Robinhood server.
func RequestLogIn(creds *Credentials) error {
	resp, err := requestAuthenticate(creds)
	if err != nil {
		return err
	}

	if resp.MFA != nil && *resp.MFA == true {
		code, err := getMFACode(resp.MFAType)
		if err != nil {
			return err
		}
		creds.addMultiFactorCode(code)
		resp, err = requestAuthenticate(creds)
		if err != nil {
			return err
		}
	}

	if resp.T != nil {
		creds.addToken(*resp.T)
		return nil
	}

	return errors.New("failed to retrieve auth token")
}

func requestAuthenticate(creds *Credentials) (*authResp, error) {
	req, err := http.NewRequest("POST", authURL, creds.generateAuthPostForm())
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
