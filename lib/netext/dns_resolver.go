/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package netext

import (
	"errors"
	"fmt"
	"net"
)

type DnsResolver struct {
	disableIpv6 bool
	cache       map[string]net.IP
}

func NewDnsResolver(disableIpv6 bool) *DnsResolver {
	return &DnsResolver{
		disableIpv6: disableIpv6,
		cache:       make(map[string]net.IP, 0),
	}
}

func (d *DnsResolver) Resolve(host string) (net.IP, error) {
	if ip, ok := d.cache[host]; ok {
		fmt.Printf("cache %s \n", ip.String())
		return ip, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	if d.disableIpv6 {
		for _, ip := range ips {
			if ip.To4() != nil {
				d.cache[host] = ip
				fmt.Printf("ipv4 found: %s \n", ip.String())

				return ip, nil
			}
		}

		newIps, err := net.LookupIP(host)
		if err != nil {
			return nil, err
		}

		for _, ip := range newIps {
			if ip.To4() != nil {
				d.cache[host] = ip
				fmt.Printf("ipv4 found in the second time: %s \n", ip.String())

				return ip, nil
			}
		}

		return nil, errors.New("ipv4 not found")
	} else {
		fmt.Printf("Fall back first ip %s \n", ips[0].String())
		return ips[0], nil
	}
}
