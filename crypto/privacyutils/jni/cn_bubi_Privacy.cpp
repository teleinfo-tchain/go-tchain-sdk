#include <iostream>

#include <stdint.h>
#include <string.h>
#include "stdio.h"
#include "privacy/privacy.hpp"
#include "cn_bubi_Privacy.h"

#ifdef WIN32
#define FMT_I64 "%I64d"
#define FMT_I64_EX(fmt) "%" #fmt "I64d"
#define FMT_U64 "%I64u"
#define FMT_X64 "%I64X"
#define FMT_SIZE "%lu"
#else
#ifdef __x86_64__
#define FMT_I64 "%ld"
#define FMT_I64_EX(fmt) "%" #fmt "ld"
#define FMT_U64 "%lu"
#define FMT_X64 "%lX"
#define FMT_SIZE "%lu"
#else
#define FMT_I64 "%lld"
#define FMT_I64_EX(fmt) "%" #fmt "lld"
#define FMT_U64 "%llu"
#define FMT_X64 "%llX"
#define FMT_SIZE "%u"
#endif

#endif

std::string ToString(const int64_t val) {
	char buf[32] = { 0 };
#ifdef WIN32    
	_snprintf(buf, sizeof(buf), FMT_I64, val);
#else
	snprintf(buf, sizeof(buf), FMT_I64, val);
#endif
	return buf;
}

/// @brief to uint64 
static uint64_t Stoui64(const std::string &str) {
	uint64_t v = 0;
#ifdef WIN32
	sscanf_s(str.c_str(), FMT_U64, &v);
#else
	sscanf(str.c_str(), FMT_U64, &v);
#endif
	return v;
}


JNIEXPORT jlong JNICALL Java_cn_bubi_Privacy_init
(JNIEnv *env, jclass classObject) {
	secp256k1::Privacy* privacy = new secp256k1::Privacy();
	if (privacy == nullptr) return 0;
	(void)classObject;
	return (uintptr_t)privacy;  
}

JNIEXPORT void JNICALL Java_cn_bubi_Privacy_destroy
(JNIEnv *env, jclass classObject, jlong obj) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	if(privacy != nullptr) delete privacy;
	(void)classObject;(void)env;
	return;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_createPedersenCommit
(JNIEnv* env, jclass classObject,jlong obj, jlong value, jstring blind) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	std::string commit_str;
	do
	{
		if (privacy == nullptr || blind == nullptr || value < 0) break;
		const char *blind_ptr = { 0 };
		blind_ptr = env->GetStringUTFChars(blind, NULL);
		std::string blind_str = blind_ptr;
	
		code = privacy->CreatePedersenCommit(value, blind_str, commit_str);
		ret = ToString(code);
	} while (false);

	jobjectArray retArray = env->NewObjectArray(3, env->FindClass("java/lang/String"), env->NewStringUTF(""));

	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	env->SetObjectArrayElement(retArray, 2, env->NewStringUTF(commit_str.c_str()));

	return retArray;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_createEcdhKey
(JNIEnv* env, jclass classObject, jlong obj, jstring seckey, jstring pubkey) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	std::string ecdh_key;
	do
	{
		if (privacy == nullptr || seckey == nullptr || pubkey == nullptr) break;

		const char *seckey_ptr = { 0 };
		seckey_ptr = env->GetStringUTFChars(seckey, NULL);
		std::string seckey_str = seckey_ptr;

		const char *pubkey_ptr = { 0 };
		pubkey_ptr = env->GetStringUTFChars(pubkey, NULL);
		std::string pubkey_str = pubkey_ptr;

		
		privacy->CreateEcdhKey(seckey_str, pubkey_str, ecdh_key);
		code = privacy->CreateEcdhKey(seckey_str, pubkey_str, ecdh_key);
		ret = ToString(code);
	} while (false);
	
	jobjectArray retArray = env->NewObjectArray(3, env->FindClass("java/lang/String"), env->NewStringUTF(""));

	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	env->SetObjectArrayElement(retArray, 2, env->NewStringUTF(ecdh_key.c_str()));

	return retArray;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_createKeyPair
  (JNIEnv* env, jclass classObject, jlong obj) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	std::string pubkey;
	std::string seckey;

	do
	{
		if (privacy == nullptr) break;
		code = privacy->CreateKeyPair(pubkey, seckey);
		ret = ToString(code);
	} while (false);
	
	jobjectArray retArray = env->NewObjectArray(4, env->FindClass("java/lang/String"), env->NewStringUTF(""));

	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	env->SetObjectArrayElement(retArray, 2, env->NewStringUTF(pubkey.c_str()));
	env->SetObjectArrayElement(retArray, 3, env->NewStringUTF(seckey.c_str()));

    return retArray;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_getPublicKey
(JNIEnv* env, jclass classObject, jlong obj, jstring seckey) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	std::string pubkey;
	do
	{
		if (privacy == nullptr || seckey == nullptr) break;

		const char *seckey_ptr = { 0 };
		seckey_ptr = env->GetStringUTFChars(seckey, NULL);
		std::string seckey_str = seckey_ptr;
	    code = privacy->GetPublicKey(seckey_str, pubkey);
		ret = ToString(code);
	} while (false);
	
	jobjectArray retArray = env->NewObjectArray(3, env->FindClass("java/lang/String"), env->NewStringUTF(""));

	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	env->SetObjectArrayElement(retArray, 2, env->NewStringUTF(pubkey.c_str()));

	return retArray;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_bpRangeproofProve
(JNIEnv* env, jclass classObject, jlong obj, jstring blind, jlong value) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	std::string proof;
	do
	{
		if (privacy == nullptr || blind == nullptr || value < 0) break;

		const char *blind_ptr = { 0 };
		blind_ptr = env->GetStringUTFChars(blind, NULL);
		std::string blind_str = blind_ptr;

		code = privacy->BpRangeproofProve(blind_str, value, proof);
		ret = ToString(code);
	} while (false);

	jobjectArray retArray = env->NewObjectArray(3, env->FindClass("java/lang/String"), env->NewStringUTF(""));

	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	env->SetObjectArrayElement(retArray, 2, env->NewStringUTF(proof.c_str()));

	return retArray;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_bpRangeproofVerify
(JNIEnv* env, jclass classObject, jlong obj, jstring commit, jstring proof) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	do
	{
		if (privacy == nullptr || commit == nullptr || proof == nullptr) break;

		const char *commit_ptr = { 0 };
		commit_ptr = env->GetStringUTFChars(commit, NULL);
		std::string commit_str = commit_ptr;

		const char *proof_ptr = { 0 };
		proof_ptr = env->GetStringUTFChars(proof, NULL);
		std::string proof_str = proof_ptr;

		code = privacy->BpRangeproofVerify(commit_str, proof_str);
		ret = ToString(code);
	} while (false);

	jobjectArray retArray = env->NewObjectArray(2, env->FindClass("java/lang/String"), env->NewStringUTF(""));

	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));

	return retArray;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_ecdsaSign
(JNIEnv* env, jclass classObject, jlong obj, jstring seckey, jstring data) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	std::string sig;
	do
	{
		if (privacy == nullptr || seckey == nullptr || data == nullptr) break;

		const char *seckey_ptr = { 0 };
		seckey_ptr = env->GetStringUTFChars(seckey, NULL);
		std::string seckey_str = seckey_ptr;

		const char *data_ptr = { 0 };
		data_ptr = env->GetStringUTFChars(data, NULL);
		std::string data_str = data_ptr;

		code = privacy->EcdsaSign(seckey_str, data_str, sig);
		ret = ToString(code);
	} while (false);
	
	jobjectArray retArray = env->NewObjectArray(3, env->FindClass("java/lang/String"), env->NewStringUTF(""));

	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	env->SetObjectArrayElement(retArray, 2, env->NewStringUTF(sig.c_str()));

	return retArray;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_ecdsaVerify
(JNIEnv* env, jclass classObject, jlong obj, jstring pubkey, jstring data, jstring sig) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	do
	{
		if (privacy == nullptr || pubkey == nullptr || data == nullptr || sig == nullptr) break;

		const char *pubkey_ptr = { 0 };
		pubkey_ptr = env->GetStringUTFChars(pubkey, NULL);
		std::string pubkey_str = pubkey_ptr;

		const char *data_ptr = { 0 };
		data_ptr = env->GetStringUTFChars(data, NULL);
		std::string data_str = data_ptr;

		const char *sig_ptr = { 0 };
		sig_ptr = env->GetStringUTFChars(sig, NULL);
		std::string sig_str = sig_ptr;

		code = privacy->EcdsaVerify(pubkey_str, data_str, sig_str);
	} while (false);
	
	jobjectArray retArray = env->NewObjectArray(2, env->FindClass("java/lang/String"), env->NewStringUTF(""));
	ret = ToString(code);
	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	return retArray;

}


JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_excessSign
(JNIEnv* env, jclass classObject, jlong obj, jobjectArray inputs, jobjectArray outputs, jstring msg) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	std::string sig;
	jobjectArray retArray = env->NewObjectArray(3, env->FindClass("java/lang/String"), env->NewStringUTF(""));
	do
	{
		if (privacy == nullptr || msg == nullptr) break;

		jsize inputs_len = env->GetArrayLength(inputs);
		std::vector<std::string> inputs_vec;
		bool invalid_elem = false;
		for (int i = 0; i < inputs_len; i++) {
			jstring objString = (jstring)env->GetObjectArrayElement(inputs, i);
			if (objString == nullptr) {
				invalid_elem = true;
				break;
			}

			const char *commit_ptr = { 0 };
			commit_ptr = env->GetStringUTFChars(objString, NULL);
			std::string commit_str = commit_ptr;
			inputs_vec.push_back(commit_str);
		}

		jsize outputs_len = env->GetArrayLength(outputs);
		std::vector<std::string> outputs_vec;
		for (int i = 0; i < outputs_len; i++) {
			jstring objString = (jstring)env->GetObjectArrayElement(outputs, i);
			if (objString == nullptr) {
				invalid_elem = true;
				break;
			}

			const char *commit_ptr = { 0 };
			commit_ptr = env->GetStringUTFChars(objString, NULL);
			std::string commit_str = commit_ptr;
			outputs_vec.push_back(commit_str);
		}
		if (invalid_elem) break;

		const char *msg_ptr = { 0 };
		msg_ptr = env->GetStringUTFChars(msg, NULL);
		std::string msg_str = msg_ptr;

		code = privacy->ExcessSign(inputs_vec, outputs_vec, msg_str, sig);
	} while (false);

	ret = ToString(code);
	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	env->SetObjectArrayElement(retArray, 2, env->NewStringUTF(sig.c_str()));

	return retArray;
}

JNIEXPORT jobjectArray JNICALL Java_cn_bubi_Privacy_pedersenTallyVerify
(JNIEnv* env, jclass classObject, jlong obj, jobjectArray inputs, jobjectArray outputs, jstring msg, jstring sig) {
	secp256k1::Privacy* privacy = (secp256k1::Privacy*)(uintptr_t)obj;
	int64_t code = BPERRORCODE::ERRCODE_INVALID_PARAMETER;
	std::string ret = ToString(code);
	jobjectArray retArray = env->NewObjectArray(3, env->FindClass("java/lang/String"), env->NewStringUTF(""));
	do
	{
		if (privacy == nullptr || msg == nullptr || sig == nullptr) break;
		jsize inputs_len = env->GetArrayLength(inputs);
		std::vector<std::string> inputs_vec;
		bool invalid_elem = false;
		for (int i = 0; i < inputs_len; i++) {

			jstring objString = (jstring)env->GetObjectArrayElement(inputs, i);
			if (objString == nullptr) {
				invalid_elem = true;
				break;
			}

			const char *commit_ptr = { 0 };
			commit_ptr = env->GetStringUTFChars(objString, NULL);
			std::string commit_str = commit_ptr;
			if (commit_str.empty()) return retArray;

			inputs_vec.push_back(commit_str);
		}

		jsize outputs_len = env->GetArrayLength(outputs);
		std::vector<std::string> outputs_vec;
		for (int i = 0; i < outputs_len; i++) {
			jstring objString = (jstring)env->GetObjectArrayElement(outputs, i);
			if (objString == nullptr) {
				invalid_elem = true;
				break;
			}

			const char *commit_ptr = { 0 };
			commit_ptr = env->GetStringUTFChars(objString, NULL);
			std::string commit_str = commit_ptr;
			if (commit_str.empty()) return retArray;

			outputs_vec.push_back(commit_str);
		}
		if (invalid_elem) break;

		const char *msg_ptr = { 0 };
		msg_ptr = env->GetStringUTFChars(msg, NULL);
		std::string msg_str = msg_ptr;

		const char *sig_ptr = { 0 };
		sig_ptr = env->GetStringUTFChars(sig, NULL);
		std::string sig_str = sig_ptr;

		code = privacy->PedersenTallyVerify(inputs_vec, outputs_vec, msg_str, sig_str);
		ret = ToString(code);
	} while (false);
	
	env->SetObjectArrayElement(retArray, 0, env->NewStringUTF(ret.c_str()));
	env->SetObjectArrayElement(retArray, 1, env->NewStringUTF(privacy->GetErrorMsg(code).c_str()));
	return retArray;
}