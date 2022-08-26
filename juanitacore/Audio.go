package juanitacore

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	CHANNELS   int = 2
	FRAME_RATE int = 48000
	FRAME_SIZE int = 960
	MAX_BYTES  int = (FRAME_SIZE * 2) * 2
)

/*
this shit is messy and i don't fully understand it yet credit to github.com/bwmarrin's voice example for the base code
*/

func (connection *JuanitaConnection) sendPCM(voice *discordgo.VoiceConnection, pcm <-chan []int16) {
	connection.Lock.Lock()
	if connection.Sendpcm || pcm == nil {
		connection.Lock.Unlock()
		return
	}
	connection.Sendpcm = true
	connection.Lock.Unlock()
	defer func() {
		connection.Sendpcm = false
	}()
	encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder error,", err)
		return
	}
	for {
		receive, ok := <-pcm
		if !ok {
			fmt.Println("PCM channel closed")
			return
		}
		opus, err := encoder.Encode(receive, FRAME_SIZE, MAX_BYTES)
		if err != nil {
			fmt.Println("Encoding error,", err)
			return
		}
		if !voice.Ready || voice.OpusSend == nil {
			fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", voice.Ready, voice.OpusSend)
			return
		}
		voice.OpusSend <- opus
	}
}

func (connection *JuanitaConnection) Play(ffmpeg *exec.Cmd) error {
	if connection.Playing {
		return errors.New("song already playing")
	}
	connection.StopRunning = false
	out, err := ffmpeg.StdoutPipe()
	if err != nil {
		return err
	}
	buffer := bufio.NewReaderSize(out, 16384)
	err = ffmpeg.Start()
	if err != nil {
		return err
	}
	connection.Playing = true
	defer func() {
		connection.Playing = false
	}()
	connection.VoiceConnection.Speaking(true)
	defer connection.VoiceConnection.Speaking(false)
	if connection.Send == nil {
		connection.Send = make(chan []int16, 2)
	}
	go connection.sendPCM(connection.VoiceConnection, connection.Send)
	for {
		if connection.StopRunning {
			ffmpeg.Process.Kill()
			break
		}
		audioBuffer := make([]int16, FRAME_SIZE*CHANNELS)
		err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			return err
		}
		connection.Send <- audioBuffer
	}
	return nil
}

func (connection *JuanitaConnection) Stop() {
	connection.StopRunning = true
	connection.Playing = false
}
