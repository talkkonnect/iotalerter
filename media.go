package iotalerter

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/talkkonnect/volume-go"
)

func PlayWavLocal(SoundFilePath string, PlayBackVolume int) error {
	origVolume, _ := volume.GetVolume(Config.Global.Settings.Outputdevice)
	cmd := exec.Command("/usr/bin/aplay", SoundFilePath)
	err := volume.SetVolume(PlayBackVolume, Config.Global.Settings.Outputdevice)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("alert: set volume failed: %+v", err))
	}
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("alert: cmd.Run() for aplay failed with %s\n", err))
	}
	err = volume.SetVolume(origVolume, Config.Global.Settings.Outputdevice)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("alert: set volume failed: %+v", err))
	}
	return nil
}

func findEventSound(findEventSound string) EventSoundStruct {
	for _, sound := range Config.Global.Sounds.Sound {
		if sound.Enabled && sound.Event == findEventSound {
			return EventSoundStruct{sound.Enabled, sound.File, sound.Volume, sound.Blocking}
		}
	}
	return EventSoundStruct{false, "", "0", false}
}

func localMediaPlayer(fileNameWithPath string, playbackvolume int, blocking bool, duration float32, loop int) {

	if loop == 0 || loop > 3 {
		log.Println("warn: Infinite Loop or more than 3 loops not allowed")
		return
	}

	CmdArguments := []string{fileNameWithPath, "-volume", strconv.Itoa(playbackvolume), "-autoexit", "-loop", strconv.Itoa(loop), "-autoexit", "-nodisp"}

	if duration > 0 {
		CmdArguments = []string{fileNameWithPath, "-volume", strconv.Itoa(playbackvolume), "-autoexit", "-t", fmt.Sprintf("%.1f", duration), "-loop", strconv.Itoa(loop), "-autoexit", "-nodisp"}
	}

	cmd := exec.Command("/usr/bin/ffplay", CmdArguments...)

	WaitForFFPlay := make(chan struct{})
	go func() {
		cmd.Run()
		if blocking {
			WaitForFFPlay <- struct{}{} // signal that the routine has completed
		}
	}()
	if blocking {
		<-WaitForFFPlay
	}
}

// func aplayLocal(fileNameWithPath string) {
// 	var player string
// 	var CmdArguments = []string{}

// 	if path, err := exec.LookPath("aplay"); err == nil {
// 		CmdArguments = []string{fileNameWithPath, "-q", "-N"}
// 		player = path
// 	} else if path, err := exec.LookPath("paplay"); err == nil {
// 		CmdArguments = []string{fileNameWithPath}
// 		player = path
// 	} else {
// 		return
// 	}

// 	log.Printf("debug: player %v CmdArguments %v", player, CmdArguments)

// 	cmd := exec.Command(player, CmdArguments...)

// 	_, err := cmd.CombinedOutput()

// 	if err != nil {
// 		return
// 	}
// }

// func PlayTone(toneFreq int, toneDuration float32, destination string, withRXLED bool) {

// 	if destination == "local" {

// 		cmdArguments := []string{"-f", "lavfi", "-i", "sine=frequency=" + strconv.Itoa(toneFreq) + ":duration=" + fmt.Sprintf("%f", toneDuration), "-autoexit", "-nodisp"}
// 		cmd := exec.Command("/usr/bin/ffplay", cmdArguments...)
// 		var out bytes.Buffer
// 		cmd.Stdout = &out

// 		err := cmd.Run()
// 		if err != nil {
// 			log.Println("error: ffplay error ", err)
// 			return
// 		}

// 		log.Printf("info: Played Tone at Frequency %v Hz With Duration of %v Seconds\n", toneFreq, toneDuration)
// 	}
// }

// func findInputEventSoundFile(findInputEventSound string) InputEventSoundFileStruct {
// 	for _, sound := range Config.Global.Software.Sounds.Input.Sound {
// 		if sound.Event == findInputEventSound {
// 			if sound.Enabled {
// 				return InputEventSoundFileStruct{sound.Event, sound.File, sound.Enabled}
// 			}
// 		}
// 	}
// 	return InputEventSoundFileStruct{findInputEventSound, "", false}
// }

// func playIOMedia(inputEvent string) {
// 	if Config.Global.Software.Sounds.Input.Enabled {
// 		var inputEventSoundFile InputEventSoundFileStruct = findInputEventSoundFile(inputEvent)
// 		if inputEventSoundFile.Enabled {
// 			go aplayLocal(inputEventSoundFile.File)
// 		}
// 	}
// }
