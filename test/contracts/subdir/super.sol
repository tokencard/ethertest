pragma solidity ^0.6.0;

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
