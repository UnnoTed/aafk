//go:generate fileb0x b0x.yml
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/MakeNowJust/hotkey"
	"github.com/UnnoTed/aafk/static"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
)

var running = false

const audioStart = "./heavy_negativevocalization06.mp3"
const audioEnd = "./heavy_goodjob01.mp3"

func main() {
	rand.Seed(time.Now().Unix())
	systray.Run(onReady, onExit)
}

func onReady() {
	icon, _ := static.ReadFile("./icon.ico")

	systray.SetIcon(icon)
	systray.SetTitle("Anti-AFK")
	systray.SetTooltip("Anti-AFK")
	mQuit := systray.AddMenuItem("Quit", "911")

	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	hkey := hotkey.New()
	hkey.Register(hotkey.Ctrl, hotkey.PRIOR, func() {
		running = !running

		if running {
			enable()
		} else {
			disable()
		}
	})

	log.Println("Press CTRL + PageUp")
}

func enable() {
	log.Println("Anti-AFK is enabled")
	play(audioStart)
	time.Sleep(2 * time.Second)

	go func() {
		for running {
			i := 9 + rand.Intn(11)
			log.Println("Next action in", i, "seconds")
			time.Sleep(time.Duration(i) * time.Second)
			randomAction()
		}
	}()
}

func disable() {
	log.Println("Anti-AFK is disabled")
	play(audioEnd)
}

func onExit() {
}

func randomAction() {
	n := rand.Intn(10)
	log.Println("action", n)

	switch n {
	case 0:
		holdKey("space")
	case 1:
		robotgo.MouseClick("left", true)
	case 2:
		holdKey("w")
	case 3:
		holdKey("a")
	case 4:
		holdKey("s")
	case 5:
		holdKey("d")
	case 6:
		holdKey("1")
	case 7:
		holdKey("2")
	case 8:
		holdKey("3")
	case 9:
		holdKey("control")
	case 10:
		holdKey("tab")

	default:
		robotgo.KeyTap("space")
	}
}

func holdKey(key string) {
	robotgo.KeyToggle(key, "down")

	i := rand.Intn(6)
	time.Sleep(time.Duration(i) * time.Second)

	robotgo.KeyToggle(key, "up")
}

func play(fpath string) {
	f, err := static.FS.OpenFile(static.CTX, fpath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Decode the .mp3 File, if you have a .wav file, use wav.Decode(f)
	s, format, _ := mp3.Decode(f)

	// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Channel, which will signal the end of the playback.
	playing := make(chan struct{})

	// Now we Play our Streamer on the Speaker
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		// Callback after the stream Ends
		close(playing)
	})))
	<-playing
}
