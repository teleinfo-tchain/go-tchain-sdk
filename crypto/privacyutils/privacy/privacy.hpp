#ifndef PRIVACY_HPP_
#define PRIVACY_HPP_

#include <string>
#include <vector>
#include <sstream>

#ifndef _WIN32
    #include <stdexcept>
#endif

#include <stdint.h>


struct secp256k1_bulletproof_generators;
typedef struct secp256k1_context_struct secp256k1_context;
typedef struct secp256k1_scratch_space_struct secp256k1_scratch;

enum BPERRORCODE {
	ERRCODE_SUCCESS = 0,
	ERRCODE_RANDOM_ERROR = 201, //Generates a random number error!
	ERRCODE_INVALID_PARAMETER = 202, //Invalid parameters
	ERRCODE_CREATE_PEDERSEN = 203,   //Failed to create pedersen commitment
	ERRCODE_PARSE_PEDERSEN = 204, //Failed to parse pedersen commitment
	ERRCODE_SERIALIZE_PEDERSEN = 205, //Failed to serialize pedersen commitment
	ERRCODE_SERIALIZE_PUBKEY = 206, //Failed to serialize pubkey
	ERRCODE_VERIFY_TALLY = 207, //Failed to verify tally
	ERRCODE_CREATE_PUBKEY = 208, //Failed to create pubkey
	ERRCODE_RANGEPROOF_PROVE = 209, //Failed to generate rangeproof
	ERRCODE_RANGEPROOF_VERIFY = 210, //Failed to verify rangeproof
	ERRCODE_OUT_RANGE = 211, //Out of range
	ERRCODE_ECDSA_CRATE = 212, //Failed to create ecdsa signature
	ERRCODE_ECDSA_SERIALIZE = 213, //Failed to serialize ecdsa signature
	ERRCODE_ECDSA_VERIFY = 214, //Failed to verify ecdsa signature
	ERRCODE_ECDSA_PARSE = 215, //Failed to parse ecdsa signature
	ERRCODE_PARSE_PUBKEY = 216, //Failed to parse pubkey
	ERRCODE_BLIND_SUM = 217, //Failed to blind sum
	ERRCODE_UNKNOWN = 218,//An unknown error
	ERRCODE_BP_LIB_INTERNAL = 219,//BP library internal error
};

namespace secp256k1 {
	class Privacy {
	public:
		Privacy();
		~Privacy();

		//用于创建Pedersen承若
		int64_t CreatePedersenCommit(uint64_t value, const std::string& blind, std::string& commit);
		//通过私钥*公钥生成ECDH共享秘钥
		int64_t CreateEcdhKey(const std::string& priv_key, const std::string& pub_key, std::string& ecdh_key);

	    //创建公私钥对
		int64_t CreateKeyPair(std::string& pub_key, std::string& priv_key);
		//获取公钥
		int64_t GetPublicKey(const std::string& priv_key, std::string& pub_key);

		/*int64_t BpRangeproofProve(const std::vector<std::string>& blinds, const std::vector<uint64_t>& values, std::string& proof);
		int64_t BpRangeproofVerify(const std::vector<std::string>& commit, const std::string& proof);*/

		//生成单个Rangeproof证明
		int64_t BpRangeproofProve(const std::string& blind, uint64_t value, std::string& proof);
		//验证单个Rangeproof证明
		int64_t BpRangeproofVerify(const std::string& commit, const std::string& proof);

		//ECDSA签名
		int64_t EcdsaSign(const std::string& priv_key, const std::string& data, std::string& sig);
		//验证ECDSA签名
		int64_t EcdsaVerify(const std::string& pub_key, const std::string& data, const std::string& sig);

		int64_t ExcessBlind(const std::vector<std::string>& inputs, const std::vector<std::string>& outputs, std::string& blind);

		//验证余项签名
		int64_t ExcessSign(const std::vector<std::string>& inputs, const std::vector<std::string>& outputs, const std::string& msg, std::string& sig);
		//验证等式平衡
		int64_t PedersenTallyVerify(const std::vector<std::string>& inputs, const std::vector<std::string>& outputs, const std::string& msg, const std::string& sig);
		
		//Combine public key
		int64_t CombinePublicKey(const std::vector<std::string>& inputs, std::string &combine_key);

		std::string GetErrorMsg(int64_t err_code);
		std::string GetErrorMsg() { return error_msg_; }
		static std::string ArrayToHexStr(unsigned char *array, int len) {
			std::string result;
			result.resize(len * 2);
			for (size_t i = 0; i < len; i++) {
				uint8_t item = array[i];
				uint8_t high = (item >> 4);
				uint8_t low = (item & 0x0f);
				result[2 * i] = (high >= 0 && high <= 9) ? (high + '0') : (high - 10 + 'a');
				result[2 * i + 1] = (low >= 0 && low <= 9) ? (low + '0') : (low - 10 + 'a');
			}
			return result;
		}

		static void HexStrToArray(const std::string& hex_str, unsigned char *array) {
			unsigned int c;
			for (int i = 0; i < hex_str.size(); i += 2) {
				std::istringstream hex_stream(hex_str.substr(i, 2));
				hex_stream >> std::hex >> c;
				array[i / 2] = c;
			}
		}

		int64_t TallyVerify(const std::vector<std::string>& inputs, const std::vector<std::string>& outputs);
	private:
		static void ErrorCallBack(const char* msg, void* data){
			std::string error_msg = "[privacy] internal consistency check failed:";
			error_msg += msg;
			throw std::runtime_error(error_msg);
		}

		static void IllegalCallBack(const char* msg, void* data){
			std::string error_msg = "[privacy] illegal argument:";
			error_msg += msg;
			throw std::runtime_error(error_msg);
		}

	private:
		secp256k1_context* ctx_;
		secp256k1_scratch* scratch_;
		secp256k1_bulletproof_generators* gens_;
		std::string error_msg_;
	};
}

#endif
