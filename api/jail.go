package api

import (
	"encoding/json"
	"net/http"

	"github.com/Strum355/log"
	"github.com/go-chi/chi"
)

func (a *API) getJail(w http.ResponseWriter, r *http.Request) {
	errHandle := func(err error) bool {
		if err != nil {
			log.WithError(err).Error("Could not get jail")
			a.handleError(w, r, err)
			return false
		}
		return true
	}

	jail := chi.URLParam(r, "jail")

	currentFailed, totalFailed, fileList, currentBanned, totalBanned, ipList, err := a.Fail2Conn.JailStatus(jail)
	if err != nil {
		log.WithError(err).Error("Could not get jail")
		a.handleError(w, r, err)
		return
	}

	failRegex, err := a.Fail2Conn.JailFailRegex(jail)
	if !errHandle(err) {
		return
	}
	findTime, err := a.Fail2Conn.JailFindTime(jail)
	if !errHandle(err) {
		return
	}
	useDNS, err := a.Fail2Conn.JailUseDNS(jail)
	if !errHandle(err) {
		return
	}
	maxRetry, err := a.Fail2Conn.JailMaxRetry(jail)
	if !errHandle(err) {
		return
	}
	actions, err := a.Fail2Conn.JailActions(jail)
	if !errHandle(err) {
		return
	}

	if ipList == nil {
		ipList = []string{}
	}
	if failRegex == nil {
		failRegex = []string{}
	}

	output := map[string]interface{}{
		"current_failed":   currentFailed,
		"total_failed":     totalFailed,
		"file_list":        fileList,
		"currently_banned": currentBanned,
		"total_banned":     totalBanned,
		"ip_list":          ipList,
		"fail_regexes":     failRegex,
		"find_time":        findTime,
		"use_dns":          useDNS,
		"max_retry":        maxRetry,
		"actions":          actions,
	}

	json.NewEncoder(w).Encode(output)
}

func (a *API) jailBanIP(w http.ResponseWriter, r *http.Request) {
	input := struct {
		IP string `json:"ip"`
	}{}
	json.NewDecoder(r.Body).Decode(&input)

	output, err := a.Fail2Conn.JailBanIP(chi.URLParam(r, "jail"), input.IP)
	if err != nil {
		log.WithError(err).Error("Error banning IP")
		a.handleError(w, r, err)
		return
	}

	log.WithFields(log.Fields{
		"ip": input.IP,
	}).Info("IP banned")
	json.NewEncoder(w).Encode(output)
}

func (a *API) jailUnbanIP(w http.ResponseWriter, r *http.Request) {
	input := struct {
		IP string `json:"ip"`
	}{}
	json.NewDecoder(r.Body).Decode(&input)

	output, err := a.Fail2Conn.JailUnbanIP(chi.URLParam(r, "jail"), input.IP)
	if err != nil {
		log.WithError(err).Error("Error unbanning IP")
		a.handleError(w, r, err)
		return
	}

	log.WithFields(log.Fields{
		"ip": input.IP,
	}).Info("IP unbanned")
	json.NewEncoder(w).Encode(output)
}

func (a *API) jailAddFailRegex(w http.ResponseWriter, r *http.Request) {
	input := struct {
		FailRegex string `json:"fail_regex"`
	}{}
	json.NewDecoder(r.Body).Decode(&input)

	output, err := a.Fail2Conn.JailAddFailRegex(chi.URLParam(r, "jail"), input.FailRegex)
	if err != nil {
		log.WithError(err).Error("Error adding fail regex")
		a.handleError(w, r, err)
		return
	}

	log.WithFields(log.Fields{
		"regex": input.FailRegex,
	}).Info("Fail Regex added")
	json.NewEncoder(w).Encode(output)
}

func (a *API) jailDeleteFailRegex(w http.ResponseWriter, r *http.Request) {
	input := struct {
		FailRegex string `json:"fail_regex"`
	}{}
	json.NewDecoder(r.Body).Decode(&input)

	output, err := a.Fail2Conn.JailDeleteFailRegex(chi.URLParam(r, "jail"), input.FailRegex)
	if err != nil {
		log.WithError(err).Error("Error removing fail regex")
		a.handleError(w, r, err)
		return
	}

	log.WithFields(log.Fields{
		"regex": input.FailRegex,
	}).Info("Fail Regex added")
	json.NewEncoder(w).Encode(output)
}
