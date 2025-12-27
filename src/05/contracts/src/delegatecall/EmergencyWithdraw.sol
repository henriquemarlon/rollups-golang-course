// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {IERC20} from "@openzeppelin-contracts-5.2.0/token/ERC20/IERC20.sol";

contract EmergencyWithdraw {
    function emergencyERC20Withdraw(IERC20 token, address to) public {
        token.transfer(to, token.balanceOf(address(this)));
    }

    function emergencyETHWithdraw(address to) public {
        uint256 balance = address(this).balance;
        require(balance > 0, "No ETH to withdraw");
        (bool success,) = to.call{value: balance}("");
        require(success, "ETH transfer failed");
    }
}
