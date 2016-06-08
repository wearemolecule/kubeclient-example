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
	"fmt"
	"strings"

	"golang.org/x/net/context"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete pod or replication controller",
	Long: `Delete a pod or replication controller given a name. It will scope to a namespace if given.

	-- delete pod/pod-name
	-- delete rc/rc-name`,
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
			return deletePod(splitArgs[1])
		case "rc":
			return deleteRc(splitArgs[1])
		default:
			return defaultErr
		}

		return nil
	},
}

func deletePod(name string) error {
	err := kubeClient.DeletePod(context.TODO(), namespace, name)
	if err != nil {
		return fmt.Errorf("Error deleting pod: %v", err)
	}

	fmt.Printf("Deleted pod/%s successfully\n", name)
	return nil
}

func deleteRc(name string) error {
	err := kubeClient.DeleteReplicationController(context.TODO(), namespace, name)
	if err != nil {
		return fmt.Errorf("Error deleting rc: %v", err)
	}

	fmt.Printf("Deleted rc/%s successfully\n", name)
	return nil
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
