package cn.bubi;

import cn.bubi.encryption.utils.hex.HexFormat;
import cn.bubi.utils.AesCbc;

public class Main {

    public static void main(String[] args) throws Exception {
        Privacy privacy = new Privacy();
        long cppObj = privacy.getCppObject();

        // test create keypair
        System.out.println("Test keypair");
        String[] keypair1 = privacy.createKeyPair(cppObj);
        System.out.println("pubkey1:" + keypair1[2] + ", seckey1:" + keypair1[3]);

        String[] keypair2 = privacy.createKeyPair(cppObj);
        System.out.println("pubkey2:" + keypair2[2] + ", seckey2:" + keypair2[3] + "\n");

        // test get public key
        System.out.println("Test get public key");
        String[] pubkey_ret_1 = privacy.getPublicKey(cppObj, "250fd4784392ed30ead5e94531227e569e4b165cca541194b3eda69f58cc11c0");
        String errCodePubRet1 = pubkey_ret_1[0];
        String errDescPubRet1 = pubkey_ret_1[1];
        String pubkey1 = pubkey_ret_1[2];
        if (Integer.valueOf(errCodePubRet1) == 0) {
            System.out.println("pubkey1 from seckey1:" + pubkey1 + "\n");
        } else {
            System.out.println("Error pubkey1 from seckey1:" + errDescPubRet1 + "\n");
        }

        // test ecdh
        System.out.println("Test ecdh");
        String[] ecdh_ret_1 = privacy.createEcdhKey(cppObj, keypair1[3], keypair2[2]);
        String errCodeEcdhRet1 = ecdh_ret_1[0];
        String errDescEcdhRet1 = ecdh_ret_1[1];
        String ecdh1 = ecdh_ret_1[2];
        if (Integer.valueOf(errCodeEcdhRet1) == 0) {
            System.out.println("ecdh r1 * R2:" + ecdh1);
        } else {
            System.out.println("Error ecdh r1 * R2:" + errDescEcdhRet1 + "\n");
        }
        String[] ecdh_ret_2 = privacy.createEcdhKey(cppObj, keypair2[3], keypair1[2]);
        String errCodeEcdhRet2 = ecdh_ret_2[0];
        String errDescecdhRet2 = ecdh_ret_2[1];
        String ecdh2 = ecdh_ret_2[2];
        if (Integer.valueOf(errCodeEcdhRet2) == 0) {
            System.out.println("ecdh r2 * R1:" + ecdh2 + "\n");
        } else {
            System.out.println("Error ecdh r2 * R1:" + errDescecdhRet2 + "\n");
        }

        // test ecdsa
        System.out.println("Test ecdsa");
        String msg = keypair1[3].substring(0, 32);
        String[] sig_ret = privacy.ecdsaSign(cppObj, keypair1[3], msg);
        String errCodeSigRet = sig_ret[0];
        String errDescSigRet = sig_ret[1];
        String sig = sig_ret[2];
        if (Integer.valueOf(errCodeSigRet) == 0) {
            System.out.println("ecdsa sign(r1, data):" + sig);
            String[] ecd_very_ret = privacy.ecdsaVerify(cppObj, keypair1[2], msg, sig);
            String errCodeEcdVeryRet = ecd_very_ret[0];
            String errDescEcdVeryRet = ecd_very_ret[1];
            if (Integer.valueOf(errCodeEcdVeryRet) == 0) {
                System.out.println("ecdsa verify(R1, sig, data) successfully\n");
            } else {
                System.out.println("Error ecdsa verify(R1, sig, data):" + errDescEcdVeryRet + " \n");
            }
        } else {
            System.out.println("Error ecdsa sign(r1, data):" + errDescSigRet);
        }

        // test commitment
        System.out.println("Test commitment");
        String blind1= ecdh1; // r = H(r1 * R2)
        String[] commit_ret = privacy.createPedersenCommit(cppObj, 1000, blind1);
        String errCodeCmtRet = commit_ret[0];
        String errDescCmtRet = commit_ret[1];
        String commit = commit_ret[2];
        if (Integer.valueOf(errCodeCmtRet) == 0) {
            System.out.println("Commit from sender with H(r1 * R2) * G + v * H:" + commit + "\n");
        } else {
            System.out.println("Error commit from sender with H(r1 * R2) * G + v * H:" + errDescCmtRet + "\n");
        }

        // test range proof
        System.out.println("Test bulletproof rangeproof");
        String[] proof_ret = privacy.bpRangeproofProve(cppObj, blind1, 1000);
        String errCodePrfRet = proof_ret[0];
        String errDescPrfRet = proof_ret[1];
        String proof = proof_ret[2];
        String[] commits = new String[1];
        commits[0] = commit;

        if (Integer.valueOf(errCodePrfRet) == 0) {
            System.out.println("Create range proof :" + proof + "\n");
            String[] bp_very_ret1 = privacy.bpRangeproofVerify(cppObj, commit, proof);
            String errCodeBpRet1 = bp_very_ret1[0];
            String errDescBpRet1 = bp_very_ret1[1];
            if (Integer.valueOf(errCodeBpRet1) == 0) {
                System.out.println("Verify range proof successful\n");
            } else {
                System.out.println("Error verify range proof:" + proof + ": " + errDescBpRet1 + "\n");
            }
        } else {
            System.out.println("Error bpRangeproofPrvoe: " + errDescPrfRet + ", " + errDescPrfRet + "\n");
        }


        // test excess
        System.out.println("Test tally");

        String[] blind_issue_ret = privacy.createEcdhKey(cppObj, keypair1[3], keypair1[2]);; // r = H(r1 * R1)
        String errCodeBlindRet = blind_issue_ret[0];
        String errDescBlindRet = blind_issue_ret[1];
        String blind_issue = blind_issue_ret[2];
        if (Integer.valueOf(errCodeBlindRet) == 0) {
            System.out.println("Issue blind :" + blind_issue + "\n");
        } else {
            System.out.println("Error issue blind :" + errDescBlindRet + "\n");
        }

        String[] commit_issue_ret = privacy.createPedersenCommit(cppObj, 2000, blind_issue);
        String errCodeCmtIssueRet = commit_issue_ret[0];
        String errDescCmtIssueRet = commit_issue_ret[1];
        String commit_issue = commit_issue_ret[2];
        if (Integer.valueOf(errCodeCmtIssueRet) == 0) {
            System.out.println("Issue commit :" + commit_issue + "\n");
        } else {
            System.out.println("Error issue commit :" + errDescCmtIssueRet + "\n");
        }

        int change_size = 1;
        String[] outputs = new String[change_size+1];
        String[] output_blinds = new String[change_size+1];
        String[] values1 = new String[change_size+1];
        outputs[0] = commit;
        output_blinds[0] = blind1;
        values1[0] = Long.toString(1000);
        long avarage_value = 1000/change_size;
        long first_one_value = 1000 % change_size + avarage_value;
        System.out.println("ava:" + avarage_value + ", first:" + first_one_value + "," + avarage_value * 98 + first_one_value);
        for (int i = 0; i < change_size; i++) {
            long change_value = avarage_value;
            if(i==0) change_value = first_one_value;
            String[] res = privacy.createPedersenCommit(cppObj, change_value, blind_issue);
            String errCodeCmtCngRet = res[0];
            String errDescCmtCngRet = res[1];
            String commit_change = res[2];
            if (Integer.valueOf(errCodeCmtCngRet) == 0) {
                System.out.println("Commit change :" + commit_change + "\n");
            } else {
                System.out.println("Error commit change :" + errDescCmtCngRet + "\n");
            }
            outputs[i+1] = commit_change;
            output_blinds[i+1] = blind_issue;
            values1[i+1] = Long.toString(change_value);
        }

        String[] inputs = new String[1];
        String[] input_blinds = new String[1];
        inputs[0] = commit_issue;
        input_blinds[0] = blind_issue;

        String data = keypair1[3].substring(0, 32);
        System.out.println("excess sign ……");
        String[] excess_sig_ret = privacy.excessSign(cppObj, input_blinds, output_blinds, data);
        String errCodeExcSigRet = excess_sig_ret[0];
        String errDescExcSigRet = excess_sig_ret[1];
        String excess_sig  = excess_sig_ret[2];
        if (Integer.valueOf(errCodeExcSigRet) == 0) {
            System.out.println("Excess sig :" + excess_sig + "\n");

            System.out.println("excess sign done, start verify tally ……");
            String[] tally_very_ret = privacy.pedersenTallyVerify(cppObj, inputs, outputs, data, excess_sig);
            String errCodeTalVeryRet = tally_very_ret[0];
            String errDescTalVeryRet = tally_very_ret[1];
            if (Integer.valueOf(errCodeTalVeryRet) == 0) {
                System.out.println("Verify tally result successfully\n");
            } else {
                System.out.println("Error failed to Verify tally result\n");
            }
        } else {
            System.out.println("Error excess sig :" + errDescExcSigRet + "\n");
        }
        System.out.println("inputs[0]:" + inputs[0]);
        System.out.println("outputs[0]:" + outputs[0] + ", outputs[1]:" + outputs[1]);
        System.out.println("data:" + data + ", excess_sig:" + excess_sig + "\n");

        // AES-CBC test
        System.out.println("Test aes");
        String aes_key = "4b01f09275b87b602acf34483b65242304d0b84b5d641195094305d2e78774f6";
        String encrypt_data = HexFormat.byteToHex(AesCbc.encrypt("100000000".getBytes(), aes_key.substring(0, 32).getBytes())).toLowerCase();
        System.out.println("encrypt data:" + encrypt_data);
        System.out.println("decrypt data:" + AesCbc.decrypt(HexFormat.hexToByte(encrypt_data), aes_key.substring(0, 32).getBytes()));

        // destory
        if(cppObj != -1) privacy.destroy(cppObj);
    }
}
