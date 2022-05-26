package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"strings"
)

func getYaml(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("request failed: %d", resp.StatusCode))
	}
	var config interface{}
	err = yaml.NewDecoder(resp.Body).Decode(&config)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return config, nil
}

func main() {
	http.HandleFunc("/ky", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		url := r.URL.Query().Get("url")
		config, err := getYaml(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		proxies := config.(map[string]interface{})["proxies"].([]interface{})
		var filtered []interface{}
		for _, p := range proxies {
			pp := p.(map[string]interface{})
			if strings.Contains(pp["name"].(string), "日本") &&
				!strings.Contains(pp["name"].(string), "仅海外用户") {
				filtered = append(filtered, p)
			}
		}
		result, err := yaml.Marshal(map[string][]interface{}{
			"proxies": filtered,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	})

	http.HandleFunc("/nex", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		url := r.URL.Query().Get("url")
		config, err := getYaml(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		proxies := config.(map[string]interface{})["proxies"].([]interface{})
		var filtered []interface{}
		for _, p := range proxies {
			pp := p.(map[string]interface{})
			if strings.Contains(pp["name"].(string), "Japan") {
				filtered = append(filtered, p)
			}
		}
		result, err := yaml.Marshal(map[string][]interface{}{
			"proxies": filtered,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
