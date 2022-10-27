package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

func main() {

	configPath := flag.String("c", "./config.json", "Specify the configuration file to use.")
	storefrontPath := flag.String("sf", "./storefrontmappings.json", "Specify the storefront mapping file to use.")
	flag.Parse()

	config, err := getConfiguration(configPath)
	if errorOccured(err) {
		log.Fatalf("Failed to find/parse required configuration file at %s: %s\n", *configPath, err)
	}

	storefronts, err := getStorefronts(storefrontPath)
	if errorOccured(err) {
		log.Fatalf("Failed to find/parse required storefront mapping file at %s: %s\n", *storefrontPath, err)
	}

	for _, app := range config.Apps {
		details, err := getAppDetails(app, storefronts)

		if errorOccured(err) {
			fmt.Println(err)
		} else {
			feed := makeFeedForApp(app, details, config)
			writeFeeds(feed, app, config)
		}
	}
}

func getAppDetails(app App, storefronts []Storefront) (details AppDetails, err error) {
	details.Name = app.ShortName

	client := &http.Client{}
	endpoint := MakeEndpointForApp(app)
	req, err := http.NewRequest("GET", endpoint, nil)
	if errorOccured(err) {
		return details, err
	}

	storefront, err := getMatchingStorefront(app, storefronts)
	if errorOccured(err) {
		return details, err
	}

	req.Header.Add("X-Apple-Store-Front", MakeStorefrontHeader(storefront))
	// req.Header.Add("apple-originating-system", "MZStore")

	if errorOccured(err) {
		return details, err
	}

	resp, err := client.Do(req)
	if errorOccured(err) {
		return details, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	response := ASResponse{}
	err = json.Unmarshal(respBody, &response)

	if errorOccured(err) {
		return details, err
	}

	details.Releases = response.Releases
	productDv, ok := response.StorePlatformData["product-dv-product"].(map[string]interface{})

	if !ok {
		return details, nil
	}

	results, ok := productDv["results"].(map[string]interface{})

	if !ok {
		return details, nil
	}

	metaData, ok := results[app.AppstoreId].(map[string]interface{})

	if !ok {
		return details, nil
	}

	details.Name = metaData["name"].(string)
	details.URL = metaData["url"].(string)

	return details, nil
}

func makeFeedForApp(app App, details AppDetails, config Configuration) *feeds.Feed {
	now := time.Now()

	feed := &feeds.Feed{
		Title:       fmt.Sprintf("App Updates: %s", details.Name),
		Link:        &feeds.Link{Href: details.URL},
		Description: fmt.Sprintf("Feed of App Updates for %s", details.Name),
		Author:      &feeds.Author{Name: "App Update Bot", Email: config.FeedAuthorEmail},
		Created:     now,
	}

	feed.Items = []*feeds.Item{}
	for _, release := range details.Releases {
		// fmt.Printf("Release for %s: %q\n\n", release.VersionString, release.ReleaseNotes)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       fmt.Sprintf("App Version %s", release.VersionString),
			Link:        &feeds.Link{Href: MakeEndpointForApp(app)},
			Description: release.ReleaseNotes,
			Created:     release.ReleaseDate,
		})
	}

	return feed
}
