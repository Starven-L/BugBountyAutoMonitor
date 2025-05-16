package util

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

var colors = []string{
	"\033[36m", // Cyan
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[35m", // Magenta
	"\033[31m", // Red
	"\033[34m", // Blue
}

//var banner = `
//____              ____                    _         __  __             _ _     ___  _ __
//| __ ) _   _  __ _| __ )  ___  _   _ _ __ | |_ _   _|  \/  | ___  _ __ (_) |_ /  _ \| '__|
//|  _ \| | | |/ _ \|  _ \ / _ \| | | | '_ \| __| | | | |\/| |/ _ \| '_ \| | __|| | | | |
//| |_) | |_| | (_| | |_) | (_) | |_| | | | | |_| |_| | |  | | (_) | | | | | |_ | |_| | |
//|____/ \__,_|\__, |____/ \___/ \__,_|_| |_|\__|\__, |_|  |_|\___/|_| |_|_|\__| \___/|_|
//             |___/                             |___/
//`

func ShowLogo() {
	myFigure := figure.NewFigure("BugBountyMonitor", "", true)
	fmt.Println(Blue + myFigure.String())
	fmt.Printf("%s【+】Author: Starven\n", Red)
	fmt.Printf("%s【+】Email: starvenl@qq.com\n", Yellow)
	fmt.Printf("%s【+】Team：Syclover三叶草安全技术小组\n", Green)
}
