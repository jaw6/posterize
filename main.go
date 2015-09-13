package main

// given name in input, return poster if one is found
// in, eg http://www.imdbapi.com/?i=&t=inception

import (
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  // "path/filepath"
  // "strconv"
)

func getNameFromArgs() (name string, err error){
  if len(os.Args) < 2 {
    err := errors.New("must have a name to search for")
    return name, err
  }

  name = os.Args[len(os.Args) - 1]
  return
}

type result struct {
  Title  string `json:"Title"`
  Poster string `json:"Poster"`
  Year   string  `json:"Year"`
}

func queryWithName(name string) (res result, err error) {
  u, err := url.Parse("http://www.imdbapi.com/")
  params := url.Values{}
  params.Add("t", name)
  params.Add("i", "")
  u.RawQuery = params.Encode()

  response, err := http.Get(u.String())
  if err != nil {
    fmt.Println("Error while downloading", u.String(), "-", err)
    return
  }

  defer response.Body.Close()

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    fmt.Println(err)
  }

  err = json.Unmarshal(body, &res)
  return
}

func main() {
  name, err := getNameFromArgs()

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  result, err := queryWithName(name)

  if err != nil {
    fmt.Println("No poster, because", err)
    os.Exit(1)
  }

  poster := result.Poster
  fmt.Println(poster)
}