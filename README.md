# CIPHER3

cipher3 v0.0.1 on going

# Purpose

Demonstrate a different encryption way respecting the following propositions:

- encryption shouldn't be decrypted using brute force. It should have no way to reverse-engineering it
- do not have complicated mathematical usage, encryption should be simple to understand
- encryption algorithm could be public without make it weaken
- safe authentication of the sender should be possible
- data flow encryption should be fast

One way to got an absolutely not decryptable algorithm is to use a random key as long as the data to encrypt and never reuse two times, so each data to encrypt has its own unique random key

For instance, if the following data:
"Hi, how are you?" with the random key,
encrypted using a simple xor with the following random key:
"48 E2 8A 23 B2 C8 12 AA 40 32 56 56 22 8D 03 00",
there is no way to decrypt the encrypted data if you do not have the key
there is also no brute force algorithm able to decrypt it no matter the time used, because xor is just not reversable.

xx ^ yy -> zz
knowing zz, there is no way to found xx and yy without having more information about xx or yy.

All the difficulties is to be able to produce a nearly infinite random key and ensure that the key is enough random, to give no informatoin on the key bytes.

Unfortunately, there is no way to produce a fully random key like that using machine code, no matter the algorithm. A random key produce by a machine is never completely random.

Hopefully we don't need a fully random key.
When data is encrypted and decrypted using the key, no matter the encryption/decryption algorithm we can think that the bytes of the key used to encrypt the data are compromised and well known by hackers.
So all the force of the encryption algorithm is that the next byte values of the key can be anticipate using the previous already used key values.
That is easier to achieve than having a really random key.

...







# Install

- prerequisite: have go installed and GOPATH set
- install glide: go get glide
- clone this project in ĜOPATH/src/github.com/freignat91/cipher3
- execute: glide update
- execute: make install
- then the command cipher3 is available

For Ubuntu, you have a pre-build cipher3.ubuntu file you can use without cloning and building the projet.

# Usage

## global options

- --help help on command
- -v verbose, display information during command execution
- -debug: display more information during command execution

## cipher createKeys [keyPath] -size [keysize]

This commande generates [keyPath].pub and [keyPath].key keys (public and private) having [keysize] bits long

## cipher encryptFile [sourceFilePath] [targetFilePath] [publicKeyPath]

This command encrypt the file [sourceFilePath] and save the  result in [targetFilePath] using the public key [publicKeyPath]

## cipher decryptFile [sourceFilePath] [targetFilePath] [privateKeyPath]

This command decrypt the file [sourceFilePath] and save the result in [targetFilePath] using the private key [privateKeyPath]
