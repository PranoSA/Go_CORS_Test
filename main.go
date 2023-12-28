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

type CorsResults struct {
	origin_problem      bool
	ajax_method         bool
	credentials_problem bool
}

func main() {
	//Read Flags and Set Variables

	//If Flags Don't Exist this, Ask User For Information from the terminal 1 by 1

	flag.StringVar(&origin, "origin", "", "Where are you sending the request from?")
	flag.StringVar(&ajax_destination, "dest", "", "Destination URL of Ajax Request")
	flag.StringVar(&ajax_methods, "methods", "", "Which METHOD are you using?")
	flag.StringVar(&ajax_headers, "headers", "", "Which Headers Are You Sending?")
	flag.BoolVar(&credentials, "credentials", false, "Do you Require Cookies?")
	flag.BoolVar(&expose_headers, "expose_headers", false, "Do You Require Particular Headers to Be Readable By Javascript?")
	flag.Parse()

	methods := strings.Split(ajax_methods, ",")
	headers := strings.Split(ajax_headers, ",")

	if origin == ajax_destination {
		fmt.Println("Origin and Destination Cannot Be The Same")
		return
	}

	//Make an HTTP Request to the origin with the following information
	/*
		Headers :
			Origin : origin
			Access-Control-Request-Method : ajax_method
			Access-Control-Request-Headers : headers

	*/

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
		fmt.Println("Origin Is Allowed For All Methods And Headers")
		return
	}

	// If Bad Cors Response -> Origin, Headers, Method doesn't line up

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

		allowedMethods := resp.Header.Get("Access-Control-Allow-Methods")

		if strings.Contains(allowedMethods, method) {
			fmt.Println("Method Allowed " + method)

			//Now Do The Checks for All Headers
			for _, header := range headers {
				req, err := http.NewRequest("OPTIONS", ajax_destination, nil)
				req.Header.Set("Origin", origin)
				req.Header.Set("Access-Control-Request-Method", method)
				req.Header.Set("Access-Control-Request-Headers", header)

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Print(err)
					return
				}
				defer resp.Body.Close()

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
