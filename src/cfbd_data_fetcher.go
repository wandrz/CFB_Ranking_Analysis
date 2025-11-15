package main

import (
	"encoding/json"
	"flag"
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

type fetcher struct {
	url     string
	api_key string
	client  *http.Client
}

func (f *fetcher) cfbdFetcher(resource string, year string) ([]byte, error) {

	urlGetByYear := f.url + "/" + resource + "?year=" + year
	fmt.Println("getting ", resource, " for year ", year, "->", urlGetByYear)

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

func saveData(body []byte, resource string, year string) error {
	dir := filepath.Join("data", resource)
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

	resourcePtr := flag.String("resource", "", "Resource to fetch: games or rankings")
	fromPtr := flag.Int("from", 2000, "Start year")
	toPtr := flag.Int("to", 2025, "End year")
	flag.Parse()

	resource := *resourcePtr
	from := *fromPtr
	to := *toPtr

	if resource != "games" && resource != "rankings" {
		log.Fatal("Resource must be 'games' or 'rankings'")
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}

	fetcher := fetcher{
		url:     getEnv("CFB_API_URL", true),
		api_key: getEnv("CFBD_API_KEY", false),
		client:  &http.Client{},
	}

	//ctx := context.Background()
	for i := from; i <= to; i++ {
		getYear := strconv.Itoa(i)
		games, err := fetcher.cfbdFetcher(resource, getYear)
		if err != nil {
			fmt.Println("Error retrieving ", getYear, "->", err)
		}

		err = saveData(games, resource, getYear)
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
