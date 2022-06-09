package cn.bubi.utils;

public enum BcTxTypeEnum {

    TRANSFER(0, "默认代币转账"),
	ASSET_ATP10(1, "ATP10"),
	CONTRACT(2, "调用合约"),
    ;


    private final Integer code;
    private final String msg;

    private BcTxTypeEnum(Integer code, String msg) {
        this.code = code;
        this.msg = msg;
    }

    public Integer getCode() {
        return code;
    }

	public String getMsg() {
		return msg;
	}
    
}
