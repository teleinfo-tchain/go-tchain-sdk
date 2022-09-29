package privacyutils

/*
#cgo CFLAGS: -D ENABLE_MODULE_RECOVERY
#cgo CFLAGS: -D ENABLE_MODULE_ECDH
#cgo CFLAGS: -D ENABLE_MODULE_RANGEPROOF
#cgo CFLAGS: -D ENABLE_MODULE_BULLETPROOF
#cgo CFLAGS: -D ENABLE_MODULE_GENERATOR
#cgo CFLAGS: -D USE_FIELD_INV_BUILTIN
#cgo CFLAGS: -D USE_NUM_NONE
#cgo CFLAGS: -D USE_SCALAR_INV_BUILTIN
#cgo CFLAGS: -D USE_FIELD_10X26
#cgo CFLAGS: -D USE_SCALAR_8X32
#cgo CFLAGS: -D HAVE_BUILTIN_EXPECT
#cgo CFLAGS: -I "./"
#cgo CFLAGS: -I "./src"
#include "privacy/purec_privacy.c"
*/
import "C"

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"unsafe"
)

var PrivErrorInfo = map[C.int64_t]error{
	C.PUREC_ERRCODE_SUCCESS:            nil,
	C.PUREC_ERRCODE_RANDOM_ERROR:       errors.New("Generates a random number error!"),
	C.PUREC_ERRCODE_INVALID_PARAMETER:  errors.New("Invalid parameters"),
	C.PUREC_ERRCODE_CREATE_PEDERSEN:    errors.New("Failed to create pedersen commitment"),
	C.PUREC_ERRCODE_PARSE_PEDERSEN:     errors.New("Failed to parse pedersen commitment"),
	C.PUREC_ERRCODE_SERIALIZE_PEDERSEN: errors.New("Failed to serialize pedersen commitment"),
	C.PUREC_ERRCODE_SERIALIZE_PUBKEY:   errors.New("Failed to serialize pubkey"),
	C.PUREC_ERRCODE_VERIFY_TALLY:       errors.New("Failed to verify tally"),
	C.PUREC_ERRCODE_CREATE_PUBKEY:      errors.New("Failed to create pubkey"),
	C.PUREC_ERRCODE_RANGEPROOF_PROVE:   errors.New("Failed to generate rangeproof"),
	C.PUREC_ERRCODE_RANGEPROOF_VERIFY:  errors.New("Failed to verify rangeproof"),
	C.PUREC_ERRCODE_OUT_RANGE:          errors.New("Out of range"),
	C.PUREC_ERRCODE_ECDSA_CRATE:        errors.New("Failed to create ecdsa signature"),
	C.PUREC_ERRCODE_ECDSA_SERIALIZE:    errors.New("Failed to serialize ecdsa signature"),
	C.PUREC_ERRCODE_ECDSA_VERIFY:       errors.New("Failed to verify ecdsa signature"),
	C.PUREC_ERRCODE_ECDSA_PARSE:        errors.New("Failed to parse ecdsa signature"),
	C.PUREC_ERRCODE_PARSE_PUBKEY:       errors.New("Failed to parse pubkey"),
	C.PUREC_ERRCODE_BLIND_SUM:          errors.New("Failed to blind sum"),
	C.PUREC_ERRCODE_UNKNOWN:            errors.New("An unknown error"),
	C.PUREC_ERRCODE_BP_LIB_INTERNAL:    errors.New("BP library internal error"),
}

func InitPrivacy() uint64 {
	var ptr C.uint64_t = C.InitPrivacy()
	return uint64(ptr)
}

func DestroyPrivacy(ptr uint64) {
	C.DestroyPrivacy(C.uint64_t(ptr))
}

func CreatePedersenCommit(ptr uint64, value uint64, blind string) (string, error) {
	cblind := C.pure_string{}
	cblind.c_str = C.CString(blind)
	cblind.length = C.int(len([]rune(blind)))
	defer C.free(unsafe.Pointer(cblind.c_str))

	ccommit := C.pure_string{}
	ret := C.PurecCreatePedersenCommit(C.uint64_t(ptr), C.uint64_t(value), &cblind, &ccommit)
	commit := C.GoStringN(ccommit.c_str, ccommit.length)
	defer C.free(unsafe.Pointer(ccommit.c_str))
	if ret == 0 {
		return commit, nil
	} else {
		return "", PrivErrorInfo[ret]
	}
}

func PedersenTallyVerify(ptr uint64, inputs []string, outputs []string, msg string, sig string) error {
	cinputs := make([](C.pure_string), 0)
	for i := range inputs {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(inputs[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(inputs[i])))
		cinputs = append(cinputs, strptr)
	}
	coutputs := make([](C.pure_string), 0)
	for i := range outputs {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(outputs[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(outputs[i])))
		coutputs = append(coutputs, strptr)
	}

	ccmsg := C.pure_string{}
	ccmsg.c_str = C.CString(msg)
	ccmsg.length = C.int(len([]rune(msg)))
	defer C.free(unsafe.Pointer(ccmsg.c_str))

	csig := C.pure_string{}
	csig.c_str = C.CString(sig)
	csig.length = C.int(len([]rune(sig)))
	defer C.free(unsafe.Pointer(csig.c_str))

	ret := C.PurecPedersenTallyVerify(C.uint64_t(ptr), (**C.pure_string)(unsafe.Pointer(&cinputs[0])), C.int(len(cinputs)),
		(**C.pure_string)(unsafe.Pointer(&coutputs[0])), C.int(len(coutputs)),
		&ccmsg, &csig)
	//    cblind := C.malloc(C.sizeof(C.pure_string{}));
	return PrivErrorInfo[ret]
}

func TallyVerify(ptr uint64, inputs []string, outputs []string) error {
	cinputs := make([](C.pure_string), 0)
	for i := range inputs {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(inputs[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(inputs[i])))
		cinputs = append(cinputs, strptr)
	}
	coutputs := make([](C.pure_string), 0)
	for i := range outputs {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(outputs[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(outputs[i])))
		coutputs = append(coutputs, strptr)
	}

	ret := C.PurecTallyVerify(C.uint64_t(ptr),
		(**C.pure_string)(unsafe.Pointer(&cinputs[0])),
		C.int(len(cinputs)),
		(**C.pure_string)(unsafe.Pointer(&coutputs[0])),
		C.int(len(coutputs)))
	//    cblind := C.malloc(C.sizeof(C.pure_string{}));
	return PrivErrorInfo[ret]
}

func CreateEcdhKey(ptr uint64, priv_key string, pub_key string) (string, error) {
	c_priv_key := C.pure_string{}
	c_priv_key.c_str = C.CString(priv_key)
	c_priv_key.length = C.int(len([]rune(priv_key)))
	defer C.free(unsafe.Pointer(c_priv_key.c_str))

	c_pub_key := C.pure_string{}
	c_pub_key.c_str = C.CString(pub_key)
	c_pub_key.length = C.int(len([]rune(pub_key)))
	defer C.free(unsafe.Pointer(c_pub_key.c_str))

	c_ecdh_key := C.pure_string{}
	ret := C.PurecCreateEcdhKey(C.uint64_t(ptr), &c_priv_key, &c_pub_key, &c_ecdh_key)
	defer C.free(unsafe.Pointer(c_ecdh_key.c_str))
	ecdh_key := C.GoStringN(c_ecdh_key.c_str, c_ecdh_key.length)

	return ecdh_key, PrivErrorInfo[ret]
}

func EcdsaSign(ptr uint64, priv_key string, data string) (string, error) {
	c_priv_key := C.pure_string{}
	c_priv_key.c_str = C.CString(priv_key)
	c_priv_key.length = C.int(len([]rune(priv_key)))
	defer C.free(unsafe.Pointer(c_priv_key.c_str))

	c_data := C.pure_string{}
	c_data.c_str = C.CString(data)
	c_data.length = C.int(len([]rune(data)))
	defer C.free(unsafe.Pointer(c_data.c_str))

	c_sig := C.pure_string{}
	ret := C.PurecEcdsaSign(C.uint64_t(ptr), &c_priv_key, &c_data, &c_sig)
	defer C.free(unsafe.Pointer(c_sig.c_str))
	sig := C.GoStringN(c_sig.c_str, c_sig.length)

	return sig, PrivErrorInfo[ret]
}

func EcdsaVerify(ptr uint64, pub_key string, data string, sig string) error {
	c_pub_key := C.pure_string{}
	c_pub_key.c_str = C.CString(pub_key)
	c_pub_key.length = C.int(len([]rune(pub_key)))
	defer C.free(unsafe.Pointer(c_pub_key.c_str))

	c_data := C.pure_string{}
	c_data.c_str = C.CString(data)
	c_data.length = C.int(len([]rune(data)))
	defer C.free(unsafe.Pointer(c_data.c_str))

	c_sig := C.pure_string{}
	c_sig.c_str = C.CString(sig)
	c_sig.length = C.int(len([]rune(sig)))
	defer C.free(unsafe.Pointer(c_sig.c_str))

	return PrivErrorInfo[C.PurecEcdsaVerify(C.uint64_t(ptr), &c_pub_key, &c_data, &c_sig)]
}

func CreateKeyPair(ptr uint64) (string, string, error) {
	c_priv_key := C.pure_string{}
	c_pub_key := C.pure_string{}
	ret := C.PurecCreateKeyPair(C.uint64_t(ptr), &c_pub_key, &c_priv_key)
	defer C.free(unsafe.Pointer(c_priv_key.c_str))
	defer C.free(unsafe.Pointer(c_pub_key.c_str))
	priv_key := C.GoStringN(c_priv_key.c_str, c_priv_key.length)
	pub_key := C.GoStringN(c_pub_key.c_str, c_pub_key.length)

	return pub_key, priv_key, PrivErrorInfo[ret]
}

func GetPublicKey(ptr uint64, priv_key string) (string, error) {
	c_priv_key := C.pure_string{}
	c_priv_key.c_str = C.CString(priv_key)
	c_priv_key.length = C.int(len([]rune(priv_key)))
	defer C.free(unsafe.Pointer(c_priv_key.c_str))

	c_pub_key := C.pure_string{}

	ret := C.PurecGetPublicKey(C.uint64_t(ptr), &c_priv_key, &c_pub_key)
	defer C.free(unsafe.Pointer(c_pub_key.c_str))
	pub_key := C.GoStringN(c_pub_key.c_str, c_pub_key.length)

	return pub_key, PrivErrorInfo[ret]
}

func CombinePublicKey(ptr uint64, pub_keys []string) (string, error) {
	cinputs := make([](C.pure_string), 0)
	for i := range pub_keys {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(pub_keys[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(pub_keys[i])))
		cinputs = append(cinputs, strptr)
	}

	c_pub_key := C.pure_string{}

	ret := C.PurecPublicKeyCombine(C.uint64_t(ptr), (**C.pure_string)(unsafe.Pointer(&cinputs[0])), C.int(len(cinputs)), &c_pub_key)
	defer C.free(unsafe.Pointer(c_pub_key.c_str))
	pub_key := C.GoStringN(c_pub_key.c_str, c_pub_key.length)

	return pub_key, PrivErrorInfo[ret]
}

func ExcessSign(ptr uint64, inputs []string, outputs []string, msg string) (sig string, ret error) {
	cinputs := make([](C.pure_string), 0)
	for i := range inputs {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(inputs[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(inputs[i])))
		cinputs = append(cinputs, strptr)
	}
	coutputs := make([](C.pure_string), 0)
	for i := range outputs {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(outputs[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(outputs[i])))
		coutputs = append(coutputs, strptr)
	}

	ccmsg := C.pure_string{}
	ccmsg.c_str = C.CString(msg)
	ccmsg.length = C.int(len([]rune(msg)))
	defer C.free(unsafe.Pointer(ccmsg.c_str))

	csig := C.pure_string{}
	defer C.free(unsafe.Pointer(csig.c_str))

	ret = PrivErrorInfo[C.PurecExcessSign(C.uint64_t(ptr), (**C.pure_string)(unsafe.Pointer(&cinputs[0])), C.int(len(cinputs)),
		(**C.pure_string)(unsafe.Pointer(&coutputs[0])), C.int(len(coutputs)),
		&ccmsg, &csig)]

	sig = C.GoStringN(csig.c_str, csig.length)
	//    cblind := C.malloc(C.sizeof(C.pure_string{}));
	return
}

func ExcessBlind(ptr uint64, inputs []string, outputs []string) (blind string, retError error) {
	cinputs := make([](C.pure_string), 0)
	for i := range inputs {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(inputs[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(inputs[i])))
		cinputs = append(cinputs, strptr)
	}
	coutputs := make([](C.pure_string), 0)
	for i := range outputs {
		strptr := C.pure_string{}
		strptr.c_str = C.CString(outputs[i])
		defer C.free(unsafe.Pointer(strptr.c_str))
		strptr.length = C.int(len([]rune(outputs[i])))
		coutputs = append(coutputs, strptr)
	}

	cblind := C.pure_string{}
	defer C.free(unsafe.Pointer(cblind.c_str))

	ret := C.PurecExcessBlind(C.uint64_t(ptr), (**C.pure_string)(unsafe.Pointer(&cinputs[0])), C.int(len(cinputs)),
		(**C.pure_string)(unsafe.Pointer(&coutputs[0])), C.int(len(coutputs)),
		&cblind)
	retError = PrivErrorInfo[ret]

	blind = C.GoStringN(cblind.c_str, cblind.length)
	return
}

func CombineBlind(ptr uint64, blinds []string) (blind string, retError error) {
	const DEFAULT_BLIND string = "0000000000000000000000000000000000000000000000000000000000000000"
	return ExcessBlind(ptr, append(blinds, DEFAULT_BLIND), []string{DEFAULT_BLIND})
}

func BpRangeproofProve(ptr uint64, blind string, value uint64) (proof string, retError error) {
	c_blind := C.pure_string{}
	c_blind.c_str = C.CString(blind)
	c_blind.length = C.int(len([]rune(blind)))
	defer C.free(unsafe.Pointer(c_blind.c_str))

	c_proof := C.pure_string{}

	ret := C.PurecBpRangeproofProve(C.uint64_t(ptr), &c_blind, C.uint64_t(value), &c_proof)
	defer C.free(unsafe.Pointer(c_proof.c_str))
	retError = PrivErrorInfo[ret]
	proof = C.GoStringN(c_proof.c_str, c_proof.length)

	return
}

func BpRangeproofVerify(ptr uint64, commit string, proof string) error {
	c_commit := C.pure_string{}
	c_commit.c_str = C.CString(commit)
	c_commit.length = C.int(len([]rune(commit)))
	defer C.free(unsafe.Pointer(c_commit.c_str))

	c_proof := C.pure_string{}
	c_proof.c_str = C.CString(proof)
	c_proof.length = C.int(len([]rune(proof)))
	defer C.free(unsafe.Pointer(c_proof.c_str))

	return PrivErrorInfo[C.PurecBpRangeproofVerify(C.uint64_t(ptr), &c_commit, &c_proof)]
}

func HideAddr(ptr uint64, dest_pubkey string, input_nonce_seckey string) (hide_pubkey string, nonce_seckey string, nonce_pubkey string, ret_err error) {
	//P = Rb + hash(rc * Rb)G  rc = nonce_seckey  Rb = dest_pubkey
	nonce_seckey = input_nonce_seckey
	if nonce_seckey == "" {
		_, nonce_seckey, _ = CreateKeyPair(ptr)
	}
	nonce_pubkey, ret_err = GetPublicKey(ptr, nonce_seckey)
	if ret_err != nil {
		return
	}

	ecdh_key, ret_err := CreateEcdhKey(ptr, nonce_seckey, dest_pubkey)
	if ret_err != nil {
		return
	}

	ecdh_key_byte, _ := hex.DecodeString(ecdh_key)
	hashed_sec := sha256.Sum256([]byte(ecdh_key_byte))
	var hashed_sec_tmp = make(chan []byte, 1)
	hashed_sec_tmp <- hashed_sec[:]
	hashed_pub, ret_err := GetPublicKey(ptr, hex.EncodeToString(<-hashed_sec_tmp))
	if ret_err != nil {
		return
	}

	hide_pubkey, ret_err = CombinePublicKey(ptr, []string{dest_pubkey, hashed_pub})

	return
}

func main1() {
	var ptr1 uint64 = InitPrivacy()
	//ret1, commit1 := CreatePedersenCommit(ptr1, 3, "8888888888888888888888888888888888888888888888888888888888888888")
	//ret1, commit2 := CreatePedersenCommit(ptr1, 1, "4444444444444444444444444444444444444444444444444444444444444444")
	//ret1, commit3 := CreatePedersenCommit(ptr1, 2, "4444444444444444444444444444444444444444444444444444444444444444")
	//fmt.Println(ret1)
	//fmt.Println(commit1)
	//fmt.Println(commit2)
	//fmt.Println(commit3)

	//var pubkey string
	//var privkey string
	//ret1, pubkey, privkey = CreateKeyPair(ptr1)
	//fmt.Println("pubkey:", pubkey, ", privkey:", privkey)

	//var pubkey1 string
	//	ret1, pubkey1 = GetPublicKey(ptr1, "63542b9f52c15f79192d88e833cabf47bb65af26e719678d4f16606dcd22c368")
	//	fmt.Println("pubkey1:", pubkey1)
	//ret1, pubkey1 = GetPublicKey(ptr1, "63542b9f52c15f79192d88e833cabf47bb65af26e719678d4f16606dcd22c368")
	//fmt.Println("pubkey1:", pubkey1)

	//inputs := []string{"0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0", "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87"}
	//outputs := []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0bad"}

	//ret2 := TallyVerify(ptr1, inputs, outputs)
	//fmt.Println("ret1:", ret2)

	//DestroyPrivacy(ptr1)
	//C.secp256k1_context_create(C.int(0));
	//fmt.Println(C.sum(C.int(12), C.int(14)))
	//return commit1

	inputs := []string{"0871235f90f81a2fc5b1afbf2805667ab4c4a03c7af94f892712a64a61547393d0", "0837b273781da2842e5d3ecc04936195783e3af4ea7c37566710fb49ef8fe18f87", "0814c5fdf7753ba8d99123ee4a09eba6c07c641bb7d9982ea7c6aa53faa0faf6c1"}
	outputs := []string{"0853c901dcca5aa5fd00115bc28607875c55d7b23e7e4c6d840c5589d4dcdd0bad"}

	ret2 := TallyVerify(ptr1, inputs, outputs)
	fmt.Println("ret1:", ret2)

}
