package auth

import "errors"
import "net/http"
import "fmt"

const logOutURL = "https://api.robinhood.com/api-token-logout/"

func RequestLogOut(creds *Credentials) error {
	if creds.token == nil {
		return errors.New("auth: unable to log out, missing token")
	}

	req, err := http.NewRequest("POST", logOutURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Host", "api.robinhood.com")
	req.Header.Add("User-Agent", "go-robinhood")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", *creds.token))

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp.Status)
	return nil
}
