package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"strings"
)

type Proxy = *yaml.Node

type Filter = func(Proxy) bool

type Mutation = func(Proxy)

func FilterStringField(key string, value string) Filter {
	return func(p Proxy) bool {
		if p.Kind != yaml.MappingNode {
			return false
		}
		kv := map[string]interface{}{}
		_ = p.Decode(&kv)
		if v, ok := kv[key]; ok {
			if strings.Contains(v.(string), value) {
				return true
			}
		}
		return false
	}
}

func MutationSet(key string, value interface{}) Mutation {
	return func(p Proxy) {
		if p.Kind != yaml.MappingNode {
			return
		}
		payload := map[string]interface{}{}
		_ = p.Decode(&payload)
		payload[key] = value
		_ = p.Encode(payload)
	}
}

func FilterNot(f Filter) Filter {
	return func(p Proxy) bool {
		return !f(p)
	}
}

func FilterAll(fs ...Filter) Filter {
	return func(p Proxy) bool {
		for _, f := range fs {
			if !f(p) {
				return false
			}
		}
		return true
	}
}

func handle(f Filter, ms ...Mutation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		resp, err := http.Get(url)
		if hasError(err, w) {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			hasError(errors.New(fmt.Sprintf("request failed: %d", resp.StatusCode)), w)
			return
		}
		data, err := io.ReadAll(resp.Body)
		if hasError(err, w) {
			return
		}
		proxies, err := getProxies(data)
		if hasError(err, w) {
			return
		}
		filtered := ProxiesFilter(proxies, f)
		for _, m := range ms {
			for _, p := range filtered {
				m(p)
			}
		}
		result, err := yaml.Marshal(map[string][]Proxy{
			"proxies": filtered,
		})
		if hasError(err, w) {
			return
		}
		w.Write(result)
	}
}

func getProxies(data []byte) ([]*yaml.Node, error) {
	var config map[string]yaml.Node
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	proxies := config["proxies"]
	return proxies.Content, nil
}

func hasError(err error, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}

func ProxiesFilter(ps []Proxy, filter Filter) []Proxy {
	var ret []Proxy
	for _, p := range ps {
		if filter(p) {
			ret = append(ret, p)
		}
	}
	return ret
}

var (
	KyFilter   Filter
	NexFilter  Filter
	GaFilter   Filter
	GaMutation Mutation
)

func init() {
	KyFilter = FilterAll(
		FilterStringField("name", "日本"),
		FilterNot(FilterStringField("name", "仅海外用户")),
		FilterNot(FilterStringField("name", "SS")),
	)
	NexFilter = FilterStringField("name", "Japan")
	GaFilter = FilterStringField("name", "日本")
	GaMutation = MutationSet("udp", true)
}

func main() {
	http.HandleFunc("/ky", handle(KyFilter))
	http.HandleFunc("/nex", handle(NexFilter))
	http.HandleFunc("/ga", handle(GaFilter, GaMutation))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
