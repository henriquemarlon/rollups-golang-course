// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {NFT} from "./NFT.sol";

contract NFTFactory {
    event NFTDeployed(address indexed nft, bytes32 salt);
    event NFTAlreadyDeployed(address indexed nft, bytes32 salt);

    function newNFT(address initialOwner, bytes32 salt, string memory name, string memory symbol)
        external
        returns (NFT)
    {
        address predicted = _computeAddress(initialOwner, salt, name, symbol);

        if (predicted.code.length > 0) {
            emit NFTAlreadyDeployed(predicted, salt);
            return NFT(predicted);
        }

        NFT nft = new NFT{salt: salt}(initialOwner, name, symbol);

        emit NFTDeployed(address(nft), salt);
        return nft;
    }

    function _computeAddress(address initialOwner, bytes32 salt, string memory name, string memory symbol)
        private
        view
        returns (address)
    {
        return address(
            uint160(
                uint256(
                    keccak256(
                        abi.encodePacked(
                            bytes1(0xff),
                            address(this),
                            salt,
                            keccak256(abi.encodePacked(type(NFT).creationCode, abi.encode(initialOwner, name, symbol)))
                        )
                    )
                )
            )
        );
    }
}
