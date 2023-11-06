// Copyright 2022 Fraunhofer AISEC
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
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package voc

func (r *Resource) Related() []string {
	if r.Parent != "" {
		return []string{string(r.Parent)}
	}

	return nil
}

// Related returns related resources for the virtual machine, e.g., its attached storage and network interfaces.
func (v VirtualMachine) Related() []string {
	list := make([]string, 0)

	for _, b := range v.BlockStorage {
		list = append(list, string(b))
	}

	if v.BootLogging != nil {
		list = append(list, v.BootLogging.Related()...)
	}

	if v.OsLogging != nil {
		list = append(list, v.OsLogging.Related()...)
	}

	list = append(list, v.Compute.Related()...)

	return list
}

func (c Compute) Related() []string {
	list := make([]string, 0)

	for _, n := range c.NetworkInterfaces {
		list = append(list, string(n))
	}

	if c.ResourceLogging != nil {
		list = append(list, c.ResourceLogging.Related()...)
	}
	list = append(list, c.Resource.Related()...)

	return list
}

func (r *Logging) Related() []string {
	list := make([]string, 0)

	for _, n := range r.LoggingService {
		list = append(list, string(n))
	}

	return list
}

// Related returns related resources for the logging service, e.g., its storage.
func (l LoggingService) Related() []string {
	list := make([]string, 0)

	for _, s := range l.Storage {
		list = append(list, string(s))
	}

	list = append(list, l.Resource.Related()...)

	return list
}

func (s StorageService) Related() []string {
	list := make([]string, 0)

	for _, s := range s.Storage {
		list = append(list, string(s))
	}

	list = append(list, s.Resource.Related()...)

	return list
}

func (s Storage) Related() []string {
	list := make([]string, 0)

	for _, b := range s.Backups {
		if b.Storage != "" {
			list = append(list, string(b.Storage))
		}
	}

	list = append(list, s.Resource.Related()...)

	return list
}

func (a Application) Related() []string {
	list := make([]string, 0)

	for _, dep := range a.Dependencies {
		list = append(list, string(dep))
	}

	for _, tu := range a.TranslationUnits {
		list = append(list, string(tu))
	}

	list = append(list, a.Resource.Related()...)

	return list
}
