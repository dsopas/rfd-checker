RFD Checker
===========
[![GPL3](https://camo.githubusercontent.com/268d96c6dd81f1fff98b19675ef5867412a2a223/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f6c6963656e73652d47504c332d627269676874677265656e2e7376673f7374796c653d666c61742d737175617265)](https://github.com/dsopas/rfd-checker/blob/master/LICENSE.md) [![Go Report Card](https://goreportcard.com/badge/github.com/dsopas/rfd-checker)](https://goreportcard.com/report/github.com/dsopas/rfd-checker)

Command line security tool to check whether a given URL is vulnerable to RFD -
Reflected File Download. This tool was developed by David Sopas [@dsopas][1]
and Paulo Silva [@pauloasilva_com][2] with the main purpose of validating and
automating the search for the RFD web attack vector.

## Usage

```
$ rfd-checker -h
RFD Checker (by @dsopas and @pauloasilva_com)

Usage: rfd-checker -target=URL
Options:
  -header value
        Request header e.g. "Cookie: SESSID=a16fb"
  -target string
        Target URL
  -h --help
        Prints this help
```

### Examples

```shell
$ go run rfd-checker.go -target="https://0xhack.com/webminar_rfd/json.php?callback=jQuery" -header="User-Agent: RFD-Checker" -header="Cookie: PHPSESSID=123"
Target URL: https://0xhack.com/webminar_rfd/json.php?callback=jQuery
Permissive query parameters: callback
Permissive URL: https://0xhack.com/webminar_rfd/json.php.bat?callback=jQuery
```

* Permissive query parameteres: a comma separated list of permissive query
  string parameters
* Permissive URL: "no" if the URL is not permissive, otherwise computed URL
  (e.g. https://0xhack.com/webminar_rfd/json.php.bat?callback=jQuery)

If you want to test a batch of URLs, exported from Burp, for example, you can
place them, one per line, on a text file and run

```shell
$ cat inputs.txt | xargs -I url go run ./rfd-checker.go -target=url
``` 

Or

```shell
$ cat inputs.txt | xargs -I url ./rfd-checker -target=url
```

![RFD checker diagram](https://0xhack.com/rfd_checker.png)

## Build

```
$ go build rfd-checker.go
```

## Resources

* [Reflected File Download - A New Web Attack Vector][3]
* [Reflected File Download Cheat Sheet][4]
* [Practical Reflected File Download and JSONP][5]
* [RFD: Still Threatening the Biggest Names on the Web][6]

[1]: https://www.twitter.com/dsopas
[2]: https://www.twitter.com/pauloasilva_com
[3]: https://www.blackhat.com/docs/eu-14/materials/eu-14-Hafif-Reflected-File-Download-A-New-Web-Attack-Vector.pdf
[4]: https://www.davidsopas.com/reflected-file-download-cheat-sheet/
[5]: http://blog.davidvassallo.me/2014/11/02/practical-reflected-file-download-and-jsonp/
[6]: https://info.checkmarx.com/resources/webinars/rfd-still-threatening-the-biggest-names-web-on-demand?hsCtaTracking=70be984d-c6b2-4eb6-a280-32ac7aa6a520%7C17df43d5-14db-4b83-ad12-09f16270754f

