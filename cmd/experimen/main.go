package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/pdcgo/common_conf/pdc_common"
)

func resolve(u string) {
	dialer := &net.Dialer{
		Timeout: 60 * time.Second,
	}
	rawConn, err := dialer.Dial("tcp", u)
	if err != nil {
		pdc_common.ReportError(err)
		return
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
		KeyLogWriter:       os.Stdout,
		VerifyConnection: func(cs tls.ConnectionState) error {
			opts := x509.VerifyOptions{
				DNSName:       cs.ServerName,
				Intermediates: x509.NewCertPool(),
				KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			}
			for _, cert := range cs.PeerCertificates[1:] {
				opts.Intermediates.AddCert(cert)
			}
			_, err := cs.PeerCertificates[0].Verify(opts)

			if err != nil {
				pdc_common.ReportError(err)
			}

			return err
		},
	}
	conn := tls.Client(rawConn, config)
	err = conn.Handshake()
	if err != nil {
		pdc_common.ReportError(err)
	}
	fmt.Println(u, err)
	conn.Close()
}

func main() {

	resolve("172.67.73.171:443")
	return

	log.SetFlags(log.Lshortfile)

	conf := &tls.Config{
		InsecureSkipVerify: true,
		KeyLogWriter:       os.Stdout,
	}

	conn, err := tls.Dial("tcp", "172.67.73.171:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	println(string(buf[:n]))

}

// func main() {
// 	log.SetFlags(log.Lshortfile)

// 	cdns := doh.Use(doh.CloudflareProvider, doh.GoogleProvider)
// 	host := "api.xeggex.com"

// 	rsp, err := cdns.Query(context.Background(), dns.Domain(host), dns.TypeA)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// addr := ""
// 	for _, a := range rsp.Answer {
// 		// addr = a.Data
// 		log.Println("asdasdasd", a.Data)
// 	}

// 	// log.Println("final", addr)
// 	// connectTls()

// 	// tlsConfig := &tls.Config{}

// 	// tr := &http.Transport{
// 	// 	TLSClientConfig: tlsConfig,
// 	// }
// 	// client := &http.Client{Transport: tr}

// 	// for range [10]int{} {
// 	// 	_, err := client.Get("https://api.xeggex.com/api/v2/market/getorderbookbysymbol/VISH_USDT")
// 	// 	log.Println(err)
// 	// }

// 	Request()

// }

// type Rotator struct {
// 	Data []string
// 	I    int
// }

// func (r *Rotator) Get() string {

// 	addr := r.Data[r.I]

// 	r.I += 1
// 	if r.I >= len(r.Data) {
// 		r.I = 0
// 	}

// 	return addr

// }

// func Request() {
// 	rotate := Rotator{
// 		Data: []string{
// 			"104.26.6.43:443",
// 			"172.67.73.171:443",
// 			"104.26.7.43:443",
// 		},
// 	}

// 	dialer := &net.Dialer{
// 		Timeout: time.Minute,

// 		KeepAlive: 60 * time.Second,
// 		// DualStack: true, // this is deprecated as of go 1.16
// 	}

// 	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
// 		ip := net.ParseIP(addr)
// 		if ip.To4() != nil || ip.To16() != nil {
// 			return dialer.DialContext(ctx, network, addr)
// 		}

// 		addr = rotate.Get()

// 		return dialer.DialContext(ctx, network, addr)
// 	}

// 	httpClient := &http.Client{

// 		Transport: &http.Transport{
// 			// DialTLSContext: tlsDialContext,
// 			DialContext:        dialContext,
// 			DisableCompression: false,
// 			// DisableKeepAlives:   true,
// 			MaxIdleConnsPerHost: 1,
// 			// TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
// 		},
// 	}

// 	for range [100]int{} {
// 		res, err := httpClient.Get("https://api.xeggex.com/api/v2/market/getorderbookbysymbol/VISH_USDT")
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}

// 		res.Body.Close()

// 		log.Println("asd")

// 	}

// }

// func connectTls() {
// 	log.SetFlags(log.Lshortfile)

// 	conf := &tls.Config{
// 		InsecureSkipVerify: true,
// 		MaxVersion:         tls.VersionTLS13,
// 	}

// 	conn, err := tls.Dial("tcp", "172.67.73.171:443", conf)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	n, err := conn.Write([]byte("hello\n"))
// 	if err != nil {
// 		log.Println(n, err)
// 		return
// 	}

// 	buf := make([]byte, 100)
// 	n, err = conn.Read(buf)
// 	if err != nil {
// 		log.Println(n, err)
// 		return
// 	}

// 	println(string(buf[:n]))
// }
