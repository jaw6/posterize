package main

// given name in input, return poster if one is found
// in, eg http://www.imdbapi.com/?i=&t=inception
// optional: given file name, save poster to file

import (
  "encoding/json"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "regexp"
)

func getNameFromArgs() (name string, err error) {
  if len(os.Args) < 2 {
    err := errors.New("must have a name to search for")
    return name, err
  }

  name = os.Args[1]

  // oddly specific requirement: if name ends in " (2001)", remove that
  r, _ := regexp.Compile(" ([[({]?\\d{4}.?)\\z")
  name = r.ReplaceAllString(name, "")
  return
}

func getFileNameFromArgs() (filename string, err error) {
  if len(os.Args) < 3 {
    err := errors.New("no filename provided")
    return filename, err
  }

  filename = os.Args[2]
  return
}

type result struct {
  Title  string `json:"Title"`
  Poster string `json:"Poster"`
  Year   string `json:"Year"`
}

func queryWithName(name string) (res result, err error) {
  u, err := url.Parse("http://www.omdbapi.com/")
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

func urlToFile(url string, file string) (err error) {
  response, err := http.Get(url)
  if err != nil {
    fmt.Printf("ERROR! %v\n", err)
  }
  defer response.Body.Close()

  target, err := os.Create(file)
  if err != nil {
    fmt.Printf("ERROR! %v\n", err)
  }
  defer target.Close()

  _, err = io.Copy(target, response.Body)
  if err != nil {
    fmt.Printf("Err! %v\n", err)
  }
  return
}

func main() {
  name, err := getNameFromArgs()
  filename, _ := getFileNameFromArgs()

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

  if poster != "" && filename != "" {
    err = urlToFile(poster, filename)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  } else {
    fmt.Println(poster)
  }
}
