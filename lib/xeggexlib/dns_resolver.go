package xeggexlib

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/ncruces/go-dns"
)

type Rotator struct {
	Data []string
	I    int
}

func (r *Rotator) Get() string {

	addr := r.Data[r.I]

	r.I += 1
	if r.I >= len(r.Data) {
		r.I = 0
	}

	return addr

}

func CreateHttpClient() *http.Client {

	// rotate := Rotator{
	// 	Data: []string{
	// 		"104.26.6.43:443",
	// 		"172.67.73.171:443",
	// 		"104.26.7.43:443",
	// 	},
	// }

	resolver, err := dns.NewDoHResolver(
		"https://dns.google/dns-query",
		dns.DoHCache())

	if err != nil {
		panic(err)
	}

	dialer := &net.Dialer{
		Timeout:  time.Minute,
		Resolver: resolver,
		// KeepAlive: 60 * time.Second,
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	// cfg := &tls.Config{}

	// tlsdialer := tls.Dialer{
	// 	Config:    cfg,
	// 	NetDialer: &net.Dialer{},
	// }

	// tlsDialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
	// 	log.Println(addr)
	// 	ip := net.ParseIP(addr)
	// 	if ip.To4() != nil || ip.To16() != nil {
	// 		return tlsdialer.DialContext(ctx, network, addr)
	// 	}

	// 	host, port, err := net.SplitHostPort(addr)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	rsp, err := cdns.Query(ctx, dns.Domain(host), dns.TypeA)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	for _, a := range rsp.Answer {
	// 		addr = a.Data + ":" + port
	// 		break
	// 	}

	// 	log.Println(addr)

	// 	return tlsdialer.DialContext(ctx, network, addr)

	// }

	httpClient := &http.Client{

		Transport: &http.Transport{
			// DialTLSContext: tlsDialContext,
			DialContext:        dialContext,
			DisableCompression: false,
			// DisableKeepAlives:  true,
			MaxIdleConnsPerHost: 1,
			// TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		},
	}

	return httpClient
}
