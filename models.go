package main

import "time"

// Go App Models
type Configuration struct {
	FeedAuthorEmail string `json:"feedAuthorEmail"`
	OutputDirectory string `json:"outputDir"`
	Apps            []App  `json:"apps"`
	Summary         *struct {
		OutputPath string `json:"outputPath"`
		BaseURL    string `json:"baseURL"`
	} `json:"summary"`
}

type App struct {
	ShortName  string `json:"shortName"`
	AppstoreId string `json:"appstoreId"`
	CountryId  string `json:"countryId"`
}

type AppDetails struct {
	Name     string
	URL      string
	Releases []AppRelease
}

type Storefront struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	StorefrontID int    `json:"storefrontId"`
}

// App Store Response Models
type ASResponse struct {
	StorePlatformData map[string]interface{} `json:"storePlatformData"`
	Releases          []AppRelease           `json:"versionHistory"`
}

type AppRelease struct {
	ReleaseNotes  string    `json:"releaseNotes"`
	VersionString string    `json:"versionString"`
	ReleaseDate   time.Time `json:"releaseDate"`
}
