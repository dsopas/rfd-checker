// Copyright © 2018 David Sopas <davidsopas@gmail.com>.
//
// Usage:
//     $ rfd-checker -target="URL" -header="X-Header-Name: value"
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
    "bytes"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "path"
    "strings"
    "time"
)

// custom flag Usage message
func Usage () {
    fmt.Printf("RFD Checker (by @dsopas and @pauloasilva_com)\n\n")
    fmt.Printf("Usage: %s -target=URL\n", os.Args[0])
    fmt.Printf("Options:\n")
    flag.PrintDefaults()
}

// Common vulnerable query parameters
var commonVulnParams = map[string]bool{
    "callback": true,
    "cb":       true,
    "jsonp":    true,
    "jsoncb":   true,
}

type reqHeaders map[string]string

func (h reqHeaders) String() string {
    return ""
}

func (h reqHeaders) Set(value string) error {
    parts := strings.Split(value, ":")

    if len(parts) == 1 {
        fmt.Fprintf(os.Stderr, "[WARN] Ignoring bad header \"%s\"...\n", value)
        return nil
    }

    hName := strings.Trim(parts[0], " ")
    hValue := strings.Trim(parts[1], " ")

    h[hName] = hValue

    return nil
}

// Performs an HTTP GET request to `targetURL`, returning both the HTTP status
// and response body.
// Any internal error is returned as is with `nil` `response` and "000 Error"
// status
func request(targetURL string, headers reqHeaders) (string, []byte, error) {
    tr := &http.Transport{IdleConnTimeout: 30 * time.Second}
    client := &http.Client{Transport: tr}

    req, err := http.NewRequest("GET", targetURL, nil)

    if err != nil {
        return "000 Error", nil, err
    }

    for k, v := range headers {
        req.Header.Add(k, v)
    }

    resp, err := client.Do(req)
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return "000 Error", nil, err
    }

    return resp.Status, body, nil
}

// Tests whether given query string paramenter `pName` from `oURL` URL is
// vulnerable to RFD.
// Any internal error is returned as is and should be handled by calling
// function
func testQueryParameter(pName string, oURL string, h reqHeaders) (bool, error) {
    tempURL, err := url.Parse(oURL)

    if err != nil {
        return false, err
    }

    q := tempURL.Query()
    if commonVulnParams[pName] {
        q.Set(pName, "calc||")
    } else {
        q.Set(pName, "\"||calc||")
    }

    tempURL.RawQuery = q.Encode()

    status, body, err := request(tempURL.String(), h)

    if status != "200 OK" {
        return false, err
    }

    return strings.Contains(string(body), "calc"), nil
}

func main() {
    var exitCode int = 0

    var targetURL string
    var targetHeaders reqHeaders = make(reqHeaders)

    var permissiveParams []string
    var permissiveURL string = "no"

    defer func() {
        os.Exit(exitCode)
    }()

    flag.Usage = Usage

    flag.StringVar(&targetURL, "target", "", "Target URL")
    flag.Var(&targetHeaders, "header",
        "Request header e.g. \"Cookie: SESSID=a16fb\"")

    flag.Parse()

    if targetURL == "" {
        fmt.Fprintf(os.Stderr, "Missing option -target <URL>\n")
        exitCode = 1

        return
    }

    u, err := url.Parse(targetURL)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Invalid -target option value\n")
        exitCode = 2

        return
    }

    // find permissive query parameters
    parameters := u.Query()
    for k := range parameters {
        isPermissive, err := testQueryParameter(k, targetURL, targetHeaders)

        if err == nil && isPermissive {
            permissiveParams = append(permissiveParams, k)
        }
    }

    // test url permissiveness
    dir, filename := path.Split(u.Path)

    if filename == "" {
        filename = "setup.bat"
    } else {
        if strings.Contains(filename, ".") {
            filename = filename + ".bat"
        } else {
            filename = filename + ";setup.bat"
        }
    }

    u.Path = path.Join(dir, filename)

    // test whether requesting the original URL and the modified one return the
    // successfully with the same response body
    oStatus, oBody, err := request(targetURL, targetHeaders)
    tStatus, tBody, err := request(u.String(), targetHeaders)
    if oStatus == "200 OK" && tStatus == "200 OK" && bytes.Equal(tBody, oBody) {
        permissiveURL = u.String()
    }

    permissiveParamsCSV := strings.Join(permissiveParams[:], ",")

    fmt.Printf("Target URL: %s\n", targetURL)
    fmt.Printf("Permissive query parameters: %v\n", permissiveParamsCSV)
    fmt.Printf("Permissive URL: %v\n", permissiveURL)
}

// vim: set ts=4 sw=4 et :
