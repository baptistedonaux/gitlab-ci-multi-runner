package network

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"gitlab.com/gitlab-org/gitlab-ci-multi-runner/common"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type client struct {
	http.Client
	url        *url.URL
	caFile     string
	skipVerify bool
	updateTime time.Time
}

func (n *client) ensureTlsConfig() {
	// certificate got modified
	if stat, err := os.Stat(n.caFile); err == nil && n.updateTime.Before(stat.ModTime()) {
		n.Transport = nil
	}

	// create or update transport
	if n.Transport == nil {
		n.updateTime = time.Now()
		n.createTransport()
	}
}

func (n *client) createTransport() {
	// create reference TLS config
	tlsConfig := tls.Config{
		MinVersion:         tls.VersionTLS10,
		InsecureSkipVerify: n.skipVerify,
	}

	// load TLS certificate
	if file := n.caFile; file != "" {
		logrus.Debugln("Trying to load", file, "...")

		data, err := ioutil.ReadFile(file)
		if err == nil {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM(data) {
				tlsConfig.RootCAs = pool
			} else {
				logrus.Errorln("Failed to parse PEM in", n.caFile)
			}
		} else {
			if !os.IsNotExist(err) {
				logrus.Errorln("Failed to load", n.caFile, err)
			}
		}
	}

	// create transport
	n.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tlsConfig,
	}
}

func (n *client) do(uri, method string, statusCode int, request interface{}, response interface{}) (int, string) {
	var body []byte

	url, err := n.url.Parse(uri)
	if err != nil {
		return -1, err.Error()
	}

	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return -1, fmt.Sprintf("failed to marshal project object: %v", err)
		}
	}

	req, err := http.NewRequest(method, url.String(), bytes.NewReader(body))
	if err != nil {
		return -1, fmt.Sprintf("failed to create NewRequest: %v", err)
	}

	if request != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	n.ensureTlsConfig()

	res, err := n.Do(req)
	if err != nil {
		return -1, fmt.Sprintf("couldn't execute %v against %s: %v", req.Method, req.URL, err)
	}
	defer res.Body.Close()

	if res.StatusCode == statusCode {
		if response != nil {
			d := json.NewDecoder(res.Body)
			err = d.Decode(response)
			if err != nil {
				return -1, fmt.Sprintf("Error decoding json payload %v", err)
			}
		}
	}

	return res.StatusCode, res.Status
}

func (n *client) fullUrl(uri string, a ...interface{}) string {
	url, err := n.url.Parse(fmt.Sprintf(uri, a...))
	if err != nil {
		return ""
	}
	return url.String()
}

func newClient(config common.RunnerCredentials) (c *client, err error) {
	url, err := url.Parse(strings.TrimRight(config.URL, "/") + "/api/v1/")
	if err != nil {
		return
	}

	c = &client{
		url:        url,
		skipVerify: config.TLSSkipVerify,
		caFile:     config.TLSCAFile,
	}

	if CertificateDirectory != "" && c.caFile == "" {
		hostAndPort := strings.Split(url.Host, ":")
		c.caFile = filepath.Join(CertificateDirectory, hostAndPort[0]+".crt")
	}

	return
}