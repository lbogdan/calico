// Copyright (c) 2016 Tigera, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package intdataplane

import (
	log "github.com/Sirupsen/logrus"
	"github.com/projectcalico/felix/go/felix/iptables"
	"github.com/projectcalico/felix/go/felix/proto"
	"github.com/projectcalico/felix/go/felix/rules"
)

type policyManager struct {
	filterTable  *iptables.Table
	ruleRenderer rules.RuleRenderer
}

func newPolicyManager(filterTable *iptables.Table, ruleRenderer rules.RuleRenderer) *policyManager {
	return &policyManager{
		filterTable:  filterTable,
		ruleRenderer: ruleRenderer,
	}
}

func (m *policyManager) OnUpdate(msg interface{}) {
	switch msg := msg.(type) {
	case *proto.ActivePolicyUpdate:
		chains := m.ruleRenderer.PolicyToIptablesChains(msg.Id, msg.Policy)
		m.filterTable.UpdateChains(chains)
	case *proto.ActivePolicyRemove:
		inName := rules.PolicyChainName(rules.PolicyInboundPfx, msg.Id)
		outName := rules.PolicyChainName(rules.PolicyOutboundPfx, msg.Id)
		m.filterTable.RemoveChainByName(inName)
		m.filterTable.RemoveChainByName(outName)
	case *proto.ActiveProfileUpdate:
		// TODO(smc) Profile updates
		log.WithField("msg", msg).Warn("Message not implemented")
	case *proto.ActiveProfileRemove:
		// TODO(smc) Profile updates
		log.WithField("msg", msg).Warn("Message not implemented")
	}
}

func (m *policyManager) CompleteDeferredWork() {
	// Nothing to do, we don't defer any work.
}
