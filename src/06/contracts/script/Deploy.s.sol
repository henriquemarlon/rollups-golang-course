pragma solidity ^0.8.27;

import {Script, console} from "forge-std-1.9.7/src/Script.sol";
import {SafeERC721Mint} from "../src/delegatecall/SafeERC721Mint.sol";
import {EmergencyWithdraw} from "../src/delegatecall/EmergencyWithdraw.sol";
import {SafeERC20Transfer} from "../src/delegatecall/SafeERC20Transfer.sol";
import {Token} from "../src/token/ERC20/Token.sol";
import {NFTFactory} from "../src/token/ERC721/NFTFactory.sol";

contract Deploy is Script {
    SafeERC721Mint public safeERC721Mint;
    EmergencyWithdraw public emergencyWithdraw;
    SafeERC20Transfer public safeERC20Transfer;

    Token public token;

    NFTFactory public nftFactory;

    function run() public returns (SafeERC721Mint, EmergencyWithdraw, SafeERC20Transfer, Token, NFTFactory) {
        vm.startBroadcast();

        safeERC721Mint = new SafeERC721Mint{salt: keccak256("1596")}();
        emergencyWithdraw = new EmergencyWithdraw{salt: keccak256("1596")}();
        safeERC20Transfer = new SafeERC20Transfer{salt: keccak256("1596")}();
        nftFactory = new NFTFactory{salt: keccak256("1596")}();

        token = new Token();

        vm.stopBroadcast();

        _saveDeploymentInfo();

        return (safeERC721Mint, emergencyWithdraw, safeERC20Transfer, token, nftFactory);
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
            '"safeERC721Mint":"',
            vm.toString(address(safeERC721Mint)),
            '","emergencyWithdraw":"',
            vm.toString(address(emergencyWithdraw)),
            '","safeERC20Transfer":"',
            vm.toString(address(safeERC20Transfer)),
            '","token":"',
            vm.toString(address(token)),
            '","nftFactory":"',
            vm.toString(address(nftFactory)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, string.concat("./deployments/", vm.toString(block.chainid), ".json"));
    }
}
