package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type WeatherData struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temperature float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not set")
	}

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		city := r.FormValue("city")
		if city == "" {
			http.Error(w, "Please specify a city", http.StatusBadRequest)
			return
		}
		weatherInfo, err := getWeatherInfo(city, apiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := fmt.Sprintf("Current weather in %s is: %s with %s, temperature is %.2f°C", city, weatherInfo.Weather[0].Main, weatherInfo.Weather[0].Description, weatherInfo.Main.Temperature)
		fmt.Fprintln(w, response)
	})

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
				<label for="city">都市名を入力してください:</label>
				<input type="text" id="city" name="city">
				<input type="submit" value="Get Weather">
			</form>
		</body>
		</html>
	`))
	tmpl.Execute(w, nil)
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
