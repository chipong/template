package core

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	guuid "github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// MoveFile ...
func MoveFile(src, dst string) error {
	//return os.Rename(src, dest)
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}

	return nil
}

// GenIDNumber ...
// func GenIDNumber() string {
// 	ID := strconv.FormatInt(time.Now().Unix(), 10)
// 	ID += strconv.FormatInt(int64(10000000000+rand.Intn(89999999999)), 10)
// 	return ID
// }

// GenUUID ...
func GenUUID() string {
	id := guuid.New()
	ID := strings.ReplaceAll(id.String(), "-", "")
	return ID
}

// GetLocalIP ...
func GetLocalIP() (string, error) {
	/*
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return "", err
		}
		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
		return "", nil
	*/
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

// FillStruct ...
func FillStruct(data map[string]interface{}, result interface{}) error {
	jsonbody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonbody, result); err != nil {
		return err
	}
	return nil
}

func FileCheckSum(path string) string {
	filename, _ := filepath.Abs(path)
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// byte array file CheckSum
func FileCheckSumByteArray(path []byte) string {
	b := md5.Sum(path)
	return hex.EncodeToString(b[:])
}

// String & []byte CheckSum Overload
func FileCheckSumOverload(arg interface{}) string {
	switch argType := arg.(type) {
	case string:
		_ = argType
		return FileCheckSum(arg.(string))
	case []byte:
		_ = argType
		return FileCheckSumByteArray(arg.([]byte))
	default:
		_ = argType
		return ""
	}
}
