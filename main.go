package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mop",
	Short: "A CLI tool to manage node_modules and target folders",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists out all node_modules and target folders",
	Run:   listFolders,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes all node_modules and target folders",
	Run:   deleteFolders,
}

func init() {
	rootCmd.AddCommand(listCmd, deleteCmd)
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.0f MB", float64(bytes)/float64(div))
}

func findFolders() ([]string, int64) {
	var folders []string
	var totalSize int64

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && (info.Name() == "node_modules" || info.Name() == "target") {
			if size, err := getDirSize(path); err == nil {
				folders = append(folders, path)
				totalSize += size
			}
		}
		return nil
	})

	return folders, totalSize
}

func listFolders(cmd *cobra.Command, args []string) {
	folders, totalSize := findFolders()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "PATH\tDIRECTORY SIZE")
	fmt.Fprintln(w, strings.Repeat("-", 50))

	for _, folder := range folders {
		size, _ := getDirSize(folder)
		fmt.Fprintf(w, "./%s\t%s\n", folder, formatSize(size))
	}

	fmt.Fprintln(w, strings.Repeat("-", 50))
	fmt.Fprintf(w, "TOTAL\t%s\n", formatSize(totalSize))
	fmt.Fprintln(w, strings.Repeat("-", 50))
	w.Flush()
}

func deleteFolders(cmd *cobra.Command, args []string) {
	folders, totalSize := findFolders()

	for _, folder := range folders {
		os.RemoveAll(folder)
		fmt.Printf("Deleted ./%s\n", folder)
	}

	fmt.Printf("\n%s freed\n", formatSize(totalSize))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
