package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/fogleman/gg"
)

const (
	width  = 1500
	height = 1480
)

type step struct {
	title   string
	desc1   string
	desc2   string
	states  []string
	l       int
	r       int
	m       int
	showM   bool
	comment string
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func setFont(dc *gg.Context, fontPath string, size float64) {
	must(dc.LoadFontFace(fontPath, size))
}

func fillRounded(dc *gg.Context, x, y, w, h, radius float64, hex string) {
	dc.DrawRoundedRectangle(x, y, w, h, radius)
	dc.SetHexColor(hex)
	dc.Fill()
}

func strokeRounded(dc *gg.Context, x, y, w, h, radius float64, hex string, line float64) {
	dc.DrawRoundedRectangle(x, y, w, h, radius)
	dc.SetHexColor(hex)
	dc.SetLineWidth(line)
	dc.Stroke()
}

func drawLabel(dc *gg.Context, fontPath string, x, y float64, text string, size float64, hex string) {
	setFont(dc, fontPath, size)
	dc.SetHexColor(hex)
	dc.DrawString(text, x, y)
}

func drawCentered(dc *gg.Context, fontPath string, x, y float64, text string, size float64, hex string) {
	setFont(dc, fontPath, size)
	dc.SetHexColor(hex)
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
}

func drawArrow(dc *gg.Context, x, fromY, toY float64, label, fontPath string) {
	dc.SetHexColor("#4b5563")
	dc.SetLineWidth(2.5)
	dc.DrawLine(x, fromY, x, toY)
	dc.Stroke()

	angle := math.Atan2(toY-fromY, 0)
	size := 9.0
	dc.NewSubPath()
	dc.MoveTo(x, toY)
	dc.LineTo(x-size*math.Cos(angle-math.Pi/6), toY-size*math.Sin(angle-math.Pi/6))
	dc.LineTo(x-size*math.Cos(angle+math.Pi/6), toY-size*math.Sin(angle+math.Pi/6))
	dc.ClosePath()
	dc.Fill()

	drawCentered(dc, fontPath, x, fromY-18, label, 22, "#111827")
}

func drawArray(dc *gg.Context, fontPath string, x, y float64, nums []int, states []string, l, r, m int, showM bool) {
	cellW := 86.0
	cellH := 64.0
	gap := 0.0
	colors := map[string]string{
		"red":   "#e8c0b7",
		"blue":  "#c9d7ef",
		"white": "#f7f2e8",
	}

	for i, v := range nums {
		cx := x + float64(i)*(cellW+gap)
		fillRounded(dc, cx, y, cellW, cellH, 0, colors[states[i]])
		strokeRounded(dc, cx, y, cellW, cellH, 0, "#70685d", 1.6)
		drawCentered(dc, fontPath, cx+cellW/2, y+cellH/2, fmt.Sprintf("%d", v), 28, "#2a2119")
	}

	if l >= 0 && l < len(nums) {
		xp := x + float64(l)*cellW + cellW/2
		drawArrow(dc, xp, y+112, y+cellH+2, "L", fontPath)
	}
	if r >= 0 && r < len(nums) {
		xp := x + float64(r)*cellW + cellW/2
		drawArrow(dc, xp, y+112, y+cellH+2, "R", fontPath)
	}
	if showM && m >= 0 && m < len(nums) {
		xp := x + float64(m)*cellW + cellW/2
		drawArrow(dc, xp, y+160, y+cellH+2, "M", fontPath)
	}
}

func drawBlock(dc *gg.Context, fontPath string, x, y, w, h float64, s step, nums []int) {
	fillRounded(dc, x, y, w, h, 18, "#f7f2e8")
	strokeRounded(dc, x, y, w, h, 18, "#d5c8b4", 1.4)
	drawLabel(dc, fontPath, x+28, y+40, s.title, 30, "#2f241c")
	drawLabel(dc, fontPath, x+28, y+78, s.desc1, 22, "#2f241c")
	drawLabel(dc, fontPath, x+28, y+108, s.desc2, 22, "#2f241c")
	drawLabel(dc, fontPath, x+w-400, y+108, s.comment, 22, "#8f4b4b")
	drawArray(dc, fontPath, x+28, y+140, nums, s.states, s.l, s.r, s.m, s.showM)
}

func drawInvariantBox(dc *gg.Context, fontPath string, x, y, w, h float64, nums []int) {
	fillRounded(dc, x, y, w, h, 18, "#f7f2e8")
	strokeRounded(dc, x, y, w, h, 18, "#d5c8b4", 1.4)
	drawLabel(dc, fontPath, x+28, y+40, "关键：循环不变量", 30, "#2f241c")
	drawLabel(dc, fontPath, x+28, y+78, "1. `L - 1` 始终是红色", 24, "#8f3a3a")
	drawLabel(dc, fontPath, x+28, y+110, "2. `R + 1` 始终是蓝色", 24, "#355f9e")
	drawLabel(dc, fontPath, x+28, y+150, "循环结束时 `L > R`，白色区间为空。", 24, "#2f241c")
	drawLabel(dc, fontPath, x+28, y+182, "红蓝分界点就在 `R` 和 `R + 1` 之间，因此答案 = `R + 1`。", 24, "#2f241c")
	drawLabel(dc, fontPath, x+28, y+214, "又因为结束时 `L = R + 1`，所以答案也可以写成 `L`。", 24, "#2f241c")

	states := []string{"red", "red", "red", "blue", "blue", "blue"}
	drawArray(dc, fontPath, x+w-640, y+68, nums, states, 3, 2, -1, false)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <font-path> <output-path>", os.Args[0])
	}
	fontPath := os.Args[1]
	outputPath := os.Args[2]

	dc := gg.NewContext(width, height)
	dc.SetHexColor("#efe6d6")
	dc.Clear()

	nums := []int{5, 7, 7, 8, 8, 10}

	drawLabel(dc, fontPath, 72, 78, "二分查找闭区间流程图", 42, "#241a14")
	drawLabel(dc, fontPath, 72, 118, "问题：返回有序数组中第一个 >= 8 的位置；如果都 < 8，则返回数组长度", 26, "#2f241c")

	drawLabel(dc, fontPath, 72, 168, "数组：", 24, "#2f241c")
	headerStates := []string{"white", "white", "white", "white", "white", "white"}
	drawArray(dc, fontPath, 150, 130, nums, headerStates, -1, -1, -1, false)

	dc.SetHexColor("#c86d76")
	dc.SetLineWidth(4)
	dc.DrawCircle(150+3*86+43, 162, 34)
	dc.Stroke()
	drawLabel(dc, fontPath, 72, 226, "暴力做法：从左到右遍历，第一个 >= 8 的位置就是答案。", 24, "#2f241c")
	drawLabel(dc, fontPath, 72, 262, "高效做法：在闭区间 [L, R] 中不断缩小答案所在范围。", 24, "#2f241c")

	fillRounded(dc, 72, 300, 650, 118, 18, "#f7f2e8")
	strokeRounded(dc, 72, 300, 650, 118, 18, "#d5c8b4", 1.4)
	drawLabel(dc, fontPath, 100, 342, "染色规则：", 28, "#2f241c")
	drawLabel(dc, fontPath, 100, 380, "红色表示 false，即 < 8", 24, "#8f3a3a")
	drawLabel(dc, fontPath, 380, 380, "蓝色表示 true，即 >= 8", 24, "#355f9e")
	drawLabel(dc, fontPath, 100, 410, "白色表示尚未确定，答案一定在白色区间里", 22, "#4b5563")

	drawLabel(dc, fontPath, 790, 342, "闭区间模板核心：", 28, "#2f241c")
	drawLabel(dc, fontPath, 790, 380, "1. `mid` 是红色：`L = mid + 1`", 24, "#2f241c")
	drawLabel(dc, fontPath, 790, 410, "2. `mid` 是蓝色：`R = mid - 1`", 24, "#2f241c")
	drawLabel(dc, fontPath, 790, 440, "3. 当 `L > R` 时停止", 24, "#2f241c")

	steps := []step{
		{
			title:   "Step 1  初始闭区间",
			desc1:   "L = 0, R = 5，mid = 2",
			desc2:   "nums[2] = 7，是红色，因此答案不在 [0, 2]",
			states:  []string{"red", "red", "red", "white", "white", "white"},
			l:       0,
			r:       5,
			m:       2,
			showM:   true,
			comment: "更新：L = mid + 1 = 3",
		},
		{
			title:   "Step 2  缩小到右半边",
			desc1:   "L = 3, R = 5，mid = 4",
			desc2:   "nums[4] = 8，是蓝色，因此答案在 [3, 4]",
			states:  []string{"red", "red", "red", "white", "blue", "blue"},
			l:       3,
			r:       5,
			m:       4,
			showM:   true,
			comment: "更新：R = mid - 1 = 3",
		},
		{
			title:   "Step 3  逼出分界点",
			desc1:   "L = 3, R = 3，mid = 3",
			desc2:   "nums[3] = 8，是蓝色，因此答案在 [3, 3]",
			states:  []string{"red", "red", "red", "blue", "blue", "blue"},
			l:       3,
			r:       3,
			m:       3,
			showM:   true,
			comment: "更新：R = mid - 1 = 2，随后 L > R，停止",
		},
	}

	drawBlock(dc, fontPath, 72, 470, 1356, 230, steps[0], nums)
	drawBlock(dc, fontPath, 72, 730, 1356, 230, steps[1], nums)
	drawBlock(dc, fontPath, 72, 990, 1356, 230, steps[2], nums)
	drawInvariantBox(dc, fontPath, 72, 1250, 1356, 190, nums)

	drawLabel(dc, fontPath, 1050, 1410, "答案：R + 1 = 3", 30, "#2f241c")

	must(dc.SavePNG(outputPath))
}
