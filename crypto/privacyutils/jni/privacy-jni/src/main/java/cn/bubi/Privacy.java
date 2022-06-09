package cn.bubi;

import cn.bubi.utils.OSinfo;

import java.io.File;

public class Privacy {

    static{
        long obj = -1;
        try {
            String path = Privacy.class.getResource("/libbuchain_secp256k1.so").getPath();
            if(OSinfo.isWindows()){
                path = Privacy.class.getResource("/secp256k1.dll").getPath();
            }
            if(OSinfo.isMacOS()){
                path = Privacy.class.getResource("/libbuchain_secp256k1.dylib").getPath();
            }
            System.load(path);
            System.out.println("Load privacy dynamic library done, path:" + path);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public long getCppObject() {
         return init();
    }

    /**
     * construct c++ object
     *
     * @return c++ object pointer
     */
    public native long init();

    /**
     * destruct c++ object
     *
     * @param obj long type of c++ object pointer
     */
    public native  void destroy(long obj);
    /**
     *  create Pedersen Commitment
     *
     * @param obj c++ object pointer
     * @param value long type of the value to hide
     * @param blind string of blind factor, must be 64 bytes (32 bytes with 4-bit hex encoding)
     * @return string array with 3 element, ret[0] is errcode, ret[1] is error message, ret[2] is pedersen commitment
     */
    public native String[] createPedersenCommit(long obj, long value, String blind);

    /**
     *  create ECDH secret key
     *
     * @param obj c++ object pointer
     * @param seckey string of secret key used in ecdh
     * @param pubkey string of public key used in ecdh
     * @return string array with 3 element, ret[0] is errcode, ret[1] is error message, ret[2] is ecdh key
     */
    public native String[] createEcdhKey(long obj, String seckey, String pubkey);

    /**
     *  create ecc-secp256k1 public key and secret key
     *
     * @param obj c++ object pointer
     * @return string array with 4 element, ret[0] is errcode, ret[1] is error message, ret[2] is public key, ret[3] is secret key
     *  public key and secret key is encode with hex every 4 bit
     */
    public native String[] createKeyPair(long obj);

    /**
     *  get public key from secret key
     *
     * @param obj c++ object pointer
     * @param seckey string of secret key
     * @return string array with 3 element, ret[0] is errcode, ret[1] is error message, ret[2] is public key
     */
    public native String[] getPublicKey(long obj, String seckey);

    /**
     *  use bulletproof to create the rangeproof of input value, prove the input value is positive
     *
     * @param obj c++ object pointer
     * @param blind string of blind factor
     * @param value long type of output value
     * @return string array with 3 element, ret[0] is errcode, ret[1] is error message, ret[2] is proof
     */
    public native String[] bpRangeproofProve(long obj, String blind, long value);

    /**
     *  use bulletproof to verify the input value is positive
     *
     * @param obj c++ object pointer
     * @param commit string of pedersen commit
     * @param proof string of rangeproof
     * @return string array with 2 element, ret[0] is errcode, ret[1] is error message, error code 0 means verify true
     */
    public native String[] bpRangeproofVerify(long obj, String commit, String proof);

    /**
     *  create ecdsa signature
     *
     * @param obj c++ object pointer
     * @param seckey string of secret key
     * @param data string of data, must be 64 bytes (32 bytes with 4-bit hex encoding)
     * @return string array with 3 element, ret[0] is errcode, ret[1] is error message, ret[2] is ecdsa signature result
     */
    public native String[] ecdsaSign(long obj, String seckey, String data);

    /**
     *  verify ecdsa signature
     *
     * @param obj c++ object pointer
     * @param pubkey string of public key
     * @param data string of the signature data, must be 64 bytes (32 bytes with 4-bit hex encoding)
     * @param sig string of the signature result
     * @return string array with 2 element, ret[0] is errcode, ret[1] is error message, error code 0 means verify true
     */
    public native String[] ecdsaVerify(long obj, String pubkey, String data, String sig);

    /**
     *  create transaction excess signature
     *
     * @param obj c++ object pointer
     * @param inputs string array of input blind factors
     * @param outputs string array of output blind factors
     * @param msg string of the signature message, must be 64 bytes (32 bytes with 4-bit hex encoding)
     * @return string array with 3 element, ret[0] is errcode, ret[1] is error message, ret[2] is excess signature result
     */
    public native String[] excessSign(long obj, String[] inputs, String[] outputs, String msg);

    /**
     *  verify inputs equal to outputs
     *
     * @param obj c++ object pointer
     * @param inputs string array of input pedersen commitments
     * @param outputs string array of output pedersen commitments
     * @param msg string of the transaction excess signature message, must be 64 bytes (32 bytes with 4-bit hex encoding)
     * @param sig string of the transaction excess signature result
     * @return string array with 2 element, ret[0] is errcode, ret[1] is error message, error code 0 means verify true
     */
    public native String[] pedersenTallyVerify(long obj, String[] inputs, String[] outputs, String msg, String sig);
}
