// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Test} from "forge-std-1.9.7/src/Test.sol";
import {Token} from "../src/token/ERC20/Token.sol";
import {MockApplication} from "./mock/MockApplication.sol";
import {EmergencyWithdraw} from "../src/delegatecall/EmergencyWithdraw.sol";
import {Outputs} from "cartesi-rollups-contracts-2.0.0/src/common/Outputs.sol";

contract EmergencyWithdrawTest is Test {
    Token public token;
    MockApplication public mockApplication;
    EmergencyWithdraw public emergencyWithdraw;

    address public user;
    address public recipient;

    event Transfer(address indexed from, address indexed to, uint256 value);

    function setUp() public {
        user = makeAddr("user");
        recipient = makeAddr("recipient");

        token = new Token();
        mockApplication = new MockApplication();
        emergencyWithdraw = new EmergencyWithdraw();

        // Fund the mock application with tokens and ETH
        token.mint(address(mockApplication), 1000e18);
        vm.deal(address(mockApplication), 10 ether);
    }

    function test_EmergencyERC20WithdrawThroughDelegateCallVoucher() public {
        uint256 initialBalance = token.balanceOf(address(mockApplication));
        uint256 recipientInitialBalance = token.balanceOf(recipient);

        bytes memory encodedWithdrawTx = abi.encodeCall(EmergencyWithdraw.emergencyERC20Withdraw, (token, recipient));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(emergencyWithdraw), encodedWithdrawTx));

        vm.expectEmit(true, true, false, true);
        emit Transfer(address(mockApplication), recipient, initialBalance);
        mockApplication.executeOutput(delegateCallVoucher);

        assertEq(token.balanceOf(address(mockApplication)), 0);
        assertEq(token.balanceOf(recipient), recipientInitialBalance + initialBalance);
    }

    function test_EmergencyETHWithdrawThroughDelegateCallVoucher() public {
        uint256 initialBalance = address(mockApplication).balance;
        uint256 recipientInitialBalance = recipient.balance;

        bytes memory encodedWithdrawTx = abi.encodeCall(EmergencyWithdraw.emergencyETHWithdraw, (recipient));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(emergencyWithdraw), encodedWithdrawTx));

        mockApplication.executeOutput(delegateCallVoucher);

        assertEq(address(mockApplication).balance, 0);
        assertEq(recipient.balance, recipientInitialBalance + initialBalance);
    }

    function test_EmergencyERC20WithdrawWithZeroBalance() public {
        bytes memory encodedWithdrawTx = abi.encodeCall(EmergencyWithdraw.emergencyERC20Withdraw, (token, recipient));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(emergencyWithdraw), encodedWithdrawTx));
        mockApplication.executeOutput(delegateCallVoucher);

        mockApplication.executeOutput(delegateCallVoucher);

        assertEq(token.balanceOf(address(mockApplication)), 0);
        assertEq(token.balanceOf(recipient), 1000e18);
    }

    function test_EmergencyETHWithdrawWithZeroBalance() public {
        bytes memory encodedWithdrawTx = abi.encodeCall(EmergencyWithdraw.emergencyETHWithdraw, (recipient));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(emergencyWithdraw), encodedWithdrawTx));
        mockApplication.executeOutput(delegateCallVoucher);

        vm.expectRevert("No ETH to withdraw");
        mockApplication.executeOutput(delegateCallVoucher);
    }
}
