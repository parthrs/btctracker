package httputils

import (
	"io/ioutil"
	"net/http"
)

// HttpGet is a utility function that makes a GET request to
// the given url and returns the response bytes and an additional error.
// It also accepts headers from the user (variadic function)
func HttpGet(url string, headers ...map[string]string) (response []byte, err error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	// Add headers to the request if present
	if len(headers) >= 1 {
		for _, header := range headers {
			// Unpack the key, value of the header using range
			for k, v := range header {
				request.Header.Add(k, v)
			}
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	return
}
