pragma solidity ^0.4.24;

import "./subdir/super.sol";

contract Test is Super {

  string public value;

  constructor(string _value) public {
    value=_value;
  }


  function setValue(string _value) external {
    value = _value;
  }

  function willFail() external {
    alwaysFails();
  }


}
