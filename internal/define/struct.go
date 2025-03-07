/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package define

type CERT struct {
	Country            string
	State              string
	Locality           string
	Organization       string
	OrganizationalUnit string
	CommonName         string
	Domains            []string

	ForK8S  string
	OwnUser string
	OwnUID  string
	OwnGID  string
}
