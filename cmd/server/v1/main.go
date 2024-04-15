package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type WeatherData struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temperature float64 `json:"temp"`
	} `json:"main"`
	City string `json:"-"`
}

func main() {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not set")
	}

	http.HandleFunc("/", basicAuth(formHandler, "tryu", "k8s"))
	http.HandleFunc("/weather", weatherHandler)
	fmt.Println("Starting server at port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("form").Parse(`
		<html>
		<head>
			<title>Weather App</title>
		</head>
		<body>
			<form action="/weather" method="post">
				<label for="cities">都市名をカンマ区切りで入力してください:</label>
				<input type="text" id="cities" name="cities">
				<input type="submit" value="天気を見る">
			</form>
		</body>
		</html>
	`))
	tmpl.Execute(w, nil)
}

// Basic認証を行うミドルウェア
func basicAuth(handler http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "v1-Unauthorized.", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	cities := r.FormValue("cities")
	if cities == "" {
		http.Error(w, "都市名を最低一つ入れてください。", http.StatusBadRequest)
		return
	}
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	cityList := strings.Split(cities, ",")
	weatherInfos, err := getWeatherInfos(cityList, apiKey)
	if err != nil {
		http.Error(w, "天気情報の取得に失敗しました。", http.StatusInternalServerError)
		return
	}
	for _, weatherInfo := range weatherInfos {
		response := fmt.Sprintf("現在の%sの天気は %s ,気温は %.2f°Cです。\n", weatherInfo.City, weatherInfo.Weather[0].Main, weatherInfo.Main.Temperature)
		fmt.Fprintln(w, response)
	}
}

func getWeatherInfos(cities []string, apiKey string) ([]*WeatherData, error) {
	var wg sync.WaitGroup
	ch := make(chan *WeatherData)
	errCh := make(chan error)

	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			data, err := getWeatherInfo(city, apiKey)
			if err != nil {
				errCh <- err
				return
			}
			data.City = city
			ch <- data
		}(city)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(errCh)
	}()

	var weatherInfos []*WeatherData
	for info := range ch {
		weatherInfos = append(weatherInfos, info)
	}

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return weatherInfos, nil
}

func getWeatherInfo(city, apiKey string) (*WeatherData, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&lang=ja&units=metric", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to OpenWeatherMap failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %v", err)
	}

	var data WeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshaling failed: %v", err)
	}

	return &data, nil
}
