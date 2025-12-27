// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Outputs} from "cartesi-rollups-contracts-2.0.0/src/common/Outputs.sol";
import {ReentrancyGuard} from "@openzeppelin-contracts-5.2.0/utils/ReentrancyGuard.sol";
import {LibError} from "cartesi-rollups-contracts-2.0.0/src/library/LibError.sol";
import {LibAddress} from "cartesi-rollups-contracts-2.0.0/src/library/LibAddress.sol";

contract MockApplication is ReentrancyGuard {
    using LibError for bytes;
    using LibAddress for address;

    error OutputNotExecutable(bytes output);
    error InsufficientFunds(uint256 value, uint256 balance);

    function executeOutput(bytes calldata output) external nonReentrant {
        bytes4 selector = bytes4(output[:4]);
        bytes calldata arguments = output[4:];

        if (selector == Outputs.Voucher.selector) {
            _executeVoucher(arguments);
        } else if (selector == Outputs.DelegateCallVoucher.selector) {
            _executeDelegateCallVoucher(arguments);
        } else {
            revert OutputNotExecutable(output);
        }
    }

    function _executeVoucher(bytes calldata arguments) internal {
        address destination;
        uint256 value;
        bytes memory payload;

        (destination, value, payload) = abi.decode(arguments, (address, uint256, bytes));

        bool enoughFunds;
        uint256 balance;

        (enoughFunds, balance) = destination.safeCall(value, payload);

        if (!enoughFunds) {
            revert InsufficientFunds(value, balance);
        }
    }

    function _executeDelegateCallVoucher(bytes calldata arguments) internal {
        address destination;
        bytes memory payload;

        (destination, payload) = abi.decode(arguments, (address, bytes));

        destination.safeDelegateCall(payload);
    }
}
