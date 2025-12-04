// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";

contract WNFT is ERC721{
    
    uint public MAX_APES = 10000; // 总量
    string public IPFS;
    constructor(string memory name, string memory symbol) ERC721 (name,symbol){} 

    // Base URI
    function _baseURI() internal view override returns (string memory) {
        string memory _ipfs = string.concat("ipfs://",IPFS,"/");
        return _ipfs;
    }
    /**
    bafkreiebo3k4z7dufrjdk4sga4vxeofr7iw3eej5irgwhmgjqgdq4knuam
    {
    "description": "Friendly OpenSea Creature that enjoys long swims in the ocean.", 
    "external_url": "https://openseacreatures.io/3", 
    "image": "https://tomato-recent-llama-633.mypinata.cloud/ipfs/bafybeida7zerku5esolwjrui4hmgbi2inhqzph5kybfimrefgy3bwg7rze", 
    "name": "Dave Starbelly",
    "attributes": [ 
        {
        "trait_type": "Eyes", 
        "value": "Big"
        }, 
        {
        "trait_type": "Mouth", 
        "value": "Confused"
        }, 
        {
        "trait_type": "Level", 
        "value": 5
        }
    ], 
    }
    */
    // 铸币
    function mintNFT(address to,uint256 tokenId,string memory ipfs) external {
        require(tokenId >= 0 && tokenId < MAX_APES, "tokenId out of range");
        IPFS = ipfs;
        _mint(to,tokenId);
        tokenURI(tokenId);
    }

    // 销毁
    function burnNFT(uint256 tokenId) external {
        require(msg.sender == ownerOf(tokenId),"Not Owner");
        _burn(tokenId);
    }
}