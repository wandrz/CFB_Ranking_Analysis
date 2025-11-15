package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	// "time"
)

// shhh.... it'll be ok
var (
	_ = json.Marshal
	_ = http.NewRequest
)

type Game = map[string]interface{}

type GameFetcher struct {
	url     string
	api_key string
	client  *http.Client
}

func (f *GameFetcher) fetchGame(year string) ([]byte, error) {

	urlGetByYear := f.url + "/games?year=" + year
	fmt.Println("getting year ", year, "->", urlGetByYear)

	req, err := http.NewRequest("GET", urlGetByYear, nil)
	if err != nil {
		return nil, fmt.Errorf("http request failed! %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+f.api_key)
	req.Header.Add("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("FetchGame failed! %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	//fmt.Println(string(body))

	return body, nil
}

func saveData(body []byte, year string) error {
	dir := filepath.Join("data", "games")
	file := filepath.Join(dir, year+".json")
	os.MkdirAll(dir, 0755)

	err := os.WriteFile(file, body, 0664)
	if err != nil {
		return fmt.Errorf("error writing file %s to dir %s", file, dir)
	}
	return nil
}

func main() {
	message := "here we go!"
	println(message)

	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}

	fetcher := GameFetcher{
		url:     getEnv("CFB_API_URL", true),
		api_key: getEnv("CFBD_API_KEY", false),
		client:  &http.Client{},
	}

	//ctx := context.Background()
	for i := 2000; i <= 2025; i++ {
		getYear := strconv.Itoa(i)
		games, err := fetcher.fetchGame(getYear)
		if err != nil {
			fmt.Println("Error retrieving ", getYear, "->", err)
		}

		err = saveData(games, getYear)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("~~~ fin ~~~")
}

func getEnv(env_var string, print_var bool) string {
	evar := os.Getenv(env_var)
	if evar == "" {
		log.Fatalf("Error: %v not found in .env", env_var)
	}

	if print_var {
		fmt.Println(evar)
	} else {
		fmt.Println("Successfully read", env_var)
	}

	return evar
}
