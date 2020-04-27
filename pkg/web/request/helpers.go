package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
)

var (
	// ErrInvalidContentType is returned by ParseRequest if it can't unmarshal it into the passed struct
	ErrInvalidContentType = errors.New("Invalid request content type")

	// gorilla/schema decoder is a shared object, as it caches information about structs
	decoder = schema.NewDecoder()
)

func parseIP(ipStr string) net.IP {
	ip := net.ParseIP(ipStr)
	return ip
}

// Parse takes the input body from the passed request and tries to unmarshal it into data
func Parse(r *http.Request, data interface{}) error {
	// Get the contentType for comparisons
	contentType := r.Header.Get("Content-Type")

	// Deterimine the passed ContentType
	if strings.Contains(contentType, "application/json") {
		// It's JSON, so read the body into a variable
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		// And then unmarshal it into the passed interface
		err = json.Unmarshal(body, data)

		return err
	} else if contentType == "" ||
		strings.Contains(contentType, "application/x-www-form-urlencoded") ||
		strings.Contains(contentType, "multipart/form-data") {
		// net/http should be capable of parsing the form data
		err := r.ParseForm()
		if err != nil {
			return err
		}

		// Unmarshal them into the passed interface
		err = decoder.Decode(data, r.PostForm)

		return err
	}

	return ErrInvalidContentType
}

// RemoteAddr extracts the remote address of the request, taking into
// account proxy headers.
func RemoteAddr(r *http.Request) string {
	if prior := r.Header.Get("X-Forwarded-For"); prior != "" {
		proxies := strings.Split(prior, ",")
		if len(proxies) > 0 {
			remoteAddr := strings.Trim(proxies[0], " ")
			if parseIP(remoteAddr) != nil {
				return remoteAddr
			}
		}
	}
	// X-Real-Ip is less supported, but worth checking in the
	// absence of X-Forwarded-For
	if realIP := r.Header.Get("X-Real-Ip"); realIP != "" {
		if parseIP(realIP) != nil {
			return realIP
		}
	}

	return r.RemoteAddr
}

// RemoteIP extracts the remote IP of the request, taking into
// account proxy headers.
func RemoteIP(r *http.Request) string {
	addr := RemoteAddr(r)

	// Try parsing it as "IP:port"
	if ip, _, err := net.SplitHostPort(addr); err == nil {
		return ip
	}

	return addr
}

// GetHTTPHost returns the currenthost in the request
func GetHTTPHost(r *http.Request) string {
	proto := "https"
	if r.TLS == nil {
		proto = "http"
	}
	return fmt.Sprintf("%s://%s", proto, r.Host)
}

// GetHTTPUrl returns the url in the request
func GetHTTPUrl(r *http.Request) string {
	proto := "https"
	if r.TLS == nil {
		proto = "http"
	}
	return fmt.Sprintf("%s://%s%s", proto, r.Host, r.URL)
}
