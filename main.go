package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"text/template"
)

//config provides configuration data for body json upload,
//and application specific information
type config struct {
	URL                string
	UserID             string
	ProcessorSessionID string
	SecurityToken      string
}

//flag provides command line parameters
type flags struct {
	startDateTime  string
	endDateTime    string
	configFilePath string
	outputFilePath string
	pr             int
	ps             int
}

//bodyParams provides parameters for POST body template
type bodyParams struct {
	UserID             string
	ProcessorSessionID string
	SecurityToken      string
	StartDateTime      string
	EndDateTime        string
}

//getFlags gets command line flags, such as config file path, start and end date
func getFlags() flags {
	startDTPtr := flag.String("s", "", "Start DateTime in 'YYYY-MM-DD' or 'YYYY-MM-DDTHH-MM-SS' format. (Required)")
	endDTPtr := flag.String("e", "", "End DateTime in 'YYYY-MM-DD' or 'YYYY-MM-DDTHH-MM-SS' format. (Required)")
	prPtr := flag.Int("pr", 1, `The number of pages requested; this is the number that should be changed in
the case there are more	than 10 queries in a search your financial institution
is trying to utilize. For example, if the financial institution has 300
applications in the month of December (using the same time constraints from
before), then the URL parameter “pr” would be set to any value from 1-30
(taking the applications 10 entries at a time – 1-10, 11-20, 21-30, etc).`)
	psPtr := flag.Int("ps", 10, `Stands for the size of the page.
Since the standard is having 10 applications
maximum per page (as seen in the admin portal search queries), we recommend
only using a maximum of 10 for this value in the URL.`)
	configPathPtr := flag.String("c", "./config.json", "Configuration file path.")
	outputFilePathPtr := flag.String("o", "./data.json", "Output file path")

	flag.Parse()

	if *startDTPtr == "" || *endDTPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return flags{*startDTPtr,
		*endDTPtr,
		*configPathPtr,
		*outputFilePathPtr,
		*prPtr, *psPtr}
}

func main() {
	var err error
	// Command line flags
	flags := getFlags()
	// configuration parameters from config file
	config := getConfigData(flags.configFilePath)
	// get body template parameters
	bodyParams := bodyParams{config.UserID, config.ProcessorSessionID,
		config.SecurityToken, flags.startDateTime, flags.endDateTime}

	template, err := template.ParseFiles(`./bodytemplate.json`)
	panicIfErr(err)

	var buff bytes.Buffer
	err = template.Execute(&buff, bodyParams)
	panicIfErr(err)
	log.Println("Created body content and ready to make a request")

	url := urlQueryString(config.URL, flags.pr, flags.ps)
	log.Printf("URL Path %s", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buff.Bytes()))
	panicIfErr(err)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Sending request to the server: %v", url)
	client := &http.Client{}
	resp, err := client.Do(req)
	panicIfErr(err)
	if resp.StatusCode == 200 {
		log.Println("Successfully got response from the server")
	} else {
		log.Fatalf("Server response: %s", resp.Status)
	}
	defer resp.Body.Close()

	out, err := os.Create(flags.outputFilePath)
	panicIfErr(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	panicIfErr(err)
	log.Printf("Successfully loaded and saved file to %s", flags.outputFilePath)
}

//getConfigData load configuration information from config file
//Configuration information includes POST body request payload
//URL address and output file location
func getConfigData(configFilePath string) config {
	file, err := os.Open(configFilePath)
	panicIfErr(err)
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := config{}
	err = decoder.Decode(&config)
	panicIfErr(err)

	return config
}

//urlQueryString creates url with query string parameters
func urlQueryString(baseURL string, pr, ps int) string {
	u, err := url.Parse(baseURL)
	panicIfErr(err)

	q := u.Query()
	q.Add("pr", strconv.Itoa(pr))
	q.Add("ps", strconv.Itoa(ps))
	u.RawQuery = q.Encode()

	return u.String()
}

//panicIfErr prints error stack and exit application with error code 1
func panicIfErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
