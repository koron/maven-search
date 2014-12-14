package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func searchArtifact(q string, n int) error {
	p := url.Values{
		"rows": []string{strconv.Itoa(n)},
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
	// parse JSON with anonymous struct.
	d := json.NewDecoder(resp.Body)
	var v struct {
		Response struct {
			Docs []struct {
				ID            string
				LatestVersion string
			}
		}
	}
	if err := d.Decode(&v); err != nil {
		return err
	}
	// print results.
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
	if err := searchArtifact(os.Args[1], 20); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
