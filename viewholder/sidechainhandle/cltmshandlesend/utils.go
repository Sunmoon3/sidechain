package cltmshandlesend

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

func Message_split(message []byte) [][]byte {
	ans := make([][]byte, 0)
	s, t := 0, 32
	mlen := len(message)
	for {
		if mlen > 32 {
			ans = append(ans, message[s:t])
			s += 32
			t += 32
			mlen -= 32
		} else {
			ans = append(ans, message[s:s+mlen])
			break
		}
	}
	return ans
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func Path() []string {
	prefix := "/root/work/chainmaker/viewholder_mult/crypto-config/"
	suffix := "/user/admin1/admin1.sign.key"
	paths := make([]string, 0)
	for i := 0; i < 16; i++ {
		path := prefix + "wx-org" + strconv.Itoa(i+1) + ".chainmaker.org" + suffix
		paths = append(paths, path)
	}
	return paths
}
