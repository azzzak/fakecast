package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/azzzak/fakecast/store"
)

type overview struct {
	Channel  *store.Channel  `json:"info"`
	Podcasts []store.Podcast `json:"podcasts"`
}

type updateChannel struct {
	Channel  *store.Channel `json:"channel"`
	OldAlias string         `json:"old_alias"`
}

type updateResponse struct {
	Cover string `json:"cover"`
	Error bool   `json:"error,omitempty"`
}

func checkChannels(cfg *Cfg, cs []store.Channel) []store.Channel {
	ix := 0
	for _, c := range cs {
		if cfg.FS.IsDirExist(c.Alias) {
			cs[ix] = c
			ix++
			continue
		}
		cfg.Store.DeleteChannel(c.ID)
	}

	for j := ix; j < len(cs); j++ {
		cs[j] = store.Channel{}
	}

	return cs[:ix]
}

func checkPodcasts(cfg *Cfg, alias string, ps []store.Podcast) []store.Podcast {
	ix := 0
	for _, p := range ps {
		if cfg.FS.IsPodcastExist(alias, p.Filename) {
			ps[ix] = p
			ix++
			continue
		}
		cfg.Store.DeletePodcast(p.ID)
	}

	for j := ix; j < len(ps); j++ {
		ps[j] = store.Podcast{}
	}

	return ps[:ix]
}

func (cfg *Cfg) createChannel(w http.ResponseWriter, r *http.Request) error {
	cid, err := cfg.Store.AddChannel()
	if err != nil {
		return err
	}

	if err = cfg.FS.CreateDir(cid); err != nil {
		return err
	}

	c, err := cfg.Store.ChannelInfo(cid)
	if err != nil {
		return err
	}

	c.ID = cid
	c.Title = fmt.Sprintf("New channel %d", cid)
	c.Alias = fmt.Sprintf("%d", cid)

	if err = cfg.Store.UpdateChannel(c); err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(c); err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) list(w http.ResponseWriter, r *http.Request) error {
	cs, err := cfg.Store.ListChannels()
	if err != nil {
		return err
	}

	cs = checkChannels(cfg, cs)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(cs); err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) overview(w http.ResponseWriter, r *http.Request) error {
	cid := r.Context().Value(CID).(int64)

	info, err := cfg.Store.ChannelInfo(cid)
	if err != nil {
		return err
	}

	info.Host = cfg.Host
	setCoverURL(cfg, info)

	ps, err := cfg.Store.ListPodcastsFrom(cid)
	if err != nil {
		return err
	}

	ps = checkPodcasts(cfg, info.Alias, ps)

	overview := overview{
		Channel:  info,
		Podcasts: ps,
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(overview); err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) updateChannel(w http.ResponseWriter, r *http.Request) error {
	var u *updateChannel

	var out interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		return err
	}

	var aliasError bool

	if u.Channel.Alias != u.OldAlias {
		err := cfg.FS.RenameDir(u.OldAlias, u.Channel.Alias)
		if err != nil {
			u.Channel.Alias = u.OldAlias
			aliasError = true
		}
	}

	err := cfg.Store.UpdateChannel(u.Channel)
	if err != nil {
		return err
	}

	if u.Channel.Cover != "" {
		setCoverURL(cfg, u.Channel)
	}

	out = updateResponse{
		Cover: u.Channel.Cover,
		Error: aliasError,
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(out); err != nil {
		return err
	}

	return nil
}

func (cfg *Cfg) deleteChannel(w http.ResponseWriter, r *http.Request) error {
	cid := r.Context().Value(CID).(int64)

	alias, err := cfg.Store.SwapCIDForAlias(cid)
	if err != nil {
		return err
	}

	if err = cfg.Store.DeleteChannel(cid); err != nil {
		return err
	}

	if err = cfg.FS.RemoveDir(alias); err != nil {
		return err
	}

	return nil
}
