package main

import (
	"image"
	_ "image/png"
	"os"
	"os/exec"
	"math/rand"
	"strconv"
	"fmt"
	"math"
	"time"
)

func main() {
	for {
		file, err := getScreen()
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		tx, ty := getTargetPos(img)
		sx, sy := getSourcePos(img)

		//img.Bounds().Max.X
		dis := math.Sqrt(math.Abs(math.Pow(float64(tx-sx), 2)) + math.Abs(math.Pow(float64(ty-sy), 2)))
		//dis := math.Abs(float64(sx - tx))
		second := dis * 1.29
		if second < 300 {
			second = 300
		}
		fmt.Println(tx, ty, sx, sy, second, dis, file.Name())

		jumpToTarget(second)

		//os.Rename(file.Name(), fmt.Sprintf("%d_%s.png", dis, file.Name()))
		time.Sleep(time.Second * 2)
	}

}

func jumpToTarget(second float64) {
	exec.Command("adb", "shell", "input", "swipe ", "10", "10", "10", "10", strconv.Itoa(int(second))).Run()
}

func getScreen() (file *os.File, err error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	fileName := fmt.Sprintf("wechat_game_%d.png", r1.Intn(10000))
	exec.Command("adb", "shell", "screencap", "/sdcard/"+fileName).Run()
	exec.Command("adb", "pull", "/sdcard/"+fileName).Run()
	file, err = os.Open(fileName)
	return
}

func getTargetPos(img image.Image) (targetX int, targetY int) {
	height := img.Bounds().Max.Y
	width := img.Bounds().Max.X

	for y := height / 4; y < height*2/3; y++ {
		sameColorLine := true
		beforeColor := 0
		for x := 0; x < width; x++ {
			nowColorR, nowColorG, nowColorB, _ := img.At(x, y).RGBA()
			nowColorRP := int(nowColorR / 257)
			nowColorGP := int(nowColorG / 257)
			nowColorBP := int(nowColorB / 257)
			nowColor := nowColorRP + nowColorGP + nowColorBP
			if beforeColor != 0 && (math.Abs(float64(nowColor - beforeColor))) >= 10 {
				sameColorLine = false
			}
			if !sameColorLine {
				targetX = x
				targetY = y + 105
				return
			}
			beforeColor = nowColor
		}
	}
	return
}

func getSourcePos(img image.Image) (sourceX int, sourceY int) {
	for y := img.Bounds().Max.Y / 3; y < img.Bounds().Max.Y*2/3; y++ {

		for x := 0; x < img.Bounds().Max.X; x++ {
			nowRedColor, nowGreenColor, nowBlueColor, _ := img.At(x, y).RGBA()
			if int(nowRedColor/257) == 59 && nowGreenColor/257 == 59 && nowBlueColor/257 == 77 {
				sourceX = x
				sourceY = y + 135
				return
			}
		}
	}
	return
}
