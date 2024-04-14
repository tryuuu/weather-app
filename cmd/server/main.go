// http://localhost:8080/?city=TOKYO

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv" // godotenv パッケージをインポート
)

// OpenWeatherMap API からの JSON 応答に合わせて設計された構造体
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
	// .env ファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("API key not set")
	}

	// HTTPサーバーの設定
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		if city == "" {
			http.Error(w, "Please specify a city", http.StatusBadRequest)
			return
		}
		weatherInfo, err := getWeatherInfo(city, apiKey)
		fmt.Printf("weatherInfo: %v\n", weatherInfo)
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

// getWeatherInfo はOpenWeatherMap APIから天気情報を取得する関数です。
func getWeatherInfo(city, apiKey string) (*WeatherData, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&lang=ja&units=metric", city, apiKey)
	resp, err := http.Get(url)
	fmt.Println(resp)
	if err != nil {
		return nil, fmt.Errorf("request to OpenWeatherMap failed: %v", err)
	}
	// 関数が終了する際に必ずレスポンスのボディを閉じる
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %v", err)
	}

	var data WeatherData
	// バイト列のbodyを受け取り適切に変換しWeatherData型のdataポインタに格納する
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshaling failed: %v", err)
	}
	fmt.Println(data)

	return &data, nil
}
