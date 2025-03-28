// Copyright 2024 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//	         $$\                           $$\ $$\   $$\
//	         $$ |                          $$ |\__|  $$ |
//	$$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
//
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//
//	\_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package assessment

import (
	"context"
	"time"

	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
)

// waitingRequest contains all information of an evidence request that still waits for
// more data
type waitingRequest struct {
	*evidence.Evidence

	started time.Time

	// waitingFor should ideally be empty at some point
	waitingFor map[string]bool

	resourceId string

	s *Service

	newResources chan string
	ctx          context.Context
}

func (l *waitingRequest) WaitAndHandle() {
	for {
		// Wait for an incoming resource
		resource := <-l.newResources

		// Check, if the incoming resource is of interest for us
		delete(l.waitingFor, resource)

		// Are we ready to assess?
		if len(l.waitingFor) == 0 {
			log.Infof("Evidence %s is now ready to assess", l.Evidence.Id)

			// Gather our additional resources
			additional := make(map[string]ontology.IsResource)

			for _, r := range l.Evidence.ExperimentalRelatedResourceIds {
				l.s.em.RLock()

				e, ok := l.s.evidenceResourceMap[r]
				l.s.em.RUnlock()

				if !ok {
					log.Errorf("Apparently, we are missing an evidence for a resource %s which we are supposed to have", r)
					break
				}

				msg := e.GetOntologyResource()
				if msg == nil {
					break
				}

				additional[r] = msg
			}

			// Let's go
			_, _ = l.s.handleEvidence(l.ctx, l.Evidence, l.Evidence.GetOntologyResource(), additional)

			duration := time.Since(l.started)

			log.Infof("Evidence %s was waiting for %s", l.Evidence.Id, duration)
			break
		}
	}

	// Lock requests for writing
	l.s.rm.Lock()
	// Remove ourselves from the list of requests
	delete(l.s.requests, l.Evidence.Id)
	// Unlock writing
	l.s.rm.Unlock()

	// Inform our wait group, that we are done
	l.s.wg.Done()
}

// informWaitingRequests informs any waiting requests of the arrival of a new resource ID, so that they might update
// their waiting decision.
func (svc *Service) informWaitingRequests(resourceId string) {
	// Lock requests for reading
	svc.rm.RLock()
	// Defer unlock at the exit of the go-routine
	defer svc.rm.RUnlock()
	for _, l := range svc.requests {
		if l.resourceId != resourceId {
			l.newResources <- resourceId
		}
	}
}
