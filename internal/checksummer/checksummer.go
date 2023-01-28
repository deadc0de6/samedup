/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package checksummer

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"

	"github.com/cespare/xxhash"
)

const (
	// UseSHA1 sha1 checksum
	UseSHA1 = iota
	// UseXXHash xxhash checksum
	UseXXHash
	// UseCRC32 CRC32 checksum
	UseCRC32
	// UseMD5 MD5 checksum
	UseMD5
)

// GetHashMethods returns all available checksum methods and index
func GetHashMethods() string {
	var out []string
	out = append(out, fmt.Sprintf("%d:sha1", UseSHA1))
	out = append(out, fmt.Sprintf("%d:xxhash", UseXXHash))
	out = append(out, fmt.Sprintf("%d:crc32", UseCRC32))
	out = append(out, fmt.Sprintf("%d:md5", UseMD5))
	return strings.Join(out, ",")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getContent(path string, n int64, h io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("checksum opening %s: %v", path, err)
	}
	defer f.Close()

	if n < 0 {
		_, err = io.Copy(h, f)
		if err != nil {
			return fmt.Errorf("checksum read content %s: %v", path, err)
		}
	} else {
		_, err = io.CopyN(h, f, n)
	}
	if err != nil && err != io.EOF {
		return fmt.Errorf("checksum read content %s: %v", path, err)
	}
	return nil
}

func checksumSHA1(path string, n int64) (string, error) {
	h := sha1.New()
	err := getContent(path, n, h)
	if err != nil {
		return "", err
	}

	hash := h.Sum(nil)
	sum := hex.EncodeToString(hash[:])
	return "sha1:" + sum, nil
}

func checksumCRC32(path string, n int64) (string, error) {
	h := crc32.NewIEEE()
	err := getContent(path, n, h)
	if err != nil {
		return "", err
	}

	hash := h.Sum(nil)
	sum := hex.EncodeToString(hash[:])
	return "crc32:" + sum, nil
}

func checksumMD5(path string, n int64) (string, error) {
	h := md5.New()
	err := getContent(path, n, h)
	if err != nil {
		return "", err
	}

	hash := h.Sum(nil)
	sum := hex.EncodeToString(hash[:])
	return "md5:" + sum, nil
}

func checksumXXHash(path string, n int64) (string, error) {
	h := xxhash.New()
	err := getContent(path, n, h)
	if err != nil {
		return "", err
	}

	hash := h.Sum64()
	return fmt.Sprintf("xxhash:%d", hash), nil
}

func checksum(path string, method int, n int64) (string, error) {
	var err error
	var chk string

	if !fileExists(path) {
		return "", fmt.Errorf("%s does not exist", path)
	}
	switch method {
	case UseMD5:
		chk, err = checksumMD5(path, n)
	case UseCRC32:
		chk, err = checksumCRC32(path, n)
	case UseSHA1:
		chk, err = checksumSHA1(path, n)
	case UseXXHash:
		chk, err = checksumXXHash(path, n)
	default:
		return "", fmt.Errorf("no such hash method")
	}
	return chk, err
}
