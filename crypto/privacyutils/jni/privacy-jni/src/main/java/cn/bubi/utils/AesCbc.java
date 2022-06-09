package cn.bubi.utils;

import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.security.AlgorithmParameters;

public class AesCbc {

    public static byte[] encrypt(byte[] plainText, byte[] key) throws Exception {
        byte[] encrypted = null;
        try {
            int len = 0;
            if (plainText.length % 16 != 0) {
                len = (plainText.length / 16 + 1) * 16;
            }
            byte[] clean = new byte[len];
            System.arraycopy(plainText, 0, clean, 0, plainText.length);

            byte[] iv = new byte[16];
            IvParameterSpec ivParameterSpec = new IvParameterSpec(iv);
            SecretKeySpec secretKeySpec = new SecretKeySpec(key, "AES");
            Cipher cipher = Cipher.getInstance("AES/CBC/NoPadding");
            AlgorithmParameters params = AlgorithmParameters.getInstance("AES");
            params.init(ivParameterSpec);
            cipher.init(Cipher.ENCRYPT_MODE, secretKeySpec, params);
            encrypted = cipher.doFinal(clean);
        }
        catch (Exception e) {
            e.printStackTrace();
        }
        return encrypted;
    }

    public static String decrypt(byte[] encryptedIvTextBytes, byte[] key) throws Exception {
        byte[] decrypted = null;
        try {
            byte[] iv = new byte[16];
            IvParameterSpec ivParameterSpec = new IvParameterSpec(iv);
            SecretKeySpec secretKeySpec = new SecretKeySpec(key, "AES");
            Cipher cipher = Cipher.getInstance("AES/CBC/NoPadding");
            AlgorithmParameters params = AlgorithmParameters.getInstance("AES");
            params.init(ivParameterSpec);
            cipher.init(Cipher.DECRYPT_MODE, secretKeySpec, params);
            decrypted = cipher.doFinal(encryptedIvTextBytes);
        }
        catch (Exception e) {
            e.printStackTrace();
        }

        String des = "";
        if (decrypted != null) {
            des = new String(decrypted);
        }

        return des.substring(0, des.indexOf("\u0000"));
    }
}