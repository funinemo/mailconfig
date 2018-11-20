package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type PostData struct {
	EMailAddress string `xml:"Request>EMailAddress"`
}
type MailData struct {
	Localpart string
	Subdomain string
}

func mytest() (string, string) {
	dataXML := `<?xml version="1.0" encoding="utf-8"?><Autodiscover xmlns="http://schemas.microsoft.com/exchange/autodiscover/outlook/requestschema/2006"><Request><EMailAddress>dummy@al.em-net.ne.jp</EMailAddress><AcceptableResponseSchema>http://schemas.microsoft.com/exchange/autodiscover/outlook/responseschema/2006a</AcceptableResponseSchema></Request></Autodiscover>`
	var p PostData
	if err := xml.Unmarshal([]byte(dataXML), &p); err != nil {
		log.Fatal()
	}
	at := strings.Index(p.EMailAddress, "@")
	localpart := p.EMailAddress[:at]
	subdomain := p.EMailAddress[at+1 : at+3]
	fmt.Printf("l/s:%s,%s\n", localpart, subdomain)
	return localpart, subdomain
}
func main() {
	http.HandleFunc("/mail/config-v1.1.xml", funcT)
	http.HandleFunc("/autodiscover/autodiscover.xml", funcM)
	http.ListenAndServe(":8080", nil)
}

func funcT(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(`./config-v1.1.xml`)
	if err != nil {
		log.Fatal()
	}
	fmt.Fprint(w, string(data))
}
func funcM(w http.ResponseWriter, r *http.Request) {

	l, s := mytest()
	md := &MailData{l, s}
	fmt.Printf("md=%#v\n", md)
	tmpl := template.Must(template.ParseFiles("./autodiscover.xml"))
	err := tmpl.Execute(w, md)
	//	data, err := ioutil.ReadFile(`./autodiscover.xml`)
	if err != nil {
		log.Fatal()
	}
	//	printRequest(r)
	//	fmt.Fprint(w, string(data))
}
func printRequest(r *http.Request) {
	d, _ := httputil.DumpRequest(r, true)
	fmt.Printf("===Dump Request[START]\n%s\n===Dump Request[END]\n", d)
}
