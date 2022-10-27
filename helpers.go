package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/feeds"
)

func errorOccured(e error) bool {
	return e != nil
}

func getConfiguration(configPath *string) (config Configuration, err error) {
	configFile, err := os.Open(*configPath)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)

	if errorOccured(err) {
		fmt.Println(err)
		return config, err
	}

	return config, nil
}

func getStorefronts(path *string) (mappings []Storefront, err error) {
	file, err := os.Open(*path)
	if err != nil {
		return mappings, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&mappings)

	if errorOccured(err) {
		fmt.Println(err)
		return mappings, err
	}

	return mappings, nil
}

func getMatchingStorefront(app App, storefronts []Storefront) (sf Storefront, err error) {
	for _, storefront := range storefronts {
		if strings.EqualFold(storefront.Code, app.CountryId) {
			return storefront, nil
		}
	}
	return sf, fmt.Errorf("no matching storefront found for country code %q", app.CountryId)
}

func writeFeeds(feed *feeds.Feed, app App, config Configuration) {
	err := writeAtomFeed(feed, app, config)
	if errorOccured(err) {
		log.Fatalf("Couldn't write Atom Feed for %s: %v", app.ShortName, err.Error())
	} else {
		log.Printf("Successfully wrote Atom Feed for %s", app.ShortName)
	}

	err = writeJSONFeed(feed, app, config)
	if errorOccured(err) {
		log.Fatalf("Couldn't write JSON Feed for %s: %v", app.ShortName, err.Error())
	} else {
		log.Printf("Successfully wrote JSON Feed for %s", app.ShortName)
	}

	err = writeRssFeed(feed, app, config)
	if errorOccured(err) {
		log.Fatalf("Couldn't write RSS Feed for %s: %v", app.ShortName, err.Error())
	} else {
		log.Printf("Successfully wrote RSS Feed for %s", app.ShortName)
	}
}

func writeAtomFeed(feed *feeds.Feed, app App, config Configuration) error {
	file, err := os.Create(makeAtomOutputPath(app, config))
	if errorOccured(err) {
		return err
	}

	writer := bufio.NewWriter(file)
	feed.WriteAtom(writer)
	writer.Flush()
	return nil
}

func writeRssFeed(feed *feeds.Feed, app App, config Configuration) error {
	file, err := os.Create(makeRssOutputPath(app, config))
	if errorOccured(err) {
		return err
	}

	writer := bufio.NewWriter(file)
	feed.WriteRss(writer)
	writer.Flush()
	return nil
}

func writeJSONFeed(feed *feeds.Feed, app App, config Configuration) error {
	file, err := os.Create(makeJSONOutputPath(app, config))
	if errorOccured(err) {
		return err
	}

	writer := bufio.NewWriter(file)
	feed.WriteJSON(writer)
	writer.Flush()
	return nil
}

func MakeEndpointForApp(app App) string {
	return fmt.Sprintf("https://apps.apple.com/%s/app/id%s", app.CountryId, app.AppstoreId)
}

func MakeStorefrontHeader(storefront Storefront) string {
	return fmt.Sprintf("%s,20", strconv.Itoa(storefront.StorefrontID))
}

func makeBaseOutputPath(app App, config Configuration) string {
	filename := fmt.Sprintf("%s-%s", app.ShortName, app.CountryId)
	return filepath.Join(config.OutputDirectory, filename)
}

func makeJSONOutputPath(app App, config Configuration) string {
	return makeBaseOutputPath(app, config) + ".json"
}

func makeAtomOutputPath(app App, config Configuration) string {
	return makeBaseOutputPath(app, config) + ".atom"
}

func makeRssOutputPath(app App, config Configuration) string {
	return makeBaseOutputPath(app, config) + ".xml"
}
