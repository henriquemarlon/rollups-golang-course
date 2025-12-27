// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {ERC20} from "@openzeppelin-contracts-5.2.0/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin-contracts-5.2.0/token/ERC20/extensions/ERC20Burnable.sol";
import {ERC20Pausable} from "@openzeppelin-contracts-5.2.0/token/ERC20/extensions/ERC20Pausable.sol";

contract Token is ERC20, ERC20Burnable, ERC20Pausable {
    constructor() ERC20("MyToken", "MTK") {}

    function pause() public {
        _pause();
    }

    function unpause() public {
        _unpause();
    }

    function mint(address to, uint256 amount) public {
        _mint(to, amount);
    }

    // The following functions are overrides required by Solidity.

    function _update(address from, address to, uint256 value) internal override(ERC20, ERC20Pausable) {
        super._update(from, to, value);
    }
}
