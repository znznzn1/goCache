package main

import (
	"fmt"
	"goCache/Cache"
)

func main() {
	cache := Cache.NewMemCache()
	cache.SetMaxMemory("200MB")
	cache.Set("int", 1, 100)
	cache.Set("bool", false, 100)
	cache.Set("data", map[string]interface{}{"a": 1}, 100)
	cache.Get("int")
	cache.Del("int")
	fmt.Println(cache.Exists("data"))
	fmt.Println(cache.Keys())
	cache.Flush()
	fmt.Println(cache.Keys())

}
