package killswitch

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/miekg/dns"
)

// WhoamiDNS return public ip by quering DNS server
func WhoamiDNS() (string, error) {
	var record *dns.TXT

	target := "o-o.myaddr.l.google.com"
	server := "ns1.google.com"

	c := dns.Client{}
	m := dns.Msg{}

	m.SetQuestion(target+".", dns.TypeTXT)
	r, _, err := c.Exchange(&m, server+":53")

	if err != nil {
		return "", err
	}

	if len(r.Answer) == 0 {
		return "", fmt.Errorf("could not find public IP")
	}

	for _, ans := range r.Answer {
		record = ans.(*dns.TXT)
	}

	return strings.TrimSpace(record.Txt[0]), nil
}

// WhoamiWWW return IP by quering http server
func WhoamiWWW() (string, error) {
	client := &http.Client{}
	// Create request
	req, _ := http.NewRequest("GET", "http://myip.country/ip", nil)
	req.Header.Set("User-Agent", "killswitch")
	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		req, _ = http.NewRequest("GET", "http://checkip.amazonaws.com/", nil)
		req.Header.Set("User-Agent", "killswitch")
		resp, err = client.Do(req)
		if err != nil {
			return "", err
		}
	}
	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)
	return string(respBody), nil
}
