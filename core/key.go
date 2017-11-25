package core

import "crypto/rand"

func getRandomKey(size int, verbose bool) *[]byte {
  key := make([]byte, size, size)
  rand.Read(key)
  return &key
}
