package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/beego/mux"
	"github.com/toukii/goutils"

	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	mx := mux.New()
	var data []*Data
	mx.Get("/1", func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		v := req.FormValue("v")
		fmt.Println(v)
		if v != "" {
			data = fetch(v)
		}
		rw.Write(goutils.ReadFile("i1.html"))
	})
	mx.Get("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(goutils.ReadFile("index.html"))
	})
	mx.Get("/v1.png", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(goutils.ReadFile("v1.png"))
	})
	mx.Get("/circle-split.min.js", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(goutils.ReadFile("circle-split.min.js"))
	})
	mx.Get("/r", func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		v := req.FormValue("v")
		fmt.Println(v)
		if v == "" {
			return
		}
		data = fetch(v)
		http.Redirect(rw, req, "/", 302)
	})
	mx.Get("/data.json", func(rw http.ResponseWriter, req *http.Request) {
		bs, _ := json.Marshal(data)
		rw.Write(bs)
	})

	http.ListenAndServe(":80", mx)
}

type Data struct {
	Value string `json:"value"`
}

func fetch(url string) []*Data {
	doc, err := goquery.NewDocument(url)
	if goutils.CheckErr(err) {
		return nil
	}
	ret := make([]*Data, 0, 10)
	wg := &sync.WaitGroup{}
	// Find the review items
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(s *goquery.Selection, wg *sync.WaitGroup) {
			defer wg.Done()
			img, e := s.Attr("src")
			// fmt.Println(img, e)
			if e {
				ret = append(ret, &Data{Value: img})
			}
		}(s, wg)
		wg.Wait()
	})
	// fmt.Println(ret)
	return ret
}
