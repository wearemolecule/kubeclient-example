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
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearemolecule/kubeclient"
)

var (
	namespace  string
	kubeClient *kubeclient.Client
	defaultErr error
)

var RootCmd = &cobra.Command{
	Use:   "kubeclient-example",
	Short: "A minimal kubectl clone using the github.com/wearemolecule/kubeclient library",
	Long: `The kubeclient library (github.com/wearemolecule/kubeclient) provides an easy way to interact with a kubernetes cluster inside of a golang app.
	This cli tool was created to demonstrate how to use the library.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	var err error
	kubeClient, err = kubeclient.GetKubeClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to kubernetes, error: ", err)
	}

	defaultErr = fmt.Errorf("Unable to process command use --help for more info")

	RootCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "If present, the namespace scope for this CLI request")
}
