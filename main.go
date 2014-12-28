//
// dnsdiff
// 2014 @oseemann
//

package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
)

type Options struct {
	dns1 string
	dns2 string
	zone string
}

func lookup(zone string, server string, record_type uint16) []string {
	m := new(dns.Msg)
	m.SetQuestion(zone, record_type)
	in, err := dns.Exchange(m, server)
	if err != nil {
		return nil
	}

	ret := make([]string, len(in.Answer))
	for i, elem := range in.Answer {
		ret[i] = elem.String()
	}
	return ret
}

func run(opt Options) {
	fmt.Println(opt)
	zone := opt.zone
	server := opt.dns1
	// SOA
	soa := lookup(zone, server, dns.TypeSOA)
	for _, elem := range soa {
		fmt.Printf("SOA: %s\n", elem)
	}

	// A
	a := lookup(zone, server, dns.TypeA)
	for _, elem := range a {
		fmt.Printf("A: %s\n", elem)
	}

	// AAAA
	aaaa := lookup(zone, server, dns.TypeAAAA)
	for _, elem := range aaaa {
		fmt.Printf("AAAA: %s\n", elem)
	}

	// CNAME
	cname := lookup(zone, server, dns.TypeCNAME)
	for _, elem := range cname {
		fmt.Printf("CNAME: %s\n", elem)
	}

	// MX
	mx := lookup(zone, server, dns.TypeMX)
	for _, elem := range mx {
		fmt.Printf("MX: %s\n", elem)
	}

	// TXT
	txt := lookup(zone, server, dns.TypeTXT)
	for _, elem := range txt {
		fmt.Printf("TXT: %s\n", elem)
	}

	// NS
	ns := lookup(zone, server, dns.TypeNS)
	for _, elem := range ns {
		fmt.Printf("NS: %s\n", elem)
	}

}

func main() {
	opt := Options{}

	flag.StringVar(&opt.dns1, "dns1", "", "DNS Server Address/Name")
	flag.StringVar(&opt.dns2, "dns2", "", "DNS Server Address/Name")
	flag.StringVar(&opt.zone, "zone", "", "DNS Zone to query")

	flag.Parse()

	run(opt)
}

// vim: set filetype=go expandtab:
