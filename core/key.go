package core

import "crypto/rand"

func getRandomKey(randomStr string, size int, verbose bool) *[]byte {
  key := make([]byte, size, size)
  rand.Read(key)
  if randomStr != "" {
    data := []byte(randomStr)
    ll := 0
    for cc, val := range key {
      key[cc] = val ^ data[ll]
      ll++
      if ll >= len(data) {
        ll = 0
      }
    }
  }
  reDistribute(&key)
  return &key
}

func reDistribute(data *[]byte) {
  count := make([]byte, 256, 256)
  for _, val := range *data {
    count[val]++
  }
  for {
    min := 256
    minIndex := 0
    max := -1
    maxIndex := 0
    for ii, valc := range count {
      val := int(valc)
      if val <= min {
        min = val
        minIndex = ii
      }
      if val >= max {
        max = val
        maxIndex = ii
      }
      if min+2 > max {
        return
      }
      for ii, val := range *data {
        if int(val) == maxIndex {
          (*data)[ii] = byte(minIndex)
        }
      }
    }
  }
}
