package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func lights(w http.ResponseWriter, r *http.Request) {
	lights := getLights()

	js, err := json.Marshal(lights)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func readConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // path to look for the config file in
	viper.SetDefault("PORT", "8080")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	readConfig()

	http.HandleFunc("/lights", lights)
	http.HandleFunc("/healthcheck", healthCheck)
	log.Fatal(http.ListenAndServe(":"+viper.GetString("PORT"), nil))
}
