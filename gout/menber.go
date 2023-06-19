package gout

import (
	"regexp"
)

type Manager struct {
	menbers *[]Entry //管理的小组所在的地址
}

type Entry struct {
	addr     *Manager       //所在小组的地址
	params   map[int]string //储存params参数
	rPattern string         //加工后的path
}

var (
	regexFactor   string = ":[\\S]+"
	ReplaceFactor string = "[\\S]+"
)

func check(pattern string, menbers []Entry) {
	for _, menber := range menbers {
		if pattern == menber.rPattern {
			panic("exist the same url ,dangerous action!")
		}
	}
}

func (m *Manager) insert(pattern string, parts []string) {
	menbers := *m.menbers
	check(pattern, menbers)
	e := Entry{
		addr:   m,
		params: make(map[int]string),
	}
	//判断是否携带params参数
	var joint string = "/"
	for i, part := range parts {
		if regexMatch(regexFactor, part) {
			e.params[i] = part
			part = ReplaceFactor
		}
		if i == len(parts)-1 {
			joint += part
		} else {
			joint += part + "/"
		}
	}
	e.rPattern = joint
	//插入
	*m.menbers = append(menbers, e)
}

func (m *Manager) search(pattern string, parts []string) (string, map[int]string) {
	menbers := *m.menbers
	//遍历整个小组，查看是否有与之一致的url
	for _, menber := range menbers {
		if menber.rPattern == "" {
			break
		}
		//查看能否精准匹配
		if pattern == menber.rPattern {
			return menber.rPattern, menber.params
		}
	}
	//匹配正则
	for _, menber := range menbers {
		if regexMatch(menber.rPattern, pattern) && menber.rPattern != "/" {
			return jointPattern(parts, menber.params), menber.params
		}
	}
	return "", nil
}

func regexMatch(regex, s string) bool {
	re := regexp.MustCompile(regex)
	return re.MatchString(s)
}

func jointPattern(parts []string, params map[int]string) string {
	joint := "/"

	for i, part := range parts {
		if params[i] != "" {
			if i == len(parts)-1 {
				joint += params[i]
				break
			}
			joint += params[i] + "/"
		} else if i == len(parts)-1 {
			joint += part
		} else {
			joint += part + "/"
		}
	}

	return joint
}
