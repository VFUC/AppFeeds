
# AppFeeds

RSS Feeds for iOS App Store Updates

## TL;DR
This project provides two things:

- a tool to generate RSS/Atom/JSON-Feeds for the version history of given iOS Apps
- a Github Action, which periodically creates new feeds for the apps specified in this repo

Example feed: https://feeds.jonas.lol/apps/duolingo-us.xml 


## How can I use this?

- You can run the tool to generate and publish your own feeds.
- You can use the feeds generated and published by this repo. The currently published feeds can be found [here](current-feeds.md). If an app is missing, feel free to add it to [the config](config.json) and open a PR.

## Getting Started

- Check if the values in `config.json` meet your expectations
- build a binary, e.g. using `go build -o feedmaker` 
- run the binary, e.g. `./feedmaker`
- find the generated feeds in the folder specified in `config.json`

### Arguments

The app needs a configuration file and a storefront mappings file to run. Examples can be found in this repo.

By default, the tool expects these in the same directory as it's being run from, but this can be controlled when ran with the following flags:

```
-c path/to/my/config.json
-sf path/to/my/storefrontmappings.json
```

## Configuration values

Based on the configuration values, the tool will fetch the version history for each of the `apps` listed in the configuration. It will generate RSS, Atom and JSON feeds and write them into the specified `outputDir`.

If the `summary` configuration is set, a markdown file with published feeds will be generated. This may be useful when combined with a publishing workflow, like the GitHub Actions of this repository.

| Configuration Key  | Value | |
| ------------- | ------------- | ------------- |
| `feedAuthorEmail`  | The feed author, used in Atom and RSS feeds  | required |
| `outputDir`  | Where to store the generated feeds. **Must exist.**   | required |
| `apps`  | List of apps to generate feeds for   | required |
| `apps.shortName`  | A name for the app, used to name the output files   | required |
| `apps.appstoreId`  | The app store identifier of the app, see [here](#adding-an-app) for details.  | required |
| `apps.countryId`  | The app store to check. App must exist in the country's store to succeed    | required |
| `summary`  | Values used to generate a summary markdown file such as [this one](current-feeds.md)    | optional |
| `summary.outputPath`  | Where to write the summary file to    |  |
| `summary.baseURL`  | Which base URL to use in the summary file    |  |

## Adding an App

To get an app's `appstoreId` and `countryId`, lookup it's App Store URL using your favorite search engine.

As an example, if the URL is https://apps.apple.com/us/app/duolingo-language-lessons/id570060128 we can derive

`countryId`: `us`

`appstoreId` : `570060128`

## References 

I used [this StackOverflow answer](https://stackoverflow.com/questions/12273811/how-do-i-check-my-ios-app-version-history-detail-on-itunesconnect/48098811#48098811) to get JSON-formatted app version histories and [this Gist](https://gist.github.com/BrychanOdlum/2208578ba151d1d7c4edeeda15b4e9b1) to map App Store countries to their "storefrontId's".

To generate the feeds, the [feeds](https://github.com/gorilla/feeds) package is being used.