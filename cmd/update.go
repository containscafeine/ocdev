package cmd

import (
	"fmt"
	"github.com/redhat-developer/ocdev/pkg/component"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var updateCmd = &cobra.Command{
	Use: "update",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 2 {
			// TODO: Improve this message
			return fmt.Errorf("Invalid arguments, maximum 1 arguments possible")
		}
		return nil
	},
	Example: `  # Change the source of a currently active component to local (use the current directory as a source)
  ocdev update --local

  # Change the source of the frontend component to local with source in ./frontend directory
  ocdev update frontend --local ./frontend

  # Change the source of a currently active component to git 
  ocdev update --git https://github.com/openshift/nodejs-ex.git

  # Change the source of the component named node-ex to git
  ocdev update node-ex --git https://github.com/openshift/nodejs-ex.git
	`,
	Short: "Change the source of a component",
	Run: func(cmd *cobra.Command, args []string) {
		client := getOcClient()

		checkFlag := 0

		if len(componentBinary) != 0 {
			checkFlag++
		}
		if len(componentGit) != 0 {
			checkFlag++
		}
		if len(componentLocal) != 0 {
			checkFlag++
		}

		if checkFlag > 1 {
			fmt.Println("The source can be either --binary or --local or --git")
			os.Exit(1)
		}

		if len(componentBinary) != 0 {
			fmt.Printf("--binary is not implemented yet\n\n")
			os.Exit(1)
		}

		var (
			componentName string
			err           error
		)
		if len(args) == 0 {
			componentName, err = component.GetCurrent(client)
			if err != nil {
				fmt.Println("Unable to get current component")
			}
		} else {
			componentName = args[0]
		}

		if len(componentGit) != 0 {
			err := component.Update(client, componentName, "git", componentGit)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("The component %s was updated successfully\n", componentName)
		} else if len(componentLocal) != 0 {
			// we want to use and save absolute path for component
			dir, err := filepath.Abs(componentLocal)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = component.Update(client, componentName, "dir", dir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("The component %s was updated successfully\n", componentName)
		} else {
			// we want to use and save absolute path for component
			dir, err := filepath.Abs("./")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = component.Update(client, componentName, "dir", dir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("The component %s was updated successfully\n", componentName)
		}
	},
}

func init() {
	updateCmd.Flags().StringVar(&componentBinary, "binary", "", "binary artifact")
	updateCmd.Flags().StringVar(&componentGit, "git", "", "git source")
	updateCmd.Flags().StringVar(&componentLocal, "local", "", "Use local directory as a source for component.")

	rootCmd.AddCommand(updateCmd)
}