package cipherserver

import (
  "bytes"
  "fmt"
  "net/http"

  "github.com/freignat91/cipher3/core"
)

//NewCipherServer .
func NewCipherServer(name string, host string, port string) (*CipherServer, error) {
  s := CipherServer{name: name, host: host, port: port}
  key, err := core.CreateKey(keyDimension, keySize, nil, false, false)
  if err != nil {
    return nil, fmt.Errorf("error creating server private key: %v\n", err)
  }
  s.privateKey = key
  s.targetMap = make(map[string]*CipherTarget)
  return &s, nil
}

//StopCipherServer .
func StopCipherServer(host string, port string) error {
  _, err := http.Get("http://" + host + ":" + port + "/api/v1/stop")
  return err
}

//Handcheck .
func Handcheck(fromHost string, fromPort string, toHost string, toPort string) error {
  _, err := http.Get("http://" + fromHost + ":" + fromPort + "/api/v1/sendHandcheck/" + toHost + "/" + toPort)
  return err
}

//SendData .
func SendData(from string, to string, data string) error {
  _, err := http.Post("http://"+from+"/api/v1/sendData/"+to, "text/plain", bytes.NewBuffer([]byte(data)))
  return err
}
