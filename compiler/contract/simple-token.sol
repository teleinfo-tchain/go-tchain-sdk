pragma solidity ^0.5.0;

contract Sample {
    int public val_0;
    uint256 public val_1;

    constructor() public {
        val_0 = 10;
        val_1 = 2;
    }

    function getName() public view returns (int) {
        return val_0;
    }
}