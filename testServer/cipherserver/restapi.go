package cipherserver

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "time"

  "github.com/freignat91/cipher3/core"
  "github.com/gorilla/mux"
)

//TargetDoc .
type TargetDoc struct {
  Name string `json:"name"`
  Host string `json:"host"`
  Port string `json:"port"`
  Data string `json:"data"`
}

//DataDoc .
type DataDoc struct {
  From string `json:"name"`
  Data string `json:"data"`
}

func (s CipherServer) apiStop(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "stopping server\n")
  fmt.Println("server stopped")
  go func() {
    time.Sleep(1 * time.Second)
    os.Exit(0)
  }()
}

func (s CipherServer) apiCreateKey(w http.ResponseWriter, r *http.Request) {
  key, err := core.CreateKey(keyDimension, keySize, nil, false, false)
  if err != nil {
    fmt.Fprintf(w, "error: %v\n", err)
    //w.WriteHeader(422)
    return
  }
  s.privateKey = key
  s.targetMap = make(map[string]*CipherTarget)
  fmt.Fprintf(w, "done\n")
}

func (s CipherServer) apiSendHandcheck(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  host := vars["host"]
  port := vars["port"]
  if err := s.sendHandcheck(host, port); err != nil {
    fmt.Printf("Error: %v\n", err)
    fmt.Fprintf(w, "error: %v\n", err)
    //w.WriteHeader(422)
    return
  }
  fmt.Fprintf(w, "send handcheck done\n")
}

func (s CipherServer) apiHandcheck(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var t TargetDoc
  decoder.Decode(&t)
  if err := s.handcheck(&t); err != nil {
    fmt.Printf("Error: %v\n", err)
    fmt.Fprintf(w, "error: %v\n", err)
    //w.WriteHeader(422)
    return
  }
  fmt.Fprintf(w, "handcheck done\n")
}

func (s CipherServer) apiSendData(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  targetName := vars["target"]
  defer r.Body.Close()
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  if err := s.sendData(targetName, data); err != nil {
    fmt.Printf("Error: %v\n", err)
    fmt.Fprintf(w, "error: %v\n", err)
    //w.WriteHeader(422)
    return
  }
  fmt.Fprintf(w, "data sent\n")
}

func (s CipherServer) apiReceiveData(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  targetName := vars["target"]
  defer r.Body.Close()
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  if err := s.receiveData(targetName, data); err != nil {
    fmt.Printf("Error: %v\n", err)
    fmt.Fprintf(w, "error: %v\n", err)
    //w.WriteHeader(422)
    return
  }
  fmt.Fprintf(w, "data received\n")
}

func (s CipherServer) apiSendAuthentication(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  targetName := vars["target"]
  if err := s.sendAuthentication(targetName); err != nil {
    fmt.Printf("Error: %v\n", err)
    fmt.Fprintf(w, "error: %v\n", err)
    //w.WriteHeader(422)
    return
  }
  fmt.Fprintf(w, "data sent\n")
}

func (s CipherServer) apiAuthenticate(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  targetName := vars["target"]
  defer r.Body.Close()
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  if err := s.authenticate(targetName, data); err != nil {
    fmt.Printf("Error: %v\n", err)
    fmt.Fprintf(w, "error: %v\n", err)
    //w.WriteHeader(422)
    return
  }
  fmt.Fprintf(w, "data sent\n")
}
