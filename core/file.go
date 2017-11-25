package core

import (
  "fmt"
  "io"
  "os"
  "strings"
)

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
    keys.Cipher(keys.index, data)
    //fmt.Printf("%d: enc: (%d):%v\n", nn, len(datac), datac)
    if _, err := fileo.Write(data); err != nil {
      return err
    }
  }
  //fmt.Printf("end: %d\n", lastN)
  keys.SaveKey(keyPath)
  return nil
}

func writeHeader(file *os.File, index *CipherIndex) error {
  for _, indexValue := range index.indexes {
    if _, err := file.WriteString(fmt.Sprintf("%x ", indexValue)); err != nil {
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
  data := make([]byte, 10000, 10000)
  nn, err := readHeader(sourcePath, keys.index, data)
  if err != nil {
    return err
  }
  fmt.Printf("nn: %d\n", nn)
  header := make([]byte, nn, nn)
  filei.Read(header)
  fmt.Printf("header (%d): [%s]\n", nn, string(header))
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
    keys.Cipher(keys.index, data)
    //fmt.Printf("%d: enc: (%d):%v\n", nn, len(datac), datac)
    if _, err := fileo.Write(data); err != nil {
      return err
    }
  }
  //fmt.Printf("end: %d\n", lastN)
  return nil
}

func readHeader(sourcePath string, index *CipherIndex, data []byte) (int, error) {
  file, errf := os.OpenFile(sourcePath, os.O_RDWR, 0666)
  if errf != nil {
    return 0, errf
  }
  defer file.Close()
  n, err := file.Read(data)
  if err != nil {
    return 0, err
  }
  buf := string(data[:n])
  nn := 0
  nt := 0
  indexValue := 0
  for ii := range index.indexes {
    fmt.Sscanf(buf, "%x-", &indexValue)
    index.indexes[ii] = indexValue
    fmt.Printf("index %d: %x\n", ii, indexValue)
    nn = strings.Index(buf, " ") + 1
    nt += nn
    if nn < 0 {
      return 0, fmt.Errorf("Header format error")
    }
    buf = buf[nn:]
  }
  return nt, nil
}
