package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/azzzak/fakecast/fs"
	"github.com/azzzak/fakecast/store"
)

func (cfg *Cfg) uploadPodcast(w http.ResponseWriter, r *http.Request) error {
	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	cid := r.Context().Value(CID).(int64)
	filename := header.Filename

	short, err := cfg.Store.SwapCIDForAlias(cid)
	if err != nil {
		return err
	}

	if cfg.FS.IsPodcastExist(short, filename) {
		t := strings.Split(filename, ".")
		if len(t) >= 2 {
			t[len(t)-2] = fmt.Sprintf("%s-%x", t[len(t)-2], time.Now().Unix())
			filename = strings.Join(t, ".")
		}
	}

	holder, err := cfg.FS.SavePodcastToDir(short, filename)
	if err != nil {
		return err
	}
	defer holder.Close()

	title, _ := fs.NameAndExtFrom(filename)

	podcast, err := cfg.Store.AddPodcastToChannel(cid, filename, title, r.FormValue("length"))
	if err != nil {
		return err
	}

	io.Copy(holder, file)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(podcast); err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) podcastInfo(w http.ResponseWriter, r *http.Request) error {
	pid := r.Context().Value(PID).(int64)

	res, err := cfg.Store.PodcastInfo(pid)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) updatePodcast(w http.ResponseWriter, r *http.Request) error {
	var p *store.Podcast

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		return err
	}

	if p.GUID == "" {
		now := time.Now()
		p.GUID = fmt.Sprintf("%x", now.Unix())
		p.PubDate = now.UTC().Format("Mon, 2 Jan 2006 15:04:05 MST")
	}

	err := cfg.Store.UpdatePodcast(p)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) deletePodcast(w http.ResponseWriter, r *http.Request) error {
	cid := r.Context().Value(CID).(int64)
	pid := r.Context().Value(PID).(int64)

	short, err := cfg.Store.SwapCIDForAlias(cid)
	if err != nil {
		return err
	}

	filename, err := cfg.Store.SwapPIDForFilename(pid)
	if err != nil {
		return err
	}

	err = cfg.Store.DeletePodcast(pid)
	if err != nil {
		return err
	}

	err = cfg.FS.RemovePodcast(short, filename)
	if err != nil {
		return err
	}

	return nil
}
