package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"
)

const HomePage = `
<html>
	<head>
		<title>Cloud Detect</title>
		<link href="https://fonts.googleapis.com/css?family=Roboto+Mono:700&display=swap" rel="stylesheet">
	</head>
	<body style="margin: 0; font-family: 'Roboto Mono', monospace;">
		<div style="display: flex; flex-direction: column; justify-content: center; min-height: 100vh; background: {{ .Color }};">
			<h1 style="text-align: center; color: white; font-size: 4em;">{{ .Name }}</h1>
		</div>
	</body>
</htlm>
`

var netClient = &http.Client{
	Timeout: time.Second * 5,
}

type cloud struct {
	Name  string
	Color string
}

type endpoint struct {
	Path  string
	Cloud cloud
}

type result struct {
	Error error
	Cloud cloud
}

var es = []endpoint{
	endpoint{Path: "/latest/meta-data", Cloud: cloud{Name: "AWS", Color: "#FF9900"}},
	endpoint{Path: "/metadata/instance", Cloud: cloud{Name: "Azure", Color: "#007FFF"}},
	endpoint{Path: "/computeMetadata/", Cloud: cloud{Name: "GCP", Color: "#DB4437"}},
}

func main() {
	log.Println("Checking cloud provider")
	c := getCloudProvider("http://169.254.169.254")

	t := template.Must(template.New("home").Parse(HomePage))
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		log.Fatalf("Could not render html page: %v", err)
	}

	log.Println("Starting web server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(buffer.Bytes())
	})
	http.ListenAndServe(":8080", nil)
}

func getCloudProvider(baseUrl string) cloud {
	c := make(chan result)

	t := func(e endpoint) result {
		resp, err := netClient.Get(baseUrl + e.Path)
		if err != nil {
			return result{Error: err}
		}

		if resp.StatusCode != 200 {
			return result{Error: errors.New("Bad response status code")}
		}

		return result{Cloud: e.Cloud}
	}

	for _, e := range es {
		go func(e endpoint) { c <- t(e) }(e)
	}

	timeout := time.After(500 * time.Millisecond)
	for {
		select {
		case res := <-c:
			if res.Error == nil {
				return res.Cloud
			}
		case <-timeout:
			return cloud{Name: "Unknown", Color: ""}
		}
	}
}
