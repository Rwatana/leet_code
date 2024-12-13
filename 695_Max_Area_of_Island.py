# n = len(grid)
# m = len(grid[0])
# time: O(n*m)
# space: O(1)
from typing import List


class Solution:
    def maxAreaOfIsland(self, grid: List[List[int]]) -> int:
        row, col = len(grid), len(grid[0])
        max_area = 0

        def count_the_width_of_island(cordinate_x: int, cordinate_y: int) -> int:
            if {0 <= cordinate_x < len(grid)
                    and 0 <= cordinate_y < len(grid[0])
                    and grid[cordinate_x][cordinate_y] == 1}:
                grid[cordinate_x][cordinate_y] = 0
                area = 1
                area += count_the_width_of_island(cordinate_x + 1, cordinate_y)
                area += count_the_width_of_island(cordinate_x - 1, cordinate_y)
                area += count_the_width_of_island(cordinate_x, cordinate_y + 1)
                area += count_the_width_of_island(cordinate_x, cordinate_y - 1)
                return area

            return 0
        
        for i in range(row):
            for j in range(col):
                if grid[i][j] == 1:
                    num_area = count_the_width_of_island(i, j)
                    max_area = max(max_area, num_area)
        
        return max_area
