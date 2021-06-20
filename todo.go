package main

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
	"time"
)

type Lists struct {
	Items  map[string]*List
	Locker sync.Mutex

	// f is a cache for the loadpath; it will not
	// persist, but is useful for saving to disk
	// (so *Lists.Save() knows where to put files)
	f string
}

func LoadLists(f string) (l *Lists, err error) {
	file, err := os.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			l = &Lists{
				Items:  make(map[string]*List),
				Locker: sync.Mutex{},
			}
			err = nil
		}
		return
	}

	b := bytes.NewBuffer(file)
	dec := gob.NewDecoder(b)

	err = dec.Decode(l)

	l.f = f // do this last

	return
}

func (l *Lists) Save() (err error) {
	b := bytes.Buffer{}
	dec := gob.NewEncoder(&b)

	err = dec.Encode(l)
	if err != nil {
		return
	}

	return os.WriteFile(l.f, b.Bytes(), 0600)
}

type List struct {
	Items []*Item
}

type Item struct {
	ID         int
	Done       bool
	Title      string
	CreatedAt  time.Time
	MarkedDone time.Time
}

func NewList() *List {
	return &List{
		Items: make([]*Item, 0),
	}
}

func (l *List) Create(title string) (i *Item) {
	i = &Item{
		ID:        len(l.Items),
		Done:      false,
		Title:     title,
		CreatedAt: time.Now(),
	}

	l.Items = append(l.Items, i)

	return
}

func (l *List) Update(id int, title string) {
	if id >= len(l.Items) {
		return
	}

	l.Items[id].Title = title
}

func (l *List) Finish(id int) {
	if id >= len(l.Items) {
		return
	}

	l.Items[id].Done = true
	l.Items[id].MarkedDone = time.Now()

	return
}

func (l *List) Delete(id int) {
	if id >= len(l.Items) {
		return
	}

	l.Items = append(l.Items[:id], l.Items[id:]...)

	for idx, item := range l.Items {
		item.ID = idx
	}
}
