package ssl

import (
	"crypto/tls"
	"errors"
	"net"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

type fakeListener struct {
	addr net.Addr
}

func (l *fakeListener) Accept() (net.Conn, error) { return nil, errors.New("not implemented") }
func (l *fakeListener) Close() error              { return nil }
func (l *fakeListener) Addr() net.Addr            { return l.addr }

func TestListenAndServeTLSDefaults(t *testing.T) {
	listener := &fakeListener{
		addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 443},
	}
	var gotNetwork string
	var gotAddr string
	listenCalled := false
	serveCalled := false

	listen := func(network, addr string, config *tls.Config) (net.Listener, error) {
		listenCalled = true
		gotNetwork = network
		gotAddr = addr
		if config == nil {
			t.Fatal("expected tls config")
		}
		return listener, nil
	}
	serve := func(l net.Listener, handler nethttp.Handler) error {
		serveCalled = true
		if l != listener {
			t.Fatalf("expected listener to be passed to serve")
		}
		if handler == nil {
			t.Fatalf("expected handler to be passed to serve")
		}
		return nil
	}

	err := listenAndServeTLS("", nethttp.NewServeMux(), &tls.Config{}, listen, serve)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !listenCalled {
		t.Fatalf("expected listen to be called")
	}
	if !serveCalled {
		t.Fatalf("expected serve to be called")
	}
	if gotNetwork != "tcp" {
		t.Fatalf("expected network tcp, got %s", gotNetwork)
	}
	if gotAddr != ":https" {
		t.Fatalf("expected default address :https, got %s", gotAddr)
	}
}

func TestListenAndServeTLSListenError(t *testing.T) {
	listen := func(network, addr string, config *tls.Config) (net.Listener, error) {
		return nil, errors.New("listen failed")
	}
	serve := func(l net.Listener, handler nethttp.Handler) error {
		t.Fatalf("serve should not be called on listen error")
		return nil
	}

	err := listenAndServeTLS("127.0.0.1:8443", nethttp.NewServeMux(), &tls.Config{}, listen, serve)
	if err == nil {
		t.Fatalf("expected error from listen")
	}
}

func TestVerifyOUsHandler(t *testing.T) {
	handler := verifyOUsHandler([]string{"db"}, func(r *nethttp.Request, validOUs []string) error {
		return errors.New("unauthorized")
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodGet, "http://example.com", nil)
	handler(recorder, req, nil)

	if recorder.Code != nethttp.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", nethttp.StatusUnauthorized, recorder.Code)
	}
}
