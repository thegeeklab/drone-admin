package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/jackspirou/syscerts"
	"golang.org/x/oauth2"
)

func New(server string, token string) (drone.Client, error) {
	s, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	if len(s.Scheme) == 0 {
		s.Scheme = "http"
	}

	// attempt to find system CA certs
	certs := syscerts.SystemRootsPool()
	tlsConfig := &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: false,
	}

	oauth := new(oauth2.Config)
	authenticator := oauth.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: token,
		},
	)

	authenticator.Timeout, _ = time.ParseDuration("60s")

	trans, _ := authenticator.Transport.(*oauth2.Transport)
	trans.Base = &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
	}

	return drone.NewClient(s.String(), authenticator), nil
}
