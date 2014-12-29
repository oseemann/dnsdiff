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
    zone := opt.zone

    var record_types = map[string]uint16{
        "SOA":   dns.TypeSOA,
        "A":     dns.TypeA,
        "AAAA":  dns.TypeAAAA,
        "CNAME": dns.TypeCNAME,
        "MX":    dns.TypeMX,
        "TXT":   dns.TypeTXT,
        "NS":    dns.TypeNS,
    }

    fmt.Printf("Comparing %s\n", zone)
    for name, rtype := range record_types {
        fmt.Printf("\t%s:", name)
        records1 := lookup(zone, opt.dns1, rtype)
        records2 := lookup(zone, opt.dns2, rtype)

        if len(records1) != len(records2) {
            fmt.Printf("ERROR: %d vs %d records\n", len(records1), len(records2))
            continue
        }
        if len(records1) == 0 {
            fmt.Printf("OK, 0 found.\n")
            continue
        }
        for i, _ := range records1 {
            a := records1[i]
            b := records2[i]
            if records1[i] != records2[i] {
                fmt.Printf("ERR: %s != %s\n", a, b)
                continue
            }
        }
        fmt.Printf("OK\n")
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

// vim: set filetype=go ts=4 sw=4 sts=4 expandtab:
