package cn.bubi.utils;

import cn.bubi.SDK;
import cn.bubi.blockchain.TransactionService;
import cn.bubi.blockchain.impl.TransactionServiceImpl;
import cn.bubi.encryption.key.PrivateKey;
import cn.bubi.encryption.utils.hex.HexFormat;
import cn.bubi.model.request.*;
import cn.bubi.model.request.operation.AssetSendOperation;
import cn.bubi.model.request.operation.BUSendOperation;
import cn.bubi.model.request.operation.ContractCreateOperation;
import cn.bubi.model.request.operation.ContractInvokeByAssetOperation;
import cn.bubi.model.response.*;
import cn.bubi.model.response.result.TransactionBuildBlobResult;
import cn.bubi.model.response.result.data.MetadataInfo;
import cn.bubi.model.response.result.data.Signature;
import cn.bubi.model.response.result.data.TransactionHistory;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

public class BlockChainManager {
//    @Value("${blockchain.node.ip}")
    private String bcUrl = "http://192.168.10.100:29333";
//    @Value("${blockchain.single.tx.fee}")
    private long fee = 10000000000L;
//    @Value("${blockchain.tx.gasprice}")
    private long gasprice = 1000;

    private Long appearLegerTime = 20L;

    public BlobDataResp buildCreateContractBlob(String initiator,Long txFee, CreateContractReq createContractReq){
        BlobDataResp blobData = null;
        ContractCreateOperation operation = new ContractCreateOperation();
        operation.setSourceAddress(createContractReq.getSourceAddress());
        operation.setInitBalance(createContractReq.getInitBalance());
        operation.setPayload(createContractReq.getPayLoad());
        operation.setInitInput(createContractReq.getInput());
        //获取交易nonce
        long accNonce = getAccountNonce(initiator)+1;
        TransactionBuildBlobRequest transactionBuildBlobRequest = new TransactionBuildBlobRequest();
        transactionBuildBlobRequest.setSourceAddress(initiator);
        transactionBuildBlobRequest.setNonce(accNonce);
        transactionBuildBlobRequest.setFeeLimit(txFee);
        transactionBuildBlobRequest.setGasPrice(gasprice);
        transactionBuildBlobRequest.addOperation(operation);
        // 获取交易BLob串
        TransactionService transactionService = new TransactionServiceImpl();
        TransactionBuildBlobResponse transactionBuildBlobResponse = transactionService.buildBlob(transactionBuildBlobRequest);
        if(transactionBuildBlobResponse.getErrorCode().equals(SdkErrorCodeEnum.SUCCESS.getCode())){
            TransactionBuildBlobResult transactionBuildBlobResult = transactionBuildBlobResponse.getResult();
            blobData = new BlobDataResp();
            blobData.setBlob(transactionBuildBlobResult.getTransactionBlob());
            blobData.setHash(transactionBuildBlobResult.getHash());
        }

        return blobData;
    }

    public BlobDataResp buildCreateAccountsBlob(String initiator, String[] accounts) {
        BlobDataResp blobData = null;
        Long accNonce = getAccountNonce(initiator)+1;

        TransactionBuildBlobRequest transactionBuildBlobRequest = new TransactionBuildBlobRequest();
        transactionBuildBlobRequest.setSourceAddress(initiator);
        transactionBuildBlobRequest.setNonce(accNonce);
        transactionBuildBlobRequest.setFeeLimit(fee);
        transactionBuildBlobRequest.setGasPrice(gasprice);
        for(String addr: accounts) {
            BUSendOperation operation = new BUSendOperation();
            operation.setSourceAddress(initiator);
            operation.setDestAddress(addr);
            operation.setAmount(100000000000L);
            transactionBuildBlobRequest.addOperation(operation);
        }
        // 获取交易BLob串
        TransactionService transactionService = new TransactionServiceImpl();
        TransactionBuildBlobResponse transactionBuildBlobResponse = transactionService.buildBlob(transactionBuildBlobRequest);
        if(transactionBuildBlobResponse.getErrorCode().equals(SdkErrorCodeEnum.SUCCESS.getCode())){
            TransactionBuildBlobResult transactionBuildBlobResult = transactionBuildBlobResponse.getResult();
            blobData = new BlobDataResp();
            blobData.setBlob(transactionBuildBlobResult.getTransactionBlob());
            blobData.setHash(transactionBuildBlobResult.getHash());
        }
        return blobData;
    }

    public BlobDataResp buildBlobData(BcTransferReq bcTransfer, String initiator){
        BlobDataResp blobData = null;
        //获取交易nonce
        long accNonce = getAccountNonce(initiator)+1;
        TransactionBuildBlobRequest transactionBuildBlobRequest = new TransactionBuildBlobRequest();
        transactionBuildBlobRequest.setSourceAddress(initiator);
        transactionBuildBlobRequest.setNonce(accNonce);
        transactionBuildBlobRequest.setFeeLimit(fee);
        transactionBuildBlobRequest.setGasPrice(gasprice);
//        Long timeOut = bcTransfer.getTimeOut() ;
//        if(!Tools.isNull(timeOut)){
//            Long ceilLedgerSeq = bcTransfer.getTimeOut()/appearLegerTime;
//            transactionBuildBlobRequest.setCeilLedgerSeq(ceilLedgerSeq);
//        }

        Integer txType = bcTransfer.getTxType();
        if(txType.equals(BcTxTypeEnum.TRANSFER.getCode())){//BU
            BUSendOperation buSendOperation = new BUSendOperation();
            buSendOperation.setSourceAddress(bcTransfer.getFromAddress());
            buSendOperation.setDestAddress(bcTransfer.getToAddress());
            buSendOperation.setAmount(bcTransfer.getAmount());
            buSendOperation.setMetadata(bcTransfer.getMetadata());
            transactionBuildBlobRequest.addOperation(buSendOperation);
        }else if(txType.equals(BcTxTypeEnum.CONTRACT.getCode())){
            ContractInvokeByAssetOperation operation = new ContractInvokeByAssetOperation();
            operation.setContractAddress(bcTransfer.getToAddress());
            operation.setSourceAddress(bcTransfer.getFromAddress());
            operation.setInput(bcTransfer.getInput());
            transactionBuildBlobRequest.addOperation(operation);
        }else if(txType.equals(BcTxTypeEnum.ASSET_ATP10.getCode())){
            AssetSendOperation assetSendOperation = new AssetSendOperation();
            assetSendOperation.setSourceAddress(bcTransfer.getFromAddress());
            assetSendOperation.setDestAddress(bcTransfer.getToAddress());
            assetSendOperation.setAmount(bcTransfer.getAmount());
            assetSendOperation.setCode(bcTransfer.getCode());
            assetSendOperation.setIssuer(bcTransfer.getIssuer());
            assetSendOperation.setMetadata(bcTransfer.getMetadata());
            transactionBuildBlobRequest.addOperation(assetSendOperation);
        }
        // 获取交易BLob串
        TransactionService transactionService = new TransactionServiceImpl();
        TransactionBuildBlobResponse transactionBuildBlobResponse = transactionService.buildBlob(transactionBuildBlobRequest);
        if(transactionBuildBlobResponse.getErrorCode().equals(SdkErrorCodeEnum.SUCCESS.getCode())){
            TransactionBuildBlobResult transactionBuildBlobResult = transactionBuildBlobResponse.getResult();
            blobData = new BlobDataResp();
            blobData.setBlob(transactionBuildBlobResult.getTransactionBlob());
            blobData.setHash(transactionBuildBlobResult.getHash());
        }
        return blobData;
    }

    public Long getAccountNonce(String address){
        try{
            SDK sdk = SDK.getInstance(bcUrl);
            AccountGetNonceRequest request = new AccountGetNonceRequest();
            request.setAddress(address);
            AccountGetNonceResponse response = sdk.getAccountService().getNonce(request);
            if(SdkErrorCodeEnum.SUCCESS.getCode().equals(response.getErrorCode())){
                return response.getResult().getNonce();
            }
            System.out.println(response.getErrorDesc());
        }catch(Exception e){
            e.printStackTrace();
        }
        return null;
    }

    public TransactionSubmitResponse submitTx(SubmitTxReq submitTxReq){
        TransactionSubmitRequest transactionSubmitRequest = new TransactionSubmitRequest();
        transactionSubmitRequest.setTransactionBlob(submitTxReq.getBlob());

        List<SignEntity> listSigner = submitTxReq.getListSigner();
        int length = listSigner.size();
        Signature[] signatures = new Signature[length];
        for(int i=0;i<length;i++){
            SignEntity signEntity = listSigner.get(i);
            Signature signature = new Signature();
            signature.setPublicKey(signEntity.getPublicKey());
            signature.setSignData(signEntity.getSignBlob());
            signatures[i]=signature;
        }

        transactionSubmitRequest.setSignatures(signatures);
        SDK sdk = SDK.getInstance(bcUrl);
        TransactionSubmitResponse transactionSubmitResponse = sdk.getTransactionService().submit(transactionSubmitRequest);
        return transactionSubmitResponse;
    }

    public TransactionHistory getTransactionByHash(String txhash) {
        try{
            SDK sdk = SDK.getInstance(bcUrl);
            TransactionGetInfoRequest request = new TransactionGetInfoRequest();
            request.setHash(txhash);
            TransactionGetInfoResponse response = sdk.getTransactionService().getInfo(request);
            if(SdkErrorCodeEnum.NOT_EXIST.getCode() == response.getErrorCode()) {
                return null;
            }
            if(SdkErrorCodeEnum.SUCCESS.getCode() == response.getErrorCode()) {
                List<TransactionHistory> listHis = Arrays.asList(response.getResult().getTransactions());
                for (TransactionHistory transactionHistory : listHis) {
                    if(txhash.equals(transactionHistory.getHash())) {
                        return transactionHistory;
                    }
                }
            }
        }catch(Exception e){
            e.printStackTrace();
        }
        return null;
    }

    public List<SignEntity> signBlob(BlobDataResp blobDataResp,List<String> privateKeys){
        List<SignEntity> listSigner = new ArrayList<SignEntity>();
        for(String privateKey : privateKeys){
            PrivateKey privateObj = new PrivateKey(privateKey);
            byte[] signByte = privateObj.sign(HexFormat.hexStringToBytes(blobDataResp.getBlob()));
            String signBlob = HexFormat.byteToHex(signByte);
            SignEntity entity = new SignEntity();
            entity.setPublicKey(privateObj.getEncPublicKey());
            entity.setSignBlob(signBlob);
            listSigner.add(entity);
        }
        return listSigner;
    }

    /**
     * 获取账户下metadata某个key的值
     * @param address
     * @param key
     * @return
     */
    public String getAccMetadataByKey(String address,String key){
        try{
            SDK sdk = SDK.getInstance(bcUrl);
            AccountGetMetadataRequest request = new AccountGetMetadataRequest();
            request.setAddress(address);
            request.setKey(key);
            AccountGetMetadataResponse response = sdk.getAccountService().getMetadata(request);
            if(SdkErrorCodeEnum.SUCCESS.getCode().equals(response.getErrorCode())){
                MetadataInfo[] metadataInfos = response.getResult().getMetadatas();
                return metadataInfos[0].getValue();
            }
        }catch(Exception e){
            e.printStackTrace();
        }
        return null;
    }

}
