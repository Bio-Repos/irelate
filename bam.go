// Implements Relatable for Bam

package irelate

import (
	"io"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"
	"github.com/brentp/xopen"
)

type Bam struct {
	*sam.Record
	source     uint32
	related    []Relatable
	Chromosome string
	_end       uint32
}

func (a *Bam) Chrom() string {
	return a.Chromosome
}

// cast to 32 bits.
func (a *Bam) Start() uint32 {
	return uint32(a.Record.Start())
}

func (a *Bam) End() uint32 {
	if a._end != 0 {
		return a._end
	}
	a._end = uint32(a.Record.End())
	return a._end
}

func (a *Bam) Source() uint32 {
	return a.source
}

func (a *Bam) SetSource(src uint32) {
	a.source = src
}

func (a *Bam) AddRelated(b Relatable) {
	if a.related == nil {
		a.related = make([]Relatable, 1, 2)
		a.related[0] = b
	} else {
		a.related = append(a.related, b)
	}
}
func (a *Bam) Related() []Relatable {
	return a.related
}

func (a *Bam) MapQ() int {
	return int(a.Record.MapQ)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func BamToRelatable(file string) RelatableChannel {

	ch := make(chan Relatable, 64)

	go func() {
		f, err := xopen.XReader(file)
		check(err)
		b, err := bam.NewReader(f, 0)
		check(err)
		for {
			rec, err := b.Read()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err)
				}
			}
			if rec.RefID() == -1 { // unmapped
				continue
			}
			// TODO: see if keeping the list of chrom names and using a ref is better.
			bam := Bam{Record: rec, Chromosome: rec.Ref.Name(), related: nil}
			ch <- &bam
		}
		close(ch)
		b.Close()
		f.(io.ReadCloser).Close()
	}()
	return ch
}
