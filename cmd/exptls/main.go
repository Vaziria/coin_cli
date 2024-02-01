package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"time"

	"github.com/ncruces/go-dns"
)

type CustomResolver struct {
	*net.Resolver
}

func main() {
	cfg := &tls.Config{
		// InsecureSkipVerify: true,
	}

	resolver, err := dns.NewDoHResolver(
		"https://dns.google/dns-query",
		dns.DoHCache())

	if err != nil {
		log.Println("resolver", err)
	}

	tcpdialer := &net.Dialer{
		Resolver: resolver,
	}

	// tlsdialer := tls.Dialer{
	// 	Config:    cfg,
	// 	NetDialer: tcpdialer,
	// }
	var con *tls.Conn

	for c := range [5]int{} {
		log.Println("open", c)
		con, err = tls.DialWithDialer(tcpdialer, "tcp", "api.xeggex.com:443", cfg)

		// con, err := tlsdialer.Dial("tcp", "api.xeggex.com:443")

		if err != nil {
			log.Println("dial error", err)
		} else {
			defer con.Close()
			break
		}

	}

	request := `GET /api/v2/market/getorderbookbysymbol/VISH_USDT HTTP/1.1
Host: api.xeggex.com
Sec-Ch-Ua: " Not A;Brand";v="99", "Chromium";v="104"
Sec-Ch-Ua-Mobile: ?0
Sec-Ch-Ua-Platform: "Windows"
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Sec-Fetch-Site: none
Sec-Fetch-Mode: navigate
Sec-Fetch-User: ?1
Sec-Fetch-Dest: document
Accept-Encoding: gzip, deflate
Accept-Language: en-US,en;q=0.9

`
	time.Sleep(time.Second * 3)
	state := con.ConnectionState()
	log.Println("state", state.HandshakeComplete, state.ServerName)

	n, err := io.WriteString(con, request)
	if err != nil {
		log.Println("SSL Write error :", err.Error(), n)
	}

	data := []byte{}
	n, err = con.Read(data)
	if err != nil {
		log.Println("SSL Read error : " + err.Error())

		return
	}

	log.Println(n)

	log.Println(string(data))

	// state := con.ConnectionState()
	// fmt.Println("SSL ServerName : " + state.ServerName)
	// fmt.Println("SSL Handshake : ", state.HandshakeComplete)
	// fmt.Println("SSL Mutual : ", state.NegotiatedProtocolIsMutual)

}
