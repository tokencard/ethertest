pragma solidity ^0.4.24;

contract Super {

  constructor() public {
  }

  function alwaysFails() pure internal {
    if (false) {
      require(2 == 2);
    }
    require(2>3, "will fail");
  }

  function neverCalled() pure internal {
    require(1 == 1);
  }

}
