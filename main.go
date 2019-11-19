package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const HomePage = `
<html>
	<head>
		<title>Cloud Detect</title>
	</head>
	<body>
		<h1>{{ .Name }}</h1>
	</body>
</htlm>
`

var netClient = &http.Client{
	Timeout: time.Second * 5,
}

func main() {
	log.Println("Starting Cloud Detect")
	name, color := getCloudProvider()

	log.Println("Starting web server")
	tmpl := template.Must(template.New("home").Parse(HomePage))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Name  string
			Color string
		}{
			name,
			color,
		}

		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":8080", nil)
}

func getCloudProvider() (string, string) {
	resp, err := netClient.Get("http://169.254.169.254")

	if err != nil {
		return "Unknown", ""
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Unknown", ""
	}

	return "AWS", "#ff9900"
}
