package main

import "testing"

func TestMakeEndpoint(t *testing.T) {

	t.Run("test de app", func(t *testing.T) {
		app := App{
			ShortName:  "moia",
			AppstoreId: "1373271535",
			CountryId:  "de",
		}
		got := MakeEndpointForApp(app)
		want := "https://apps.apple.com/de/app/id1373271535"

		if got != want {
			t.Errorf("got %q but expected %q", got, want)
		}
	})

	t.Run("test us app", func(t *testing.T) {
		app := App{
			ShortName:  "duolingo",
			AppstoreId: "570060128",
			CountryId:  "us",
		}
		got := MakeEndpointForApp(app)
		want := "https://apps.apple.com/us/app/id570060128"

		if got != want {
			t.Errorf("got %q but expected %q", got, want)
		}
	})
}
