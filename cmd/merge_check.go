package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/LiZeC123/gmh/util"
	"github.com/urfave/cli/v3"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
)

func MergeCheckCommand() *cli.Command {
	return &cli.Command{
		Name:  "mck",
		Usage: "Check if local git branches are merged into master",
		Action: func(ctx context.Context, c *cli.Command) error {
			// 检查当前目录是否是git仓库
			if !isGitRepository() {
				return errors.New("not a git repository")
			}

			// 获取所有本地分支
			branches, err := getLocalBranches()
			if err != nil {
				return errors.New("get local branches failed")
			}

			// 检查master分支是否存在
			hasMaster, err := checkMasterBranchExists()
			if err != nil || !hasMaster {
				return errors.New("check master branch failed")
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
			for _, branch := range branches {
				// 跳过master分支自身
				if branch == "master" {
					continue
				}

				hasUnmerged, aheadCount, behindCount, err := checkBranchStatus(branch)
				if err != nil {
					fmt.Printf("❌ 检查分支 %s 失败: %v\n", branch, err)
					continue
				}
				util.PrintToFile(w, "%s", makeDisplayBranchStatusString(branch, hasUnmerged, aheadCount, behindCount))
			}

			err = w.Flush()
			if err != nil {
				return err
			}

			return nil
		},
	}
}

// 检查当前目录是否是git仓库
func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// 获取所有本地分支
func getLocalBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--list", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var branches []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		branch := strings.TrimSpace(scanner.Text())
		if branch != "" {
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

// 检查master分支是否存在
func checkMasterBranchExists() (bool, error) {
	cmd := exec.Command("git", "show-ref", "--verify", "refs/heads/master")
	err := cmd.Run()
	return err == nil, nil
}

// 检查分支状态
func checkBranchStatus(branch string) (bool, int, int, error) {
	// 获取分支在master之前的所有提交
	cmd := exec.Command("git", "log", "master.."+branch, "--oneline")
	aheadOutput, _ := cmd.Output()
	aheadCommits := strings.Split(strings.TrimSpace(string(aheadOutput)), "\n")
	aheadCount := 0
	for _, commit := range aheadCommits {
		if strings.TrimSpace(commit) != "" {
			aheadCount++
		}
	}

	// 获取分支落后master的提交数
	cmd = exec.Command("git", "log", branch+"..master", "--oneline")
	behindOutput, _ := cmd.Output()
	behindCommits := strings.Split(strings.TrimSpace(string(behindOutput)), "\n")
	behindCount := 0
	for _, commit := range behindCommits {
		if strings.TrimSpace(commit) != "" {
			behindCount++
		}
	}

	// 如果有在master之前的提交，说明有未合并的内容
	hasUnmerged := aheadCount > 0

	return hasUnmerged, aheadCount, behindCount, nil
}

// 显示分支状态
func makeDisplayBranchStatusString(branch string, hasUnmerged bool, aheadCount, behindCount int) string {
	statusText := ""
	if hasUnmerged {
		statusText = "🔴 未合入"
	} else {
		statusText = "🟢 已合入"
	}

	if strings.HasPrefix(branch, "*") {
		branch = branch[1:] // 移除*号
	}

	countStr := ""
	if aheadCount > 0 {
		countStr = fmt.Sprintf("🔼 领先 %d 个提交 ", aheadCount)
	}
	if behindCount > 0 {
		countStr = countStr + fmt.Sprintf("🔽 落后 %d 个提交", behindCount)
	}

	return fmt.Sprintf("%s\t%s\t%s\t\n", branch, statusText, countStr)

}
