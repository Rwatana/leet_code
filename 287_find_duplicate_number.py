# time complexity O(nlogn)
# space complexity O(1)
from typing import List


class Solution:
    def findDuplicate(self, nums: List[int]) -> int:
        n = len(nums) - 1
        left, right = 1, n
        while left < right:
            mid = (right + left) // 2
            count = 0
            for num in nums:
                if num <= mid:
                    count += 1
            if count > mid:
                right = mid
            else:
                left = mid + 1
        return right
