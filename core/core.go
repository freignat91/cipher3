package core

import (
  "fmt"
  "io/ioutil"
  "os"
  "strings"
  "time"
)

//CipherCore .
type CipherCore struct {
  size  int
  space []*CipherKey
  index *CipherIndex
}

//CipherKey .
type CipherKey struct {
  key []byte
}

//CipherIndex .
type CipherIndex struct {
  values  []byte
  indexes []int
}

//CreateKey .
func CreateKey(dimension int, pSize int, randomList []string, verbose bool, debug bool) (*CipherCore, error) {
  if dimension <= 0 {
    return nil, fmt.Errorf("invalide number of dimension, shoud be >0")
  }
  if verbose {
    fmt.Printf("Compute key dimension: %d size: %d bits\n", dimension, pSize)
  }
  if randomList == nil {
    randomList = make([]string, dimension, dimension)
  }

  //find two random prime size sqrt(keyBitsize)
  space := make([]*CipherKey, dimension, dimension)
  done := 0
  for ii := range space {
    key := getRandomKey(randomList[ii], pSize, verbose)
    space[ii] = &CipherKey{key: *key}
    done++
    if verbose {
      fmt.Printf("prime: %x\n", key)
    }
  }
  for done < dimension {
    time.Sleep(1 * time.Second)
  }
  keys := CipherCore{space: space, size: pSize}
  keys.index = keys.getNewCipherIndex()
  return &keys, nil
}

func (keys *CipherCore) getNewCipherIndex() *CipherIndex {
  return &CipherIndex{
    indexes: make([]int, len(keys.space), len(keys.space)),
    values:  make([]byte, len(keys.space), len(keys.space)),
  }
}

//SaveKey .
func (keys *CipherCore) SaveKey(path string) error {
  file, errc := os.Create(path)
  if errc != nil {
    return errc
  }
  skey := keys.ToString()
  if _, err := file.WriteString(skey); err != nil {
    return err
  }
  file.Close()
  return nil
}

//Copy .
func (keys *CipherCore) Copy() *CipherCore {
  skey := keys.ToString()
  key, _ := ReadKey(skey)
  return key
}

//ToString .
func (keys *CipherCore) ToString() string {
  ret := ""
  for ii := range keys.space {
    if ii > 0 {
      ret += "\n"
    }
    ret += keys.keyToString(ii)
  }
  return ret
}

func (keys *CipherCore) keyToString(num int) string {
  return fmt.Sprintf("%05x%x", keys.index.indexes[num], keys.space[num].key)
}

//GetKey .
func GetKey(path string) (*CipherCore, error) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  return ReadKey(string(data))
}

//ReadKey .
func ReadKey(skey string) (*CipherCore, error) {
  list := strings.Split(skey, "\n")
  dimension := len(list)
  space := make([]*CipherKey, dimension, dimension)
  indexes := make([]int, dimension, dimension)
  for ii, buf := range list {
    index := 0
    fmt.Sscanf(buf[0:5], "%x", &index)
    indexes[ii] = index
    key := make([]byte, (len(buf)-5)/2, (len(buf)-5)/2)
    fmt.Sscanf(buf[5:], "%x", &key)
    space[ii] = &CipherKey{key: key}
  }
  kindex := &CipherIndex{indexes: indexes, values: make([]byte, dimension, dimension)}
  keys := CipherCore{space: space, index: kindex, size: len(space[0].key)}
  keys.initIndexesValues()
  return &keys, nil
}

func (keys *CipherCore) initIndexesValues() {
  for ii, index := range keys.index.indexes {
    keys.index.values[ii] = keys.space[ii].key[index]
  }
}

//Cipher .
func (keys *CipherCore) Cipher(index *CipherIndex, data []byte) error {
  if index == nil {
    index = keys.index
  }
  for ii := range data {
    keyc, err := keys.getNextComputedKeyByte(index)
    if err != nil {
      return err
    }
    data[ii] = data[ii] ^ keyc
  }
  return nil
}

func (keys *CipherCore) getNextComputedKeyByte(index *CipherIndex) (byte, error) {
  nkey := 0
  for {
    index.indexes[nkey]++
    if index.indexes[nkey] < keys.size {
      break
    }
    index.indexes[nkey] = 0
    nkey++
    if nkey == len(index.indexes) {
      return 0, fmt.Errorf("The key has been completly used")
    }
  }
  var ret byte
  for ii, key := range keys.space {
    ret = ret ^ key.key[index.indexes[ii]]
  }
  return ret, nil
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
  return keys.index.indexes[num]
}

//GetKeyBytes .
func (keys *CipherCore) GetKeyBytes(num int) []byte {
  return keys.space[num].key
}

//DisplayIndex .
func (keys *CipherCore) DisplayIndex() {
  fmt.Printf("Index: ")
  for _, val := range keys.index.indexes {
    fmt.Printf("%x-", val)
  }
  fmt.Printf("\n")
}
