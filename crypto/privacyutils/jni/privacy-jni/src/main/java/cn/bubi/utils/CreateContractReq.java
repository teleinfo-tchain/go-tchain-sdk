package cn.bubi.utils;

public class CreateContractReq {

    private String sourceAddress;
    private Long initBalance;
    private String payLoad;
    private String input;

    public String getSourceAddress() {
        return sourceAddress;
    }

    public void setSourceAddress(String sourceAddress) {
        this.sourceAddress = sourceAddress;
    }

    public Long getInitBalance() {
        return initBalance;
    }

    public void setInitBalance(Long initBalance) {
        this.initBalance = initBalance;
    }

    public String getPayLoad() {
        return payLoad;
    }

    public void setPayLoad(String payLoad) {
        this.payLoad = payLoad;
    }

    public String getInput() {
        return input;
    }

    public void setInput(String input) {
        this.input = input;
    }
}
