package tests

import (
  "crypto/rand"
  "fmt"
  "os"
  "testing"

  "github.com/freignat91/cipher3/core"
)

var (
  keyPath    = "./test.key"
  testNumber = 10000
)

func TestMain(m *testing.M) {
  //api.SetLogLevel("info")
  os.Exit(m.Run())
}

//TestKeySaveReload .
func TestKeySaveReload(t *testing.T) {
  keys1, err := core.CreateKey(3, 256, nil, false, false)
  if err != nil {
    t.Fatalf("Error creating key: %v\n", err)
    return
  }
  if errs := keys1.SaveKey(keyPath); errs != nil {
    t.Fatalf("Error saving key: %v\n", errs)
    return
  }
  keys2, err := core.GetKey(keyPath)
  if err != nil {
    t.Fatalf("Error reading key: %v\n", err)
    return
  }
  if keys1.GetKeySize() != keys2.GetKeySize() {
    t.Fatalf("Error key size %d not equal\n", keys1.GetKeySize())
    return
  }
  for ii := 0; ii < keys1.GetDimension(); ii++ {
    if keys1.GetKeyIndex(ii) != keys2.GetKeyIndex(ii) {
      t.Fatalf("Error index %d not equal\n", ii)
      return
    }
    same := true
    buf1 := keys1.GetKeyBytes(ii)
    buf2 := keys2.GetKeyBytes(ii)
    for jj := 0; jj < keys1.GetKeySize(); jj++ {
      if buf1[jj] != buf2[jj] {
        same = false
        break
      }
    }
    if !same {
      t.Fatalf("Error key %d not equal\n", ii)
      return
    }
  }
  fmt.Println("keys ok")
}

//TestCipher .
func TestCipher(t *testing.T) {
  keySend, err := core.GetKey(keyPath)
  if err != nil {
    t.Fatalf("Error reading key: %v\n", err)
    return
  }
  keyReceive, err := core.GetKey(keyPath)
  if err != nil {
    t.Fatalf("Error reading key: %v\n", err)
    return
  }
  list := []byte{15, 211, 218, 155, 207, 209, 212, 102, 241, 192, 130, 92, 10, 92, 213, 236, 172, 190, 189, 213, 116, 66, 8, 33, 132, 16, 66, 8, 33, 132, 16}
  //list := []byte{200, 166, 141, 215, 211, 66, 55, 245, 7, 183, 83, 15, 192, 57, 118, 110, 186, 145, 209, 127, 28, 54, 180, 68, 13, 122, 155, 105, 108, 110, 239}
  if err := testEncriptDecript(keySend, keyReceive, list, true); err != nil {
    t.Fatalf("Error: %v\n", err)
    return
  }
  keySend.SaveKey(keyPath)
}

//Testn .
func TestCypherN(t *testing.T) {
  keySend, err := core.GetKey(keyPath)
  if err != nil {
    t.Fatalf("Error reading key: %v\n", err)
    return
  }
  keyReceive, err := core.GetKey(keyPath)
  if err != nil {
    t.Fatalf("Error reading key: %v\n", err)
    return
  }
  fmt.Printf("start test")
  nn := 0
  list := make([]byte, keySend.GetKeySize(), keySend.GetKeySize())
  for {
    rand.Read(list)
    if err := testEncriptDecript(keySend, keyReceive, list, false); err != nil {
      t.Fatalf("Error: %v\n", err)
      return
    }
    nn++
    if nn%1000 == 0 {
      fmt.Println(nn)
    }
    if testNumber == nn {
      break
    }
  }
  keySend.SaveKey(keyPath)
}

func testEncriptDecript(keySend *core.CipherCore, keyReceive *core.CipherCore, list []byte, verbose bool) error {
  if verbose {
    fmt.Printf("list: %v\n", list)
  }
  listm := make([]byte, len(list), len(list))
  for ii, val := range list {
    listm[ii] = val
  }
  if err := keySend.Cipher(nil, list); err != nil {
    return err
  }
  if verbose {
    fmt.Printf("encr: %v\n", list)
  }
  keyReceive.Cipher(nil, list)
  if verbose {
    fmt.Printf("decr: %v\n", list)
  }
  same := true
  for jj := 0; jj < len(list); jj++ {
    if list[jj] != listm[jj] {
      same = false
      break
    }
  }
  if !same {
    fmt.Println("Error")
    fmt.Printf("decr: %v\n", list)
    fmt.Printf("init: %v\n", listm)
    return fmt.Errorf("Error plain versus decrypt on data\n")
  }
  return nil
}

func TestEncryptFile(t *testing.T) {
  os.Remove("./filee")
  if err := core.EncryptFile("./fileTest", "./filee", keyPath); err != nil {
    t.Fatalf("Error: %v\n", err)
  }
}

func TestDecryptFile(t *testing.T) {
  os.Remove("./filed")
  if err := core.DecryptFile("./filee", "./filed", keyPath); err != nil {
    t.Fatalf("Error: %v\n", err)
  }
}
