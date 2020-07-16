// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// 
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package groups

import (
	"github.com/kiprotect/kiprotect/actions/anonymize/aggregate"
	"sync"
)

type InMemoryGroup struct {
	mutex         sync.Mutex
	state         aggregate.State
	shard         *InMemoryShard
	groupByValues map[string]interface{}
	triggers      []*aggregate.Trigger
	hash          []byte
}

func MakeInMemoryGroup(hash []byte,
	groupByValues map[string]interface{},
	triggers []*aggregate.Trigger,
	shard *InMemoryShard) *InMemoryGroup {
	return &InMemoryGroup{
		hash:          hash,
		shard:         shard,
		groupByValues: groupByValues,
		triggers:      triggers,
	}
}

func (g *InMemoryGroup) Lock() {
	g.mutex.Lock()
}

func (g *InMemoryGroup) Unlock() {
	g.mutex.Unlock()
}

func (g *InMemoryGroup) Clone() (aggregate.Group, error) {
	clonedState, err := g.state.Clone()
	if err != nil {
		return nil, err
	}
	return &InMemoryGroup{
		state:         clonedState,
		shard:         g.shard,
		groupByValues: g.groupByValues,
		hash:          g.hash,
		triggers:      g.triggers,
	}, nil
}

// Returns whether a given group is initialized
func (g *InMemoryGroup) Initialized() bool {
	return g.state == nil
}

// Initialize the group
func (g *InMemoryGroup) Initialize(state aggregate.State) error {
	if g.state != nil {
		return aggregate.AlreadyInitialized
	}
	g.state = state
	return nil
}

// Return the state of the group
func (g *InMemoryGroup) State() aggregate.State {
	return g.state
}

// Return the group-by fields of a given group
func (g *InMemoryGroup) GroupByValues() map[string]interface{} {
	return g.groupByValues
}

// Get the upper time boundary for the group
func (g *InMemoryGroup) Triggers() []*aggregate.Trigger {
	return g.triggers
}

// Get the hash for the group
func (g *InMemoryGroup) Hash() []byte {
	return g.hash
}
