// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std-1.9.7/src/Script.sol";
import {NFTFactory} from "../src/token/ERC721/NFTFactory.sol";

contract Deploy is Script {
    NFTFactory public nftFactory;

    function run() public returns (NFTFactory) {
        vm.startBroadcast();
        nftFactory = new NFTFactory{salt: keccak256("1596")}();
        vm.stopBroadcast();

        _saveDeploymentInfo();

        return nftFactory;
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"deployer":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"nftFactory":"',
            vm.toString(address(nftFactory)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, string.concat("./deployments/", vm.toString(block.chainid), ".json"));
    }
}
