package main

import "sort"

type timelineItem struct {
	Event                Event
	Depth                int
	MissingParent, Cycle bool
}

func parentID(ev Event) string {
	var id string
	for _, tag := range ev.Tags {
		if len(tag) >= 2 && tag[0] == "e" {
			id = tag[1]
		}
	}
	return id
}

func buildTimeline(events []Event) []timelineItem {
	nodes := map[string]Event{}
	for _, ev := range events {
		if _, exists := nodes[ev.ID]; !exists {
			nodes[ev.ID] = ev
		}
	}
	sorted := make([]Event, 0, len(nodes))
	for _, ev := range nodes {
		sorted = append(sorted, ev)
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].CreatedAt == sorted[j].CreatedAt {
			return sorted[i].ID < sorted[j].ID
		}
		return sorted[i].CreatedAt > sorted[j].CreatedAt
	})
	children := map[string][]Event{}
	roots := []Event{}
	for _, ev := range sorted {
		if _, ok := nodes[parentID(ev)]; ok {
			children[parentID(ev)] = append(children[parentID(ev)], ev)
		} else {
			roots = append(roots, ev)
		}
	}
	result := []timelineItem{}
	visited := map[string]bool{}
	var walk func(Event, int, bool)
	walk = func(ev Event, depth int, cycle bool) {
		if visited[ev.ID] {
			return
		}
		visited[ev.ID] = true
		parent := parentID(ev)
		_, found := nodes[parent]
		result = append(result, timelineItem{ev, depth, parent != "" && !found, cycle})
		for _, child := range children[ev.ID] {
			walk(child, depth+1, false)
		}
	}
	for _, ev := range roots {
		walk(ev, 0, false)
	}
	for _, ev := range sorted {
		if !visited[ev.ID] {
			walk(ev, 0, true)
		}
	}
	return result
}
