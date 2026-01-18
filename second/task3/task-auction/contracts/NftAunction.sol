// SPDX-License-Identifier:MIT
pragma solidity ^0.8.22;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "hardhat/console.sol";

contract NftAuction is Initializable, UUPSUpgradeable {
    //
    struct Auction {
        address seller; // 卖家
        uint256 startTime; // 开始时间
        uint256 duration; // 持续时间
        bool ended; //是否结束
        uint256 startPrice; // 起拍价
        address highestBidder; //最高出价者
        uint256 highestBid; // 最高价
        address nftAddr; // token 地址
        uint256 tokenID; // token ID
        address payTokenAddress; // 0x:ETH，其他:ERC20
    }

    Auction auction; //拍卖
    uint256 nextAuctionID; //下一个拍卖ID
    mapping(uint256 => Auction) public auctions; //拍卖ID => 拍卖
    address admin; // 管理员
    mapping(address => AggregatorV3Interface) priceFeeds;

    function initialize() public initializer {
        admin = msg.sender;
    }

    /*
        ETH / USD  : 0x694AA1769357215DE4FAC081bf1f309aDC325306
        USDC / USD : 0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E
    */
    function setAggV3s(address _payTokenAddress, address _priceETHFeed) public {
        priceFeeds[_payTokenAddress] = AggregatorV3Interface(_priceETHFeed);
    }

    function getChainLinkDataFeedLatestAnswer(
        address _payTokenAddress
    ) public view returns (int) {
        AggregatorV3Interface priceFeed = priceFeeds[_payTokenAddress];

        (
            ,
            /* uint80 roundId */ int256 answer /*uint256 startedAt*/ /*uint256 updatedAt*/ /*uint80 answeredInRound*/,
            ,
            ,

        ) = priceFeed.latestRoundData();
        return answer;
    }

    // 创建拍卖
    function createAuction(
        uint256 _duration,
        uint256 _startPrice,
        address _nftAddr,
        uint256 _tokenID
    ) public {
        require(msg.sender == admin, "Must be admin.");
        require(_duration >= 10, "Duration must be greater than 10s.");
        require(_startPrice > 0, "Start price must be greater than 0.");

        IERC721(_nftAddr).approve(msg.sender, _tokenID); // 将token授权给拍卖方
        auction.seller = msg.sender;
        auction.startTime = block.timestamp;
        auction.duration = _duration;
        auction.startPrice = _startPrice;
        auction.nftAddr = _nftAddr;
        auction.tokenID = _tokenID;
        auctions[nextAuctionID] = auction;
        nextAuctionID++;
    }

    // 买家参与拍卖
    function placeBid(
        uint _auctionID,
        uint amount,
        address _payTokenAddress
    ) external payable {
        Auction storage auction = auctions[_auctionID];
        // 判断当前拍卖是否结束
        require(
            !auction.ended &&
                auction.startTime + auction.duration >= block.timestamp,
            "bid: Auction has ended."
        );

        // 判断出价是否大于最高价
        // 当前拍卖的价格汇率转换
        if (_payTokenAddress == address(0)) amount = msg.value; //ETH时 amount=msg.value
        uint payValue = amount *
            uint(getChainLinkDataFeedLatestAnswer(_payTokenAddress));

        // 当前拍卖前的价格汇率转换。
        uint pariValue = uint(
            getChainLinkDataFeedLatestAnswer(auction.payTokenAddress)
        );
        uint startPriceValue = auction.startPrice * pariValue;
        uint hightestBidValue = auction.highestBid * pariValue;

        // 基于汇率转换后的价格进行比较
        require(
            payValue > hightestBidValue && payValue > startPriceValue,
            "Bid must be heigher than current heighest bid"
        );

        // 退回之前的最高价
        if (auction.highestBidder != address(0)) {
            if (auction.payTokenAddress == address(0)) {
                //payToken为0地址，直接退回
                (bool succ, ) = payable(auction.highestBidder).call{
                    value: auction.highestBid
                }("");
                if (succ) {}
            } else {
                //payToken为非0地址， ERC20 退回
                ERC20(auction.payTokenAddress).transfer(
                    auction.highestBidder,
                    hightestBidValue
                );
            }
        }
        auction.payTokenAddress = _payTokenAddress;
        auction.highestBidder = msg.sender;
        auction.highestBid = amount;
        // console.log("endPlaceBid address", auction.highestBidder);
    }

    // 结束拍卖
    function endAuction(uint256 _auctionID) external {
        Auction storage auction = auctions[_auctionID];
        // 判断当前拍卖是否结束
        require(
            !auction.ended &&
                auction.startTime + auction.duration >= block.timestamp,
            "end: Auction has ended."
        );

        // 转移NFT到最高价者
        IERC721(auction.nftAddr).safeTransferFrom(
            admin,
            auction.highestBidder,
            auction.tokenID
        );

        // 转移剩余的资金到卖家
        // payable(address(this)).transfer(address(this).balance);
        (bool buss, ) = payable(address(this)).call{
            value: address(this).balance
        }("");
        if (buss) auction.ended = true;
    }

    //
    function _authorizeUpgrade(
        address newImplementation
    ) internal view override {
        // 只有管理员可以升级合约
        require(msg.sender == admin, "Only admin can upgrade");
    }
}
