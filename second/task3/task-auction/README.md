
1：delegatecall 跟 call 的区别是什么
主要区别是合约的上下文不同。调用链是:A->B->C的情况下，那么call时，B和C合约的上下文对应各自合约；delegatecall时，B和C合约的上下文都是B合约的。

2：可升级合约的执行流程是什么（user -> proxy -> implementation）
1）、升级，user调用proxy合约的升级函数，将implementation合约地址状态变量修改为新逻辑合约地址；
2）、调用：
a、user通过proxy合约调用implementation合约的函数，
b、由于proxy合约上没有对应的函数签名，会执行proxy的fallback函数，
c、proxy的fallback函数中使用delegatecall调用implementation合约的函数。

3：代理合约上本身是有存储的，怎么避免跟逻辑合约上的存储产生冲突
1）、保证代理合约和逻辑合约的存储槽位保持一致，
2）、可以采用预留的方式设置状态变量的类型大小，防止合约升级时的槽位冲突。
3）、基于ERC1967标准处理。

4： 逻辑合约升级的存储冲突问题
1）、delegatecall使得被调用的合约和调用合约使用的是同一个上下文，
2）、EVM的槽位分配规则是按顺序从0开始的，
3）、当逻辑合约的槽位和代理合约的槽位不同时，逻辑合约会将代理合约槽位上的数据破坏。

为了解决这个问题，有两种常见的方法：
方法一：不使用常规的 Solidity 存储布局机制来存储代理数据，而是利用 EVM 的 sstore 和 sload 指令，在伪随机数字槽位中读写数据。例如，通过 keccak256(my.proxy.version) 这样的哈希函数返回的槽位来存储数据，从而避免冲突。
方法二：使用相同的存储布局，并结合高级的数据争议解决方法。

5： 可以在逻辑合约的构造函数中初始化变量吗？为什么
不可以。
原因：
1）、逻辑合约处理的变量是在代理合约的上下文中的，构造函数初始化的变量是逻辑合约自身的，和代理合约无关，
2）、逻辑合约的构造函数不会在代理合约的上下文中执行，
3）、逻辑合约部署时会初始化变量，使得代理合约和逻辑合约的槽位发生冲突。


# Sample Hardhat Project

This project demonstrates a basic Hardhat use case. It comes with a sample contract, a test for that contract, and a Hardhat Ignition module that deploys that contract.

Try running some of the following tasks:

```shell
npx hardhat help
npx hardhat test
REPORT_GAS=true npx hardhat test
npx hardhat node
npx hardhat ignition deploy ./ignition/modules/Lock.js
```
代理合约地址： 0x0E8C19C974D40b987B89aABDf811d5F31781Cd7C

逻辑合约地址： 0xa31F0eCdeD6FdF1126b071d02f547462F8B0cc27

项目结构：
task-auction/
├── contracts/
│   ├── NftAuction.sol         # 拍卖合约
│   └── test/TestERC721.sol    # NTF合约
├── deploy/
|   ├── 00_deploy_nft_auction.js  # 部署脚本
├── test/
│   └── NftAuction.js  # 测试文件
├── .env                      # Infura key
├── hardhat.config.js         # 配置
├── package.json              # 依赖
└── README.md                 # 说明
