+++
title = "二分查找怎么写才不出错：红蓝染色法 + LeetCode 34"
date = 2026-04-23T11:20:00+08:00
draft = false
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

把下标区间想成一条线：

- 蓝色：满足条件（`true`）
- 红色：不满足条件（`false`）

如果颜色分布是单调的（例如 `红红红蓝蓝蓝`），那么二分就是找“第一块蓝色”。

这比“在数组里找目标值”更通用。很多题都能转成：

- 找第一个满足 `check(i)` 的位置
- 或找最后一个不满足 `check(i)` 的位置

## 3. 最稳模板：找第一个满足条件的位置

我们先写一个最小模板：`[l, r)` 左闭右开。

```cpp
int lower_bound_custom(const vector<int>& a, int target) {
    int l = 0, r = (int)a.size(); // [l, r)
    while (l < r) {
        int mid = l + (r - l) / 2;
        if (a[mid] >= target) {
            r = mid;      // mid 是蓝色，收右边界
        } else {
            l = mid + 1;  // mid 是红色，排除 mid
        }
    }
    return l; // 第一个 >= target 的下标
}
```

循环不变量是：

1. `[0, l)` 一定是红色（`< target`）
2. `[r, n)` 一定是蓝色（`>= target`）

当 `l == r`，分界点就是答案。

## 4. LeetCode 34：第一个和最后一个位置

题意：在有序数组中找目标值 `target` 的起止下标。

核心转化：

1. `left = lower_bound(nums, target)`
2. `right = lower_bound(nums, target + 1) - 1`

如果 `left` 越界或 `nums[left] != target`，说明不存在。

```cpp
class Solution {
public:
    int lowerBound(const vector<int>& nums, int target) {
        int l = 0, r = (int)nums.size();
        while (l < r) {
            int mid = l + (r - l) / 2;
            if (nums[mid] >= target) r = mid;
            else l = mid + 1;
        }
        return l;
    }

    vector<int> searchRange(vector<int>& nums, int target) {
        int left = lowerBound(nums, target);
        if (left == (int)nums.size() || nums[left] != target) {
            return {-1, -1};
        }
        int right = lowerBound(nums, target + 1) - 1;
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
    long long l = 0, r = 1e12; // 假设答案区间
    while (l < r) {
        long long mid = l + (r - l) / 2;
        if (check(mid)) r = mid; // mid 可行，试更小
        else l = mid + 1;        // mid 不可行，增大
    }
    return l;
}
```

关键不是二分本身，而是 `check` 必须单调。

## 6. 高频坑位清单

1. `mid = (l + r) / 2` 可能溢出，改成 `l + (r - l) / 2`
2. 忘记区分区间类型（`[l, r]` vs `[l, r)`）
3. 题目要边界，你却只判断“存在”
4. `check` 不单调，还强行二分
5. 写完不测极端样例：空数组、全相同、目标不存在

## 7. 一句话记忆

先把条件染成红蓝，再去找分界点。

二分代码只是外壳，真正决定对错的是：

1. 判定条件是否单调
2. 循环不变量是否被每一步维护

## 小结

如果你也经常“觉得会写，但总差一行”，建议固定一套模板长期复用：

1. 统一用 `[l, r)`
2. 统一先写 `lower_bound`
3. 题目都往“找第一块蓝色”转

这样二分会从“玄学”变成“机械正确”。
