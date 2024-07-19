package arista

import (
	"regexp"
	"strings"
	"time"

	"github.com/stellaraf/go-parselog/types"
)

var patternISIS = regexp.MustCompile(`^L2 Neighbor State Change .+ SystemID (?P<remote>\S+) on (?P<iface>\S+).*to (?P<state>\S+)(: (?P<reason>.+))?$`)
var patternBGP = regexp.MustCompile(`^peer (?P<remote>\S+) \(VRF (?P<table>\S+) AS (?P<remote_as>\S+)\) old .+ new state (?P<state>\S+)$`)

const (
	bgpLen     int = 5
	isisMinLen int = 5
	isisMaxLen int = 6
)

type Parser func(string, string, time.Time, map[string]any) (types.Log, error)

var parseMap = map[*regexp.Regexp]Parser{
	regexp.MustCompile(`^L[12] Neighbor.+`):      ParseISIS,
	regexp.MustCompile(`^peer [0-9a-f\.\:]+.*$`): ParseBGP,
}

func ParseISIS(msg, src string, ts time.Time, extra map[string]any) (types.Log, error) {
	names := patternISIS.SubexpNames()
	matches := patternISIS.FindStringSubmatch(msg)

	if len(matches) < isisMinLen || len(matches) > isisMaxLen {
		return nil, types.ErrIncompleteMatch
	}

	iRemote := patternISIS.SubexpIndex(names[1])
	iIf := patternISIS.SubexpIndex(names[2])
	iState := patternISIS.SubexpIndex(names[3])

	state := strings.TrimSpace(matches[iState])
	remote := strings.TrimSpace(matches[iRemote])
	iface := strings.TrimSpace(matches[iIf])

	iReason := patternISIS.SubexpIndex(names[5])
	reason := ""
	if len(names) == 6 && iReason != -1 {
		reason = strings.TrimSpace(matches[iReason])
	}

	if strings.Contains(strings.ToLower(state), "init") {
		return nil, nil
	}

	l := &types.ISISLog{
		Base:      types.Base{Type: types.ISIS, Original: msg, Extra: extra},
		Local:     src,
		Timestamp: ts,
		Remote:    remote,
		Interface: iface,
		State:     types.DOWN,
		Reason:    reason,
	}
	if strings.Contains(strings.ToLower(state), "up") {
		l.State = types.UP
	}
	return l, nil
}

func ParseBGP(msg, src string, ts time.Time, extra map[string]any) (types.Log, error) {
	names := patternBGP.SubexpNames()
	matches := patternBGP.FindStringSubmatch(msg)
	if len(matches) != bgpLen {
		return nil, types.ErrIncompleteMatch
	}

	iRemote := patternBGP.SubexpIndex(names[1])
	iTable := patternBGP.SubexpIndex(names[2])
	iASN := patternBGP.SubexpIndex(names[3])
	iState := patternBGP.SubexpIndex(names[4])

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
		for pattern, parser := range parseMap {
			if pattern.MatchString(msg) {
				l, err := parser(msg, req.Source, req.Timestamp, req.Extra)
				if err != nil {
					return nil, err
				}
				if l != nil {
					logs = append(logs, l)
				}
			}
		}
	}
	if len(logs) == 0 {
		return nil, types.ErrNoMatchingParser
	}
	return logs, nil
}
