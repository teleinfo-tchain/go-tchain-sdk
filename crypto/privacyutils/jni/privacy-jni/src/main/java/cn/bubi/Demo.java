package cn.bubi;

import cn.bubi.common.Tools;
import cn.bubi.encryption.utils.hex.HexFormat;
import cn.bubi.exception.SDKException;
import cn.bubi.model.request.*;
import cn.bubi.model.response.AccountGetMetadataResponse;
import cn.bubi.model.response.TransactionSubmitResponse;
import cn.bubi.model.response.result.AccountGetMetadataResult;
import cn.bubi.model.response.result.data.TransactionHistory;
import cn.bubi.protobuf.PrivacyProtobuf;
import cn.bubi.crypto.Keypair;
import cn.bubi.utils.*;
import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.googlecode.protobuf.format.JsonFormat;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;

import java.util.*;

/**
 * @author kuangkai@bubi.cn
 * @since 19/9/29 下午3:34.
 */
public class Demo{

    //  genesis account
    private static String genesis_address = "adxSiukTQWSCnFPoTrYLke5aXhdVp2wDxosFp";
    private static String genesis_privateKey = "privbzmruHFUHWD5T7ZHCKUbzcDXER7jihgDSdzeMEzdPY7WrK8tHrZm";
    private static String contract_address;
    private static Keypair from_keypair = Keypair.generator();
    private static Keypair to_keypair = Keypair.generator();
    private static Keypair to_keypair1 = Keypair.generator();
    private static CtEntity ct_from_keypair = CocoonUtil.buildCtEntity();
    private static CtEntity ct_to_keypair = CocoonUtil.buildCtEntity();
    private static CtEntity ct_from_keypair1 = CocoonUtil.buildCtEntity();
    private static CtEntity ct_to_keypair1 = CocoonUtil.buildCtEntity();
    private static CtEntity ct_to_keypair2 = CocoonUtil.buildCtEntity();
    private static BlockChainManager bm = new BlockChainManager();

    private static String contract_file = "D:\\bubichain-v4\\src\\secp256k1-bulletproof\\privacy\\ct_token.js";
    private static SDK bubiSdk;

    public static JSONObject buildToken(String from_pubkey, String to, String encrypt_value, String rangeproof, String commit) {
        JSONObject token = new JSONObject();
        token.put("from_pubkey", from_pubkey);
        token.put("to", to);
        token.put("encrypt_value", encrypt_value);
        if (rangeproof != null) {
            token.put("range_proof", rangeproof);
        }
        token.put("commit", commit);
        return token;
    }

    public static JSONObject buildTx(String[] inputs_ids, JSONArray outputs, String msg, String sig) {
        JSONObject tx = new JSONObject();

        JSONArray ips = new JSONArray();
        for (String id: inputs_ids) {
            JSONObject token = new JSONObject();
            token.put("id", id);
            ips.add(token);
        }
        tx.put("inputs", ips);

        JSONArray ops = new JSONArray();
        for (Object token: outputs) {
            ops.add(token);
        }
        tx.put("outputs", ops);
        tx.put("excess_sig", sig);
        tx.put("excess_msg", msg);
        return tx;
    }

    public static Map getTokenByValue(String address, String[] privs, long value, String[] ids) throws Exception {
        AccountGetMetadataRequest request = new AccountGetMetadataRequest();
        request.setAddress(contract_address);
        request.setKey(address);
        AccountGetMetadataResponse response = bubiSdk.getAccountService().getMetadata(request);
        AccountGetMetadataResult result = null;
        if(response.getErrorCode() == 0) {
            result = response.getResult();
        } else {
            System.out.println("error: " + response.getErrorDesc());
        }

        JSONObject dataJson = JSONObject.parseObject(result.getMetadatas()[0].getValue());

        Long total_balance = 0L;
        Map<String, JSONObject> map = new HashMap<String, JSONObject>();
        for (Object token: dataJson.getJSONArray("tokens")) {
            JSONObject t = (JSONObject) token;
            if (ids.length > 0) {
                boolean got = false;
                for (String id : ids) {
                    if (t.getString("id").equals(id)) {
                        got = true;
                        break;
                    }
                }
                if(!got) continue;
            }

            Long balance = 0L;
            boolean decrypt_done = false;
            for(String priv: privs) {
                String decrypt_key = CocoonUtil.buildEcdhKey(priv, t.getString("from_pubkey"));
                String encrypt_value = t.getString("encrypt_value");
                try {
                    balance = Long.parseLong(AesCbc.decrypt(HexFormat.hexToByte(encrypt_value), decrypt_key.substring(0, 32).getBytes()));
                    t.put("spend_key", priv);
                    t.put("value", balance);
                    decrypt_done = true;
                    break;
                }catch (Exception e) {
                    continue;
                }
            }
            if(!decrypt_done) {
                throw new UnsatisfiedLinkError("Failed to decrypt token of id: " + t.getString("id"));
            };

            total_balance += balance;
            map.put(t.getString("id"), t);
        }

        if (total_balance < value) {
            //throw new UnsatisfiedLinkError("Balance not enough, want:" + Long.toString(value) + ", got:" + Long.toString(total_balance));
            System.out.println("Balance not enough, want:" + Long.toString(value) + ", got:" + Long.toString(total_balance));
        }

        if (ids.length > 0) return map;

        Long tmp_value = 0L;
        Map<String, JSONObject> new_map = new HashMap<String, JSONObject>();
        for (String key: map.keySet()) {
            JSONObject t = map.get(key);
            tmp_value += t.getLong("value");
            new_map.put(key, map.get(key));
            if(tmp_value >= value) break;
        }
        return new_map;
    }

    // bubi address mapping only one confidential transcation account
    public static JSONObject createTx(String from, String to, CtEntity ct_from, String ct_to, long value, String[] ids) throws Exception {

        // get input token from contract
        String[] priv_list = new String[1];
        priv_list[0] = ct_from.getCtPrivateKey();

        Map input_tokens = getTokenByValue(from, priv_list, value, ids);

        String[] input_blinds = new String[input_tokens.size()];
        String[] inputs = new String[input_tokens.size()];
        String[] input_ids = new String[input_tokens.size()];

        int i = 0;
        Long total_balance = 0L;
        for (Object key: input_tokens.keySet()) {
            JSONObject token = (JSONObject) input_tokens.get(key);
            total_balance += token.getLong("value");
            String from_pubkey = token.getString("from_pubkey");
            input_blinds[i] = CocoonUtil.buildEcdhKey(token.getString("spend_key"), from_pubkey);
            inputs[i] = token.getString("commit");
            input_ids[i] = key.toString();
            i += 1;
        }

        // build the spend token
        PrivacyProtobuf.ConfidentialAsset.Builder toAsset = CocoonUtil.buildAsset(to, value, ct_from.getCtPrivateKey(), ct_to);
        String[] outputs = null;
        String[] output_blinds = null;
        JSONArray output_tokens = new JSONArray();
        JsonFormat jsonFormat = new JsonFormat();
        String toAssetStr = jsonFormat.printToString(toAsset.build());
        output_tokens.add(JSONObject.parseObject(toAssetStr));

        // build the change token
        if (total_balance - value > 0) {
            PrivacyProtobuf.ConfidentialAsset.Builder changeAsset = CocoonUtil.buildAsset(from, total_balance - value, ct_from.getCtPrivateKey(), ct_from.getCtPublicKey());
            output_blinds = new String[2];
            outputs = new String[2];
            jsonFormat = new JsonFormat();
            outputs[1] = changeAsset.getCommit();
            output_blinds[1] = changeAsset.getId(); // store blinds in parameter id
            String changeAssetStr = jsonFormat.printToString(changeAsset.build());
            output_tokens.add(JSONObject.parseObject(changeAssetStr));
        } else {
            output_blinds = new String[1];
            outputs = new String[1];
        }
        outputs[0] = toAsset.getCommit();
        output_blinds[0] = toAsset.getId(); // store blind factor in parameter id;

        // build the excess signature
        String data = "64e4a2fcf36693904a0d549303c6f35c";
        String excess_sig = CocoonUtil.excessSign(input_blinds, output_blinds, data);

        // verify tx data
        int res = CocoonUtil.pedersenTallyVerify(inputs, outputs, data, excess_sig);
        if(res != 0) {
            System.out.println("Failed to verify tally");
        } else {
            System.out.println("Verify tally done\n");
        }

        // build tx
        JSONObject txJson = buildTx(input_ids, output_tokens, data, excess_sig);

        JSONObject input = new JSONObject();
        input.put("method", "transfer");
        input.put("params", txJson);
        System.out.println("The contrat input is:");
        System.out.println(JSON.toJSONString(txJson, true));
        return input;
    }

    // bubi address mapping two or more confidential address
    // ct_from_keypairs: {public_key1: private_key1, public_key2: private_key2}, transfer_list: {1:[to1, public_key1, value1], 2:[to2, public_key2, value2]}
//    public static JSONObject createTx1(String from, Map<String, String> ct_from_keypairs, Map<Integer, String[]> transfer_list, String[] ids) throws Exception {
//
//        // get input token from contract
//        String[] privs = new String[ct_from_keypairs.size()];
//        String change_pub = null;
//        int i = 0;
//        for (String key: ct_from_keypairs.keySet()) {
//            if(i==0) change_pub = key;
//            privs[i++] = ct_from_keypairs.get(key);
//        }
//
//        long total_value = 0;
//        for(int key: transfer_list.keySet()) {
//            String[] item = transfer_list.get(key);
//            total_value += Long.parseLong(item[2]);
//        }
//        Map input_tokens = getTokenByValue(from, privs, total_value, ids);
//
//        String[] input_blinds = new String[input_tokens.size()];
//        String[] inputs = new String[input_tokens.size()];
//        String[] input_ids = new String[input_tokens.size()];
//
//        i = 0;
//        Long total_balance = 0L;
//        for (Object key: input_tokens.keySet()) {
//            JSONObject token = (JSONObject) input_tokens.get(key);
//            total_balance += token.getLong("value");
//            String from_pubkey = token.getString("from_pubkey");
//            String priv = token.getString("spend_key");
//
//            String[] blind_ret = cn.bubi.Privacy.createEcdhKey(cpp_obj, priv, from_pubkey);
//            String errCodeBlindRet = blind_ret[0];
//            String errDescBlindRet = blind_ret[1];
//            if (Integer.valueOf(errCodeBlindRet) != 0) {
//                System.out.println("Failed to create blind: " + errDescBlindRet);
//                i += 1;
//                continue;
//            }
//            String blind = blind_ret[2];
//
//            input_blinds[i] = blind;
//            inputs[i] = token.getString("commit");
//            input_ids[i] = key.toString();
//            i += 1;
//        }
//
//        // build the spend token
//        JSONArray output_tokens = null;
//        if(total_balance - total_value > 0) {
//            String[] item = new String[3];
//            item[0] = from;
//            item[1] = change_pub;
//            item[2] = Long.toString(total_balance - total_value);
//            transfer_list.put(transfer_list.size() + 1, item);
//        }
//
//        output_tokens = new JSONArray(transfer_list.size());
//        String[] output_blinds = new String[transfer_list.size()];
//        String[] output_values = new String[transfer_list.size()];
//        String[] output_commits = new String[transfer_list.size()];
//
//        i = 0;
//        for (int key: transfer_list.keySet()) {
//            String[] item = transfer_list.get(key);
//            String ct_to = item[1];
//            String addr = item[0];
//            long value = Long.parseLong(item[2]);
//
//            // get random keypair for blind
//            CtEntity random_keypair = buildCtEntity();
//
//            // blind factor can create from random private key or the private key of from account
//            String ecdh_key = CocoonUtil.buildEcdhKey(random_keypair.getCtPrivateKey(), ct_to);
//            String commit = CocoonUtil.buildPedersenCommitment(value, random_keypair.getCtPrivateKey(), ct_to);
//
//            String encrypt_data = HexFormat.byteToHex(AesCbc.encrypt(Long.toString(value).getBytes(), ecdh_key.substring(0, 32).getBytes())).toLowerCase();
//
//            String[] proof_ret = CocoonUtil.bpRangeproofProve(cpp_obj, ecdh_key, value);
//            String errCodePrfRet = proof_ret[0];
//            String errDescPrfRet = proof_ret[1];
//            String proof = proof_ret[2];
//            if (Integer.valueOf(errCodePrfRet) != 0) {
//                System.out.println("Error proof change: " + errDescPrfRet);
//            } else {
//                String[] bp_very_ret1 = cn.bubi.Privacy.bpRangeproofVerify(cpp_obj, commit, proof);
//                String errCodeBpRet1 = bp_very_ret1[0];
//                String errDescBpRet1 = bp_very_ret1[1];
//                if (Integer.valueOf(errCodeBpRet1) != 0) {
//                    System.out.println("Failed to verify range proof: " + errDescPrfRet + "\n");
//                }
//            }
//
//            JSONObject token = buildToken(random_pub, addr, encrypt_data, proof, commit);
//            output_tokens.add(token);
//            output_values[i] = Long.toString(value);
//            output_blinds[i] = ecdh_key;
//            output_commits[i] = commit;
//            i += 1;
//        }
//
//        // create excess
//        String data = "64e4a2fcf36693904a0d549303c6f35c";
//        String[] excess_sig_ret = cn.bubi.Privacy.excessSign(cpp_obj, input_blinds, output_blinds, data);
//        String errCodeExcSigRet = excess_sig_ret[0];
//        String errDescExcSigRet = excess_sig_ret[1];
//        String excess_sig = excess_sig_ret[2];
//        if (Integer.valueOf(errCodeExcSigRet) != 0) {
//            System.out.println("Error excess sig: " + errDescExcSigRet);
//        }
//
//        // verify tx data
//        String[] tally_very_ret = cn.bubi.Privacy.pedersenTallyVerify(cpp_obj, inputs, output_commits, data, excess_sig);
//        String errCodeTalVeryRet = tally_very_ret[0];
//        String errDescTalVeryRet = tally_very_ret[1];
//        if(Integer.valueOf(errCodeTalVeryRet) != 0) {
//            throw new UnsatisfiedLinkError("Failed to verify tally, " + errDescTalVeryRet);
//        } else {
//            System.out.println("Verify tally done\n");
//        }
//
//        // build tx
//        JSONObject txJson = buildTx(input_ids, output_tokens, data, excess_sig);
//
//        JSONObject input = new JSONObject();
//        input.put("method", "transfer");
//        input.put("params", txJson);
//        System.out.println("The contrat input is:");
//        System.out.println(JSON.toJSONString(txJson, true));
//        return input;
//    }

    public static String getMetadata(String address, String key) throws Exception {
        AccountGetMetadataRequest request = new AccountGetMetadataRequest();
        request.setAddress(address);
        request.setKey(key);
        AccountGetMetadataResponse response = bubiSdk.getAccountService().getMetadata(request);

        AccountGetMetadataResult result = null;
        if(response.getErrorCode() == 0) {
            result = response.getResult();
            return result.getMetadatas()[0].getValue();
        } else {
            System.out.println("error: " + response.getErrorDesc());
        }
        return null;
    }
    public static void getBalance(String address, String[] privs) throws Exception {
        String result = getMetadata(contract_address, address);
        if (result == null) {
            return;
        }
        JSONObject dataJson = JSONObject.parseObject(result);

        for (Object token: dataJson.getJSONArray("tokens")) {
            JSONObject t = (JSONObject) token;

            long balance = 0;
            boolean decrypt_done = false;
            for(String priv: privs) {
                String dec_key = CocoonUtil.buildEcdhKey(priv, t.getString("from_pubkey"));
                String encrypt_value = t.getString("encrypt_value");
                try {
                    balance = Long.parseLong(AesCbc.decrypt(HexFormat.hexToByte(encrypt_value), dec_key.substring(0, 32).getBytes()));
                    decrypt_done = true;
                    break;
                }catch (Exception e) {
                    continue;
                }
            }
            if(!decrypt_done) {
                throw new UnsatisfiedLinkError("Failed to decrypt token of id: " + t.getString("id"));
            };

            System.out.println("account:" + address + ", id: " + t.getString("id") + ", value:" + balance);
        }
        return;
    }

    public static void main(String[] args) throws Exception {

//        String eventUtis = "ws://127.0.0.1:7053";
//        String ips = "127.0.0.1:19333";
//        String eventUtis = "ws://192.168.3.65:7153";
//        String ips = "192.168.3.65:29333";

        SDKConfigure sdkCfg = new SDKConfigure();
        sdkCfg.setHttpConnectTimeOut(5000);
        sdkCfg.setHttpReadTimeOut(5000);
        sdkCfg.setUrl("http://192.168.10.100:29333");
        bubiSdk = SDK.getInstance(sdkCfg);

        createAccounts();

        // create contract and transaction account
        createContract();

        Thread.sleep(2000);

        // issue
        issueCtAsset((long) 100000000);
        Thread.sleep(2000);

        // query issue asset
        String[] privs = new String[1];
        privs[0] = ct_from_keypair.getCtPrivateKey();
        getBalance(from_keypair.getAddress(), privs);

        // split
        splitCtAsset((long) 50000000);
        Thread.sleep(2000);

        getBalance(from_keypair.getAddress(), privs);

        // transfer asset
        transferCtAsset((long) 10000000);
        Thread.sleep(2000);

        getBalance(from_keypair.getAddress(), privs);
        privs[0] = ct_to_keypair.getCtPrivateKey();
        getBalance(to_keypair.getAddress(), privs);

        // this section test the input and output build by two or more confidential transaction account(ct),
        // which means one bubi address mapping two or more ct account,
        // the output can also be set to more than one dest bubi address.

//        // add a new ct account
//        createMultSplit(operationService, (long) 20000000);
//        getBalance(to_keypair.getBubiAddress(), privs);
//        privs = new String[2];
//        privs[0] = ct_from_keypair[1];
//        privs[1] = ct_from_keypair1[1];
//        getBalance(from_keypair.getBubiAddress(), privs);
//
//        // transfer from two ct account to two ct account
//        createMultTransfer(operationService, (long) 10000000, (long) 10000000, (long) 10000000);
//        getBalance(from_keypair.getBubiAddress(), privs);
//
//        privs[0] = ct_to_keypair[1];
//        privs[1] = ct_to_keypair1[1];
//        getBalance(to_keypair.getBubiAddress(), privs);
//
//        privs[0] = ct_to_keypair2[1];
//        getBalance(to_keypair1.getBubiAddress(), privs);
    }

    private static void splitCtAsset(Long value) throws Exception {
        String[] ids = {};

        JSONObject input = createTx(from_keypair.getAddress(), from_keypair.getAddress(), ct_from_keypair, ct_from_keypair.getCtPublicKey(), value, ids);
        try {
            BcTransferReq bcTransfer = new BcTransferReq();
            bcTransfer.setFromAddress(from_keypair.getAddress());
            bcTransfer.setToAddress(contract_address);
            bcTransfer.setTxType(BcTxTypeEnum.CONTRACT.getCode());
            bcTransfer.setAmount(value);
            bcTransfer.setInput(input.toJSONString());
            BlobDataResp blobDataResp = bm.buildBlobData(bcTransfer, from_keypair.getAddress());
            //签名
            List<String> privateKeys = new ArrayList<>();
            privateKeys.add(from_keypair.getPrivateKey());
            List<SignEntity> signBlobs = bm.signBlob(blobDataResp,privateKeys);
            //提交
            SubmitTxReq submitTxReq = new SubmitTxReq();
            submitTxReq.setBlob(blobDataResp.getBlob());
            submitTxReq.setListSigner(signBlobs);
            submitTxReq.setHash(blobDataResp.getHash());
            submitTxReq.setInitiator(from_keypair.getAddress());
            TransactionSubmitResponse bcResponse = bm.submitTx(submitTxReq);
            if(bcResponse.getErrorCode().equals(0)){
                //System.out.println(bcResponse.getResult()); // print issue output
            }else{
                System.out.println("提交交易到blockchain异常" + bcResponse.getErrorCode() + " hash:" + submitTxReq.getHash());
                throw new RuntimeException(bcResponse.getErrorCode() +"-" +bcResponse.getErrorDesc());//抛出异常
            }

        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static void transferCtAsset(Long value) throws Exception {
        String[] ids = {};
        JSONObject input = createTx(from_keypair.getAddress(), to_keypair.getAddress(), ct_from_keypair, ct_to_keypair.getCtPublicKey(), value, ids);
        try {
            BcTransferReq bcTransfer = new BcTransferReq();
            bcTransfer.setFromAddress(from_keypair.getAddress());
            bcTransfer.setToAddress(contract_address);
            bcTransfer.setTxType(BcTxTypeEnum.CONTRACT.getCode());
            bcTransfer.setAmount(value);
            bcTransfer.setInput(input.toJSONString());
            BlobDataResp blobDataResp = bm.buildBlobData(bcTransfer, from_keypair.getAddress());
            //签名
            List<String> privateKeys = new ArrayList<>();
            privateKeys.add(from_keypair.getPrivateKey());
            List<SignEntity> signBlobs = bm.signBlob(blobDataResp,privateKeys);
            //提交
            SubmitTxReq submitTxReq = new SubmitTxReq();
            submitTxReq.setBlob(blobDataResp.getBlob());
            submitTxReq.setListSigner(signBlobs);
            submitTxReq.setHash(blobDataResp.getHash());
            submitTxReq.setInitiator(from_keypair.getAddress());
            TransactionSubmitResponse bcResponse = bm.submitTx(submitTxReq);
            if(bcResponse.getErrorCode().equals(0)){
                //System.out.println(bcResponse.getResult()); // print issue output
            }else{
                System.out.println("提交交易到blockchain异常" + bcResponse.getErrorCode() + " hash:" + submitTxReq.getHash());
                throw new RuntimeException(bcResponse.getErrorCode() +"-" +bcResponse.getErrorDesc());//抛出异常
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
//
//    private static void createMultSplit(BcOperationService operationService, Long value) throws Exception {
//        String[] ids = {};
//        JSONObject input = createTx(from_keypair.getBubiAddress(), from_keypair.getBubiAddress(), ct_from_keypair, ct_from_keypair1[0], value, ids);
//        try {
//            Transaction transaction = operationService.newTransaction(from_keypair.getBubiAddress());
//            InvokeContractOperation invokeContractOperation = new InvokeContractOperation.Builder()
//                    .buildDestAddress(contract_address)
//                    .buildInputData(input.toJSONString())
//                    .build();
//
//            TransactionCommittedResult result = transaction.buildAddOperation(invokeContractOperation)
//                    .buildGasPrice(1000L).buildFeeLimit(100000000L)
//                    .buildTxMetadata("交易")
//                    .buildAddSigner(from_keypair.getPubKey(), from_keypair.getPriKey())
//                    .commit();
//
//            System.out.println("\n------------------------------------------------");
//            System.out.println(GsonUtil.toJson(result));
//        } catch (SdkException e) {
//            e.printStackTrace();
//        } catch (Exception e) {
//            e.printStackTrace();
//        }
//    }
//
//    private static void createMultTransfer(BcOperationService operationService, Long value, Long value1, Long value3) throws Exception {
//        String[] ids = {};
//
//        Map<Integer, String[]> transfer_list = new HashMap<>();
//        String[] item = new String[3];
//        item[0] = to_keypair.getBubiAddress();
//        item[1] = ct_to_keypair[0];
//        item[2] = Long.toString(value);
//        transfer_list.put(1, item);
//
//        String[] item1 = new String[3];
//        item1[0] = to_keypair.getBubiAddress();
//        item1[1] = ct_to_keypair1[0];
//        item1[2] = Long.toString(value1);
//        transfer_list.put(2, item1);
//
//        String[] item2 = new String[3];
//        item2[0] = to_keypair1.getBubiAddress();
//        item2[1] = ct_to_keypair2[0];
//        item2[2] = Long.toString(value3);
//        transfer_list.put(3, item2);
//
//        // params from_keypairs is a map, key is ct public key, value is ct private key
//        Map<String, String> from_keypairs = new HashMap<String, String>();
//        from_keypairs.put(ct_from_keypair[0], ct_from_keypair[1]);
//        from_keypairs.put(ct_from_keypair1[0], ct_from_keypair1[1]);
//
//        JSONObject input = createTx1(from_keypair.getBubiAddress(), from_keypairs, transfer_list, ids);
//        try {
//            Transaction transaction = operationService.newTransaction(from_keypair.getBubiAddress());
//            InvokeContractOperation invokeContractOperation = new InvokeContractOperation.Builder()
//                    .buildDestAddress(contract_address)
//                    .buildInputData(input.toJSONString())
//                    .build();
//
//            TransactionCommittedResult result = transaction.buildAddOperation(invokeContractOperation)
//                    .buildGasPrice(1000L).buildFeeLimit(100000000L)
//                    .buildTxMetadata("交易")
//                    .buildAddSigner(from_keypair.getPubKey(), from_keypair.getPriKey())
//                    .commit();
//
//            System.out.println("\n------------------------------------------------");
//            System.out.println(GsonUtil.toJson(result));
//        } catch (SdkException e) {
//            e.printStackTrace();
//        } catch (Exception e) {
//            e.printStackTrace();
//        }
//    }
//
    public static String getContractCodeFromFile(String fileName) {
        StringBuilder result = new StringBuilder();
        try{
            File file = new File(fileName);
            BufferedReader br = new BufferedReader(new FileReader(file));//构造一个BufferedReader类来读取文件

            String s = null;
            while((s = br.readLine())!=null){//使用readLine方法，一次读一行
                result.append(System.lineSeparator() + s);
            }
            br.close();
        }catch(Exception e){
            e.printStackTrace();
        }
        return result.toString();
    }

    private static TransactionHistory queryTxResult(String hash) {
        int count = 10;
        TransactionHistory transactionHistory = null;
        while (count > 0) {
            try {
                // 循环查询交易结果的次数
                Thread.sleep(1000L);
                transactionHistory = bm.getTransactionByHash(hash);
                System.out.println("----------------------->:查询交易返回的结果:" + JSON.toJSONString(transactionHistory));
                if (!Tools.isNULL(transactionHistory) ){
                    break;//查询交易成功则跳出
                }
            } catch (Throwable e) {
                e.printStackTrace();
            }
            count--;
        }
        return transactionHistory;
    }

    private static void createAccounts(){
        System.out.println("from account: " + from_keypair.getAddress() + ", " + from_keypair.getPrivateKey());
        System.out.println("ct from account: " + ct_from_keypair.getCtPublicKey() + ", " + ct_from_keypair.getCtPrivateKey());
        System.out.println("ct from account1: " + ct_from_keypair1.getCtPublicKey() + ", " + ct_from_keypair1.getCtPrivateKey());
        System.out.println("to account: " + to_keypair.getAddress() + ", " + to_keypair.getPrivateKey());
        System.out.println("ct to account: " + ct_to_keypair.getCtPublicKey() + ", " + ct_to_keypair.getCtPrivateKey());
        System.out.println("ct to account1: " + ct_to_keypair1.getCtPublicKey() + ", " + ct_to_keypair1.getCtPrivateKey());
        System.out.println("ct to account2: " + ct_to_keypair2.getCtPublicKey() + ", " + ct_to_keypair2.getCtPrivateKey());

        String[] addrs = new String[2];
        addrs[0] = from_keypair.getAddress();
        addrs[1] = to_keypair.getAddress();

        try {
            // Create blob data
            BlobDataResp blobDataResp = bm.buildCreateAccountsBlob(genesis_address, addrs);

            // Sign
            String []signerPrivateKeyArr = {genesis_privateKey};
            TransactionSignRequest signRequest = new TransactionSignRequest();
            signRequest.setBlob(blobDataResp.getBlob());
            for (int i = 0; i < signerPrivateKeyArr.length; i++) {
                signRequest.addPrivateKey(signerPrivateKeyArr[i]);
            }
            List<SignEntity> signBlobs = bm.signBlob(blobDataResp, Arrays.asList(signerPrivateKeyArr));

            // Submit tx
            SubmitTxReq submitTxReq = new SubmitTxReq();
            submitTxReq.setBlob(blobDataResp.getBlob());
            submitTxReq.setListSigner(signBlobs);
            submitTxReq.setHash(blobDataResp.getHash());
            submitTxReq.setInitiator(genesis_address);
            TransactionSubmitResponse bcResponse = bm.submitTx(submitTxReq);
            if(bcResponse.getErrorCode().equals(SdkErrorCodeEnum.SUCCESS.getCode())){
                //同步查询交易结果
                TransactionHistory transactionHistory = queryTxResult(submitTxReq.getHash());
                if(!Tools.isNULL(transactionHistory) && ConstantsUtil.BC_SUCCESS.equals(transactionHistory.getErrorCode())) {
                    //成功
                    JSONArray arrayJson = JSON.parseArray(transactionHistory.getErrorDesc());
                } else {
                    System.out.println("查询结果交易失败：" + "errorCode:" + transactionHistory.getErrorCode());
                    throw new RuntimeException(transactionHistory.getErrorCode() +"-" +transactionHistory.getErrorDesc());//抛出异常
                }
            }else{
                System.out.println("提交交易到blockchain异常" + bcResponse.getErrorCode() + " Desc:" + bcResponse.getErrorCode());
                throw new RuntimeException(bcResponse.getErrorCode() +"-" +bcResponse.getErrorDesc());//抛出异常
            }

        } catch (SDKException e) {
            e.printStackTrace();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static void createContract() {
        try {
            String sourceAddress = genesis_address;
            String privateKey = genesis_privateKey;
            CreateContractReq createContractReq = new CreateContractReq();
            createContractReq.setInitBalance(10000000L);
            createContractReq.setSourceAddress(genesis_address);
            createContractReq.setPayLoad(getContractCodeFromFile(contract_file));
            //生成blob
            BlobDataResp blobDataResp = bm.buildCreateContractBlob(genesis_address, 1100000000L, createContractReq);
            //签名
            List<String> privateKeys = new ArrayList<>();
            privateKeys.add(privateKey);
            List<SignEntity> signBlobs = bm.signBlob(blobDataResp, privateKeys);
            //提交
            SubmitTxReq submitTxReq = new SubmitTxReq();
            submitTxReq.setBlob(blobDataResp.getBlob());
            submitTxReq.setListSigner(signBlobs);
            submitTxReq.setHash(blobDataResp.getHash());
            submitTxReq.setInitiator(sourceAddress);
            TransactionSubmitResponse bcResponse = bm.submitTx(submitTxReq);
            if (bcResponse.getErrorCode().equals(SdkErrorCodeEnum.SUCCESS.getCode())) {
                //同步查询交易结果
                TransactionHistory transactionHistory = queryTxResult(submitTxReq.getHash());
                if (!Tools.isNULL(transactionHistory) && ConstantsUtil.BC_SUCCESS.equals(transactionHistory.getErrorCode())) {
                    //成功
                    JSONArray arrayJson = JSON.parseArray(transactionHistory.getErrorDesc());
                    contract_address = arrayJson.getJSONObject(0).getString("contract_address");

                } else {
                    System.out.println("查询结果交易失败：" + "errorCode:" + transactionHistory.getErrorCode());
                    throw new RuntimeException(transactionHistory.getErrorCode() + "-" + transactionHistory.getErrorDesc());//抛出异常
                }
            } else {
                System.out.println("提交交易到blockchain异常" + bcResponse.getErrorCode() + " Desc:" + bcResponse.getErrorCode());
                throw new RuntimeException(bcResponse.getErrorCode() + "-" + bcResponse.getErrorDesc());//抛出异常
            }
        } catch (SDKException e) {
            e.printStackTrace();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * 创建账户操作
     */
    private static void issueCtAsset(Long value){
        try {

            PrivacyProtobuf.ConfidentialAsset.Builder issueAsset = CocoonUtil.buildAsset(from_keypair.getAddress(), value,ct_from_keypair.getCtPrivateKey(), ct_from_keypair.getCtPublicKey());

            JsonFormat jsonFormat = new JsonFormat();
            String tokenStr = jsonFormat.printToString(issueAsset.build());

            JSONObject input = new JSONObject();
            input.put("method", "issue");
            JSONObject params = new JSONObject();
            params.put("name", "Confidential");
            params.put("symbol", "CT");
            params.put("token", JSONObject.parseObject(tokenStr));
            input.put("params", params);
            System.out.println("input: " + JSON.toJSONString(input, true));

            try {
                BcTransferReq bcTransfer = new BcTransferReq();
                bcTransfer.setFromAddress(from_keypair.getAddress());
                bcTransfer.setToAddress(contract_address);
                bcTransfer.setTxType(BcTxTypeEnum.CONTRACT.getCode());
                bcTransfer.setAmount(0L);
                bcTransfer.setInput(input.toJSONString());
                BlobDataResp blobDataResp = bm.buildBlobData(bcTransfer, from_keypair.getAddress());
                //签名
                List<String> privateKeys = new ArrayList<>();
                privateKeys.add(from_keypair.getPrivateKey());
                List<SignEntity> signBlobs = bm.signBlob(blobDataResp,privateKeys);
                //提交
                SubmitTxReq submitTxReq = new SubmitTxReq();
                submitTxReq.setBlob(blobDataResp.getBlob());
                submitTxReq.setListSigner(signBlobs);
                submitTxReq.setHash(blobDataResp.getHash());
                submitTxReq.setInitiator(from_keypair.getAddress());
                TransactionSubmitResponse bcResponse = bm.submitTx(submitTxReq);
                if(bcResponse.getErrorCode().equals(0)){
                    //System.out.println(bcResponse.getResult()); // print issue output
                }else{
                    System.out.println("提交交易到blockchain异常" + bcResponse.getErrorCode() + " hash:" + submitTxReq.getHash());
                    throw new RuntimeException(bcResponse.getErrorCode() +"-" +bcResponse.getErrorDesc());//抛出异常
                }


            } catch (SDKException e) {
                e.printStackTrace();
            } catch (Exception e) {
                e.printStackTrace();
            }

        } catch (SDKException e) {
            e.printStackTrace();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
