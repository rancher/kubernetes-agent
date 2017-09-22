package kubernetesclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

const (
	caLocation    = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	tokenLocation = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

var (
	token  []byte
	caData []byte
)

func Init() error {
	var err error
	caData, err = ioutil.ReadFile(caLocation)
	if err != nil {
		return fmt.Errorf("Failed to read CA cert %s: %v", caLocation, err)
	}

	token, err = ioutil.ReadFile(tokenLocation)
	if err != nil {
		return fmt.Errorf("Failed to read token %s: %v", tokenLocation, err)
	}

	return nil
}

func GetAuthorizationHeader() string {
	return fmt.Sprintf("Bearer %s", string(token))
}

func GetTLSClientConfig() *tls.Config {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)
	return &tls.Config{
		RootCAs: certPool,
	}
}
