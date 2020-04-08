pragma solidity ^0.6.0;

import "./subdir/super.sol";

contract Test is Super {

  string public value;

  constructor(string memory _value) public {
    value = _value;
  }


  function setValue(string calldata _value) external {
    value = _value;
  }

  function willFail() external pure {
    alwaysFails();
  }


}
