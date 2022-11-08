package voiceAlert

import (
	"golang.org/x/exp/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	success  = iota + 1 // 单次转码成功
	failed              // 转码失败,程序退出
	complete            // 转码进程完成
)

func Voice(msg int) {
	//slog.SetDefault(slog.New(slog.NewJSONHandler(io.MultiWriter(logf, os.Stdout))))
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)))
	slog.Info("voice start!")
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// 查询发音人 `say -v ?`
		voice := "Serena"
		//if !HasSpoker(voice) {
		//	voice = "Victoria"
		//}
		switch msg {
		case success:
			cmd = exec.Command("say", "-v", voice, "Rocket was launched successfully")
			if err := cmd.Start(); err != nil {
				slog.Warn("执行发声命令发生错误", err)
			}
		case failed:
			cmd = exec.Command("say", "-v", voice, "Rocket launch failed")
			if err := cmd.Start(); err != nil {
				slog.Warn("执行发声命令发生错误", err)
			}
		case complete:
			cmd = exec.Command("say", "-v", voice, "mission complete!")
			if err := cmd.Start(); err != nil {
				slog.Warn("执行发声命令发生错误", err)
			}
		}
	case "linux":
		cmd = exec.Command("echo", "-e", "\\a")
		switch msg {
		case success:
			for i := 0; i < 2; i++ {
				if err := cmd.Start(); err != nil {
					slog.Warn("执行发声命令发生错误", err)
				}
			}
		case failed:
			for i := 0; i < 50; i++ {
				if err := cmd.Start(); err != nil {
					slog.Warn("执行发声命令发生错误", err)
				}
			}
		case complete:
			for i := 0; i < 100; i++ {
				if err := cmd.Start(); err != nil {
					slog.Warn("执行发声命令发生错误", err)
				}
			}
		}
	}
	if err := cmd.Wait(); err != nil {
		slog.Warn("执行命令过程中发生错误", err)
	}

}

func HasSpoker(key string) bool {
	cmd := exec.Command("say", "-v", "?")
	slog.Info("查询发音人")
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		slog.Error("cmd.StdoutPipe产生错误", err)
	}
	if err = cmd.Start(); err != nil {
		slog.Error("cmd.Run产生的错误", err)
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		//写成输出日志
		t := string(tmp)
		t = strings.Replace(t, "\u0000", "", -1)
		if strings.Contains(t, key) {
			slog.Info("found!", "voice", key)
			return true
		} else {
			slog.Warn("not found!", "voice", key)
		}
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		slog.Error("命令执行中有错误产生", err)
	}
	return false
}
