// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.7.0;

contract Multiply {
    event Print(uint);
    event Log(string);

    uint public weight;

    constructor(uint data) public {
        weight = data;

    }

   function multiply(uint input) public returns (uint) {
       weight = input*2 + weight;
       emit Print(weight);
       return weight;
   }

   function getWeight(uint input) public returns (uint) {
      emit Print(weight+1+input);
      return weight+1+input;

   }

   function testStr(string memory input) public returns (string memory) {
      emit Log(input);
      return input;
   }

  function testInt(string memory input) public returns (string memory) {
      emit Print(7);
      return input;
   }
}