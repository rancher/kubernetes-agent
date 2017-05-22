package kubernetesclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	caLocation = "/etc/kubernetes/ssl/ca.pem"
)

var (
	token  string
	caData []byte
)

func init() {
	if err := func() error {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Failed to read token from stdin: %v", err)
		}
		token = strings.TrimSpace(string(bytes))
		if token == "" {
			return errors.New("No token passed in from stdin")
		}

		caData, err = ioutil.ReadFile(caLocation)
		if err != nil {
			return fmt.Errorf("Failed to read CA cert %s: %v", caLocation, err)
		}
		return nil
	}(); err != nil {
		log.Fatal(err)
	}
}

func GetAuthorizationHeader() string {
	return fmt.Sprintf("Bearer %s", token)
}

func GetTLSClientConfig() *tls.Config {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)
	return &tls.Config{
		RootCAs: certPool,
	}
}
