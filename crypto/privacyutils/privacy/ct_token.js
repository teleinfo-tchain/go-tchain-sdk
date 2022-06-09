'use strict';

let globalAttribute = {};
const globalAttributeKey = 'global_attribute';
const maxIdKey = 'max_id';
let max_id = '0';

function increaseMaxId() {
    if (max_id === '0') {
        let json = Chain.getAccountMetadata(Chain.thisAddress, maxIdKey);
        if (json === false) {
            max_id = '0';
        } else {
            max_id = json;
        }
    }
    max_id = Utils.int64Add(max_id, '1');
}

function validate_token(token) {
    Utils.log('token is:' + JSON.stringify(token));
    Utils.assert(typeof token.commit === 'string' && token.commit.length > 0 &&
           typeof token.encrypt_value === 'string' && token.encrypt_value.length > 0 &&
           typeof token.from_pubkey === 'string' && token.from_pubkey.length > 0 &&
		   typeof token.range_proof === 'string' && token.range_proof.length > 0 &&
           typeof token.to === 'string' && token.to.length > 0,
           'Invalid token format');
}

function transfer(params) {
    Utils.assert(typeof params.inputs === 'object' && params.inputs.length > 0 &&
           typeof params.outputs === 'object' && params.outputs.length > 0 &&
           typeof params.excess_sig === 'string' && params.excess_sig.length > 0 &&
           typeof params.excess_msg === 'string' && params.excess_msg.length > 0,
           'Failed to check args');

    // check token existence
    let src_json = Chain.getAccountMetadata(Chain.thisAddress, Chain.msg.sender);
    Utils.assert(src_json !== false, 'Failed to get source account ' + Chain.msg.sender + ' from metadata.');   
    
    let src_account = JSON.parse(src_json);

    let i = 0;
    let j = 0;
    let input_commits = [];
    for (i = 0 ; i < params.inputs.length; i += 1) {
        let exists = false;
        for (j = 0; j < src_account.tokens.length; j += 1) {
            if (src_account.tokens[j].id === params.inputs[i].id) {
                exists = true;
                break;
            }
        }
        Utils.assert(exists !== false, 'No such token:' + params.inputs[i].id);
        input_commits.push(src_account.tokens[j].commit);
        src_account.tokens.splice(j, 1); // consume the input token
    }

    // validate params
    let output_commits = [];
    for (i = 0 ; i < params.outputs.length; i += 1) {
        validate_token(params.outputs[i]);
		
		Utils.assert(Utils.bpRangeProofVerify(params.outputs[i].commit, params.outputs[i].range_proof) !== false, 'Failed to verify range proof');
		delete params.outputs[i].range_proof;
		
        increaseMaxId();
        params.outputs[i].id = max_id;
		params.outputs[i].hash = Chain.tx.hash;

        if (Chain.msg.sender === params.outputs[i].to) {
            delete params.outputs[i].to;
            src_account.tokens.push(params.outputs[i]);
        }
        
        output_commits.push(params.outputs[i].commit);
    }
	
    // verify excess signature
    Utils.assert(Utils.pedersenTallyVerify(JSON.stringify(input_commits), JSON.stringify(output_commits), params.excess_msg, params.excess_sig) !== false, 'Failed to verify excess');
    
    // update source and dest account
    for (i = 0 ; i < params.outputs.length; i += 1) {
        let dest_json = Chain.getAccountMetadata(Chain.thisAddress, params.outputs[i].to);
        //assert(dest_json !== false, 'Failed to get dest account ' + params.outputs[i].to + ' from metadata.');
		let dest_account = {};
		if (dest_json === false) {
		    dest_account.tokens = [];
		} else {
		    dest_account = JSON.parse(dest_json);
		}

        if (params.outputs[i].to !== undefined) {
			// to already delete if src adress equal to dest address
			let to_addr = params.outputs[i].to;
			delete params.outputs[i].to;
			dest_account.tokens.push(params.outputs[i]);
			Chain.store(to_addr, JSON.stringify(dest_account));
		}
    }

    Chain.store(Chain.msg.sender, JSON.stringify(src_account));
    Chain.store(maxIdKey, max_id);
}

function issue(params) {
    let json = Chain.getAccountMetadata(Chain.thisAddress, globalAttributeKey);
    Utils.assert(json === false || json === 'undefined', 'Already issued');


    Utils.assert(typeof params.name === 'string' && params.name.length > 0 &&
           typeof params.symbol === 'string' && params.symbol.length > 0 &&
           typeof params.token.commit === 'string' && params.token.commit.length > 0 &&
           typeof params.token.range_proof === 'string' && params.token.range_proof.length > 0 &&
           typeof params.token.from_pubkey === 'string' && params.token.from_pubkey.length > 0,
           'Failed to check args');

    globalAttribute.name = params.name;
    globalAttribute.symbol = params.symbol;
    globalAttribute.version = 'ETP10';

    Chain.store(globalAttributeKey, JSON.stringify(globalAttribute));

    // create pedersen commitment
    let account_data = {};
    account_data.tokens = [];
    let mwtoken = {};
    mwtoken.commit = params.token.commit;
    mwtoken.from_pubkey = params.token.from_pubkey;
    mwtoken.encrypt_value = params.token.encrypt_value;
    increaseMaxId();
    mwtoken.id = max_id;
	mwtoken.hash = Chain.tx.hash;

    Utils.assert(Utils.bpRangeProofVerify(mwtoken.commit, params.token.range_proof) !== false, 'Failed to verify rangeproof');
    account_data.tokens.push(mwtoken);
    
    Chain.store(Chain.msg.sender, JSON.stringify(account_data));
    Chain.store(maxIdKey, max_id);
}
function rangeproofVerify(inputInfo) {
	Utils.log('bpRangeproofVerifyInfo:' + JSON.stringify(inputInfo));
	let ret = Utils.bpRangeProofVerify(inputInfo.commit, inputInfo.proof);
	Chain.store('bpRangeproofVerify', JSON.stringify(ret));
}

function tallyVerify(inputInfo) {
	Utils.log("pedersenTallyVerify:" + JSON.stringify(inputInfo));
	let ret = Utils.pedersenTallyVerify(JSON.stringify(inputInfo.inputs), JSON.stringify(inputInfo.outputs),inputInfo.excess_msg,inputInfo.excess_sig);
	Chain.store('pedersenTallyVerify', JSON.stringify(ret));
}

function init() {
	return;
}

function main(input_str) {
    let input = JSON.parse(input_str);

    if (input.method === 'transfer') {
        transfer(input.params);
    }
    else if (input.method === 'issue') {
        Utils.log('call issue method');
        issue(input.params);
    }
    else if (input.method === 'tallyVerify') {
        Utils.log('call tall verify method');
        tallyVerify(input.params);
    }
    else if (input.method === 'rangeproofVerify') {
        Utils.log('call rangeproof verify method');
        rangeproofVerify(input.params);
    }
    else {
        throw '<Main interface passes an invalid operation type>';
    }
}
