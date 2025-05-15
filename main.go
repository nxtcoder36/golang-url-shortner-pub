package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/nxtcoder36/golang-url-shortner-pub/cache"
	"github.com/nxtcoder36/golang-url-shortner-pub/db"
)

func main() {
	dbImpl := db.UrlShortnerImpl()
	cacheImpl, err := cache.RedisCacheImpl()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/url", createUrlShortner(dbImpl))
	mux.HandleFunc("/", redirectLongUrl(dbImpl, cacheImpl))
	mux.HandleFunc("/testing-long-url-working-or-not-locally", testLocally)

	fmt.Println("Server is running on port " + os.Getenv("PORT"))
	if err := http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), mux); err != nil {
		panic(err)
	}
}

func createUrlShortner(dbImpl db.UrlShortnerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// if long_url received in params
		// long_url := r.URL.Query().Get("long_url")
		// if long_url == "" {
		// 	http.Error(w, "long_url is missing", http.StatusBadRequest)
		// 	return
		// }
		// fmt.Println(long_url)

		reader, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if len(reader) == 0 {
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}

		data := struct {
			LongUrl string `json:"long_url"`
		}{}

		if err := json.Unmarshal(reader, &data); err != nil {
			http.Error(w, "request body is missing or invalid long_url", http.StatusBadRequest)
			return
		}

		url, err := dbImpl.Insert(data.LongUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(url)
		json.NewEncoder(w).Encode(map[string]string{
			"short_url": url,
		})
	}
}

func redirectLongUrl(dbImpl db.UrlShortnerInterface, cacheImpl cache.RedisCacheInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		longUrl, err := cacheImpl.Get(r.URL.Path)
		if err == nil || longUrl != "" {
			http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
			return
		}
		fmt.Println("Couldn't find the key in redis", r.URL.Path)
		longUrl, err = dbImpl.Find(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if longUrl == "" {
			http.Error(w, "short URL not found", http.StatusNotFound)
			return
		}

		_ = cacheImpl.Set(r.URL.Path, longUrl)

		http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
	}
}

func testLocally(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("working"))
}
