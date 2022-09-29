#ifndef PUREC_PRIVACY_H_
#define PUREC_PRIVACY_H_

#ifdef __cplusplus
extern "C" {
#endif

enum PUREC_BPERRORCODE {
	PUREC_ERRCODE_SUCCESS = 0,
	PUREC_ERRCODE_RANDOM_ERROR = 201, //Generates a random number error!
	PUREC_ERRCODE_INVALID_PARAMETER = 202, //Invalid parameters
	PUREC_ERRCODE_CREATE_PEDERSEN = 203,   //Failed to create pedersen commitment
	PUREC_ERRCODE_PARSE_PEDERSEN = 204, //Failed to parse pedersen commitment
	PUREC_ERRCODE_SERIALIZE_PEDERSEN = 205, //Failed to serialize pedersen commitment
	PUREC_ERRCODE_SERIALIZE_PUBKEY = 206, //Failed to serialize pubkey
	PUREC_ERRCODE_VERIFY_TALLY = 207, //Failed to verify tally
	PUREC_ERRCODE_CREATE_PUBKEY = 208, //Failed to create pubkey
	PUREC_ERRCODE_RANGEPROOF_PROVE = 209, //Failed to generate rangeproof
	PUREC_ERRCODE_RANGEPROOF_VERIFY = 210, //Failed to verify rangeproof
	PUREC_ERRCODE_OUT_RANGE = 211, //Out of range
	PUREC_ERRCODE_ECDSA_CRATE = 212, //Failed to create ecdsa signature
	PUREC_ERRCODE_ECDSA_SERIALIZE = 213, //Failed to serialize ecdsa signature
	PUREC_ERRCODE_ECDSA_VERIFY = 214, //Failed to verify ecdsa signature
	PUREC_ERRCODE_ECDSA_PARSE = 215, //Failed to parse ecdsa signature
	PUREC_ERRCODE_PARSE_PUBKEY = 216, //Failed to parse pubkey
	PUREC_ERRCODE_BLIND_SUM = 217, //Failed to blind sum
	PUREC_ERRCODE_UNKNOWN = 218,//An unknown error
	PUREC_ERRCODE_BP_LIB_INTERNAL = 219,//BP library internal error
};

typedef struct tag_purec_privacy {
	secp256k1_context* ctx_;
	secp256k1_scratch* scratch_;
	secp256k1_bulletproof_generators* gens_;
	char error_msg_[64];
}purec_privacy;

typedef struct tag_purec_string {
	char *c_str;
	int length;
}pure_string;


uint64_t InitPrivacy();
void DestroyPrivacy(uint64_t obj_ptr);
void IllegalCallBack(const char* msg, void* data);
void ErrorCallBack(const char* msg, void* data);
int64_t PurecCreatePedersenCommit(uint64_t obj_ptr, uint64_t value, const pure_string *blind, pure_string *commit);
int64_t PurecTallyVerify(uint64_t obj_ptr, const pure_string **inputs, int inputs_length, const pure_string **outputs, int outputs_length);

int64_t PurecEcdsaVerify(uint64_t obj_ptr, const pure_string *pub_key, const pure_string *data, const pure_string *sig);
int64_t PurecCreateEcdhKey(uint64_t obj_ptr, const pure_string *priv_key, const pure_string *pub_key, pure_string *ecdh_key);
int64_t PurecEcdsaSign(uint64_t obj_ptr, const pure_string *priv_key, const pure_string *data, pure_string *sig);
int64_t PurecPublicKeyCombine(uint64_t obj_ptr, const pure_string **inputs, int inputs_length, pure_string *pubCombine);

int64_t PurecCreateKeyPair(uint64_t obj_ptr, pure_string *pub_key, pure_string *priv_key);
int64_t PurecGetPublicKey(uint64_t obj_ptr, const pure_string *priv_key, pure_string *pub_key);
int64_t PurecExcessSign(uint64_t obj_ptr, const pure_string **inputs, int inputs_length,
	const pure_string **outputs, int outputs_length,
	const pure_string *msg, pure_string *sig);

int64_t PurecBpRangeproofProve(uint64_t obj_ptr, const pure_string *blind, uint64_t value, pure_string *proof);
int64_t PurecBpRangeproofVerify(uint64_t obj_ptr, const  pure_string *commit, const  pure_string *proof);

int64_t PurecRand(unsigned char* rand);

int PurecHexStrToArray(const pure_string *hex_str, unsigned char *array);
void PurecArrayToHexStr(unsigned char *array, int len, pure_string *raw_str);

void InitPurecString(pure_string *raw_str);
pure_string *PurecNewString(pure_string *raw_str, size_t length);
void PurecDelString(pure_string *raw_str);

#ifdef __cplusplus
}
#endif

#endif