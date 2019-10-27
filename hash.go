package echor

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
)

// Hash function
func Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

// SaltyHash function
func SaltyHash(s string, salt string) string {
	_s := fmt.Sprintf("%s%s", s, salt)
	h := sha1.New()
	h.Write([]byte(_s))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

// UniqueID func
func UniqueID(domain string) string {
	now := NowMillis()
	r1 := rand.Float64()
	r := fmt.Sprintf("%d%f", now, r1)
	hd := Hash(domain)
	return fmt.Sprintf("%s%s-%s", hd[0:3], hd[len(hd)-4:], Hash(r))
}

// Password func
func Password(pure string) string {
	return SaltyHash(SaltyHash(pure, "S4XHS9cYKbkFpaxR"), "Bc68VZGMpSWncAaV")
}

// Sig hash
func Sig(txt string, sigKey string) string {
	return SaltyHash(txt, sigKey)
}

// SigArr func
func SigArr(parts []string, sigKey string) string {
	return Sig(strings.Join(parts, ""), sigKey)
}

// IsValidSig : check if sig is valid
func IsValidSig(parts []string, ts int64, token string, sigKey string, checksig string) bool {
	txt := strings.Join(parts, ",")
	_s := fmt.Sprintf("%s,%d,%s", txt, ts, token)
	_sig := Sig(_s, sigKey)
	return _sig == checksig
}
