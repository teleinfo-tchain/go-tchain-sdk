#include <stdlib.h>
#include <stdio.h>

#include <iostream>
#include <sstream>
#include <string>
#include <map>
#include <chrono>

#include "privacy/privacy.hpp" 
#include "aes_tool.hpp"

extern int64_t rand(unsigned char* rand);

static int64_t ITER_NUM = 1000;
static std::string send_pub;
static std::string send_priv;
static std::string recv_pub;
static std::string recv_priv;
static secp256k1::Privacy privacy; 

struct Result {
	int64_t create_time;
	int64_t prove_time;
	int64_t verify_time;
};

void create_account()
{
	int ret = 0;
	do {
		ret = privacy.CreateKeyPair(send_pub, send_priv);
		if(ret != 0) break; 

		ret = privacy.CreateKeyPair(recv_pub, recv_priv);
		if(ret != 0) break; 

	}while(false);

	if(ret!=0) {
		std::cerr << "Failed to create sender and receiver account" << std::endl;
	}
}

void function_test()
{
	int ret = 0;
	//secp256k1::Privacy privacy; 
	// function test
	do {
		std::cout << "Function Test" << std::endl;
		/*std::cout << "Generate Keypair ... " << std::endl;
		
		//std::map<std::string, std::string> account_map;
		std::string send_pub;
		std::string send_priv;
		//for(int i = 0 ; i < USER_NUM; i++) {
		ret = privacy.CreateKeyPair(send_pub, send_priv);
		if(ret != 0) break; 

		std::cout << "Sender's pub:" << send_pub << ", Sender's priv:" << send_priv << std::endl;
		//	account_map.insert(std::make_pair(pub, priv));
		//}
		
		std::string recv_pub;
		std::string recv_priv;
		//for(int i = 0 ; i < USER_NUM; i++) {
		ret = privacy.CreateKeyPair(recv_pub, recv_priv);
		if(ret != 0) break; 

		std::cout << "Receiver's pub:" << recv_pub << ", Receiver's priv:" << recv_priv << std::endl;
		*/

		std::cout << "Create ecdh key" << std::endl;
		std::string ecdh_key;
		ret = privacy.CreateEcdhKey(send_priv, send_pub, ecdh_key);
		if(ret != 0) break; 
		std::cout << "Ecdh key:" << ecdh_key << std::endl;
		
		std::cout << "Create Pedersen Commit" << std::endl;
		std::string input_commit;
		uint64_t input_value = 10000;
		
		std::string input_blind = ecdh_key;
		ret = privacy.CreatePedersenCommit(input_value, input_blind, input_commit);
		if(ret != 0) break; 
		std::cout << "Input commit:" << input_commit << std::endl;

		uint64_t output_value = 5000;
		uint64_t change_value = input_value - output_value;
		std::string output_commit;
		std::string change_commit;
		
		std::string output_blind;
		std::string change_blind;
		ret = privacy.CreateEcdhKey(send_priv, recv_pub, output_blind);
		if(ret != 0) break; 
		ret = privacy.CreateEcdhKey(send_priv, send_pub, change_blind);
		if(ret != 0) break; 
		ret = privacy.CreatePedersenCommit(output_value, output_blind, output_commit);
		if(ret != 0) break; 
		std::cout << "output_commit:" << output_commit << std::endl;
		ret = privacy.CreatePedersenCommit(change_value, change_blind, change_commit);
		if(ret != 0) break; 
		std::cout << "change_commit:" << change_commit << std::endl;

		// Range proof
		std::string output_proof;
		std::string change_proof;
		ret = privacy.BpRangeproofProve(output_blind, output_value, output_proof);
		if(ret != 0) break;
		std::cout << "output_proof:" << output_proof << std::endl;
		ret = privacy.BpRangeproofProve(change_blind, change_value, change_proof);
		if(ret != 0) break;
		std::cout << "change_proof:" << output_proof << std::endl;

		// Excess sign
		std::string msg = send_pub.substr(0, 32);
		std::string sig;
		std::vector<std::string> input_blinds;
		std::vector<std::string> output_blinds;
		input_blinds.push_back(input_blind);
		output_blinds.push_back(output_blind);
		output_blinds.push_back(change_blind);
		ret = privacy.ExcessSign(input_blinds, output_blinds, msg, sig);
		if(ret != 0) break;
		std::cout << "Excess_msg:" << msg << std::endl;
		std::cout << "Excess_sig:" << sig << std::endl;

		// verify
		ret = privacy.BpRangeproofVerify(output_commit, output_proof);
		if(ret != 0) break;
		ret = privacy.BpRangeproofVerify(change_commit, change_proof);
		if(ret != 0) break;
		std::cout << "Verify range proof done" << std::endl;

		std::vector<std::string> input_commits;
		std::vector<std::string> output_commits;
		input_commits.push_back(input_commit);
		output_commits.push_back(output_commit);
		output_commits.push_back(change_commit);

		ret = privacy.PedersenTallyVerify(input_commits, output_commits, msg, sig);
		if(ret != 0) break;
		std::cout << "Pedersen tally verify done" << std::endl;

		//public combine
		std::vector<std::string>  input_pubs;
		std::string combine_pubs;

		input_pubs.push_back("03435091d48b13056a3a1c63fec9909eaaf7c290d4179cb7a1362a653b2d1cbce6");
		input_pubs.push_back("03435091d48b13056a3a1c63fec9909eaaf7c290d4179cb7a1362a653b2d1cbce6");
		ret = privacy.PublicKeyCombine(input_pubs, combine_pubs);
		if(ret != 0 || combine_pubs != "0320b0be5eb417e0d227f43285f0a98fbca0a67f24985beff68e8a29b477195b44")
		{
			std::cout << "Failed to combine public keys" << std::endl;
		}
		
		std::cout << "Combine public keys done" << std::endl;
		
	} while(false);

	if(ret != 0) {
		std::cout << privacy.GetErrorMsg(ret) << std::endl;
	}

}

int tx_test(Result& res)
{
	int ret = 0;
	do {
		auto t1 = std::chrono::high_resolution_clock::now();
		
		// 2 input commit
		std::string ecdh_key;
		ret = privacy.CreateEcdhKey(send_priv, send_pub, ecdh_key);
		if(ret != 0) break; 
		
		uint64_t input_value1 = 10000;
		uint64_t input_value2 = 20000;
		std::string input_commit1;
		std::string input_commit2;
		
		ret = privacy.CreatePedersenCommit(input_value1, ecdh_key, input_commit1); // use r1*R1 as blind factor
		if(ret != 0) break; 
		ret = privacy.CreatePedersenCommit(input_value2, ecdh_key, input_commit2);
		if(ret != 0) break; 

		// 2 output commit
		uint64_t output_value = 25000;
		uint64_t change_value = input_value1 + input_value2 - output_value;
		std::string output_commit;
		std::string change_commit;
		
		std::string output_blind;
		std::string change_blind;
		ret = privacy.CreateEcdhKey(send_priv, recv_pub, output_blind);
		if(ret != 0) break; 
		ret = privacy.CreateEcdhKey(send_priv, send_pub, change_blind);
		if(ret != 0) break; 
		ret = privacy.CreatePedersenCommit(output_value, output_blind, output_commit);
		if(ret != 0) break; 
		ret = privacy.CreatePedersenCommit(change_value, change_blind, change_commit);
		if(ret != 0) break; 
		
		// encryt values
		std::string output_enc = Aes::HexDecrypto(std::to_string(output_value), output_blind.substr(0, 32));
		std::string change_enc = Aes::HexDecrypto(std::to_string(change_value), change_blind.substr(0, 32));
		
		auto t2 = std::chrono::high_resolution_clock::now(); 
		res.create_time = std::chrono::duration_cast<std::chrono::nanoseconds>(t2-t1).count();

		// Range proof
		std::string output_proof;
		std::string change_proof;
		ret = privacy.BpRangeproofProve(output_blind, output_value, output_proof);
		if(ret != 0) break;
		ret = privacy.BpRangeproofProve(change_blind, change_value, change_proof);
		if(ret != 0) break;

		// Excess sign
		std::string msg = send_pub.substr(0, 32);
		std::string sig;
		std::vector<std::string> input_blinds;
		std::vector<std::string> output_blinds;
		input_blinds.push_back(ecdh_key);
		input_blinds.push_back(ecdh_key);
		output_blinds.push_back(output_blind);
		output_blinds.push_back(change_blind);
		ret = privacy.ExcessSign(input_blinds, output_blinds, msg, sig);
		if(ret != 0) break;

		auto t3 = std::chrono::high_resolution_clock::now(); 
		res.prove_time = std::chrono::duration_cast<std::chrono::nanoseconds>(t3-t2).count();

		// verify
		ret = privacy.BpRangeproofVerify(output_commit, output_proof);
		if(ret != 0) break;
		ret = privacy.BpRangeproofVerify(change_commit, change_proof);
		if(ret != 0) break;

		std::vector<std::string> input_commits;
		std::vector<std::string> output_commits;
		input_commits.push_back(input_commit1);
		input_commits.push_back(input_commit2);
		output_commits.push_back(output_commit);
		output_commits.push_back(change_commit);

		ret = privacy.PedersenTallyVerify(input_commits, output_commits, msg, sig);
		if(ret != 0) break;
		auto t4 = std::chrono::high_resolution_clock::now();
		res.verify_time = std::chrono::duration_cast<std::chrono::nanoseconds>(t4-t3).count();
	} while(false);

	if(ret != 0) {
		std::cout << privacy.GetErrorMsg(ret) << std::endl;
	}
	return ret;
}

void performance_test()
{
	int ret = 0;
	std::vector<Result> res_vec;

	do{
		std::cout << "Performance test" << std::endl;
		for (int i = 0; i < ITER_NUM; i++) {
			Result res;
			res.create_time = 0;
			res.prove_time = 0;
			res.verify_time = 0;
			ret = tx_test(res);
			res_vec.push_back(res);
		}
		
		uint64_t total_prove_time = 0;
		uint64_t total_build_time = 0;
		uint64_t total_verify_time = 0;
		uint64_t max_prove_time = 0;
		uint64_t min_prove_time = -1;
		uint64_t max_build_time = 0;
		uint64_t min_build_time = -1;
		uint64_t max_verify_time = 0;
		uint64_t min_verify_time = -1;
		for(int i = 0; i < ITER_NUM; i++) {
			uint64_t build_time = res_vec[i].create_time + res_vec[i].prove_time;
			total_prove_time += res_vec[i].prove_time;
			total_build_time += build_time;
			total_verify_time += res_vec[i].verify_time;
			max_prove_time = (res_vec[i].prove_time > max_prove_time) ? res_vec[i].prove_time : max_prove_time;
			min_prove_time = (res_vec[i].prove_time < min_prove_time) ? res_vec[i].prove_time : min_prove_time;
			max_build_time = (build_time > max_build_time) ? build_time : max_build_time;
			min_build_time = (build_time < min_build_time) ? build_time : min_build_time;
			max_verify_time = (res_vec[i].verify_time > max_verify_time) ? res_vec[i].verify_time : max_verify_time;
			min_verify_time = (res_vec[i].verify_time < min_verify_time) ? res_vec[i].verify_time : min_verify_time;
		}
		
		std::cout << "Performance test result:" << std::endl;
		std::cout << "Total prove time:" << total_prove_time/1000000 << " ms, avarage: " << total_prove_time/ITER_NUM/1000000 << " ms" << std::endl;
		std::cout << "Total build time:" << total_build_time/1000000 << " ms, avarage: " << total_build_time/ITER_NUM/1000000 << " ms" << std::endl;
		std::cout << "Total verify time:" << total_verify_time/1000000 << " ms, avarage: " << total_verify_time/ITER_NUM/1000000 << " ms" << std::endl;
		std::cout << "Max prove time:" << max_prove_time/1000000 << " ms" << std::endl;
		std::cout << "Min prove time:" << min_prove_time/1000000 << " ms" << std::endl;
		std::cout << "Max build time:" << max_build_time/1000000 << " ms" << std::endl;
		std::cout << "Min build time:" << min_build_time/1000000 << " ms" << std::endl;
		std::cout << "Max verify time:" << max_verify_time/1000000 << " ms" << std::endl;
		std::cout << "Min verify time:" << min_verify_time/1000000 << " ms" << std::endl;


	} while(false);

	if(ret != 0) {
		std::cout << privacy.GetErrorMsg(ret) << std::endl;
	}
}

void stable_test() {
	int ret = 0;
	secp256k1::Privacy privacy; 

	do{
		std::cout << "stable test" << std::endl;
	} while(false);

	if(ret != 0) {
		std::cout << privacy.GetErrorMsg(ret) << std::endl;
	}

}

void usage()
{
	printf(
		"Usage:bubitest [OPTIONS]\n"
		"OPTIONS:\n"
		"  -f|--function-test			function test\n"
		"  -p|--perf_test				performance test\n"
		"		-i						total tx numbers to test\n"
	);
}
int main(int argc, char* argv[])
{
	if(argc < 2) {
		usage();
		exit(1);
	}
	std::string cmd = argv[1];

	
	create_account();

	if(cmd == "--function-test") {
		function_test();
	} else if(cmd == "--perf-test") {
		if(argc == 3) {
			int64_t num = 0;
			std::istringstream iss(argv[2]);
			iss >> num;
			ITER_NUM = num;
		}
		performance_test();
	} else if(cmd == "--stable test") {
		stable_test();
	} else if(cmd == "-h" || cmd == "--help") {
		usage();
		exit(1);
	}

	return 0;
}
