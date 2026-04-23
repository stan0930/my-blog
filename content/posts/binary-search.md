+++
title = "二分查找怎么写才不出错：红蓝染色法 + LeetCode 34"
date = 2026-04-23T11:20:00+08:00
draft = true
description = "用红蓝染色法理解二分查找，稳定写出边界正确的模板。"
summary = "从循环不变量到 lower_bound，再到 LeetCode 34，一次讲透二分边界。"
tags = ["算法", "二分", "C++", "LeetCode"]
categories = ["算法"]
series = ["算法基础"]
ShowToc = true
TocOpen = false
+++

## 写在前面

这篇是把我和 ChatGPT 聊二分时的关键点整理成一篇可复用笔记，主线就三件事：

1. 二分的本质不是“猜数字”，而是找分界点
2. 用红蓝染色法维护循环不变量
3. 用统一模板处理边界题（特别是 LeetCode 34）

## 1. 为什么二分总写错

常见报错点不是复杂逻辑，而是边界细节：

1. `while (l <= r)` 和 `while (l < r)` 混用
2. `mid` 取值后区间收缩不一致，导致死循环
3. 题目要“第一个/最后一个”，却写成“是否存在”

这类问题的统一解法是：先定义判定条件，再固定循环不变量。

## 2. 红蓝染色法：先定义真假，再二分

这一版我按视频里的习惯来记颜色：

- 红色：不满足条件（`false`）
- 蓝色：满足条件（`true`）
- 白色：当前还没判断，答案一定藏在白色区间里

对于“找第一个 `>= target`”这类题，数组会被染成：

- 左边一段红色：`< target`
- 右边一段蓝色：`>= target`

所以二分本质上是在找“第一块蓝色”。

这比“在数组里找目标值”更通用。很多题都能转成：

- 找第一个满足 `check(i)` 的位置
- 或找最后一个不满足 `check(i)` 的位置

### 2.1 例子流程图：闭区间 `[L, R]`

拿视频里最经典的例子：

- `nums = [5, 7, 7, 8, 8, 10]`
- `target = 8`
- `check(i) = nums[i] >= 8`

于是颜色就是：

- `5, 7, 7` 是红色，因为它们 `< 8`
- `8, 8, 10` 是蓝色，因为它们 `>= 8`

图里用的是闭区间思路：`[L, R]` 表示“答案还可能出现的范围”，`M` 是当前看的位置。

![二分查找闭区间流程图](../../images/binary-search-closed-interval.png)

按图里的流程走一遍：

1. 初始 `L = 0`，`R = 5`
2. 第一次二分，`M = 2`，`nums[2] = 7` 是红色，说明答案不在 `[0, 2]`，所以更新 `L = M + 1`
3. 第二次二分，`M = 4`，`nums[4] = 8` 是蓝色，说明答案在左边，更新 `R = M - 1`
4. 第三次二分，`M = 3`，`nums[3] = 8` 还是蓝色，继续更新 `R = M - 1`
5. 最后 `L = 3`，`R = 2`，循环结束，此时答案是 `R + 1 = 3`

### 2.2 循环不变量：为什么答案是 `R + 1`

这是这套写法最关键的知识点。

在整个二分过程中，我们始终维护：

1. `L - 1` 一定是红色，也就是 `< target`
2. `R + 1` 一定是蓝色，也就是 `>= target`
3. `[L, R]` 是还没有确定颜色的白色区间

所以每做一次更新，本质上都在缩小白色区间：

1. 如果 `mid` 是红色，说明 `[L, mid]` 都可以排除，更新 `L = mid + 1`
2. 如果 `mid` 是蓝色，说明 `[mid, R]` 都可以排除，更新 `R = mid - 1`

当循环结束时，`L > R`，白色区间已经空了。  
此时红蓝分界点正好在 `R` 和 `R + 1` 之间，所以：

```cpp
答案 = R + 1
```

又因为结束时一定有 `L = R + 1`，所以很多代码里也直接返回 `L`。  
但如果你想把“循环不变量”这个逻辑看得更清楚，`return R + 1` 更直观。

## 3. 最稳模板：找第一个满足条件的位置

如果你更喜欢左右都闭区间，我建议代码直接固定成这一版。它的好处是：

1. `l` 和 `r` 都是真实下标
2. 直接兼容“所有数都小于 `target`”的情况
3. 写 LeetCode 34 这种边界题最稳

```cpp
int lower_bound_closed(const vector<int>& nums, int target) {
    int l = 0, r = (int)nums.size() - 1;
    while (l <= r) {
        int mid = l + (r - l) / 2;
        if (nums[mid] >= target) {
            r = mid - 1;   // mid 是蓝色，答案在左边或就是 mid
        } else {
            l = mid + 1;   // mid 是红色，答案只能在右边
        }
    }
    return r + 1;          // 也等于 l
}
```

这个模板和上面的流程图是一一对应的：

1. 遇到红色，`L` 右移
2. 遇到蓝色，`R` 左移
3. 最后用循环不变量推出答案是 `R + 1`

## 4. LeetCode 34：第一个和最后一个位置

题意：在有序数组中找目标值 `target` 的起止下标。

核心转化：

1. `left = lowerBoundClosed(nums, target)`
2. `right = lowerBoundClosed(nums, target + 1) - 1`

如果 `left` 越界或 `nums[left] != target`，说明不存在。

```cpp
class Solution {
public:
    int lowerBoundClosed(const vector<int>& nums, int target) {
        int l = 0, r = (int)nums.size() - 1;
        while (l <= r) {
            int mid = l + (r - l) / 2;
            if (nums[mid] >= target) {
                r = mid - 1;
            } else {
                l = mid + 1;
            }
        }
        return r + 1;
    }

    vector<int> searchRange(vector<int>& nums, int target) {
        int left = lowerBoundClosed(nums, target);
        if (left == (int)nums.size() || nums[left] != target) {
            return {-1, -1};
        }
        int right = lowerBoundClosed(nums, target + 1) - 1;
        return {left, right};
    }
};
```

时间复杂度 `O(log n)`，空间复杂度 `O(1)`。

## 5. 答案二分：把值域当作“下标”

很多题不是在数组里二分，而是在答案范围二分。例如：

- 最小化最大值
- 最大化最小值

模板完全一样，只是 `check(mid)` 变成“这个答案是否可行”：

```cpp
long long solve() {
    long long l = 0, r = 1e12;
    while (l <= r) {
        long long mid = l + (r - l) / 2;
        if (check(mid)) {
            r = mid - 1;    // mid 可行，答案在左边
        } else {
            l = mid + 1;
        }
    }
    return r + 1;          // 也等于 l
}
```

关键不是二分本身，而是 `check` 必须单调。

这里有个前提：你给定的搜索区间右端，必须已经在“蓝色区域”右侧，也就是答案一定存在于这个区间内。  
这样循环结束后，`R + 1` 才有明确含义。

## 6. 高频坑位清单

1. `mid = (l + r) / 2` 可能溢出，改成 `l + (r - l) / 2`
2. 忘记区分区间类型（`[l, r]` vs `[l, r)`）
3. 闭区间版里，`nums[mid] >= target` 后应该收成 `r = mid - 1`
4. 题目要边界，你却只判断“存在”
5. `check` 不单调，还强行二分
6. 写完不测极端样例：空数组、全相同、目标不存在

## 7. 一句话记忆

先把条件染成红蓝，再去找第一块蓝色。

二分代码只是外壳，真正决定对错的是：

1. 判定条件是否单调
2. 循环不变量是否被每一步维护

## 小结

如果你也经常“觉得会写，但总差一行”，建议固定一套模板长期复用：

1. 统一用一套闭区间模板
2. 统一先写“找第一个 `>= target`”
3. 题目都往“找第一块蓝色”转

这样二分会从“玄学”变成“机械正确”。
