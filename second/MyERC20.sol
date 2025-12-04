// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MyERC20 is IERC20{
    address public owner; // 合约所有者
    mapping(address => uint256) public override balanceOf; // 账户余额(balanceOf())
    mapping(address => mapping(address => uint256)) public override allowance; //`owner`账户授权给`spender`账户的额度 
    uint256 public override totalSupply; // the value of tokens in existence.
    
    string public name; // 名称
    string public symbol;// 代号

    constructor(string memory _name,string memory _symbol){
        name = _name;
        symbol = _symbol;
        owner = msg.sender;
    }

    // 直接转账
    function transfer(address to, uint256 amount) external override returns (bool) {    
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        emit Transfer(msg.sender,to,amount);
        return true;
    }

    // 授权额度
    function approve(address spender, uint256 amount) external override returns (bool) {
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }
    
    // 基于额度转账
    function transferFrom(address from, address to, uint256 amount) external override returns (bool) {
        require(allowance[from][to] >= amount,"Insufficient allowance");
        allowance[from][to] -= amount;
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        emit Transfer(from, to, amount);
        return true;
    }

    // 铸币
    function mint(address to,uint256 amount) external {
        require(msg.sender == owner, "NEED OWNER");
        balanceOf[to] += amount;
        totalSupply += amount;
        emit Transfer(address(0), to, amount);
    }

}