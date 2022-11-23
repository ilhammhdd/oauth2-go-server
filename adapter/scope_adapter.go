package adapter

import (
	"sort"
	"strings"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

// TODO: need additional checker whether the parent of a child scope actually exists or not
func ParseScopeWithRel(rawScopes, callTraceFunc string) ([]*entity.ScopeWithRel, []*errorkit.DetailedError) {
	/* a:6 b:3 b.c:7 b.g:7 d:4 d.e:6 d.e.f:4
	read  :0b100
	write :0b10
	delete:0b1
	for the youngest child like "c", it doesn't make sense if it is a scope of its own without any parent.
	If the relations between two or more scope is many-to-many then the above statement is incorrect.
	The children of any scope is only the scopes that is directly parented by it, meaning the grandchildren doesn't count as its children, so "d"'s children is only "e". */
	var scopeWithRels []*entity.ScopeWithRel
	var detailedErrs []*errorkit.DetailedError
	scopes := sort.StringSlice(strings.Split(rawScopes, " "))
	var scopeChildrenMap map[string]*[]string = make(map[string]*[]string)
	var scopeMap map[string]*entity.ScopeWithRel = make(map[string]*entity.ScopeWithRel)

	for i := range scopes {
		scopeBytes := []byte(scopes[i])
		if len(scopeBytes) < 3 {
			detailedErrs = append(detailedErrs, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrInvalidScope, DetailedErrDescGen))
			continue
		}
		if scopeBytes[len(scopeBytes)-1] > 55 || scopeBytes[len(scopeBytes)-1] < 48 || scopeBytes[len(scopeBytes)-2] != ':' {
			detailedErrs = append(detailedErrs, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrInvalidScope, DetailedErrDescGen, string(scopeBytes[:len(scopeBytes)-2])))
			continue
		}
		var parent string
		var child string
		for j := len(scopeBytes) - 3; j >= 0 && child == ""; j-- {
			if scopeBytes[j] == '.' {
				child = string(scopeBytes[j+1 : len(scopeBytes)-2])
			}
			for k := len(scopeBytes) - 3; k >= 0 && parent == ""; k-- {
				if scopeBytes[k] == '.' {
					parent = string(scopeBytes[k+1 : j])
				}
			}
		}
		if parent == "" && child == "" {
			parent = string(scopeBytes[:len(scopeBytes)-2])
		}

		scopeWithRel := entity.ScopeWithRel{Scope: &entity.Scope{Scope: string(scopeBytes[:len(scopeBytes)-2]), Permission: scopeBytes[len(scopeBytes)-1] - 48}}
		if parent != "" && child != "" {
			if children, ok := scopeChildrenMap[parent]; ok {
				*children = append(*children, child)
				scopeChildrenMap[parent] = children
			} else {
				scopeChildrenMap[parent] = &[]string{child}
			}
			scopeMap[child] = &scopeWithRel
		} else if parent != "" {
			scopeChildrenMap[parent] = &[]string{}
			scopeMap[parent] = &scopeWithRel
		}
	}

	for key, val := range scopeMap {
		if children, ok := scopeChildrenMap[key]; ok {
			var valChildren []*entity.ScopeWithRel
			for i := range *children {
				scopeMap[(*children)[i]].ParentScope = val.Scope
				valChildren = append(valChildren, scopeMap[(*children)[i]])
			}
			val.ChildrenScopes = valChildren
			scopeMap[key] = val
		}
		scopeWithRels = append(scopeWithRels, val)
	}

	return scopeWithRels, detailedErrs
}
