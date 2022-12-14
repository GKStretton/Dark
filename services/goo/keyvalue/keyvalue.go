package keyvalue

import (
	"fmt"
	"sync"

	"github.com/gkstretton/dark/services/goo/mqtt"
)

const TOPIC_ROOT = "asol/kv/"
const TOPIC_SET = TOPIC_ROOT + "set/"
const TOPIC_SET_RESP = TOPIC_ROOT + "set-resp/"
const TOPIC_GET = TOPIC_ROOT + "get/"
const TOPIC_GET_RESP = TOPIC_ROOT + "get-resp/"

var (
	lock = &sync.Mutex{}
)

func Start() {
	// set
	mqtt.Subscribe(TOPIC_SET+"#", setCallback)
	// req
	mqtt.Subscribe(TOPIC_GET+"#", reqCallback)
}

func Get(key string) []byte {
	b, err := getKeyValue(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func Set(key string, value []byte) {
	err := setKeyValue(key, value)
	if err != nil {
		fmt.Println(err)
	}
}
