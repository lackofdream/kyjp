package main

import (
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
)

const data = `
port: 7890
socks-port: 7891
redir-port: 7892
allow-lan: false
mode: rule
log-level: silent
external-controller: '0.0.0.0:9090'
proxies:
  -
    name: '🇭🇰 香港 01 [ C1 ] V2Ray'
    type: vmess
    server: 1551.com
    port: 3005
    uuid: uuid1551
    alterId: '0'
    cipher: auto
    udp: false
  -
    name: '🇭🇰 香港 02 [ C1 ] V2Ray'
    type: vmess
    server: 1552.com
    port: 3005
    uuid: uuid1552
    alterId: '0'
    cipher: auto
    udp: false`

func TestFilterStringField(t *testing.T) {
	ps, err := getProxies([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	for idx, n := range ps {
		res := FilterStringField("name", "01")(n)
		if idx == 0 && res != true {
			t.Error()
		}
		if idx == 1 && res != false {
			t.Error()
		}
	}
}

func TestMutationSet(t *testing.T) {

	ps, err := getProxies([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	p := ps[0]
	p = MutationSet("udp", true)(p)
	data, err := yaml.Marshal(p )
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "udp: true") {
		t.Error()
	}
}