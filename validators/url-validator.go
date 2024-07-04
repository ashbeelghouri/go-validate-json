package validators

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func StatusCodeCheck(i interface{}, attr map[string]interface{}) error {
	if err := IsString(i, attr); err != nil {
		return err
	}

	if err := IsURL(i, attr); err != nil {
		return err
	}

	defaultTimeout := time.Duration(5)
	if to, ok := attr["timeout"].(float64); ok {
		defaultTimeout = time.Duration(to)
	}

	defaultStatusCodeCheck := http.StatusOK
	if sc, ok := attr["status_code"].(float64); ok {
		defaultStatusCodeCheck = int(sc)
	}
	url := i.(string)
	client := &http.Client{
		Timeout: defaultTimeout * time.Second,
	}
	resp, err := client.Head(url)
	if err != nil {
		return errors.New("failed to perform HEAD request")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("warning: body-closer", err)
		}
	}(resp.Body)
	if resp.StatusCode != defaultStatusCodeCheck {
		msg := fmt.Sprintf("url is not throwing status code: %d, it is sending %d", defaultStatusCodeCheck, resp.StatusCode)
		return errors.New(msg)
	}
	return nil
}
