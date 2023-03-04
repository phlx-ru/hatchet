package progress

import (
	"fmt"

	"github.com/cheggaaa/pb/v3"

	"github.com/phlx-ru/hatchet/cli"
)

type Bar interface {
	Start(message string, total int)
	Increment()
	Finish(message string)
}

type Progress struct {
	bar *pb.ProgressBar
}

func New() *Progress {
	return &Progress{
		bar: pb.New(0),
	}
}

func (p *Progress) Start(message string, total int) {
	fmt.Println(cli.ColorBlue + message)
	p.bar = pb.New(total)
	p.bar.Start()
}

func (p *Progress) Increment() {
	p.bar.Increment()
}

func (p *Progress) Finish(message string) {
	p.bar.Finish()
	echo := cli.ColorReset
	if message != "" {
		echo = message + cli.ColorReset
	}
	fmt.Println(echo)
}
