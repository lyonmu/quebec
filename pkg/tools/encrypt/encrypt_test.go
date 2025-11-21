package encrypt

import (
	"log"
	"testing"
)

func TestHashWithBcryptBytes(t *testing.T) {

	u := []byte("$2a$10$1WSPlPr.1D4Wbw0JbRiCt.Rp7B5bwxt2Zpt1BXacu8b7nIZ5bh1a.")
	p := []byte("$2a$10$SrT6rMl2mHM8xx4rGEzyPO4h6nynw7s1nukqI44dPIaxKE82mR2li")

	log.Println(CompareWithBcryptBytes(HashWithSHA256Bytes([]byte("system")), u))
	log.Println(CompareWithBcryptBytes(HashWithSHA256Bytes([]byte("Quebec@123456")), p))
}
