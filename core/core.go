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
  if pSize%64 != 0 {
    return nil, fmt.Errorf("number of bits should be a multiple of 64")
  }
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
  for ii := range keys.space {
    if _, err := file.WriteString(keys.keyToString(ii)); err != nil {
      return err
    }
  }
  file.Close()
  return nil
}

func (keys *CipherCore) keyToString(num int) string {
  return fmt.Sprintf("%x-%x-", keys.index.indexes[num], keys.space[num].key)
}

//GetKey .
func GetKey(path string) (*CipherCore, error) {
  data, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  list := strings.Split(string(data), "-")
  space := make([]*CipherKey, len(list)/2, len(list)/2)
  indexes := make([]int, len(list)/2, len(list)/2)
  nn := 0
  for ii := 0; ii < len(list); ii = ii + 2 {
    if len(list[ii]) == 0 {
      break
    }
    index := 0
    fmt.Sscanf(list[ii], "%x", &index)
    indexes[nn] = index
    key := make([]byte, len(list[ii+1])/2, len(list[ii+1])/2)
    fmt.Sscanf(list[ii+1], "%x", &key)
    space[nn] = &CipherKey{key: key}
    nn++
  }
  kindex := &CipherIndex{indexes: indexes, values: make([]byte, len(list)/2, len(list)/2)}
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
