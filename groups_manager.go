package main

// GroupsManager manages all groups by storing them in a map, using the group.id as key.
// "register" and "unregister" are both channels, that adds and removes a group from the
// the map.
type GroupsManager struct {
	groups     map[string]*Group
	register   chan *Group
	unregister chan *Group
}

// NewGroupsManager creates the only GroupsManager for the entire application
func NewGroupsManager() *GroupsManager {
	return &GroupsManager{
		groups:     make(map[string]*Group),
		register:   make(chan *Group),
		unregister: make(chan *Group),
	}
}

// Run runs in a goroutine for the entirety of the application. It does the management work i.e
/*
 - It listens for messages in gm.register channel to add a group to gm.groups
 - It listens for messages in gm.unregister channel to remove a group from gm.groups
*/
func (gm *GroupsManager) Run() {
	for {
		select {
		case group := <-gm.register:
			gm.groups[group.id] = group
		case group := <-gm.unregister:
			if _, ok := gm.Contains(group.id); ok {
				delete(gm.groups, group.id)
				// after deleting the group from gm.groups, the group's broadcast channel is closed,
				// to signify to Clients in the group that the group can't broadcast anymore, and
				// it's time to close the Client's connection.
				close(group.broadcast)
			}
		}
	}
}

// Contains checks if a group is in the gm, using id. It returns the found group and a bool that
// indicates if it found the group or not.
func (gm *GroupsManager) Contains(id string) (*Group, bool) {
	if group, ok := gm.groups[id]; ok {
		return group, true
	}

	return nil, false
}
