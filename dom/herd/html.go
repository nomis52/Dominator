package herd

import (
	"fmt"
	"io"
	"sort"
	"time"
)

func (herd *Herd) writeHtml(writer io.Writer) {
	numSubs := herd.countSelectedSubs(nil)
	fmt.Fprintf(writer, "Time since current cycle start: %s<br>\n",
		time.Since(herd.currentScanStartTime))
	if numSubs < 1 {
		fmt.Fprintf(writer, "Duration of previous cycle: %s<br>\n",
			herd.previousScanDuration)
	} else {
		fmt.Fprintf(writer, "Duration of previous cycle: %s (%s/sub)<br>\n",
			herd.previousScanDuration,
			herd.previousScanDuration/time.Duration(numSubs))
	}
	fmt.Fprintf(writer, "Image server: <a href=\"http://%s/\">%s</a><br>\n",
		herd.imageServerAddress, herd.imageServerAddress)
	fmt.Fprintf(writer,
		"Number of <a href=\"listSubs\">subs</a>: <a href=\"showAllSubs\">%d</a><br>\n",
		numSubs)
	numSubs = herd.countSelectedSubs(selectAliveSub)
	fmt.Fprintf(writer,
		"Number of alive subs: <a href=\"showAliveSubs\">%d</a><br>\n",
		numSubs)
	numSubs = herd.countSelectedSubs(selectDeviantSub)
	fmt.Fprintf(writer,
		"Number of deviant subs: <a href=\"showDeviantSubs\">%d</a><br>\n",
		numSubs)
	numSubs = herd.countSelectedSubs(selectCompliantSub)
	fmt.Fprintf(writer,
		"Number of compliant subs: <a href=\"showCompliantSubs\">%d</a><br>\n",
		numSubs)
	subs := herd.getSelectedSubs(nil)
	connectDurations := getConnectDurations(subs)
	shortPollDurations := getPollDurations(subs, false)
	fullPollDurations := getPollDurations(subs, true)
	showDurationStats(writer, connectDurations, "Connect")
	showDurationStats(writer, shortPollDurations, "Short poll")
	showDurationStats(writer, fullPollDurations, "Full poll")
	fmt.Fprintf(writer, "Connection slots: %d out of %d<br>\n",
		len(herd.connectionSemaphore), cap(herd.connectionSemaphore))
	fmt.Fprintf(writer, "RPC slots: %d out of %d<br>\n",
		len(herd.pollSemaphore), cap(herd.pollSemaphore))
}

func selectAliveSub(sub *Sub) bool {
	switch sub.status {
	case statusUnknown:
		return false
	case statusConnecting:
		return false
	case statusDNSError:
		return false
	case statusFailedToConnect:
		return false
	case statusFailedToPoll:
		return false
	}
	return true
}

func selectDeviantSub(sub *Sub) bool {
	if sub.status == statusUpdating {
		return true
	}
	return false
}

func selectCompliantSub(sub *Sub) bool {
	if sub.status == statusSynced {
		return true
	}
	return false
}

func getConnectDurations(subs []*Sub) []time.Duration {
	durations := make([]time.Duration, 0, len(subs))
	for _, sub := range subs {
		if sub.lastConnectDuration > 0 {
			durations = append(durations, sub.lastConnectDuration)
		}
	}
	sort.Sort(durationList(durations))
	return durations
}

func getPollDurations(subs []*Sub, full bool) []time.Duration {
	durations := make([]time.Duration, 0, len(subs))
	for _, sub := range subs {
		var duration time.Duration
		if full {
			duration = sub.lastFullPollDuration
		} else {
			duration = sub.lastShortPollDuration
		}
		if duration > 0 {
			durations = append(durations, duration)
		}
	}
	sort.Sort(durationList(durations))
	return durations
}

type durationList []time.Duration

func (dl durationList) Len() int {
	return len(dl)
}

func (dl durationList) Less(i, j int) bool {
	return dl[i] < dl[j]
}

func (dl durationList) Swap(i, j int) {
	dl[i], dl[j] = dl[j], dl[i]
}

func showDurationStats(writer io.Writer, durations []time.Duration,
	durationType string) {
	if len(durations) < 1 {
		return
	}
	var avgDuration time.Duration
	for _, duration := range durations {
		avgDuration += duration
	}
	avgDuration /= time.Duration(len(durations))
	medDuration := durations[len(durations)/2]
	unit := "ns"
	scale := 1.0
	switch {
	case medDuration > 1e9:
		unit = "s"
		scale = 1e-9
	case medDuration > 1e6:
		unit = "ms"
		scale = 1e-6
	case medDuration > 1e3:
		unit = "µs"
		scale = 1e-3
	}
	fmt.Fprintf(writer,
		"%s durations: %.3f/%.3f/%.3f/%.3f %s (avg/med/min/max)<br>\n",
		durationType,
		float64(avgDuration)*scale, float64(medDuration)*scale,
		float64(durations[0])*scale, float64(durations[len(durations)-1])*scale,
		unit)
}
