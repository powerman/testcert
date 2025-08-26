package testcert_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/powerman/check"

	"github.com/powerman/testcert"
)

func Test(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	l, err := new(net.ListenConfig).Listen(t.Context(), "tcp", "127.0.0.1:0")
	t.Nil(err)

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("ok"))
		}),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{testcert.LocalhostCertificate()},
		},
	}
	go srv.ServeTLS(l, "", "")
	defer srv.Close()
	t.Nil(waitTCPPort(l.Addr()))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: testcert.LocalhostCertPool(),
			},
		},
	}
	res, err := client.Get("https://" + l.Addr().String()) //nolint:noctx // Test.
	t.Nil(err)
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	t.Nil(err)
	t.Equal(string(body), "ok")
}

func waitTCPPort(addr fmt.Stringer) error {
	const delay = time.Second / 20
	var dialer net.Dialer
	for {
		conn, err := dialer.Dial("tcp", addr.String())
		if err == nil {
			return conn.Close()
		}
		time.Sleep(delay)
	}
}
