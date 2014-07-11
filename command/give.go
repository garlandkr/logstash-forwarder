// Licensed to Elasticsearch under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package command

import (
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf"
)

const cmdGive lsf.CommandCode = "give"

var giveOptions0 = struct{ n uint }{n: 3}
var Give *lsf.Command

func init() {
	Give = &lsf.Command{
		Name:        cmdGive,
		About:       "Ask and you shall receive!",
		Run:         runGive,
		Flag:        FlagSet(cmdGive),
		Initializer: true,
	}
	Give.Flag.UintVar(&giveOptions0.n, "me", giveOptions0.n, "how many")
}

func runGive(env *lsf.Environment, args ...string) (err error) {
	defer panics.Recover(&err)

	for n := 0; n < int(giveOptions0.n); n++ {
		log.Printf("h%c gs\n", rune('\u2661'))
	}
	return
}