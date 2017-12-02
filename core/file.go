package core

import (
  "fmt"
  "io"
  "os"
)

const bufferSize = 500000

//EncryptFile .
func EncryptFile(sourcePath string, targetPath string, keyPath string) error {
  keys, err := GetKey(keyPath)
  if err != nil {
    return fmt.Errorf("Error reading key: %v\n", err)
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
  if err := writeHeader(fileo, keys.index); err != nil {
    return err
  }
  data := make([]byte, bufferSize, bufferSize)
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
    keys.Cipher(keys.index, data)
    if _, err := fileo.Write(data); err != nil {
      return err
    }
  }
  keys.SaveKey(keyPath)
  return nil
}

func writeHeader(file *os.File, index *CipherIndex) error {
  for _, indexValue := range index.indexes {
    if _, err := file.WriteString(fmt.Sprintf("%05x", indexValue)); err != nil {
      return err
    }
  }
  return nil
}

//DecryptFile .
func DecryptFile(sourcePath string, targetPath string, keyPath string) error {
  keys, err := GetKey(keyPath)
  if err != nil {
    return fmt.Errorf("Error reading key: %v\n", err)
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
  errh := readHeader(filei, keys.index)
  if errh != nil {
    return errh
  }
  data := make([]byte, bufferSize, bufferSize)
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
    keys.Cipher(keys.index, data)
    if _, err := fileo.Write(data); err != nil {
      return err
    }
  }
  return nil
}

func readHeader(filei *os.File, index *CipherIndex) error {
  indexValue := 0
  data := make([]byte, 5, 5)
  for ii := range index.indexes {
    _, err := filei.Read(data)
    if err != nil {
      return err
    }
    fmt.Sscanf(string(data), "%5x", &indexValue)
    index.indexes[ii] = indexValue
  }
  return nil
}

func (keys *CipherCore) displayIndexes(label string) {
  fmt.Printf(label)
  for _, val := range keys.index.indexes {
    fmt.Printf("%x ", val)
  }
  fmt.Println("")
}
