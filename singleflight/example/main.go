package main

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"golang.org/x/sync/singleflight"
)

var errorNotExist = errors.New("not exist")
var g singleflight.Group

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("key")
			if err != nil {
				fmt.Println(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}

func getData(key string) (string, error) {
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		v, err, _ := g.Do(key, func() (interface{}, error) {
			return getDataFromDB(key)
		})
		if err != nil {
			log.Println(err)
			return "", err
		}
		data = v.(string)
	}
	return data, nil
}

func getDataFromCache(key string) (string, error) {
	return "", errorNotExist
}

func getDataFromDB(key string) (string, error) {
	fmt.Printf("get %s from database\n", key)
	return "I Love You!", nil
}
