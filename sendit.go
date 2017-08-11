// A utility program to send HTTP POST and GET requests.
package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	prog = "[sendit]"
)

var (
	url     = flag.String("url", "", "url of where data is sent")
	data    = flag.String("data", "", "data to send in the body of the request")
	file    = flag.String("file", "", "file whose contents will be sent as data in the body of the request")
	help    bool
	verbose bool
)

func init() {
	const (
		helpUsage   = "prints the help for sendit"
		verboseHelp = "print sendit's log output"
		progUsage   = `Sendit can send arbitrary data via a web request and then prints the
response from the listening program.

Sendit will send the data specified to the program listening at the given
'url'.  If 'url' is not specified then sendit will dump the data out to the
console.  If the options 'data' and 'file' are both specified then both will
be sent as separate requests to the same url.
`
	)
	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage)
	flag.BoolVar(&verbose, "verbose", false, verboseHelp)
	flag.BoolVar(&verbose, "v", false, verboseHelp)

	flag.Usage = func() {
		fmt.Println(progUsage)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if *url == "" {
		if *data != "" {
			vlog("Output data --")
			fmt.Println(*data, "\n")
		}

		if *file != "" {
			withFile(*file, func(content []byte) {
				vlog("File output --")
				fmt.Printf("%s\n", content)
			})
		}
	} else {
		var client = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 4000 * time.Millisecond,
		}

		if *data != "" {
			resp, err := client.Post(*url, "text/plain", strings.NewReader(*data))
			processResponse(resp, err)
		}

		if *file != "" {
			withFile(*file, func(content []byte) {
				resp, err := client.Post(*url, "text/plain", bytes.NewReader(content))
				processResponse(resp, err)
			})
		}

		if *data == "" && *file == "" {
			resp, err := client.Get(*url)
			processResponse(resp, err)
		}
	}
}

func processResponse(resp *http.Response, err error) {
	if err != nil {
		elog("Error from web request --\n    ", err)
	} else {
		vlog("Received response status:", resp.Status)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			elog("Could not read body of the response --\n    ", err)
		} else if len(body) > 0 {
			vlog("Reponse body --\n")
			fmt.Printf("%s\n", body)
		}
	}
}

func withFile(name string, fn func([]byte)) {
	content, err := ioutil.ReadFile(*file)
	if err != nil {
		elog(prog, "Could not read contents of the file:", *file)
	} else {
		fn(content)
	}
}

// For logging errors.
func elog(args ...interface{}) {
	log.Println(join(prog, args...)...)
}

// For logging verbose output.
func vlog(args ...interface{}) {
	if verbose {
		log.Println(join(prog, args...)...)
	}
}

func join(first interface{}, args ...interface{}) []interface{} {
	a := make([]interface{}, 0)
	a = append(a, first)
	return append(a, args...)
}
