package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/ky", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		url := r.URL.Query().Get("url")
		resp, err := http.Get(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(resp.StatusCode)
			return
		}
		var config interface{}
		err = yaml.NewDecoder(resp.Body).Decode(&config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()
		proxies := config.(map[string]interface{})["proxies"].([]interface{})
		var filtered []interface{}
		for _, p := range proxies {
			pp := p.(map[string]interface{})
			if strings.Contains(pp["name"].(string), "日本") &&
				!strings.Contains(pp["name"].(string), "仅海外用户"){
				filtered = append(filtered, p)
			}
		}
		result, err := yaml.Marshal(filtered)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
