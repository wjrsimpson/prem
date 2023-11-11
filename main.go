package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	baseUrl      = "https://fantasy.premierleague.com/api"
	bootstrapUrl = baseUrl + "/bootstrap-static"
	fixturesUrl  = baseUrl + "/fixtures"
)

type bootstrapData struct {
	Teams []team
}

type team struct {
	Id   int
	Name string
}

type fixture struct {
	Id          int
	HomeTeamId  int       `json:"team_h"`
	AwayTeamId  int       `json:"team_a"`
	KickOffTime time.Time `json:"kickoff_time"`
}

func main() {
	// Get bootstrap data from FPL API and write to cache
	bootstrapFilePath := getBootstrapFilePath()
	populateBootstrapCacheFile(bootstrapFilePath)
	// Read bootstrap json data into structs
	bootstrapData := readBootstrapData(bootstrapFilePath)
	// Put the teams into a map for easy lookup
	teams := make(map[int]string)
	for _, team := range (*bootstrapData).Teams {
		teams[team.Id] = team.Name
	}
	// Get fixtures data from FPL API and write to cache
	fixturesFilePath := getFixturesFilePath()
	populateFixturesCacheFile(fixturesFilePath)
	// Read fixtures json data into structs
	fixturesData := readFixturesData(fixturesFilePath)
	fmt.Println("Number of fixtures:", len(fixturesData))
	for _, fixture := range fixturesData {
		fmt.Println(fixture, teams[fixture.HomeTeamId], teams[fixture.AwayTeamId], fixture.KickOffTime)
	}
}

func readFixturesData(filePath string) []fixture {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening fixtures file:", err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data []fixture
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding fixtures file:", err)
		os.Exit(1)
	}
	return data
}

func readBootstrapData(filePath string) *bootstrapData {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening bootstrap file:", err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data bootstrapData
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding bootstrap file:", err)
		os.Exit(1)
	}
	return &data
}

func getBootstrapFilePath() string {
	cacheDir := createCacheDir()
	bootstrapFilePath := filepath.Join(cacheDir, "bootstrap.json")
	return bootstrapFilePath
}

func createCacheDir() string {
	rootCacheDir, _ := os.UserCacheDir()
	cacheDir := filepath.Join(rootCacheDir, "prem")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, 0755)
		if err != nil {
			fmt.Println("Error creating cache directory:", err)
			os.Exit(1)
		}
	}
	return cacheDir
}

func populateBootstrapCacheFile(bootstrapFilePath string) {
	if _, err := os.Stat(bootstrapFilePath); os.IsNotExist(err) {
		bootstrapFile := createBootstrapFile(bootstrapFilePath)
		res := getBootstrapData()
		copyBootstrapDataToFile(bootstrapFile, res)
		fmt.Println("bootstrap file written to cache:", bootstrapFile.Name())
	} else {
		fmt.Println("bootstrap file already exists in cache:", bootstrapFilePath)
	}
}

func createBootstrapFile(bootstrapFilePath string) *os.File {
	bootstrapFile, err := os.Create(bootstrapFilePath)
	if err != nil {
		fmt.Println("Error creating bootstrap file:", err)
		os.Exit(1)
	}
	return bootstrapFile
}

func getBootstrapData() *http.Response {
	res, err := http.Get(bootstrapUrl)
	if err != nil {
		fmt.Println("Error get bootstrap static:", err)
		os.Exit(1)
	}
	return res
}

func copyBootstrapDataToFile(bootstrapFile *os.File, res *http.Response) {
	_, err := io.Copy(bootstrapFile, res.Body)
	if err != nil {
		fmt.Println("Error writing bootstrap file:", err)
		os.Exit(1)
	}
}

func getFixturesFilePath() string {
	cacheDir := createCacheDir()
	fixturesFilePath := filepath.Join(cacheDir, "fixtures.json")
	return fixturesFilePath
}

func populateFixturesCacheFile(fixturesFilePath string) {
	if _, err := os.Stat(fixturesFilePath); os.IsNotExist(err) {
		fixturesFile := createFixturesFile(fixturesFilePath)
		res := getFixturesData()
		copyFixturesDataToFile(fixturesFile, res)
		fmt.Println("fixtures file written to cache:", fixturesFile.Name())
	} else {
		fmt.Println("fixtures file already exists in cache:", fixturesFilePath)
	}
}

func createFixturesFile(fixturesFilePath string) *os.File {
	fixturesFile, err := os.Create(fixturesFilePath)
	if err != nil {
		fmt.Println("Error creating fixtures file:", err)
		os.Exit(1)
	}
	return fixturesFile
}

func getFixturesData() *http.Response {
	res, err := http.Get(fixturesUrl)
	if err != nil {
		fmt.Println("Error get fixtures:", err)
		os.Exit(1)
	}
	return res
}

func copyFixturesDataToFile(fixturesFile *os.File, res *http.Response) {
	_, err := io.Copy(fixturesFile, res.Body)
	if err != nil {
		fmt.Println("Error writing fixtures file:", err)
		os.Exit(1)
	}
}
