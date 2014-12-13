package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// JSONResponse solrsearch response.
type JSONResponse struct {
	Response Response
}

// Response body.
type Response struct {
	Docs []Doc
}

// Doc shows an artifact.
type Doc struct {
	ID            string
	LatestVersion string
}

func searchArtifact(q string) error {
	p := url.Values{
		"rows": []string{"20"},
		"wt":   []string{"json"},
		"q":    []string{q},
	}
	u := "http://search.maven.org/solrsearch/select?" + p.Encode()
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	d := json.NewDecoder(resp.Body)
	v := JSONResponse{}
	if err := d.Decode(&v); err != nil {
		return err
	}
	for _, d := range v.Response.Docs {
		fmt.Println(d.ID + ":" + d.LatestVersion)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: need a query string")
		os.Exit(1)
	}
	if err := searchArtifact(os.Args[1]); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
