// Copyright © 2016 Molecule Software <devs@molecule.io>
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/build/kubernetes/api"
	"golang.org/x/net/context"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pod or replication controller from JSON config",
	Long: `Create a pod or replication controller from config. Given a local json file name it will read/parse the config.

	-- create pod/pod.json
	-- create rc/replication-controller.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return defaultErr
		}

		splitArgs := strings.Split(args[0], "/")
		if len(splitArgs) < 2 {
			return defaultErr
		}

		switch splitArgs[0] {
		case "pod":
			return createPod(splitArgs[1])
		case "rc":
			return createRc(splitArgs[1])
		default:
			return defaultErr
		}

		return nil
	},
}

func createPod(fileName string) error {
	ctx := context.TODO()

	podData, err := ioutil.ReadFile("pod.json")
	if err != nil {
		return fmt.Errorf("Error reading pod: %v", err)
	}

	pod := api.Pod{}
	err = json.Unmarshal(podData, &pod)
	if err != nil {
		return fmt.Errorf("Error parsing pod: %v", err)
	}

	pod.ObjectMeta.Namespace = namespace

	newPod, err := kubeClient.CreatePod(ctx, &pod)
	if err != nil {
		return fmt.Errorf("Error creating pod: %v", err)
	}

	statuses, err := kubeClient.WatchPod(ctx, newPod.Namespace, newPod.Name, newPod.ResourceVersion)
	if err != nil {
		return fmt.Errorf("Error watching task pod: %v", err)
	}

	for status := range statuses {
		podStatus := status.Pod.Status
		if podStatus.Phase == "Failed" {
			_ = kubeClient.DeletePod(ctx, newPod.Namespace, newPod.Name)
			fmt.Printf("Task pod failed (attempting to delete pod).\n")
			return nil
		}
		if podStatus.Phase == "Succeeded" {
			fmt.Printf("Task pod succeeded.\n")
			return nil
		}
	}

	return nil
}

func createRc(fileName string) error {
	ctx := context.TODO()

	rcData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Error reading rc: %v", err)
	}

	repc := api.ReplicationController{}
	err = json.Unmarshal(rcData, &repc)
	if err != nil {
		return fmt.Errorf("Error parsing rc: %v", err)
	}

	repc.ObjectMeta.Namespace = namespace

	_, err = kubeClient.CreateReplicationController(ctx, &repc)
	if err != nil {
		return fmt.Errorf("Error creating rc: %v", err)
	}

	fmt.Println("Successfully created replication controller")

	return nil
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
