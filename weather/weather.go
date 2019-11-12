package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "log"
  "net/http"
  "encoding/json"
)

// Structures

type temperature struct {
  Value float64
  Unit string
}

type weather struct {
  DateTime, IconPhrase string
  Temperature temperature
  PrecipitationProbability int
}

// Accuweather API

func APIKey() string {
  file, err := os.Open("key.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  b, err := ioutil.ReadAll(file)
  return string(b[:len(b)-1])
}

func generateURL() string {
  key := APIKey()
  req, _ := http.NewRequest(
    "GET",
    "http://dataservice.accuweather.com/forecasts/v1/hourly/12hour/56186",
    nil)
  q := req.URL.Query()
  q.Add("apikey", key)
  q.Add("metric", "true")
  q.Add("details", "true")
  req.URL.RawQuery = q.Encode()
  return req.URL.String()
}

func formatHour(h weather) string {
  hour := h.DateTime[11:16]
  return fmt.Sprintf("%s \t %s \t %.1f%s \t %d%%\n",
    hour,
    h.IconPhrase,
    h.Temperature.Value,
    h.Temperature.Unit,
    h.PrecipitationProbability,
  )
}

// Local webserver
func handler (w http.ResponseWriter, r *http.Request) {
  // Get result
  url := generateURL()
  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  body, err := ioutil.ReadAll(res.Body)
  var hours []weather
  e := json.Unmarshal(body, &hours)
  if e != nil {
    log.Fatal(err)
  }
  w.Header().Set("Content-Type", "text/html")
  fmt.Fprintf(w, "<h1>Weather</h1><ul>")
  for _, hour := range hours {
    fmt.Fprintf(w, "<li>%s</li>", formatHour(hour))
  }
}

func main () {
  http.HandleFunc("/", handler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
