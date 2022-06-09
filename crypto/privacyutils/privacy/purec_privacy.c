#include <stdio.h>

#if 0
#include "include/secp256k1.h"
#include "util.h"
#include "num_impl.h"
#include "field_impl.h"
#include "scalar_impl.h"
#include "group_impl.h"
#include "ecmult_impl.h"
#include "ecmult_const_impl.h"
#include "ecmult_gen_impl.h"
#include "ecdsa_impl.h"
#include "eckey_impl.h"
#include "hash_impl.h"
#include "scratch_impl.h"

#ifdef ENABLE_MODULE_GENERATOR
# include "include/secp256k1_generator.h"
#endif

#ifdef ENABLE_MODULE_COMMITMENT
# include "include/secp256k1_commitment.h"
#endif

#ifdef ENABLE_MODULE_RANGEPROOF
# include "include/secp256k1_rangeproof.h"
#endif

#ifdef ENABLE_MODULE_BULLETPROOF
# include "include/secp256k1_bulletproofs.h"
#endif

#include "purec_privacy.h"
#else
#include "src/secp256k1.c"
#include "src/testrand_impl.h"
#include "purec_privacy.h"

#ifdef MAIN
void main() {
	char blind_str[65];
	memset(blind_str, 0, 65);
	strcpy(blind_str, "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234");
	pure_string blind, commitp;
	blind.c_str = (char *)blind_str;
	blind.length = 64;
	uint64_t obj_ptr = InitPrivacy();
	PurecCreatePedersenCommit(obj_ptr, 1, &blind, &commitp);
	printf("purec commit:%s\n", commitp.c_str);
}
#endif

 #endif


uint64_t InitPrivacy(){
    purec_privacy *ptr = ( purec_privacy *)malloc(sizeof(purec_privacy));
    ptr->ctx_ = secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
	ptr->scratch_ = secp256k1_scratch_space_create(ptr->ctx_, 1024 * 1024);
	ptr->gens_ = secp256k1_bulletproof_generators_create(ptr->ctx_, &secp256k1_generator_const_g, 256);
	secp256k1_context_set_illegal_callback(ptr->ctx_, IllegalCallBack, NULL);
	secp256k1_context_set_error_callback(ptr->ctx_, ErrorCallBack, NULL);

    if(ptr == NULL) return 0;
    return (uint64_t)ptr;
}

void DestroyPrivacy(uint64_t obj_ptr){
	purec_privacy *ptr =  ( purec_privacy *)obj_ptr;
	if (ptr == NULL) return ;
	if (!ptr->ctx_) secp256k1_context_destroy(ptr->ctx_);
	if (!ptr->scratch_) secp256k1_scratch_destroy(ptr->scratch_);
	if (!ptr->gens_) secp256k1_bulletproof_generators_destroy(ptr->ctx_, ptr->gens_);
	free(ptr);
}


void IllegalCallBack(const char* msg, void* data) {
	//std::string error_msg = "[privacy] illegal argument:";
	//error_msg += msg;
	//throw std::runtime_error(error_msg);
}

void ErrorCallBack(const char* msg, void* data) {
	//std::string error_msg = "[privacy] internal consistency check failed:";
	//error_msg += msg;
	//throw std::runtime_error(error_msg);
}

int64_t PurecCreatePedersenCommit(uint64_t obj_ptr, uint64_t value, const pure_string  *blind, pure_string *commit) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;
	if (blind->length != 64) {
		return PUREC_ERRCODE_INVALID_PARAMETER;
	}

	secp256k1_pedersen_commitment pc;
	unsigned char output[33];
	unsigned char blind_data[32];
	PurecHexStrToArray(blind, blind_data);
	//printf("blind_data[31]:%d\n", blind_data[31]);
	//printf("blind_data[16]:%d\n", blind_data[16]);
	//try {
	if (!secp256k1_pedersen_commit(ptr->ctx_, &pc, blind_data, value, &secp256k1_generator_const_h, &secp256k1_generator_const_g)) {
		return PUREC_ERRCODE_CREATE_PEDERSEN;
		}

	if (!secp256k1_pedersen_commitment_serialize(ptr->ctx_, output, &pc)) // return 1 always
	{
		return PUREC_ERRCODE_SERIALIZE_PEDERSEN;
	}
	//}
	//catch (std::exception& e) {
	//	error_msg_ = e.what();
	//	return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
	//}

	PurecArrayToHexStr(output, 33, commit);
	return PUREC_ERRCODE_SUCCESS;
}

int64_t PurecPedersenTallyVerify(uint64_t obj_ptr, const pure_string **inputs, int inputs_length,
	const pure_string **outputs, int outputs_length, 
	const pure_string *msg, 
	const pure_string *sig) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;
	//printf("inputlen:%d, output:%d, msgleng:%d\n", inputs_length, outputs_length, msg->length);
	if (inputs_length == 0 || inputs_length > 100 || outputs_length == 0 || outputs_length > 100 || msg->length != 32) {
		return PUREC_ERRCODE_INVALID_PARAMETER;
	}
	const pure_string *inputs_array = (const pure_string *)inputs;
	const pure_string *outputs_array = (const pure_string *)outputs;


	/** A group element of the secp256k1 curve, in jacobian coordinates. */
	//typedef struct {
	//	secp256k1_fe x; /* actual X: x/z^2 */
	//	secp256k1_fe y; /* actual Y: y/z^3 */
	//	secp256k1_fe z;
	//	int infinity; /* whether this represents the point at infinity */
	//} secp256k1_gej;
	secp256k1_gej accj;
	/** A group element of the secp256k1 curve, in affine coordinates. */
	//typedef struct {
	//	secp256k1_fe x;
	//	secp256k1_fe y;
	//	int infinity; /* whether this represents the point at infinity */
	//} secp256k1_ge;
	secp256k1_ge add;

//	try {
		secp256k1_gej_set_infinity(&accj);
		int i = 0;
		for (; i < outputs_length; i++) {
			if (outputs_array[i].length != 66) return PUREC_ERRCODE_INVALID_PARAMETER;

			unsigned char commit_data[33];
			secp256k1_pedersen_commitment output_commit;
			PurecHexStrToArray(&outputs_array[i], commit_data);
			
			if (!secp256k1_pedersen_commitment_parse(ptr->ctx_, &output_commit, commit_data)) {
				return PUREC_ERRCODE_PARSE_PEDERSEN;
			}
			
			secp256k1_ge_clear(&add);
			secp256k1_pedersen_commitment_load(&add, &output_commit);
			secp256k1_gej_add_ge_var(&accj, &accj, &add, NULL);
		}


		secp256k1_gej_neg(&accj, &accj);
		i = 0;
		for (; i < inputs_length; i++) {
			if (inputs_array[i].length != 66) return PUREC_ERRCODE_INVALID_PARAMETER;

			unsigned char commit_data[33];
			secp256k1_pedersen_commitment input_commit;
			PurecHexStrToArray(&inputs_array[i], commit_data);

			if (!secp256k1_pedersen_commitment_parse(ptr->ctx_, &input_commit, commit_data)) {
				return PUREC_ERRCODE_PARSE_PEDERSEN;
			}
			
			secp256k1_ge_clear(&add);
			secp256k1_pedersen_commitment_load(&add, &input_commit);

			secp256k1_gej_add_ge_var(&accj, &accj, &add, NULL);
		}

		// create pubkey from excess
		size_t pub_len = 33;
		unsigned char pub_compress[33];
		secp256k1_pubkey pubkey;

		secp256k1_ge_clear(&add);
		secp256k1_ge_set_gej(&add, &accj);
		secp256k1_pubkey_save(&pubkey, &add);

		if (!secp256k1_ec_pubkey_serialize(ptr->ctx_, pub_compress, &pub_len, &pubkey, SECP256K1_EC_COMPRESSED)) {
			return PUREC_ERRCODE_SERIALIZE_PUBKEY;
		}

		pure_string pubkey_str;
		PurecArrayToHexStr(pub_compress, 33, &pubkey_str);

		// create commit from excess
		unsigned char excess_data[33];
		secp256k1_pedersen_commitment excess_commit;

		secp256k1_pedersen_commitment_save(&excess_commit, &add);

		if (!secp256k1_pedersen_commitment_serialize(ptr->ctx_, excess_data, &excess_commit)) {
			PurecDelString(&pubkey_str);
			return PUREC_ERRCODE_SERIALIZE_PEDERSEN;
		}
		pure_string excess_str;
		PurecArrayToHexStr(excess_data, 33, &excess_str);

		pure_string **new_outputs;
		new_outputs = malloc(sizeof(pure_string) * (outputs_length+1));
		pure_string *new_outputs_array = (pure_string *)new_outputs;
		int counter = 0;
		for (; counter < outputs_length; counter++) {
			new_outputs_array[counter] = outputs_array[counter];
		}
		new_outputs_array[counter] = excess_str;

		// verify input - output = 0
		int64_t ret = PurecTallyVerify(obj_ptr, (const pure_string **)inputs, inputs_length, (const pure_string **)new_outputs, counter+1);
		//delete the new_outputs immediately and excess_str
		free(new_outputs);
		PurecDelString(&excess_str);

		if (ret != PUREC_ERRCODE_SUCCESS) {
			PurecDelString(&pubkey_str);
			return ret;
		}

		// verify blinds_sum * G = Excess
		int64_t veresult = PurecEcdsaVerify(obj_ptr, &pubkey_str, msg, sig);
		PurecDelString(&pubkey_str);
		return veresult;
// 	}
// 	catch (std::exception& e) {
// 		//error_msg_ = e.what();
// 		return PUREC_ERRCODE_BP_LIB_INTERNAL;
// 	}
	return PUREC_ERRCODE_SUCCESS;
}

//��֤��ʽƽ��
int64_t PurecTallyVerify(uint64_t obj_ptr,
	const pure_string **inputs, int inputs_length,
	const pure_string **outputs, int outputs_length) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;

	const pure_string *inputs_array = (const pure_string *)inputs;
	const pure_string *outputs_array = (const pure_string *)outputs;

	//printf("inputlen:%d, output:%d\n", inputs_length, outputs_length);
	int64_t input_size = inputs_length;
	int64_t output_size = outputs_length;
	if ((input_size > 100) || (output_size > 101)) { // the extra one is excess
		return PUREC_ERRCODE_OUT_RANGE;
	}

//	try {
		secp256k1_pedersen_commitment input_list[100];
		secp256k1_pedersen_commitment output_list[101];
		const secp256k1_pedersen_commitment *plist[100];
		const secp256k1_pedersen_commitment *nlist[101];
		int i = 0;
		for (; i < input_size; i++) {
			unsigned char commit_data[33];
			PurecHexStrToArray(&inputs_array[i], commit_data);
			//printf("input:%s\n", inputs_array[i].c_str);

			if (!secp256k1_pedersen_commitment_parse(ptr->ctx_, &input_list[i], commit_data)) {
				return PUREC_ERRCODE_PARSE_PEDERSEN;
			}
			plist[i] = &input_list[i];
		}

		i = 0;
		for (; i < output_size; i++) {
			unsigned char commit_data[33];
			//printf("output:%s\n", outputs_array[i].c_str);
			PurecHexStrToArray(&outputs_array[i], commit_data);

			if (!secp256k1_pedersen_commitment_parse(ptr->ctx_, &output_list[i], commit_data)) {
				return PUREC_ERRCODE_PARSE_PEDERSEN;
			}
			nlist[i] = &output_list[i];
		}
		
		if (!secp256k1_pedersen_verify_tally(ptr->ctx_, plist, input_size, nlist, output_size)) {
			return PUREC_ERRCODE_VERIFY_TALLY;
		}
// 	}
// 	catch (std::exception& e) {
// 		//error_msg_ = e.what();
// 		return PUREC_ERRCODE_BP_LIB_INTERNAL;
// 	}
	return PUREC_ERRCODE_SUCCESS;
}

int64_t PurecEcdsaVerify(uint64_t obj_ptr, const pure_string *pub_key, const pure_string *data, const pure_string *sig) {
	if (pub_key->length != 66 || data->length != 32 || sig->length > 256) return PUREC_ERRCODE_INVALID_PARAMETER;

	purec_privacy *ptr = (purec_privacy *)obj_ptr;

	unsigned char pub[33];
	unsigned char sig_data[128];
	PurecHexStrToArray(pub_key, pub);
//	try {
		secp256k1_pubkey pubkey;
		if (!secp256k1_ec_pubkey_parse(ptr->ctx_, &pubkey, pub, 33)) return PUREC_ERRCODE_PARSE_PUBKEY;

		secp256k1_ecdsa_signature ecdsa_sig;
//		unsigned char* sig_data = new unsigned char[sig->length / 2];
		PurecHexStrToArray(sig, sig_data);

		if (!secp256k1_ecdsa_signature_parse_der(ptr->ctx_, &ecdsa_sig, sig_data, sig->length / 2)) return PUREC_ERRCODE_ECDSA_PARSE;

		if (!secp256k1_ecdsa_verify(ptr->ctx_, &ecdsa_sig, (unsigned char*)data->c_str, &pubkey)) return PUREC_ERRCODE_ECDSA_VERIFY;
// 	}
// 	catch (std::exception& e) {
// 		//error_msg_ = e.what();
// 		return PUREC_ERRCODE_BP_LIB_INTERNAL;
// 	}
	return PUREC_ERRCODE_SUCCESS;
}

int64_t PurecCreateEcdhKey(uint64_t obj_ptr, const pure_string *priv_key, const pure_string *pub_key, pure_string *ecdh_key) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;
	if (priv_key->length != 64 || pub_key->length != 66) return PUREC_ERRCODE_INVALID_PARAMETER;

	unsigned char priv[32];
	PurecHexStrToArray(priv_key, priv);

	unsigned char pub[33];
	PurecHexStrToArray(pub_key, pub);
	//try {
		secp256k1_pubkey pubkey;

		if (!secp256k1_ec_pubkey_parse(ptr->ctx_, &pubkey, pub, 33)) return PUREC_ERRCODE_PARSE_PEDERSEN;

		unsigned char key[32];
		if (!secp256k1_ecdh(ptr->ctx_, key, &pubkey, priv)) return PUREC_ERRCODE_CREATE_PUBKEY;

		PurecArrayToHexStr(key, 32, ecdh_key);
	//}
	//catch (std::exception& e) {
	//	error_msg_ = e.what();
	//	return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
	//}
		return PUREC_ERRCODE_SUCCESS;
}

int64_t PurecEcdsaSign(uint64_t obj_ptr, const pure_string *priv_key, const pure_string *data, pure_string *sig) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;
	if (priv_key->length != 64 || data->length != 32) return PUREC_ERRCODE_INVALID_PARAMETER;

	unsigned char priv[32];
	PurecHexStrToArray(priv_key, priv);

	size_t siglen = 72;
	unsigned char sig_data[72];
	//try {
		secp256k1_ecdsa_signature ecdsa_sig;
		if (!secp256k1_ecdsa_sign(ptr->ctx_, &ecdsa_sig, (unsigned char*)data->c_str, priv, NULL, NULL)) return PUREC_ERRCODE_ECDSA_CRATE;

		if (!secp256k1_ecdsa_signature_serialize_der(ptr->ctx_, sig_data, &siglen, &ecdsa_sig)) return PUREC_ERRCODE_ECDSA_SERIALIZE;

		PurecArrayToHexStr(sig_data, (int)siglen, sig);
	//}
	//catch (std::exception& e) {
	//	error_msg_ = e.what();
	//	return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
	//}
		return PUREC_ERRCODE_SUCCESS;
}

int64_t PurecCreateKeyPair(uint64_t obj_ptr, pure_string *pub_key, pure_string *priv_key) {
	// set random seed
	unsigned char priv[32];
	//try {
		int64_t ret = PurecRand(priv);
		if (ret != PUREC_ERRCODE_SUCCESS) {
			return ret;
		}

		PurecArrayToHexStr(priv, 32, priv_key);
		return PurecGetPublicKey(obj_ptr, priv_key, pub_key);
	//}
	//catch (std::exception& e) {
	//	error_msg_ = e.what();
	//	return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
	//}
		return PUREC_ERRCODE_SUCCESS;
}

int64_t PurecGetPublicKey(uint64_t obj_ptr, const pure_string *priv_key, pure_string *pub_key) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;
	if (priv_key->length != 64) return PUREC_ERRCODE_INVALID_PARAMETER;

	unsigned char priv[32];
	PurecHexStrToArray(priv_key, priv);
	//printf("prikey:%s\n, prikey[1]:%d, prikey[10]:%d\n", priv_key->c_str, priv[1], priv[10]);
	//try {
		secp256k1_pubkey pub;
		if (!secp256k1_ec_pubkey_create(ptr->ctx_, &pub, priv)) return PUREC_ERRCODE_CREATE_PUBKEY;
		
		//printf("pub[2]:%d, pub[10]:%d\n", pub.data[2], pub.data[10]);

		size_t pub_len = 33;
		unsigned char pub_compress[33];
		if (!secp256k1_ec_pubkey_serialize(ptr->ctx_, pub_compress, &pub_len, &pub, SECP256K1_EC_COMPRESSED)) return PUREC_ERRCODE_SERIALIZE_PUBKEY;

		PurecArrayToHexStr(pub_compress, 33, pub_key);
	//}
	//catch (std::exception& e) {
	//	error_msg_ = e.what();
	//	return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
	//}
		return PUREC_ERRCODE_SUCCESS;
}

//
int64_t PurecExcessSign(uint64_t obj_ptr, const pure_string **inputs, int inputs_length,
	const pure_string **outputs, int outputs_length,
	const pure_string *msg, pure_string *sig) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;
	unsigned char* blinds[200];
	unsigned char blinds_temp[200][32];
	unsigned char blind_out[32];
	int64_t input_size = inputs_length;
	int64_t output_size = outputs_length;
	if ((input_size > 100) || (output_size > 100)) {
		return PUREC_ERRCODE_OUT_RANGE;
	}
	const pure_string *inputs_array = (const pure_string *)inputs;
	const pure_string *outputs_array = (const pure_string *)outputs;

	int i = 0;
	for (; i < input_size; i++) {
		if (inputs_array[i].length != 64) return PUREC_ERRCODE_INVALID_PARAMETER;
		PurecHexStrToArray(&inputs_array[i], blinds_temp[i]);
		blinds[i] = blinds_temp[i];
	}

	i = 0;
	for (; i < output_size; i++) {
		if (outputs_array[i].length != 64) return PUREC_ERRCODE_INVALID_PARAMETER;
		PurecHexStrToArray(&outputs_array[i], blinds_temp[i + input_size]);
		blinds[i + input_size] = blinds_temp[i + input_size];
	}

	//try {
	if (secp256k1_pedersen_blind_sum(ptr->ctx_, blind_out, (const unsigned char * const*)blinds, input_size + output_size, input_size) != 1) {
			return PUREC_ERRCODE_BLIND_SUM;
		}

		pure_string priv_key;
		PurecArrayToHexStr(blind_out, 32, &priv_key);
		int64_t ret = PurecEcdsaSign(obj_ptr, &priv_key, msg, sig);
		PurecDelString(&priv_key);

		return ret;
	//}
	//catch (std::exception& e) {
	//	error_msg_ = e.what();
	//	return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
	//}
	//	return PUREC_ERRCODE_SUCCESS;
}

int64_t PurecBpRangeproofProve(uint64_t obj_ptr, const pure_string *blind, uint64_t value, pure_string *proof) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;
	if (blind->length != 64) return PUREC_ERRCODE_INVALID_PARAMETER;

	unsigned char* blind_list[1];
	uint64_t value_list[1];

	unsigned char blind_data[32];
	PurecHexStrToArray(blind, blind_data);
	blind_list[0] = blind_data;
	value_list[0] = value;

	size_t plen = 2000;
	unsigned char proof_data[2000];
	//try {
		unsigned char random[32];
		int64_t ret = PurecRand(random);
		if (ret != PUREC_ERRCODE_SUCCESS) {
			return ret;
		}

		if (!secp256k1_bulletproof_rangeproof_prove(ptr->ctx_, ptr->scratch_, ptr->gens_, proof_data, &plen, value_list, NULL, (const unsigned char * const*)blind_list, 1, &secp256k1_generator_const_h, 64, random, NULL, 0)) {
			return PUREC_ERRCODE_RANGEPROOF_PROVE;
		}

		PurecArrayToHexStr(proof_data, plen, proof);
	//}
	//catch (std::exception& e) {
	//	error_msg_ = e.what();
	//	return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
	//}

	//proof = ArrayToHexStr(proof_data, plen);

		return PUREC_ERRCODE_SUCCESS;
}

int64_t PurecBpRangeproofVerify(uint64_t obj_ptr, const  pure_string *commit, const  pure_string *proof) {
	purec_privacy *ptr = (purec_privacy *)obj_ptr;
	if (commit->length != 66 || proof->length == 0) return PUREC_ERRCODE_INVALID_PARAMETER;

	unsigned char commit_data[33];
	PurecHexStrToArray(commit, commit_data);

	unsigned char proof_data[2000];
	PurecHexStrToArray(proof, proof_data);

	secp256k1_pedersen_commitment pc;
	//try {
	int ret = secp256k1_pedersen_commitment_parse(ptr->ctx_, &pc, commit_data);
		if (!ret) {
			return PUREC_ERRCODE_CREATE_PEDERSEN;
		}

		ret = secp256k1_bulletproof_rangeproof_verify(ptr->ctx_, ptr->scratch_, ptr->gens_, proof_data, proof->length / 2, NULL, &pc, 1, 64, &secp256k1_generator_const_h, NULL, 0);
		if (!ret) {
			return PUREC_ERRCODE_RANGEPROOF_VERIFY;
		}
	//}
	//catch (std::exception& e) {
	//	error_msg_ = e.what();
	//	return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
	//}
		return PUREC_ERRCODE_SUCCESS;
}

int PurecHexStrToArray(const pure_string *hex_str, unsigned char *array) {
	size_t i = 0;
	if (hex_str->length % 2 != 0 || hex_str->length == 0) {
		return 0;
	}

	for (; i <hex_str->length - 1; i = i + 2) {
		uint8_t high = 0;
		if (hex_str->c_str[i] >= '0' && hex_str->c_str[i] <= '9')
			high = (hex_str->c_str[i] - '0');
		else if (hex_str->c_str[i] >= 'a' && hex_str->c_str[i] <= 'f')
			high = (hex_str->c_str[i] - 'a' + 10);
		else if (hex_str->c_str[i] >= 'A' && hex_str->c_str[i] <= 'F') {
			high = (hex_str->c_str[i] - 'A' + 10);
		}
		else {
			return 0;
		}

		uint8_t low = 0;
		if (hex_str->c_str[i + 1] >= '0' && hex_str->c_str[i + 1] <= '9')
			low = (hex_str->c_str[i + 1] - '0');
		else if (hex_str->c_str[i + 1] >= 'a' && hex_str->c_str[i + 1] <= 'f')
			low = (hex_str->c_str[i + 1] - 'a' + 10);
		else  if (hex_str->c_str[i + 1] >= 'A' && hex_str->c_str[i + 1] <= 'F') {
			low = (hex_str->c_str[i + 1] - 'A' + 10);
		}
		else {
			return 0;
		}

		int valuex = (high << 4) + low;
		//sscanf(hex_string.substr(i, 2).c_str(), "%x", &valuex);
		array[i / 2] = (char)valuex;
	}

	return 1;
}

void PurecArrayToHexStr(unsigned char *array, int len, pure_string *raw_str) {
	size_t i = 0;
	PurecNewString(raw_str, len * 2);
	for (; i < len; i++) {
		uint8_t item = array[i];
		uint8_t high = (item >> 4);
		uint8_t low = (item & 0x0f);
		raw_str->c_str[2 * i] = (high >= 0 && high <= 9) ? (high + '0') : (high - 10 + 'a');
		raw_str->c_str[2 * i + 1] = (low >= 0 && low <= 9) ? (low + '0') : (low - 10 + 'a');
	}
}

void InitPurecString(pure_string *raw_str) {
	raw_str->c_str = NULL;
	raw_str->length = 0;
}

pure_string *PurecNewString(pure_string *raw_str, size_t length) {
	raw_str->c_str = malloc(length);
	raw_str->length = length;
	return raw_str;
}

void PurecDelString(pure_string *p_str) {
	free(p_str->c_str);
	p_str->length = 0;
}

int64_t PurecRand(unsigned char* rand) {
	int64_t ret = 0;
	unsigned char seed16[16] = { 0 };
// #ifdef  WIN32
// 	std::random_device rd;
// 	for (int i = 0; i < 16; i++) {
// 		seed16[i] = (uint8_t)std::uniform_int_distribution<uint16_t>(0, 255)(rd);
// 	}
// #else
	FILE* frand = fopen("/dev/urandom", "r");
	int result = (frand == NULL) || fread(&seed16, sizeof(seed16), 1, frand) != sizeof(seed16);
	if (!result) {
		return PUREC_ERRCODE_RANDOM_ERROR;
	}
	if (frand) {
		fclose(frand);
	}
//#endif 
	uint64_t t = time(NULL) * (uint64_t)1337;
	seed16[0] ^= t;
	seed16[1] ^= t >> 8;
	seed16[2] ^= t >> 16;
	seed16[3] ^= t >> 24;
	seed16[4] ^= t >> 32;
	seed16[5] ^= t >> 40;
	seed16[6] ^= t >> 48;
	seed16[7] ^= t >> 56;
	secp256k1_rand_seed(seed16);

	secp256k1_rand256(rand);
	return ret;
}