pragma solidity ^0.4.24;

contract Test {

  string public value;

  constructor(string _value) public {
    value=_value;
  }


  function setValue(string _value) external {
    value = _value;
  }


}
