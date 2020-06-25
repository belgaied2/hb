package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

type step struct {
	Title   string
	Content string
}
type scenario struct {
	Id                 string
	Name               string
	Description        string
	Steps              []step
	Keepalive_duration string
	Virtualmachines    []map[string]string
	Pause_duration     string
	Pauseable          bool
}

type simplifiedScenario struct {
	Id   string
	Name string
}

var hb_url string

func dlCommand() *cli.Command {

	dlFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "email",
			Usage: "email address to connect to HobbyFarm",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "password to connect to HobbyFarm",
		},
		&cli.StringFlag{
			Name:  "scenario",
			Value: "ILT - Session 1",
			Usage: "Name of a scenario to download from HobbyFarm",
		},
		&cli.StringFlag{
			Name:  "region",
			Value: "eu1",
			Usage: "Region of HobbyFarm on which to look for scenario",
		},
	}

	return &cli.Command{
		Name:   "dl",
		Usage:  "Download a scenario",
		Action: dlScenMarkdown,
		Flags:  dlFlags,
	}
}

func dlScenMarkdown(ctx *cli.Context) error {

	email := ctx.String("email")
	password := ctx.String("password")
	region := ctx.String("region")
	scenName := ctx.String("scenario")

	// intermediate variables
	hb_url = "https://api." + region + ".hobbyfarm.io"
	var token = getToken(email, password)
	var bearer = "Bearer " + token

	// var baseURL = hb_url + "/a"

	scenContentBytes := getScenarioContentFromAPI(bearer, scenName)
	var scenDetail scenario
	marschallErr := json.Unmarshal(scenContentBytes, &scenDetail)

	if marschallErr != nil {
		log.Fatal("Unable to read json response: ", marschallErr)
	}

	if len(scenDetail.Steps) > 0 {
		os.Mkdir(scenName, os.ModePerm)
		os.Chdir(scenName)
	}
	return createStepFiles(scenDetail)
}

func getToken(username string, password string) string {
	form := url.Values{}

	form.Add("email", username)
	form.Add("password", password)

	resp, err := http.PostForm(hb_url+"/auth/authenticate", form)

	if err != nil {
		log.Fatal("Error during authentication request :", err)
	}

	var stringResp map[string]string
	respBody, ioErr := ioutil.ReadAll(resp.Body)

	if ioErr != nil {
		log.Fatal("Error reading the response's Body during authentication", ioErr)
	}

	marshallErr := json.Unmarshal(respBody, &stringResp)

	if marshallErr != nil {
		log.Fatal("Error parsing the response during authentication :", marshallErr)
	}

	result := stringResp["message"]

	// Test the result
	fmt.Println("The Token is: " + result)

	return result
}

func createStepFiles(scenDetail scenario) error {

	var err error

	for i, s := range scenDetail.Steps {
		stepName, nameError := base64.StdEncoding.DecodeString(s.Title)
		stepContent, contentErr := base64.StdEncoding.DecodeString(s.Content)

		if nameError != nil {
			log.Fatal(nameError)
		}

		if contentErr != nil {
			log.Fatal(contentErr)
		}
		stepNameCleaned := strings.ReplaceAll(strings.ReplaceAll(string(stepName), " ", "_"), ":", "")
		// stepNameCleaned := "step" + strconv.Itoa(i)
		fileName := strconv.Itoa(i) + "_" + stepNameCleaned + ".md"
		file, fileErr := os.Create(fileName)
		if fileErr != nil {
			log.Fatal("File could not be created ", fileErr)
		}

		_, err = file.Write(stepContent)

	}
	return err
}

func getHFContentFromURL(url string, bearer string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)

	// Check if response is not error
	if err != nil {
		log.Fatal("The HTTP request failed with error:", err)

	}

	// Read the response
	data, _ := ioutil.ReadAll(resp.Body)

	// check response body
	//fmt.Println(string(data))

	//Unmarshalling
	var jsonData map[string]interface{}
	json.Unmarshal(data, &jsonData)

	// get node content
	content := jsonData["content"].(string)
	decodedContent, err1 := base64.StdEncoding.DecodeString(content)
	if err1 != nil {

		fmt.Printf("The Base64 decoding failed! : %s\n", err1)
		// log.Fatal(err1)
	}
	return decodedContent
}
func getScenarioContentFromAPI(bearer string, scenName string) []byte {

	// Get List of scenarios from API
	scenListData := getHFContentFromURL(hb_url+"/a/scenario/list", bearer)

	// Unmarshal content to object simplifiedScenario
	scenarioData := []simplifiedScenario{}
	err2 := json.Unmarshal(scenListData, &scenarioData)

	fmt.Println("Scenario Data is : " + scenarioData[0].Id)

	if err2 != nil {
		log.Fatal(err2)
	}

	scenLink := hb_url + "/a/scenario/" + getScenarioIdFromName(scenarioData, scenName)
	fmt.Println("Link to access scenario data is: " + scenLink)
	scenData := getHFContentFromURL(scenLink, bearer)

	return scenData

}

func getScenarioIdFromName(scen []simplifiedScenario, name string) string {
	matchId := ""
	for _, s := range scen {
		decodedName, err := base64.StdEncoding.DecodeString(s.Name)
		if err == nil && string(decodedName) == name {
			matchId = s.Id
			break
		}
	}

	if matchId == "" {
		log.Fatal("No scenario found with the defined name in this region")
	}

	return matchId
}
