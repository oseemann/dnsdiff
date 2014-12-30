dnsdiff
=======

A small program to compare DNS Resource Records of a number of (sub-)domains on 2 different (authoritative) name servers.

Main use case is to check consistency when moving zones between name servers.

I used it when moving a domain with dozens of subdomains from a plain old shared hoster with a bad DNS management interface to AWS Route 53.

Usage
-----
```sh
$ dnsdiff  -dns1 ns.old-nameserver.net -dns2 ns.new-nameserver.net -name mydomain.com
```
