package junos

import (
	"regexp"
	"strings"
	"time"

	"github.com/stellaraf/go-parselog/types"
)

var patternISIS = regexp.MustCompile(`^IS-IS (?P<state>.+) .+ to (?P<remote>.+) on (?P<iface>[\S\.]+)(, reason: (?P<reason>.+))?$`)
var patternBGP = regexp.MustCompile(`^BGP peer (?P<remote>.+) \(.+AS (?P<asn>\d+).+changed state from \S+ to (?P<state>\S+).*\(instance (?P<instance>\S+)\).*$`)

const (
	bgpMinLen  int = 5
	isisMinLen int = 5
	isisMaxLen int = 6
)

type Parser func(string, string, time.Time, map[string]any) (types.Log, error)

var parseMap = map[string]Parser{
	"IS-IS":    ParseISIS,
	"BGP peer": ParseBGP,
}

func ParseISIS(msg, src string, ts time.Time, extra map[string]any) (types.Log, error) {
	names := patternISIS.SubexpNames()
	matches := patternISIS.FindStringSubmatch(msg)

	if len(matches) < isisMinLen || len(matches) > isisMaxLen {
		return nil, types.ErrIncompleteMatch
	}

	iState := patternISIS.SubexpIndex(names[1])
	iRemote := patternISIS.SubexpIndex(names[2])
	iIf := patternISIS.SubexpIndex(names[3])

	state := strings.TrimSpace(matches[iState])
	remote := strings.TrimSpace(matches[iRemote])
	iface := strings.TrimSpace(matches[iIf])

	iReason := patternISIS.SubexpIndex(names[5])
	reason := ""
	if len(names) == 6 && iReason != -1 {
		reason = strings.TrimSpace(matches[iReason])
	}

	l := &types.ISISLog{
		Base:      types.Base{Type: types.ISIS, Original: msg, Extra: extra},
		Local:     src,
		Timestamp: ts,
		Remote:    remote,
		Interface: iface,
		Reason:    reason,
		State:     types.DOWN,
	}
	if strings.Contains(strings.ToLower(state), "new") {
		l.State = types.UP
	}
	return l, nil
}

func ParseBGP(msg, src string, ts time.Time, extra map[string]any) (types.Log, error) {
	names := patternBGP.SubexpNames()
	matches := patternBGP.FindStringSubmatch(msg)
	if len(matches) != bgpMinLen {
		return nil, types.ErrIncompleteMatch
	}

	iRemote := patternBGP.SubexpIndex(names[1])
	iASN := patternBGP.SubexpIndex(names[2])
	iState := patternBGP.SubexpIndex(names[3])
	iTable := patternBGP.SubexpIndex(names[4])

	remote := strings.TrimSpace(matches[iRemote])
	asn := strings.TrimSpace(matches[iASN])
	state := strings.TrimSpace(matches[iState])
	table := strings.TrimSpace(matches[iTable])

	l := &types.BGPLog{
		Base:      types.Base{Type: types.BGP, Original: msg, Extra: extra},
		Timestamp: ts,
		Local:     src,
		Remote:    remote,
		State:     types.DOWN,
		RemoteAS:  asn,
		Table:     table,
	}
	if strings.Contains(strings.ToLower(state), "established") {
		l.State = types.UP
	}
	return l, nil
}

func Parse(req *types.Request) ([]types.Log, error) {
	logs := make([]types.Log, 0, len(req.Messages))
	for _, msg := range req.Messages {
		for prefix, parser := range parseMap {
			if strings.HasPrefix(msg, prefix) {
				l, err := parser(msg, req.Source, req.Timestamp, req.Extra)
				if err != nil {
					return nil, err
				}
				logs = append(logs, l)
			}
		}
	}

	if len(logs) != 0 {
		return logs, nil
	}
	return nil, types.ErrNoMatchingParser
}
