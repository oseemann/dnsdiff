//
// dnsdiff
// 2014 @oseemann
//

package main

import (
    "os"
    "flag"
    "fmt"
    "sort"
    "bufio"
    "github.com/miekg/dns"
)

type Options struct {
    dns1 string
    dns2 string
    name string
    namelist string
}

type Record struct {
    recordtype string
    value string
    ttl uint32
    mx_preference uint16
}

type Records []Record

func (r Records) Len() int {
    return len(r)
}

func (r Records) Less(i, j int) bool {
    a := r[i]
    b := r[j]

    if a.recordtype != b.recordtype {
        return a.recordtype < b.recordtype
    } else {
        return a.value < b.value
    }
}

func (r Records) Swap(i, j int) {
    r[i], r[j] = r[j], r[i]
}

var record_types = map[string]uint16{
    "SOA":   dns.TypeSOA,
    "A":     dns.TypeA,
    "AAAA":  dns.TypeAAAA,
    "CNAME": dns.TypeCNAME,
    "MX":    dns.TypeMX,
    "TXT":   dns.TypeTXT,
    "NS":    dns.TypeNS,
}

func lookup(name string, server string, record_type string) Records {
    m := new(dns.Msg)
    rt := record_types[record_type]
    m.SetQuestion(name, rt)
    in, err := dns.Exchange(m, server)
    if err != nil {
        return nil
    }

    // TODO: check Authoritative

    ret := make(Records, len(in.Answer))
    for i, elem := range in.Answer {
        r := Record{}
        r.ttl = elem.Header().Ttl
        switch record_type {
            case "A":
	            if t, ok := elem.(*dns.A); ok {
                    r.value = t.A.String()
                }
            case "AAAA":
	            if t, ok := elem.(*dns.AAAA); ok {
                    r.value = t.AAAA.String()
                }
            case "MX":
	            if t, ok := elem.(*dns.MX); ok {
                    r.value = t.Mx
                    r.mx_preference = t.Preference
                }
            case "TXT":
	            if t, ok := elem.(*dns.TXT); ok {
                    r.value = t.Txt[0]
                }
            case "SOA":
	            if t, ok := elem.(*dns.SOA); ok {
                    r.value = t.Ns
                }
        }
        ret[i] = r
    }
    sort.Sort(ret)
    return ret
}

func check(name, dns1, dns2 string) {
    fmt.Printf("Comparing %s\n", name)
    for rt := range record_types {
        fmt.Printf("\t%s:", rt)
        records1 := lookup(name, dns1, rt)
        records2 := lookup(name, dns2, rt)

        if len(records1) != len(records2) {
            fmt.Printf("ERROR: %d vs %d records\n", len(records1), len(records2))
            continue
        }
        if len(records1) == 0 {
            fmt.Printf("OK, 0 found.\n")
            continue
        }
        equals := make([]string, len(records1))
        e := 0
        for i, _ := range records1 {
            a := records1[i]
            b := records2[i]
            if a.value != b.value {
                fmt.Printf("ERR: %s != %s\n", a, b)
                continue
            } else {
                equals[e] = a.value
                e++
            }
        }
        fmt.Printf("OK all equal (%s)\n", equals)
    }
}

func read_name_list(filename string) []string {
    ret := make([]string, 0, 128)

    file, err := os.Open(filename)
    if err != nil {
        return ret
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        ret = append(ret, scanner.Text())
    }
    return ret
}

func run(opt Options) {
    if opt.name != "" {
        check(opt.name, opt.dns1, opt.dns2)
    } else if opt.namelist != "" {
        names := read_name_list(opt.namelist)
        for _, name := range names {
            check(name, opt.dns1, opt.dns2)
        }
    }
}

func parse_flags() Options {
    opt := Options{}

    flag.StringVar(&opt.dns1, "dns1", "", "DNS Server Address/Name")
    flag.StringVar(&opt.dns2, "dns2", "", "DNS Server Address/Name")
    flag.StringVar(&opt.name, "name", "", "Single host name to check")
    flag.StringVar(&opt.namelist, "namelist", "", "File with host names to check")

    flag.Parse()

    // TODO: check option consistency

    return opt
}

func main() {
    opt := parse_flags()
    run(opt)
}

// vim: set filetype=go ts=4 sw=4 sts=4 expandtab:
