// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

contract Voting{
    address[] allCandidates; //候选人
    mapping(address => uint256) public allVotes;// 用户投票给某个候选人

    // 用户投票给某个候选人
    function vote(address candidate) external {
        allCandidates.push(candidate);
        allVotes[candidate] += 1;
    }

    // 重置 mapping
    function resetVotes() external{
        require(allCandidates.length > 0 , "Mapping is null");
        // clear mapping
        for (uint i = 0;i < allCandidates.length; i++){
            delete allVotes[allCandidates[i]];
        }
    }

    // 2.	反转字符串 (Reverse String)
    function reverseString(string memory str) external pure returns (string memory){
        bytes memory strBytes = bytes(str);
        uint len = strBytes.length;
        for(uint i = 0; i < len/2; i++){
            (strBytes[i],strBytes[len-i-1]) = (strBytes[len-i-1],strBytes[i]);
        }
        return string(strBytes);
    }

    mapping(uint8=>bytes2) mapR;
    // 3.	用 solidity 实现整数转罗马数字
    function integerToRoman(uint num) external returns(string memory){
        require(num>0 && num<=3999," 0< n <= 3999");
        mapR[1] = bytes2("IV");
        mapR[2] = bytes2("XL");
        mapR[3] = bytes2("CD");
        mapR[4] = bytes2("M");
        bytes memory res = new bytes(15);
        uint i = 0;
        // 千位 
        while(num > 1000){res[i++] = 'M';num -= 1000;}

        // 百位
        if( num >= 900){res[i++] = 'C';res[i++] = 'M';num -= 900;}
        if( num >= 500){res[i++] = 'D';num -= 500;}
        if( num >= 400){res[i++] = 'C';res[i++] = 'D';num -= 400;}
        while(num >= 100){res[i++] = 'C';num -= 100;}

        // 十位
        if( num >= 90){res[i++] = 'X';res[i++] = 'C';num -= 90;}
        if( num >= 50){res[i++] = 'L';num -= 50;}
        if( num >= 40){res[i++] = 'X';res[i++] = 'L';num -= 40;}
        while(num >= 10){res[i++] = 'X';num -= 10;}

        // 个位
        if( num >= 9){res[i++] = 'I';res[i++] = 'X';num -= 9;}
        if( num >= 5){res[i++] = 'V';num -= 5;}
        if( num >= 4){res[i++] = 'I';res[i++] = 'V';num -= 4;}
        while(num >= 1){res[i++] = 'I';num -= 1;}
        return string(res);
    }

    // 4.	用 solidity 实现罗马数字转数整数
    function romanToInteger(string memory s) external pure returns(int num){
        bytes memory bStr = bytes(s);
        uint len = bStr.length;
        for(uint i=0; i < len; i++){
            if (i == len-1){
                if(bStr[i] == 'I'){num += 1;}
                if(bStr[i] == 'V'){num += 5;}
                if(bStr[i] == 'X'){num += 10;}
                return num;
            }
            if (bStr[i] == 'M'){num += 1000;}
            // 900 400
            if (bStr[i] == 'C' && (bStr[i+1] == 'M'||bStr[i+1] == 'D')){
                num -= 100;
            }else if(bStr[i] == 'C'){
                num += 100;
            }

            if (bStr[i] == 'D') num += 500;

            // 90 40
            if (bStr[i] == 'X' && (bStr[i+1] == 'C'||bStr[i+1] == 'L')){
                num -= 10;
            }else if(bStr[i] == 'X'){
                num += 10;
            }

            if (bStr[i] == 'L') num += 50;
            
            // 9 4
            if (bStr[i] == 'I' && (bStr[i+1] == 'X'||bStr[i+1] == 'V')){
                num -= 1;
            }else if(bStr[i] == 'I'){
                num += 1;
            }
            if (bStr[i] == 'V') num += 5;
        }
    }
    
    // 5.	合并两个有序数组 (Merge Sorted Array)
    function mergeArray(uint128[] memory arr1,uint128[] memory arr2) external pure returns(uint256[] memory){
        uint len1 = arr1.length;
        uint len2 = arr2.length;
        uint256[] memory mArr = new uint256[](len1 + len2);
        uint i = 0;
        uint j = 0;
        uint k = 0;
        while(i < len1 && j < len2){
            mArr[k++] = arr1[i] <= arr2[j] ? arr1[i++]:arr2[j++];
        }
        // 将两个数组剩下的写入目标数组
        while(i < len1){
            mArr[k++] = arr1[i++];
        }
        while(j < len2){
            mArr[k++] = arr1[j++];
        }
        return mArr;
    }

    // 6.二分查找 (Binary Search)
    function binarySearch(uint[] memory arr,uint dest) external pure returns(uint,bool){
        uint len;
        require((len = arr.length) > 0, "Arrays is null");

        uint left = 0;
        uint right = len-1;
        
        while(left <= right){
            uint idx = left + (right - left)/2;

            if (arr[idx] == dest){
                return (idx,true);
            }else if (arr[idx] < dest){
                left  = idx + 1;
            }else{
                if (idx == 0) break; // 已经到最左边
                right = idx - 1;
            }
        }
        return (0,false);
    }

    // 6.二分查找 (Binary Search)。 deepseek 深度优化版本（Gas 效率最高）
    function binarySearchOptimized(uint[] memory arr, uint dest) external pure returns (uint256, bool) {
        assembly {
            // arr 参数在内存中的布局： [长度, 元素1, 元素2, ...]
            let length := mload(arr)          // 获取数组长度
            if iszero(length) {
                // 空数组直接返回
                mstore(0, 0)                  // 索引 = 0
                mstore(0x20, 0)               // bool = false
                return(0, 0x40)               // 返回64字节
            }
            
            let left := 0
            // 注意：right = length - 1，但length是uint，需要检查length>0
            let right := sub(length, 1)
            
            // 使用内联汇编进行二分查找
            for {} iszero(gt(left, right)) {} {
                // 计算中间索引：mid = left + (right - left) / 2
                let diff := sub(right, left)
                let mid := add(left, shr(1, diff)) // shr(1, diff) = diff / 2
                
                // 获取 arr[mid] 的值
                // 每个元素占32字节，起始位置是 arr + 0x20
                let midValue := mload(add(arr, add(0x20, mul(mid, 0x20))))
                
                switch eq(midValue, dest)
                case 1 {
                    // 找到目标，返回结果
                    mstore(0, mid)            // 索引
                    mstore(0x20, 1)           // true
                    return(0, 0x40)           // 返回64字节
                }
                default {
                    // 未找到，调整边界
                    if lt(midValue, dest) {
                        // dest > arr[mid]，搜索右半部分
                        left := add(mid, 1)
                    }
                    // dest < arr[mid]，搜索左半部分
                    // 需要检查 mid > 0 防止下溢
                    if gt(mid, 0) {
                        right := sub(mid, 1)
                    } {
                        // mid == 0，直接结束循环
                        break
                    }
                }
            }
            
            // 未找到
            mstore(0, 0)  // 索引 = 0
            mstore(0x20, 0) // bool = false
            return(0, 0x40)
        }
    }
}