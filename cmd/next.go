/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

const (
	baseUrl      = "https://fantasy.premierleague.com/api"
	bootstrapUrl = baseUrl + "/bootstrap-static"
	fixturesUrl  = baseUrl + "/fixtures"
)

var refreshCache *bool

type bootstrapData struct {
	Teams []team
}

type team struct {
	Id   int
	Name string
}

type teamBucket struct {
	team
	fixtures []fixture
	points   int
}

type fixture struct {
	Id             int
	HomeTeamId     int       `json:"team_h"`
	AwayTeamId     int       `json:"team_a"`
	KickOffTime    time.Time `json:"kickoff_time"`
	Finished       bool
	HomeScore      int `json:"team_h_score"`
	AwayScore      int `json:"team_a_score"`
	HomeDifficulty int `json:"team_h_difficulty"`
	AwayDifficulty int `json:"team_a_difficulty"`
}

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Prints out the next 5 fixtures for each team",
	Long: `Prints out the next 5 fixtures for each team, along with the difficulty of each fixture.
	
The fixtures will be retrieved from the FPL API and cached in the user's cache directory. You can force a refresh of the cache by using the -r flag.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		teams := getTeamMap()
		fixturesList := getFixturesList()
		processFixtures(fixturesList, teams)
		printTeams(teams)
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	refreshCache = nextCmd.Flags().BoolP("refresh", "r", false, "refresh the cache")
}

func getTeamMap() map[int]*teamBucket {
	// Get bootstrap data from FPL API and write to cache
	bootstrapFilePath := getBootstrapFilePath()
	populateBootstrapCacheFile(bootstrapFilePath)
	// Read bootstrap json data into structs
	bootstrapData := readBootstrapData(bootstrapFilePath)
	// Put the teams into a map for easy lookup
	teams := make(map[int]*teamBucket)
	for _, team := range (*bootstrapData).Teams {
		teams[team.Id] = &teamBucket{team: team}
	}
	return teams
}

func getBootstrapFilePath() string {
	cacheDir := createCacheDir()
	bootstrapFilePath := filepath.Join(cacheDir, "bootstrap.json")
	return bootstrapFilePath
}

func populateBootstrapCacheFile(bootstrapFilePath string) {
	_, err := os.Stat(bootstrapFilePath)
	fileNotExists := os.IsNotExist(err)
	if fileNotExists || *refreshCache {
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

func getFixturesList() []fixture {
	// Get fixtures data from FPL API and write to cache
	fixturesFilePath := getFixturesFilePath()
	populateFixturesCacheFile(fixturesFilePath)
	// Read fixtures json data into structs
	fixturesData := readFixturesData(fixturesFilePath)
	return fixturesData
}

func getFixturesFilePath() string {
	cacheDir := createCacheDir()
	fixturesFilePath := filepath.Join(cacheDir, "fixtures.json")
	return fixturesFilePath
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

func populateFixturesCacheFile(fixturesFilePath string) {
	_, err := os.Stat(fixturesFilePath)
	fileNotExists := os.IsNotExist(err)
	if fileNotExists || *refreshCache {
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

func processFixtures(fixturesData []fixture, teams map[int]*teamBucket) {
	fmt.Println("Number of fixtures:", len(fixturesData))
	for _, fixture := range fixturesData {
		homeTeamBucket := teams[fixture.HomeTeamId]
		awayTeamBucket := teams[fixture.AwayTeamId]
		if fixture.Finished {
			if fixture.HomeScore > fixture.AwayScore {
				homeTeamBucket.points += 3
			} else if fixture.HomeScore < fixture.AwayScore {
				awayTeamBucket.points += 3
			} else {
				homeTeamBucket.points += 1
				awayTeamBucket.points += 1
			}
		} else if fixture.KickOffTime.Year() > 1 {
			if len(homeTeamBucket.fixtures) < 5 {
				homeTeamBucket.fixtures = append(homeTeamBucket.fixtures, fixture)
			}
			if len(awayTeamBucket.fixtures) < 5 {
				awayTeamBucket.fixtures = append(awayTeamBucket.fixtures, fixture)
			}
		}
	}
}

func printTeams(teams map[int]*teamBucket) {
	// Print out the teams in order of points
	teamBucketList := make([]teamBucket, 0)
	for _, teamBucket := range teams {
		teamBucketList = append(teamBucketList, *teamBucket)
	}
	sort.Slice(teamBucketList, func(a, b int) bool {
		return teamBucketList[a].points > teamBucketList[b].points
	})
	for _, teamBucket := range teamBucketList {
		fmt.Println()
		fmt.Println(teamBucket.team.Name, ":")
		var totalDifficulty int
		for _, fixture := range teamBucket.fixtures {
			if teamBucket.team.Id == fixture.HomeTeamId {
				totalDifficulty += fixture.HomeDifficulty
				fmt.Println("H ", teams[fixture.AwayTeamId].team.Name, fixture.HomeDifficulty, fixture.KickOffTime)
			} else {
				totalDifficulty += fixture.AwayDifficulty
				fmt.Println("A ", teams[fixture.HomeTeamId].team.Name, fixture.AwayDifficulty, fixture.KickOffTime)
			}
		}
		fmt.Println("Total difficulty:", totalDifficulty)
	}
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
