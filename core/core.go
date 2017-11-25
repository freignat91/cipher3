package core

import (
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "strings"
  "time"
)

//CipherCore .
type CipherCore struct {
  size  int
  space []*CipherKey
}

//CipherKey .
type CipherKey struct {
  index int
  value byte
  key   []byte
}

//CreateKey .
func CreateKey(dimension int, pSize int, verbose bool, debug bool) (*CipherCore, error) {
  if pSize%64 != 0 {
    return nil, fmt.Errorf("number of bits should be a multiple of 64")
  }
  if dimension <= 0 {
    return nil, fmt.Errorf("invalide number of dimension, shoud be >0")
  }
  if verbose {
    fmt.Printf("Compute key dimension: %d size: %d bits\n", dimension, pSize)
  }

  //find two random prime size sqrt(keyBitsize)
  space := make([]*CipherKey, dimension, dimension)
  done := 0
  for ii := range space {
    key := getRandomKey(pSize, verbose)
    space[ii] = &CipherKey{index: 0, key: *key}
    done++
    if verbose {
      fmt.Printf("prime: %x\n", key)
    }
  }
  for done < dimension {
    time.Sleep(1 * time.Second)
  }
  return &CipherCore{space: space, size: pSize}, nil
}

//SaveKey .
func (keys *CipherCore) SaveKey(path string) error {
  file, errc := os.Create(path)
  if errc != nil {
    return errc
  }
  for _, key := range keys.space {
    if _, err := file.WriteString(key.toString()); err != nil {
      return err
    }
  }
  file.Close()
  return nil
}

func (k *CipherKey) toString() string {
  return fmt.Sprintf("%x-%x-", k.index, k.key)
}

//GetKey .
func GetKey(path string) (*CipherCore, error) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  list := strings.Split(string(data), "-")
  space := make([]*CipherKey, len(list)/2, len(list)/2)
  nn := 0
  for ii := 0; ii < len(list); ii = ii + 2 {
    if len(list[ii]) == 0 {
      break
    }
    index := 0
    fmt.Sscanf(list[ii], "%x", &index)
    key := make([]byte, len(list[ii+1])/2, len(list[ii+1])/2)
    fmt.Sscanf(list[ii+1], "%x", &key)
    space[nn] = &CipherKey{index: index, key: key}
    nn++
  }
  return &CipherCore{space: space, size: len(space[0].key)}, nil
}

//CipherFile .
func CipherFile(sourcePath string, targetPath string, keyPath string) error {
  core, err := GetKey(keyPath)
  if err != nil {
    return err
  }
  filei, errf := os.OpenFile(sourcePath, os.O_RDWR, 0666)
  if errf != nil {
    return errf
  }
  defer filei.Close()
  fileo, errf := os.Create(targetPath)
  if errf != nil {
    return errf
  }
  defer fileo.Close()
  data := make([]byte, 10000, 10000)
  for {
    data = data[:cap(data)]
    n, err := filei.Read(data)
    if err != nil {
      if err == io.EOF {
        break
      }
      return err
    }
    data = data[:n]
    //fmt.Printf("%d: read: (%d):%v\n\n", nn, len(data), data)
    core.Cipher(data)
    //fmt.Printf("%d: enc: (%d):%v\n", nn, len(datac), datac)
    if _, err := fileo.Write(data); err != nil {
      return err
    }
  }
  //fmt.Printf("end: %d\n", lastN)
  return nil
}

//Cipher .
func (keys *CipherCore) Cipher(data []byte) {
  for ii := range data {
    data[ii] = data[ii] ^ keys.getNextComputedKeyByte()
  }
}

func (keys *CipherCore) getNextComputedKeyByte() byte {
  nkey := 0
  for {
    keys.space[nkey].index++
    if keys.space[nkey].index < keys.size {
      break
    }
    keys.space[nkey].index = 0
    nkey++
  }
  var ret byte
  for _, key := range keys.space {
    ret = ret ^ key.key[key.index]
  }
  return ret
}

//GetKeySize .
func (keys *CipherCore) GetKeySize() int {
  return keys.size
}

//GetDimension .
func (keys *CipherCore) GetDimension() int {
  return len(keys.space)
}

//GetKeyIndex .
func (keys *CipherCore) GetKeyIndex(num int) int {
  return keys.space[num].index
}

//GetKeyBytes .
func (keys *CipherCore) GetKeyBytes(num int) []byte {
  return keys.space[num].key
}

//ResetKeyIndexes .
func (keys *CipherCore) ResetKeyIndexes() {
  for _, key := range keys.space {
    key.index = 0
  }
}
