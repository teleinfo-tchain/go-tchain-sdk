#include <iostream>
#include <random>
#include <time.h>


#include "src/secp256k1.c"
#include "src/testrand_impl.h"
#include "privacy.hpp"

int64_t rand(unsigned char* rand){
	int64_t ret = 0;
	unsigned char seed16[16] = { 0 };
#ifdef  WIN32
	std::random_device rd;
	for (int i = 0; i < 16; i++) {
		seed16[i] = (uint8_t)std::uniform_int_distribution<uint16_t>(0, 255)(rd);
	}
#else
	FILE* frand = fopen("/dev/urandom", "r");
	bool result = (frand == NULL) || fread(&seed16, sizeof(seed16), 1, frand) != sizeof(seed16);
	if (!result) {
		return BPERRORCODE::ERRCODE_RANDOM_ERROR;
	}
	if (frand) {
		fclose(frand);
	}
#endif 
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

namespace secp256k1 {

	Privacy::Privacy(){
		ctx_ = secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
		scratch_ = secp256k1_scratch_space_create(ctx_, 1024 * 1024);
		gens_ = secp256k1_bulletproof_generators_create(ctx_, &secp256k1_generator_const_g, 256);
		secp256k1_context_set_illegal_callback(ctx_, IllegalCallBack, nullptr);
		secp256k1_context_set_error_callback(ctx_, ErrorCallBack, nullptr);
	};

	Privacy::~Privacy(){
		if (ctx_) secp256k1_context_destroy(ctx_);
		if (scratch_) secp256k1_scratch_destroy(scratch_);
		if (gens_) secp256k1_bulletproof_generators_destroy(ctx_, gens_);
	};

	std::string Privacy::GetErrorMsg(int64_t err_code){
		std::string err_msg = "[privacy] illegal argument:";
		if (err_code == BPERRORCODE::ERRCODE_SUCCESS){
			err_msg = "[privacy] info:";
		}
		std::string msg;
		switch (err_code)
		{
		case BPERRORCODE::ERRCODE_SUCCESS:{
			msg = "Success!";
				break; 
		}
		case BPERRORCODE::ERRCODE_RANDOM_ERROR:{
			msg = "Generates a random number error!";
			break;
		}
		case BPERRORCODE::ERRCODE_INVALID_PARAMETER:{
			msg = "Invalid parameters!";
			break;
		}
		case BPERRORCODE::ERRCODE_CREATE_PEDERSEN:{
			msg = "Failed to create pedersen commitment!";
			break;
		}
		case BPERRORCODE::ERRCODE_PARSE_PEDERSEN:{
			msg = "Failed to parse pedersen commitment!";
			break;
		}
		case BPERRORCODE::ERRCODE_SERIALIZE_PEDERSEN:{
			msg = "Failed to serialize pedersen commitment!";
			break;
		}
		case BPERRORCODE::ERRCODE_SERIALIZE_PUBKEY:{
			msg = "Failed to serialize pubkey!";
			break;
		}
		case BPERRORCODE::ERRCODE_VERIFY_TALLY:{
			msg = "Failed to verify tally!";
			break;
		}
		case BPERRORCODE::ERRCODE_CREATE_PUBKEY:{
			msg = "Failed to create pubkey!";
			break;
		}
		case BPERRORCODE::ERRCODE_RANGEPROOF_PROVE:{
			msg = "Failed to generate rangeproof!";
			break;
		}
		case BPERRORCODE::ERRCODE_RANGEPROOF_VERIFY:{
			msg = "Failed to verify rangeproof!";
			break;
		}
		case BPERRORCODE::ERRCODE_OUT_RANGE:{
			msg = "Out of range!";
			break;
		}
		case BPERRORCODE::ERRCODE_ECDSA_CRATE:{
			msg = "Failed to create ecdsa signature!";
			break;
		}
		case BPERRORCODE::ERRCODE_ECDSA_SERIALIZE:{
			msg = "Failed to serialize ecdsa signature!";
			break;
		}
		case BPERRORCODE::ERRCODE_ECDSA_VERIFY:{
			msg = "Failed to verify ecdsa signature!";
			break;
		}
		case BPERRORCODE::ERRCODE_ECDSA_PARSE:{
			msg = "Failed to parse ecdsa signature!";
			break;
		}
		case BPERRORCODE::ERRCODE_PARSE_PUBKEY:{
			msg = "Failed to parse pubkey!";
			break;
		}
		case BPERRORCODE::ERRCODE_BLIND_SUM:{
			msg = "Failed to blind sum!";
			break;
		}
		case BPERRORCODE::ERRCODE_UNKNOWN:{
			msg = "An unknown error!";
			break;
		}
		case BPERRORCODE::ERRCODE_BP_LIB_INTERNAL:{
			msg = "BP library internal error,"+error_msg_;
			break;
		}
		default:
			break;
		}
		err_msg += msg;
		return err_msg;
	}

	//用于创建Pedersen承若
	int64_t Privacy::CreatePedersenCommit(uint64_t value, const std::string& blind, std::string& commit){
		if (blind.size() != 64) {
			return BPERRORCODE::ERRCODE_INVALID_PARAMETER;
		}

		secp256k1_pedersen_commitment pc;
		unsigned char output[33];
		unsigned char blind_data[32];
		HexStrToArray(blind, blind_data);
		try{
			//创建Pedersen承若
			if (!secp256k1_pedersen_commit(ctx_, &pc, blind_data, value, &secp256k1_generator_const_h, &secp256k1_generator_const_g)){
				return BPERRORCODE::ERRCODE_CREATE_PEDERSEN;
			}
			//将一个64字节secp256k1_pedersen_commitment序列化成一个33字节的承若
			secp256k1_pedersen_commitment_serialize(ctx_, output, &pc); // return 1 always

		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}

		commit = ArrayToHexStr(output, 33);
		return BPERRORCODE::ERRCODE_SUCCESS;
	}

	int64_t Privacy::PedersenTallyVerify(const std::vector<std::string>& inputs, const std::vector<std::string>& outputs, const std::string& msg, const std::string& sig){
		if (inputs.size() == 0 || inputs.size() > 100 || outputs.size() == 0 || outputs.size() > 100 || msg.size() != 32) {
			return BPERRORCODE::ERRCODE_INVALID_PARAMETER;
		}

		//secp256k1曲线的一组元素，可比矩阵坐标系
		/** A group element of the secp256k1 curve, in jacobian coordinates. */
		//typedef struct {
		//	secp256k1_fe x; /* actual X: x/z^2 */
		//	secp256k1_fe y; /* actual Y: y/z^3 */
		//	secp256k1_fe z;
		//	int infinity; /* whether this represents the point at infinity */
		//} secp256k1_gej;
		secp256k1_gej accj;
		//secp256k1曲线的一组元素，在仿射坐标下
		/** A group element of the secp256k1 curve, in affine coordinates. */
		//typedef struct {
		//	secp256k1_fe x;
		//	secp256k1_fe y;
		//	int infinity; /* whether this represents the point at infinity */
		//} secp256k1_ge;
		secp256k1_ge add;

		try{
			//注意需要进行初始化，否则会造成不同平台内容不一样
			secp256k1_gej_set_infinity(&accj);
			for (int i = 0; i < outputs.size(); i++){
				if (outputs[i].length() != 66) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;

				unsigned char commit_data[33];
				secp256k1_pedersen_commitment output_commit;
				HexStrToArray(outputs[i], commit_data);
				//将一个33字节解成64字节secp256k1_pedersen_commitment
				if (!secp256k1_pedersen_commitment_parse(ctx_, &output_commit, commit_data)) {
					return BPERRORCODE::ERRCODE_PARSE_PEDERSEN;
				}
				//所有的输出相加
				secp256k1_ge_clear(&add);
				secp256k1_pedersen_commitment_load(&add, &output_commit);
				secp256k1_gej_add_ge_var(&accj, &accj, &add, NULL);
			}

			//将输出取反
			secp256k1_gej_neg(&accj, &accj);

			for (int i = 0; i < inputs.size(); i++){
				if (inputs[i].length() != 66) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;

				unsigned char commit_data[33];
				secp256k1_pedersen_commitment input_commit;
				HexStrToArray(inputs[i], commit_data);

				if (!secp256k1_pedersen_commitment_parse(ctx_, &input_commit, commit_data)) {
					return BPERRORCODE::ERRCODE_PARSE_PEDERSEN;
				}
			    //将所有的输入相加和输出取反后相加
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

			if (!secp256k1_ec_pubkey_serialize(ctx_, pub_compress, &pub_len, &pubkey, SECP256K1_EC_COMPRESSED)){
				return BPERRORCODE::ERRCODE_SERIALIZE_PUBKEY;
			}
			std::string pubkey_str = ArrayToHexStr(pub_compress, 33);

			// create commit from excess
			unsigned char excess_data[33];
			secp256k1_pedersen_commitment excess_commit;

			secp256k1_pedersen_commitment_save(&excess_commit, &add);

			if (!secp256k1_pedersen_commitment_serialize(ctx_, excess_data, &excess_commit)){
				return BPERRORCODE::ERRCODE_SERIALIZE_PEDERSEN;
			}
			std::string excess_str = ArrayToHexStr(excess_data, 33);

			std::vector<std::string> new_outputs;
			new_outputs.assign(outputs.begin(), outputs.end());
			new_outputs.push_back(excess_str);

			// verify input - output = 0
			int64_t ret = TallyVerify(inputs, new_outputs);
			if (ret != BPERRORCODE::ERRCODE_SUCCESS) return ret;
			
			// verify blinds_sum * G = Excess
			return EcdsaVerify(pubkey_str, msg, sig);
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}

	//验证等式平衡
	int64_t Privacy::TallyVerify(const std::vector<std::string>& inputs, const std::vector<std::string>& outputs) {

		int64_t input_size = inputs.size();
		int64_t output_size = outputs.size();
		if ((input_size > 100) || (output_size > 101)){ // the extra one is excess
			return BPERRORCODE::ERRCODE_OUT_RANGE;
		}

		try{
			secp256k1_pedersen_commitment input_list[100];
			secp256k1_pedersen_commitment output_list[101];
			const secp256k1_pedersen_commitment *plist[100];
			const secp256k1_pedersen_commitment *nlist[101];
			for (int i = 0; i < input_size; i++){
				unsigned char commit_data[33];
				HexStrToArray(inputs[i], commit_data);

				if (!secp256k1_pedersen_commitment_parse(ctx_, &input_list[i], commit_data)){
					return BPERRORCODE::ERRCODE_PARSE_PEDERSEN;
				}
				plist[i] = &input_list[i];
			}

			for (int i = 0; i < output_size; i++){
				unsigned char commit_data[33];
				HexStrToArray(outputs[i], commit_data);

				if (!secp256k1_pedersen_commitment_parse(ctx_, &output_list[i], commit_data)){
					return BPERRORCODE::ERRCODE_PARSE_PEDERSEN;
				}
				nlist[i] = &output_list[i];
			}
			//验证等式平衡
			if (!secp256k1_pedersen_verify_tally(ctx_, plist, input_size, nlist, output_size)){
				return BPERRORCODE::ERRCODE_VERIFY_TALLY;
			}
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}

	int64_t Privacy::CreateEcdhKey(const std::string& priv_key, const std::string& pub_key, std::string& ecdh_key){
		if (priv_key.size() != 64 || pub_key.size() != 66) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;

		unsigned char priv[32];
		HexStrToArray(priv_key, priv);

		unsigned char pub[33];
		HexStrToArray(pub_key, pub);
		try{
			secp256k1_pubkey pubkey;

			if (!secp256k1_ec_pubkey_parse(ctx_, &pubkey, pub, 33)) return BPERRORCODE::ERRCODE_PARSE_PEDERSEN;

			unsigned char key[32];
			if (!secp256k1_ecdh(ctx_, key, &pubkey, priv)) return BPERRORCODE::ERRCODE_CREATE_PUBKEY;

			ecdh_key = ArrayToHexStr(key, 32);
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}

	//ECDSA签名
	int64_t Privacy::EcdsaSign(const std::string& priv_key, const std::string& data, std::string& sig){
		if (priv_key.size() != 64 || data.size() != 32) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;

		unsigned char priv[32];
		HexStrToArray(priv_key, priv);

		size_t siglen = 72;
		unsigned char sig_data[72];
		try{
			secp256k1_ecdsa_signature ecdsa_sig;
			if (!secp256k1_ecdsa_sign(ctx_, &ecdsa_sig, (unsigned char*)data.c_str(), priv, NULL, NULL)) return BPERRORCODE::ERRCODE_ECDSA_CRATE;

			if (!secp256k1_ecdsa_signature_serialize_der(ctx_, sig_data, &siglen, &ecdsa_sig)) return BPERRORCODE::ERRCODE_ECDSA_SERIALIZE;

			sig = ArrayToHexStr(sig_data, siglen);
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}

	//验证ECDSA签名
	int64_t Privacy::EcdsaVerify(const std::string& pub_key, const std::string& data, const std::string& sig){
		if (pub_key.size() != 66 || data.size() != 32) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;

		unsigned char pub[33];
		HexStrToArray(pub_key, pub);
		try{
			secp256k1_pubkey pubkey;
			if (!secp256k1_ec_pubkey_parse(ctx_, &pubkey, pub, 33)) return BPERRORCODE::ERRCODE_PARSE_PUBKEY;

			secp256k1_ecdsa_signature ecdsa_sig;
			unsigned char* sig_data = new unsigned char[sig.size() / 2];
			HexStrToArray(sig, sig_data);

			if (!secp256k1_ecdsa_signature_parse_der(ctx_, &ecdsa_sig, sig_data, sig.size() / 2)) return BPERRORCODE::ERRCODE_ECDSA_PARSE;

			if (!secp256k1_ecdsa_verify(ctx_, &ecdsa_sig, (unsigned char*)data.c_str(), &pubkey)) return BPERRORCODE::ERRCODE_ECDSA_VERIFY;
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}

	//创建公私钥对
	int64_t Privacy::CreateKeyPair(std::string& pub_key, std::string& priv_key){
		// set random seed
		unsigned char priv[32];
		try{
			int64_t ret = rand(priv);
			if (ret != BPERRORCODE::ERRCODE_SUCCESS){
				return ret;
			}

			priv_key = ArrayToHexStr(priv, 32);
			return GetPublicKey(priv_key, pub_key);
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}

	//获取公钥
	int64_t Privacy::GetPublicKey(const std::string& priv_key, std::string& pub_key){
		if (priv_key.size() != 64) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;

		unsigned char priv[32];
		HexStrToArray(priv_key, priv);
		try{
			secp256k1_pubkey pub;
			if (!secp256k1_ec_pubkey_create(ctx_, &pub, priv)) return BPERRORCODE::ERRCODE_CREATE_PUBKEY;

			size_t pub_len = 33;
			unsigned char pub_compress[33];
			if (!secp256k1_ec_pubkey_serialize(ctx_, pub_compress, &pub_len, &pub, SECP256K1_EC_COMPRESSED)) return BPERRORCODE::ERRCODE_SERIALIZE_PUBKEY;

			pub_key = ArrayToHexStr(pub_compress, 33);
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}

	//
	int64_t Privacy::ExcessSign(const std::vector<std::string>& inputs, const std::vector<std::string>& outputs, const std::string& msg, std::string& sig){
		unsigned char* blinds[200];
		unsigned char blinds_temp[200][32];
		unsigned char blind_out[32];
		int64_t input_size = inputs.size();
		int64_t output_size = outputs.size();
		if ((input_size > 100) || (output_size > 100)){
			return BPERRORCODE::ERRCODE_OUT_RANGE;
		}

		for (int i = 0; i < input_size; i++){
			if (inputs[i].length() != 64) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;
			HexStrToArray(inputs[i], blinds_temp[i]);
			blinds[i] = blinds_temp[i];
		}

		for (int i = 0; i < output_size; i++){
			if (outputs[i].length() != 64) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;
			HexStrToArray(outputs[i], blinds_temp[i + input_size]);
			blinds[i + input_size] = blinds_temp[i + input_size];
		}

		try{
			//对所有私钥求和
			if (secp256k1_pedersen_blind_sum(ctx_, blind_out, blinds, input_size+output_size, input_size) != 1){
				return BPERRORCODE::ERRCODE_BLIND_SUM;
			}

			std::string priv_key = ArrayToHexStr(blind_out, 32);
			return EcdsaSign(priv_key, msg, sig);	
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}
	//生成单个Rangeproof证明
	int64_t Privacy::BpRangeproofProve(const std::string& blind, uint64_t value, std::string& proof){
		if (blind.size() != 64) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;

		unsigned char* blind_list[1];
		uint64_t value_list[1];

		unsigned char blind_data[32];
		HexStrToArray(blind, blind_data);
		blind_list[0] = blind_data;
		value_list[0] = value;

		size_t plen = 2000;
		unsigned char proof_data[2000];
		try{
			unsigned char random[32];
			int64_t ret = rand(random);
			if (ret != BPERRORCODE::ERRCODE_SUCCESS){
				return ret;
			}

			if (!secp256k1_bulletproof_rangeproof_prove(ctx_, scratch_, gens_, proof_data, &plen, value_list, NULL, blind_list, 1, &secp256k1_generator_const_h, 64, random, NULL, 0)) {
				return BPERRORCODE::ERRCODE_RANGEPROOF_PROVE;
			}

			proof = ArrayToHexStr(proof_data, plen);
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}

		proof = ArrayToHexStr(proof_data, plen);

		return BPERRORCODE::ERRCODE_SUCCESS;
	}
	//验证单个Rangeproof证明
	int64_t Privacy::BpRangeproofVerify(const std::string& commit, const std::string& proof){
		if (commit.size() != 66 || proof.empty()) return BPERRORCODE::ERRCODE_INVALID_PARAMETER;

		unsigned char commit_data[33];
		HexStrToArray(commit, commit_data);

		unsigned char proof_data[2000];
		HexStrToArray(proof, proof_data);

		secp256k1_pedersen_commitment pc;
		try{
			int ret = secp256k1_pedersen_commitment_parse(ctx_, &pc, commit_data);
			if (!ret) {
				return BPERRORCODE::ERRCODE_CREATE_PEDERSEN;
			}

			ret = secp256k1_bulletproof_rangeproof_verify(ctx_, scratch_, gens_, proof_data, proof.size() / 2, NULL, &pc, 1, 64, &secp256k1_generator_const_h, NULL, 0);
			if (!ret) {
				return BPERRORCODE::ERRCODE_RANGEPROOF_VERIFY;
			}
		}
		catch (std::exception& e){
			error_msg_ = e.what();
			return BPERRORCODE::ERRCODE_BP_LIB_INTERNAL;
		}
		return BPERRORCODE::ERRCODE_SUCCESS;
	}
}
