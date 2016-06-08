// Copyright © 2016 Stephen Meriwether <stephen@molecule.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the pods or replication controllers",
	Long: `List all pod's or all replication controller's. It will scope to a namespace if given.

	-- list pods
	-- list rcs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return defaultErr
		}

		switch args[0] {
		case "pods":
			return listPods()
		case "rcs":
			return listRcs()
		default:
			return defaultErr
		}

		return nil
	},
}

func listPods() error {
	pods, err := kubeClient.PodList(context.TODO(), namespace, "")
	if err != nil {
		return fmt.Errorf("Error listing pods: %v", err)
	}

	msg := fmt.Sprintf("| %-30s | %-10s | %-15s | %-10s\n", "Name", "Phase", "IP", "StartTime")
	for _, pod := range pods {
		msg += fmt.Sprintf("| %-30s | %-10s | %-15s | %-10s\n",
			pod.Name, pod.Status.Phase, pod.Status.PodIP, pod.Status.StartTime,
		)
	}
	fmt.Println(msg)

	return nil
}

func listRcs() error {
	rcs, err := kubeClient.ReplicationControllerList(context.TODO(), namespace, "")
	if err != nil {
		return fmt.Errorf("Error listing replication controllers: %v", err)
	}

	msg := fmt.Sprintf("| %-30s | %-10s | %-30s\n", "Name", "Replicas", "Labels")
	for _, rc := range rcs {
		var labels []string
		for _, l := range rc.Labels {
			labels = append(labels, strings.Trim(l, " "))
		}

		msg += fmt.Sprintf("| %-30s | %-10d | %-30s\n",
			rc.Name, rc.Status.Replicas, strings.Join(labels, ", "),
		)
	}
	fmt.Println(msg)

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}
