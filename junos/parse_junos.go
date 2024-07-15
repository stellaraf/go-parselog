package junos

import (
	"log"
	"regexp"
	"strings"

	"github.com/stellaraf/go-parselog/types"
)

var patternISIS = regexp.MustCompile(`^IS-IS (?P<state>.+) .+ to (?P<remote>.+) on (?P<iface>[\w\.]+)(, reason: (?P<reason>.+))?$`)
var patternBGP = regexp.MustCompile(`^BGP peer (?P<remote>.+) \(.+AS (?P<asn>\d+).+changed state from \S+ to (?P<state>\S+).*\(instance (?P<instance>\S+)\).*$`)

const (
	bgpMinLen  int = 5
	isisMinLen int = 5
	isisMaxLen int = 6
)

var parseMap = map[string]types.Parser{
	"IS-IS":    ParseISIS,
	"BGP peer": ParseBGP,
}

func ParseISIS(req *types.Request) (types.Log, error) {
	msg := strings.TrimSpace(req.Message)
	names := patternISIS.SubexpNames()
	matches := patternISIS.FindStringSubmatch(msg)

	if len(matches) < isisMinLen || len(matches) > isisMaxLen {
		log.Println(matches, len(matches))
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
		Base:      types.Base{Type: types.ISIS, Original: req.Message, Extra: req.Extra},
		Local:     req.Source,
		Timestamp: req.Timestamp.Time,
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

func ParseBGP(req *types.Request) (types.Log, error) {
	msg := strings.TrimSpace(req.Message)
	names := patternBGP.SubexpNames()
	matches := patternBGP.FindStringSubmatch(msg)
	if len(matches) != bgpMinLen {
		log.Println(matches, len(matches))
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
		Base:      types.Base{Type: types.BGP, Original: req.Message, Extra: req.Extra},
		Timestamp: req.Timestamp.Time,
		Local:     req.Source,
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

func Parse(req *types.Request) (types.Log, error) {
	for prefix, parser := range parseMap {
		if strings.HasPrefix(req.Message, prefix) {
			return parser(req)
		}
	}
	return nil, types.ErrNoMatchingParser
}
