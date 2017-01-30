// Copyright 2017, Google
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

package b2

import "fmt"

// StatusInfo reports information about a client.
type StatusInfo struct {
	// Writers report the status of each upload in progress.
	Writers map[string]*WriterStatus
}

// WriterStatus reports the status for each writer.
type WriterStatus struct {
	// Progress maps chunk IDs to a completion ratio between 0 and 1.
	Progress map[int]float64
}

// Status returns information about the current state of the client.
func (c *Client) Status(addr string) *StatusInfo {
	c.slock.Lock()
	defer c.slock.Unlock()

	si := &StatusInfo{
		Writers: make(map[string]*WriterStatus),
	}

	for name, w := range c.sWriters {
		si.Writers[name] = w.status()
	}

	return si
}

func (c *Client) addWriter(w *Writer) {
	c.slock.Lock()
	defer c.slock.Unlock()

	if c.sWriters == nil {
		c.sWriters = make(map[string]*Writer)
	}

	c.sWriters[fmt.Sprintf("%s/%s", w.o.b.Name, w.name)] = w
}

func (c *Client) removeWriter(w *Writer) {
	c.slock.Lock()
	defer c.slock.Unlock()

	if c.sWriters == nil {
		return
	}

	delete(c.sWriters, fmt.Sprintf("%s/%s", w.o.b.Name, w.name))
}
