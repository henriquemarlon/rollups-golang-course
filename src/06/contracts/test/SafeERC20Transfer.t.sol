// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Test} from "forge-std-1.9.7/src/Test.sol";
import {Token} from "../src/token/ERC20/Token.sol";
import {MockApplication} from "./mock/MockApplication.sol";
import {SafeERC20Transfer} from "../src/delegatecall/SafeERC20Transfer.sol";
import {Outputs} from "cartesi-rollups-contracts-2.0.0/src/common/Outputs.sol";

contract SafeERC20TransferTest is Test {
    Token public token;
    SafeERC20Transfer public safeERC20Transfer;
    MockApplication public mockApplication;

    address public user;
    address public recipient;

    event Transfer(address indexed from, address indexed to, uint256 value);

    function setUp() public {
        user = makeAddr("user");
        recipient = makeAddr("recipient");

        token = new Token();
        mockApplication = new MockApplication();
        safeERC20Transfer = new SafeERC20Transfer();

        // Fund the mock application with tokens
        token.mint(address(mockApplication), 1000e18);
    }

    function test_SafeTransferThroughDelegateCallVoucher() public {
        uint256 transferAmount = 100e18;
        uint256 initialBalance = token.balanceOf(address(mockApplication));
        uint256 recipientInitialBalance = token.balanceOf(recipient);

        bytes memory encodedTransferTx =
            abi.encodeCall(SafeERC20Transfer.safeTransfer, (token, recipient, transferAmount));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC20Transfer), encodedTransferTx));

        vm.expectEmit(true, true, false, true);
        emit Transfer(address(mockApplication), recipient, transferAmount);
        mockApplication.executeOutput(delegateCallVoucher);

        assertEq(token.balanceOf(address(mockApplication)), initialBalance - transferAmount);
        assertEq(token.balanceOf(recipient), recipientInitialBalance + transferAmount);
    }

    function test_SafeTransferTargetedThroughDelegateCallVoucher() public {
        uint256 transferAmount = 100e18;
        uint256 initialBalance = token.balanceOf(address(mockApplication));
        uint256 recipientInitialBalance = token.balanceOf(recipient);

        bytes memory encodedTransferTx =
            abi.encodeCall(SafeERC20Transfer.safeTransferTargeted, (token, recipient, recipient, transferAmount));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC20Transfer), encodedTransferTx));

        vm.prank(recipient);
        vm.expectEmit(true, true, false, true);
        emit Transfer(address(mockApplication), recipient, transferAmount);
        mockApplication.executeOutput(delegateCallVoucher);

        assertEq(token.balanceOf(address(mockApplication)), initialBalance - transferAmount);
        assertEq(token.balanceOf(recipient), recipientInitialBalance + transferAmount);
    }

    function test_SafeTransferTargetedWithWrongTarget() public {
        uint256 transferAmount = 100e18;

        bytes memory encodedTransferTx =
            abi.encodeCall(SafeERC20Transfer.safeTransferTargeted, (token, user, recipient, transferAmount));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC20Transfer), encodedTransferTx));

        vm.expectRevert(abi.encodeWithSelector(SafeERC20Transfer.NotTarget.selector, user));
        mockApplication.executeOutput(delegateCallVoucher);
    }

    function test_SafeTransferWithInsufficientBalance() public {
        uint256 transferAmount = 2000e18;

        bytes memory encodedTransferTx =
            abi.encodeCall(SafeERC20Transfer.safeTransfer, (token, recipient, transferAmount));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC20Transfer), encodedTransferTx));

        vm.expectRevert();
        mockApplication.executeOutput(delegateCallVoucher);
    }

    function test_SafeTransferToZeroAddress() public {
        uint256 transferAmount = 100e18;

        bytes memory encodedTransferTx =
            abi.encodeCall(SafeERC20Transfer.safeTransfer, (token, address(0), transferAmount));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC20Transfer), encodedTransferTx));

        vm.expectRevert();
        mockApplication.executeOutput(delegateCallVoucher);
    }

    function test_SafeTransferExactBalance() public {
        uint256 exactBalance = token.balanceOf(address(mockApplication));
        uint256 recipientInitialBalance = token.balanceOf(recipient);

        bytes memory encodedTransferTx =
            abi.encodeCall(SafeERC20Transfer.safeTransfer, (token, recipient, exactBalance));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC20Transfer), encodedTransferTx));

        vm.expectEmit(true, true, false, true);
        emit Transfer(address(mockApplication), recipient, exactBalance);
        mockApplication.executeOutput(delegateCallVoucher);

        assertEq(token.balanceOf(address(mockApplication)), 0);
        assertEq(token.balanceOf(recipient), recipientInitialBalance + exactBalance);
    }
}
