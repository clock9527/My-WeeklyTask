const { deployments, upgrades, ethers } = require("hardhat");

const fs = require("fs");
const path = require("path");

module.exports = async ({ getNamedAccounts, deployments }) => {

    const { save } = deployments;
    const { deployer } = await getNamedAccounts();

    // 通过可升级代理，部署合约
    const nftAuctionContract = await ethers.getContractFactory("NftAuction");
    const nftAuctionContractProxy = await upgrades.deployProxy(nftAuctionContract, [], { initializer: "initialize", });
    await nftAuctionContractProxy.waitForDeployment();
    const proxAddress = await nftAuctionContractProxy.getAddress();
    console.log("代理合约地址：", proxAddress);
    const implAddress = await upgrades.erc1967.getImplementationAddress(proxAddress);
    console.log("逻辑合约地址：", implAddress);

    //
    const storePath = path.resolve(__dirname, "./.cache/proxyNftAuction.json");
    fs.writeFileSync(storePath, JSON.stringify({
        proxAddress,
        implAddress,
        abi: nftAuctionContract.interface.format("json")
    }));
    await save("NftAuctionProxy", {
        abi: nftAuctionContract.interface.format("json"),
        address: proxAddress
    });

};

module.exports.tags = ['deployNftAuction'];