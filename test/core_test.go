package tests

import (
  "crypto/rand"
  "fmt"
  "os"
  "testing"

  "github.com/freignat91/cipher2/core"
)

func TestMain(m *testing.M) {
  //api.SetLogLevel("info")
  os.Exit(m.Run())
}

//TestKeySaveReload .
func TestKeySaveReload(t *testing.T) {
  keys, err := core.CreateKey(3, 256, false, false)
  if err != nil {
    t.Fatalf("Error creating key: %v\n", err)
    return
  }
  if errs := keys.SaveKey("test"); errs != nil {
    t.Fatalf("Error saving key: %v\n", errs)
    return
  }
  keys2, err := core.GetKey("test")
  if err != nil {
    t.Fatalf("Error reading key: %v\n", err)
    return
  }
  for ii, key := range keys.space {
    if key.index != keys.space[ii].index {
      t.Fatalf("Error index %d not equal\n", ii)
      return
    }
    same := true
    for jj, val := range key.key {
      if val != key2.space[ii].key[jj] {
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

//Test1 .
func test() error {
  fmt.Println("Load keys")
  core, err := GetKey("k256")
  if err != nil {
    return fmt.Errorf("Error reading key: %v\n", err)
  }
  list := []byte{15, 211, 218, 155, 207, 209, 212, 102, 241, 192, 130, 92, 10, 92, 213, 236, 172, 190, 189, 213, 116, 66, 8, 33, 132, 16, 66, 8, 33, 132, 16}
  //list := []byte{200, 166, 141, 215, 211, 66, 55, 245, 7, 183, 83, 15, 192, 57, 118, 110, 186, 145, 209, 127, 28, 54, 180, 68, 13, 122, 155, 105, 108, 110, 239}
  if err := testEncriptDecript(core, list, true); err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  fmt.Println("done")
  return nil
}

//Testn .
func testn(num int, size int) error {
  fmt.Println("Load keys")
  core, err := GetKey("k256")
  if err != nil {
    return fmt.Errorf("Error reading key: %v\n", err)
  }
  fmt.Printf("start test")
  nn := 0
  list := make([]byte, size, size)
  for {
    rand.Read(list)
    if err := testEncriptDecript(core, list, false); err != nil {
      fmt.Printf("Error: %v\n", err)
    }
    nn++
    if nn%1000 == 0 {
      fmt.Println(nn)
    }
    if num == nn {
      fmt.Println("done")
      return nil
    }
  }
}

func testEncriptDecript(core *CipherCore, list []byte, verbose bool) error {
  if verbose {
    fmt.Printf("list: %v\n", list)
  }
  listm := make([]byte, len(list), len(list))
  for ii, val := range list {
    listm[ii] = val
  }
  core.cipher(list)
  if verbose {
    fmt.Printf("encr: %v\n", list)
  }
  core.cipher(list)
  if verbose {
    fmt.Printf("decr: %v\n", list)
  }
  for i, val := range list {
    if val != listm[i] {
      fmt.Println("Error")
      fmt.Printf("decr: %v\n", list)
      fmt.Printf("init: %v\n", listm)
      return fmt.Errorf("Error plain versus decrypt on data %d\n", i)
    }
  }
  return nil
}
