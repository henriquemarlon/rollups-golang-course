// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

interface IERC721Mintable {
    function safeMint(address to, string memory uri) external;
}

contract SafeERC721Mint {
    error NotAContract(address target);

    function safeMint(IERC721Mintable nft, address to, string memory uri) public {
        if (address(nft).code.length == 0) {
            revert NotAContract(address(nft));
        }

        nft.safeMint(to, uri);
    }
}
