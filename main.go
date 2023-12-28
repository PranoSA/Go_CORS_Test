package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//Remember -> CORS is a per-origin configuration
// But -> It is served by a server

var origin string
var ajax_destination string
var ajax_methods string
var credentials bool //Cookies or Not ???//which Headers Are Sent
var ajax_headers string
var expose_headers bool
var simple_cors bool

func main() {
	//Read Flags and Set Variables

	//If Flags Don't Exist this, Ask User For Information from the terminal 1 by 1

	flag.StringVar(&origin, "origin", "", "Where are you sending the request from?")
	flag.StringVar(&ajax_destination, "dest", "", "Destination URL of Ajax Request")
	flag.StringVar(&ajax_methods, "methods", "", "Which METHOD are you using?")
	flag.StringVar(&ajax_headers, "headers", "", "Which Headers Are You Sending?")
	flag.BoolVar(&credentials, "credentials", false, "Do you Require Cookies?")
	flag.BoolVar(&expose_headers, "expose_headers", false, "Do You Require Particular Headers to Be Readable By Javascript?")
	flag.BoolVar(&simple_cors, "simple", false, "No Custom Headers or Methods")
	flag.Parse()

	methods := strings.Split(ajax_methods, ",")
	headers := strings.Split(ajax_headers, ",")

	if origin == ajax_destination {
		fmt.Println("Origin and Destination Cannot Be The Same")
		return
	}

	req, err := http.NewRequest("OPTIONS", ajax_destination, nil)
	if err != nil {
		log.Print(err)
		return
	}

	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", ajax_methods)
	req.Header.Set("Access-Control-Request-Headers", ajax_headers)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	//if Origin is Allowed
	//if Access-Control-Allow-Origin is present
	if resp.Header.Values("Access-Control-Allow-Origin") != nil {
		fmt.Println("Origin Is Allowed For Headers Specified Method")
		return
	}

	// If Bad Cors Response -> Origin, Headers, Method doesn't line up

	//Simple CORS ???
	if simple_cors == true {
		test_basic_cors(headers, methods, true)
		return
	}

	if test_basic_cors(headers, methods, false) == true {
		fmt.Println("Request Qualfies For Simple CORS, Run Again with --simple Enabled")
	}

	//Now Test for Method
	for _, method := range methods {
		req, err := http.NewRequest("OPTIONS", ajax_destination, nil)
		req.Header.Set("Origin", origin)
		req.Header.Set("Access-Control-Request-Method", method)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Print(err)
			return
		}

		defer resp.Body.Close()

		fmt.Println("Testing Method ----------------------------------------" + method)

		allowedMethods := resp.Header.Get("Access-Control-Allow-Methods")

		if strings.Contains(allowedMethods, method) {
			fmt.Println("Method Allowed " + method)

			checkedForMethodCredentials := false

			//Now Do The Checks for All Headers
			for _, headerU := range headers {
				header := strings.ToLower(headerU)
				fmt.Println("Testing ----------------------------------------" + method + "-" + header)
				req, err := http.NewRequest("OPTIONS", ajax_destination, nil)
				req.Header.Set("Origin", origin)
				req.Header.Set("Access-Control-Request-Method", method)
				req.Header.Set("Access-Control-Request-Headers", header)

				per_header_credentials := false

				if header == "cookie" && credentials == false {
					per_header_credentials = true
					fmt.Println("Cookies Are Subject to X-Access-Control-Allow-Credentials, Turning on Credentials Option")
				}
				if header == "authorization" && credentials == false {
					fmt.Println("Authorization Header Is Subject to Access-Control-Allow-Headers, Turning on Credentials Option")
					per_header_credentials = true
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Print(err)
					return
				}
				defer resp.Body.Close()
				allowedCredentials := resp.Header.Get("Access-Control-Allow-Credentials")

				if allowedCredentials == "false" && (credentials == true || per_header_credentials == true) {
					fmt.Println("Credentials Not Allowed For Method " + method)
				}
				if allowedCredentials == "true" && (credentials == true || per_header_credentials) && checkedForMethodCredentials == false {
					fmt.Println("Credentials Allowed For Method " + method)
					checkedForMethodCredentials = true
				}
				if allowedCredentials == "" && resp.Header.Get("Access-Control-Allow-Origin") == "*" {
					fmt.Println("WARNING !!!! Credentials Often Not Allowed for Wildcard Domains")
				}

				allowedHeaders := resp.Header.Get("Access-Control-Allow-Headers")
				//If The Header Exists
				if allowedHeaders != "" {
					fmt.Println("Header " + header + "  Allowed For Method " + method)
					continue
				}
				fmt.Println("Header " + header + "  Not Allowed For Method " + method)
			}
			continue
		}
		fmt.Println("Method Not Allowed " + method)
	}
}

func test_basic_cors(headers []string, methods []string, print bool) bool {
	/*
		Accept
		Accept-Language
		Content-Language
		Content-Type (with certain restrictions)
			application/x-www-form-urlencoded:
			multipart/form-data
			text/plain
	*/
	//Make Sure Only These Headers Are Specified

	if print {
		fmt.Println("Making Simple CORS Request")
	}

	for _, Uheader := range headers {

		header := strings.ToLower(Uheader)

		if header != "Accept" && header != "accept-language" && header != "content-language" && header != "content-type" {

			if print {
				fmt.Println("Simple CORS Headers Not Allowed")
			}
			return false
		}
		//Figure Out How to Tell User About This
		if header == "content-type" {
			if print {
				fmt.Println("Remember only x-www-form-urlencoded, multipart/form-data, text/plain are allowed for Content-Type")
			}
		}
	}

	/*
		anything in the "Authorization" header is subject to "Access-Control-Allow-Credentials",
		 but "Authentication" header is subject to Access-Control-Allow-Headers listing it
	*/

	/*
		Methods:
			Also Only Get, HEAD, or POST are allowed
			Make Sure Only These Methods Are Specified
	*/

	for _, method := range methods {
		if method != "GET" && method != "HEAD" && method != "POST" {
			if print {
				fmt.Println("Simple CORS Methods Not Allowed")
			}
			return false
		}
	}

	//Simple Cors  -> Only Per origin
	req, err := http.NewRequest("GET", ajax_destination, nil)
	req.Header.Set("Origin", origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return false
	}

	defer resp.Body.Close()

	if resp.Header.Get("Access-Control-Allow-Methods") == "" {
		if print {
			fmt.Println("Origin ALlowed for Simple CORS Request")
		}
		return true
	}
	return true
}
