package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

//config provides configuration data for body json upload,
//and application specific information
type config struct {
	URL                string
	UserID             string
	ProcessorSessionID string
	SecurityToken      string
	OutputFilePath     string
}

//flag provides command line parameters
type flags struct {
	ConfigFilePath string
	StartDateTime  string
	EndDateTime    string
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
	configPathPtr := flag.String("config", "./config.json", "Configuration file path.")
	startDTPtr := flag.String("s", "", "Start DateTime in 'YYYY-MM-DDTHH-MM-SS' format. (Required)")
	endDTPtr := flag.String("e", "", "End DateTime in 'YYYY-MM-DDTHH-MM-SS' format. (Required)")
	flag.Parse()

	if *startDTPtr == "" || *endDTPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return flags{*configPathPtr, *startDTPtr, *endDTPtr}
}

func main() {
	var err error
	// Command line flags
	flags := getFlags()
	// configuration parameters from config file
	config := getConfigData(flags.ConfigFilePath)
	// get body template parameters
	bodyParams := bodyParams{config.UserID, config.ProcessorSessionID,
		config.SecurityToken, flags.StartDateTime, flags.EndDateTime}

	template, err := template.ParseFiles(`./bodytemplate.json`)
	panicIfErr(err)

	var buff bytes.Buffer
	err = template.Execute(&buff, bodyParams)
	panicIfErr(err)
	log.Println("Created body content and ready to make a request")

	req, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(buff.Bytes()))
	panicIfErr(err)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Sending request to the server: %v", config.URL)
	client := &http.Client{}
	resp, err := client.Do(req)
	panicIfErr(err)
	if resp.StatusCode == 200 {
		log.Println("Successfully got response from the server")
	} else {
		log.Fatalf("Server response: %s", resp.Status)
	}
	defer resp.Body.Close()

	out, err := os.Create(config.OutputFilePath)
	panicIfErr(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	panicIfErr(err)
	log.Printf("Successfully loaded and saved file to %s", config.OutputFilePath)
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

//panicIfErr prints error stack and exit application with error code 1
func panicIfErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
