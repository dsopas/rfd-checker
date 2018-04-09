RFD Checker
===========

Command line security tool to check whether a given URL is vulnerable to RFD -
Reflected File Download. This CLI tool was developed by David Sopas [@dsopas][1]
and Paulo Silva [@pauloasilva_com][2] with the main purpose of validate and
automate the search for RFD web attack vector.

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

## License

Copyright (c) 2018 David Sopas <davidsopas@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

[1]: https://www.twitter.com/dsopas
[2]: https://www.twitter.com/pauloasilva_com
[3]: https://www.blackhat.com/docs/eu-14/materials/eu-14-Hafif-Reflected-File-Download-A-New-Web-Attack-Vector.pdf
[4]: https://www.davidsopas.com/reflected-file-download-cheat-sheet/
[5]: http://blog.davidvassallo.me/2014/11/02/practical-reflected-file-download-and-jsonp/
[6]: https://info.checkmarx.com/resources/webinars/rfd-still-threatening-the-biggest-names-web-on-demand?hsCtaTracking=70be984d-c6b2-4eb6-a280-32ac7aa6a520%7C17df43d5-14db-4b83-ad12-09f16270754f

