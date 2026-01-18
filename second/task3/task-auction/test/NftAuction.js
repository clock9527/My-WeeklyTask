const hre = require("hardhat") // Hardhat Runtime Enviroment
const { expect } = require("chai");

describe("Test aunction", async () => {
    const { ethers, upgrades, deployments } = hre;
    it("Should be OK", async () => {
        [deployer, bider] = await ethers.getSigners();

        // 1、部署NFT合约
        const NFTFactory = await ethers.getContractFactory("TestERC721");
        const NFTContract = await NFTFactory.deploy();
        await NFTContract.waitForDeployment();
        const NFTAddress = await NFTContract.getAddress();
        console.log("NFT合约地址", NFTAddress);

        // 铸造10个NFT
        for (let i = 0; i < 10; i++) {
            NFTContract.mint(deployer.address, i + 1);
        }

        // 2、部署Auction合约
        await deployments.fixture("deployNftAuction");
        const NftAuctionProxy = await deployments.get("NftAuctionProxy");

        // 给代理合约赋权
        await NFTContract.connect(deployer).setApprovalForAll(NftAuctionProxy.address, true);

        // 3、创建拍卖
        const tokenId = 1;
        const nftAuction = await ethers.getContractAt("NftAuction", NftAuctionProxy.address);
        await nftAuction.createAuction(
            10,
            ethers.parseEther("0.01"),
            NFTAddress,
            tokenId);
        const currAuction = await nftAuction.auctions(0);
        console.log("创建拍卖成功：：", currAuction);

        // 4、购买者参与拍卖
        const bidPrice = ethers.parseEther("0.0101");
        await nftAuction.connect(bider).placeBid(0, { value: bidPrice });

        const auctionResult0 = await nftAuction.auctions(0);
        console.log("结束拍卖后读取拍卖0：：", auctionResult0);

        // 5、等待8s后结束拍卖
        await new Promise((resolve) => setTimeout(resolve, 8 * 1000));
        await NftAuction.endAuction(0);

        // 6、验证结果
        // 拍卖情况
        const auctionResult1 = await NftAuction.auctions(0);
        console.log("结束拍卖后读取拍卖1：：", auctionResult1);
        expect(auctionResult1.highestBidder).to.equal(bider.address);
        expect(auctionResult1.highestBid).to.equal(bidPrice);

        // 验证 NFT 所有权
        const owner = await NFTContract.ownerOf(tokenId);
        console.log("owner::", owner);
        expect(owner).to.equal(bider.address);

    });
})