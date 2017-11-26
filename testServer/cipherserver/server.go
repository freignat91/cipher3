package cipherserver

import (
  "bytes"
  "encoding/json"
  "fmt"
  "log"
  "net/http"

  "github.com/freignat91/cipher3/core"
  "github.com/gorilla/mux"
)

var (
  keyDimension = 3
  keySize      = 16384
)

//CipherServer .
type CipherServer struct {
  name       string
  host       string
  port       string
  privateKey *core.CipherCore
  targetMap  map[string]*CipherTarget
}

//CipherTarget .
type CipherTarget struct {
  name           string
  host           string
  port           string
  exchangeKey    *core.CipherCore
  publicKey      *core.CipherCore
  handcheckDone  bool
  handcheckPhase int
}

//Start .
func (s *CipherServer) Start(version string) error {
  router := mux.NewRouter().StrictSlash(true)
  s.handleAPIFunctions(router)
  log.Printf("test server [%s] starting on %s\n", s.name, s.port)
  return http.ListenAndServe(":"+s.port, router)
}

func (s *CipherServer) handleAPIFunctions(r *mux.Router) {
  r.HandleFunc("/api/v1/stop", s.apiStop).Methods("GET")
  r.HandleFunc("/api/v1/createKey", s.apiCreateKey).Methods("GET")
  r.HandleFunc("/api/v1/sendHandcheck/{host}/{port}", s.apiSendHandcheck).Methods("GET")
  r.HandleFunc("/api/v1/handcheck", s.apiHandcheck).Methods("POST")
  r.HandleFunc("/api/v1/sendData/{target}", s.apiSendData).Methods("POST")
  r.HandleFunc("/api/v1/receiveData/{target}", s.apiReceiveData).Methods("POST")
  r.HandleFunc("/api/v1/sendAuthentication/{target}", s.apiSendAuthentication).Methods("GET")
  r.HandleFunc("/api/v1/receiveData/authentication/{target}", s.apiAuthenticate).Methods("POST")
}

func (s *CipherServer) sendHandcheck(host string, port string) error {
  doc := TargetDoc{Name: s.name, Host: s.host, Port: s.port}
  jsonValue, _ := json.Marshal(doc)
  _, err := http.Post("http://"+host+":"+port+"/api/v1/handcheck", "application/json", bytes.NewBuffer(jsonValue))
  return err
}

func (s *CipherServer) handcheck(doc *TargetDoc) error {
  target, ok := s.targetMap[doc.Name]
  if !ok {
    target = &CipherTarget{name: doc.Name, host: doc.Host, port: doc.Port}
    s.targetMap[doc.Name] = target
  }
  //command handcheck
  if doc.Data == "" {
    fmt.Printf("%s: received command handcheck to %s\n", s.name, target.name)
    puk, errc := s.createPublicKey(target)
    if errc != nil {
      target.handcheckDone = false
      target.handcheckPhase = 0
      return errc
    }
    target.publicKey = puk.Copy()
    target.handcheckPhase = 1
    puk.Xor(s.privateKey)
    return s.sendKey(target, puk)
  }
  //time.Sleep(1 * time.Second)
  //Handcheck protocol
  fmt.Printf("%s: handcheck phase %d with %s\n", s.name, target.handcheckPhase, target.name)
  tmp, errk := core.ReadKey(doc.Data)
  if errk != nil {
    target.handcheckDone = false
    target.handcheckPhase = 0
    return errk
  }
  //Phase 1 and 2
  if target.handcheckPhase == 1 || target.handcheckPhase == 2 {
    tmp.Xor(s.privateKey)
    target.handcheckPhase++
    return s.sendKey(target, tmp)
    //Phase 3
  } else if target.handcheckPhase == 3 {
    tmp.Xor(s.privateKey)
    fmt.Printf("%s: PartnerPublicKey: %s...\n", s.name, tmp.ToString()[5:32])
    target.exchangeKey = target.publicKey.Copy()
    target.exchangeKey.Xor(tmp)
    target.handcheckPhase = 0
    target.handcheckDone = true
    fmt.Printf("%s: handcheck with %s done, private exchangeKey: %s...\n", s.name, target.name, target.exchangeKey.ToString()[5:32])
    //Phase 0
  } else if target.handcheckPhase == 0 {
    tmp.Xor(s.privateKey)
    target.handcheckPhase = 4
    return s.sendKey(target, tmp)
    //Phase 4
  } else if target.handcheckPhase == 4 {
    tmp.Xor(s.privateKey)
    fmt.Printf("%s: PartnerPublicKey: %s...\n", s.name, tmp.ToString()[5:32])
    puk, errc := s.createPublicKey(target)
    if errc != nil {
      target.handcheckDone = false
      target.handcheckPhase = 0
      return errc
    }
    target.publicKey = puk.Copy()
    target.exchangeKey = puk.Copy()
    target.exchangeKey.Xor(tmp)
    fmt.Printf("%s: handcheck with %s done, private exchangeKey: %s...\n", s.name, target.name, target.exchangeKey.ToString()[5:32])
    target.handcheckPhase = 5
    puk.Xor(s.privateKey)
    return s.sendKey(target, puk)
    //Phase 5
  } else if target.handcheckPhase == 5 {
    tmp.Xor(s.privateKey)
    target.handcheckDone = true
    target.handcheckPhase = 0
    return s.sendKey(target, tmp)
  }
  return nil
}

func (s *CipherServer) createPublicKey(target *CipherTarget) (*core.CipherCore, error) {
  tmp, errc := core.CreateKey(keyDimension, keySize, nil, false, false)
  if errc != nil {
    return nil, errc
  }
  tmp.Xor(s.privateKey)
  fmt.Printf("%s: PublicKey: %s...\n", s.name, tmp.ToString()[5:32])
  return tmp, nil
}

func (s *CipherServer) sendData(targetName string, data []byte) error {
  target, ok := s.targetMap[targetName]
  if !ok {
    return fmt.Errorf("target %s is unknown", targetName)
  }
  if !target.handcheckDone {
    return fmt.Errorf("no handcheck has been done with target %s", targetName)
  }
  target.exchangeKey.Cipher(nil, data)
  fmt.Printf("Send data to %s (%s:%s)\n", target.name, target.host, target.port)
  _, err := http.Post("http://"+target.host+":"+target.port+"/api/v1/receiveData/"+s.name, "text/plain", bytes.NewBuffer(data))
  return err
}

func (s *CipherServer) receiveData(targetName string, data []byte) error {
  target, ok := s.targetMap[targetName]
  if !ok {
    return fmt.Errorf("target %s is unknown", targetName)
  }
  if !target.handcheckDone {
    return fmt.Errorf("no handcheck has been done with target %s", targetName)
  }
  target.exchangeKey.Cipher(nil, data)
  fmt.Printf("Received data from %s (%s:%s): %s\n", target.name, target.host, target.port, string(data))
  return nil
}

func (s *CipherServer) sendKey(target *CipherTarget, key *core.CipherCore) error {
  skey := key.ToString()
  doc := TargetDoc{Name: s.name, Host: s.host, Port: s.port, Data: skey}
  jsonValue, _ := json.Marshal(doc)
  _, err := http.Post("http://"+target.host+":"+target.port+"/api/v1/handcheck", "application/json", bytes.NewBuffer(jsonValue))
  return err
}

func (s *CipherServer) sendAuthentication(targetName string) error {
  return nil
}

func (s *CipherServer) authenticate(targetName string, data []byte) error {
  return nil
}
