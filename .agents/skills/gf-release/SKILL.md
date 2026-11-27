---
name: gf-release
description: >-
  为 GoFrame（gogf/gf）仓库执行正式发版：用户只说发版时先给出提交统计并用结构化选项让用户选择修正版本、子版本或主版本，再建发版分支、统一提升所有模块版本、提交并发起发版 PR，合并后在 master 上打注解 tag，核对自动发布出来的 GitHub Release 与子模块 tag，合并机器人自动提交的 CLI 依赖 PR，以及整理 Release Notes 与站点发布文档。
  依赖 Makefile 的 branch、version、tag 三个目标与 .make_version.sh 脚本，以及 release.yml、tag.yml 两条自动发布流水线。
  必须由用户手动触发，禁止自动运行。
compatibility: 需要 git、已登录的 GitHub CLI `gh`、可用的 Go 工具链；刷新内嵌模板需要已安装 `gf` CLI 并能访问 github.com；`make version` 会在每个模块执行 `go mod tidy`，需要模块代理可用。
---

# GF Release

按仓库自带的发版机制发布一个新版本。这是执行型技能：触发后要真的建分支、改版本号、提 PR、打 tag，而不是只给命令建议。

## 核心规则

1. 默认仓库是`gogf/gf`，发版基线是`master`，tag 形如`v2.10.6`。
2. 只能从最新的`master`出发，工作区必须干净；不在`master`或存在未提交改动时先处理再继续。
3. 用户只说「发版」「发布新版本」而没有指定版本类型时，必须先给出提交统计并用**结构化选项**让用户在修正版本、子版本、主版本中选择。选项文案必须写成按最新 tag 算出的真实版本号（例如`修正版本 v2.10.6`、`子版本 v2.11.0`、`主版本 v3.0.0`），不能用占位符；用户选中后把该版本号用于后续所有步骤。用户已明确给出类型或版本号时直接采用，不要重复提问。不要自行替用户决定升次版本或主版本。
4. 版本号只允许通过`make version`（内部`.make_version.sh`）修改，禁止手工散改各模块版本号；脚本漏改时先定位原因，再做最小补齐。
5. 禁止在功能分支上打 tag，禁止用`git tag -f`覆盖已有 tag，禁止删除或移动已发布的 tag。
6. 每一步都要先核对上一步的真实输出再进入下一步。任何一步失败就停下报告，不要带着失败继续打 tag。
7. 发版 PR 的标题固定为`chore(release): bump version to vX.Y.Z`，标题遵守`.github/PULL_REQUEST_TEMPLATE.MD`并用英文书写。
8. 打 tag 后必须合并机器人自动提交的 CLI 依赖 PR（`fix: update gf cli to vX.Y.Z`，主要改动是`cmd/gf/go.sum`），不要手工改`cmd/gf/go.sum`代替这一步，也不要把「PR 已出现」当作「已合并」。

## 前置检查

```bash
git rev-parse --show-toplevel
git branch --show-current
git status --short
git fetch origin --tags --prune
git describe --tags --abbrev=0
git log --oneline "$(git describe --tags --abbrev=0)..origin/master"
gh auth status
```

同时确认：

- `gf -v`可用。不可用且本次需要刷新内嵌模板时，执行`cd cmd/gf && go install .`。
- `go env GOPROXY`指向可用的模块代理，脚本会对每个模块跑`go mod tidy`。
- 记录上个 tag、`origin/master`相对它的提交总数，以及提交类型分布（`feat`、`fix`、`refactor`、`test`、`docs`、`ci`、`chore` 各多少条），供下一步提问使用。

```bash
git log --oneline "v<上个版本>..origin/master" | grep -oE "^[a-f0-9]+ (feat|fix|refactor|test|docs|ci|chore|perf|build)" | awk '{print $2}' | sort | uniq -c
```

## 工作流

### 1. 确认版本类型与版本号

用户没有指定版本类型时，**先提问再动手**，不要先建分支或改版本号。使用当前环境的结构化选择题工具（Claude 用`AskUserQuestion`，其他环境用对应的提问工具或可选择的工具）；工具不可用时列出三个选项等待用户回复编号。

提问前必须先按最新 tag 算出三个候选版本号，**选项文案里直接用算好的真实版本号**，禁止出现`vX.Y.(Z+1)`、`vX.Y.Z`这类占位符：

```bash
latest=$(git describe --tags --abbrev=0)                # 例如 v2.10.5
base=${latest#v}
major=${base%%.*}; rest=${base#*.}; minor=${rest%%.*}; patch=${rest#*.}; patch=${patch%%-*}
echo "修正版本 v${major}.${minor}.$((patch+1))"
echo "子版本 v${major}.$((minor+1)).0"
echo "主版本 v$((major+1)).0.0"
```

问题：`本次发布哪种版本？`

选项（以最新 tag 为`v2.10.5`为例，实际必须替换成算出来的真实版本号）：

1. `修正版本 v2.10.6`（推荐）— 只有`fix`、`test`、`docs`、`chore`、`refactor`、`perf`，或向后兼容的小特性
2. `子版本 v2.11.0` — 成批新功能、较大重构或存在兼容性调整
3. `主版本 v3.0.0` — 破坏性变更或重大架构调整

提问前先用一两句话给出依据：上个 tag、提交总数、类型分布，以及三个选项对应的具体版本号。推荐项放在第一项并标注`(Recommended)`/`（推荐）`，默认推荐修正版本。

判断参考：

| 自上个 tag 以来的改动 | 建议版本 |
| --- | --- |
| 只有`fix`、`test`、`docs`、`chore`、`refactor`、`perf` | 修正版本，补丁号加一 |
| 含成批的`feat`新能力 | 子版本，次版本号加一、补丁号归零 |
| 含破坏性变更 | 主版本，需要维护者额外确认，`v2`主线通常不升主版本 |

用户选定后，把选项里的那个真实版本号作为后续命令中的`vX.Y.Z`使用，并在同一轮回复里说明所选类型与依据。用户直接给了版本号或类型时跳过提问。后续回复中也必须写具体版本号，不要用占位符代替。

### 2. 建发版分支

```bash
make branch to=vX.Y.Z
```

该目标会切到`master`、拉取最新，再创建`fix/vX.Y.Z`。省略`to=`时自动取最新 tag 的补丁号加一，但次版本和主版本必须显式传入。

### 3. 提升版本号

```bash
make version to=vX.Y.Z
```

等价于`make -C cmd/gf pack`加`./.make_version.sh ./ vX.Y.Z`：

- `make -C cmd/gf pack`会克隆`gogf/template-single`与`gogf/template-mono`，重新生成`cmd/gf/internal/packed/template-*.go`。本次不需要刷新模板时，可以跳过这一半，直接执行`./.make_version.sh ./ vX.Y.Z`。
- `.make_version.sh`改写`version.go`里的`VERSION`、`README.MD`与`README.zh_CN.MD`的 contributors 徽章；遍历除`testdata/`、`examples/`之外的所有`go.mod`，删除`// indirect`与`toolchain`行、执行`go mod tidy`、把非间接的`github.com/gogf/gf*`依赖`go get`到新版本、再 tidy 一次。
- 处理`cmd/gf`时会临时加上`cmd/gf/go.work`里的 replace 指令，结束后自动移除。

必须显式传`to=vX.Y.Z`。漏传时`.make_version.sh`会因参数数量不符直接报错退出，此时检查命令后重跑即可，不要手工补版本号。

### 4. 核对改动

```bash
git status --short
git diff --stat
git diff version.go README.MD README.zh_CN.MD
grep -rn "v<旧版本>" --include="go.mod" --include="*.go" . | grep -v "/testdata/"
cd cmd/gf && go build ./... && cd -
go build ./...
```

核对清单：

- `version.go`的`VERSION`与新版本一致。
- `README.MD`、`README.zh_CN.MD`徽章里的`?version=`已更新。
- 根模块、`contrib/**`与`cmd/gf`的`go.mod`依赖都指向新版本，`go.sum`只有合理变动。
- 需要刷新模板时，`cmd/gf/internal/packed/template-*.go`已更新。
- 除`cmd/gf/internal/cmd/testdata/`下的固定测试夹具外，没有旧版本残留。
- 改动只包含版本相关文件。`go mod tidy`带入的无关依赖升级要逐条确认，可疑就回退该模块后重跑。

核对不通过时修正后重跑对应步骤，不要直接提交。

### 5. 提交并发起发版 PR

```bash
git add -A
git commit -m "chore(release): bump version to vX.Y.Z"
git push -u origin HEAD
gh pr create -R gogf/gf --base master --head "fix/vX.Y.Z" \
  --title "chore(release): bump version to vX.Y.Z" --body-file <正文文件>
```

正文沿用历史发版 PR 的两段结构：

```markdown
## What this does

Prepares the `vX.Y.Z` release by bumping every module's inter-dependent version requirement.

- `version.go`: `VERSION` → `vX.Y.Z`
- `contrib/*/go.mod`: framework and sibling `contrib` requirements → `vX.Y.Z`
- `cmd/gf/go.mod` / `go.sum`: framework and driver requirements → `vX.Y.Z`
- `cmd/gf/internal/packed/template-*.go`: refreshed embedded templates
- `README.MD` / `README.zh_CN.MD`: contributor badge `?version=vX.Y.Z`

## Verification

- `go build ./...` in `cmd/gf` passes
- no remaining `v<旧版本>` references outside intentional test fixtures
```

正文语言按用户偏好，默认英文以与历史发版 PR 保持一致。CI 未通过时不要催促合并。

### 6. 合并后打 tag

确认发版 PR 已合并、`master`指向该合并提交后执行：

```bash
git checkout master && git pull origin master && git log --oneline -1
make tag to=vX.Y.Z
```

`make tag`会自己切到`master`、拉取最新，再创建注解 tag（消息为`Release vX.Y.Z`）并推送到`origin`。

### 7. 等待发布流水线完成

打 tag 后不要立刻继续，先确认`release.yml`与`tag.yml`都执行完成：

```bash
gh run list -R gogf/gf --limit 10 --json databaseId,workflowName,status,conclusion,headBranch
gh run watch <release.yml 的 run id> -R gogf/gf --exit-status
```

`gh run watch`不可用或长时间无输出时，改为轮询`gh run list`，直到对应运行的`status`为`completed`并记录`conclusion`。全平台二进制构建通常需要数分钟，期间不要重复推送 tag，也不要提前编辑 Release 页面。

必须同时核对：

| 产物 | 来源 | 核对点 |
| --- | --- | --- |
| GitHub Release | `.github/workflows/release.yml` | 结论为`success`，非草稿，含`gf_linux_amd64`、`gf_windows_amd64.exe`等全平台二进制 |
| 子模块 tag | `.github/workflows/tag.yml` | `contrib/**/vX.Y.Z`与`cmd/gf/vX.Y.Z`都已创建 |
| CLI 依赖 PR | `.github/workflows/tag.yml` | 自动出现`fix: update gf cli to vX.Y.Z`，必须按步骤 8 合并；改动为`cmd/gf/go.sum`补齐新版本校验和 |

```bash
gh release view vX.Y.Z -R gogf/gf --json tagName,isDraft,publishedAt,assets
git ls-remote --tags origin | grep "vX.Y.Z$" | wc -l
gh pr list -R gogf/gf --state open --json number,title,headRefName
```

任一流水线失败就停下，用`gh run view <id> -R gogf/gf --log-failed`定位；流水线未完成或失败时，不要进入后续步骤，也不要把未完成描述成已完成。

### 8. 合并机器人 CLI 依赖 PR

`tag.yml`除了补子模块 tag，还会用机器人账号（`TagRobot` / `app/github-actions`）自动创建标题为`fix: update gf cli to vX.Y.Z`、头分支`fix/vX.Y.Z`、base 为`master`的 PR。它的实际改动通常只有`cmd/gf/go.sum`（补齐新 tag 的校验和）和`cmd/gf`下`testdata`里的固定`go.mod`夹具：发版 PR 已经改过`cmd/gf/go.mod`，但那时子模块 tag 还不存在，`go.sum`无法记录`vX.Y.Z`的校验和，只能等`tag.yml`脚本（`.github/workflows/scripts/update_version.sh ./cmd/gf vX.Y.Z`）跑完`go mod tidy`后补上。**这一步不能省**：不合并的话`cmd/gf/go.sum`缺少新版本校验和，后续发版会反复出现同样的 PR，干净环境构建也会失败。发版人自己合并，不要指派给其他贡献者。

```bash
gh pr list -R gogf/gf --state open --json number,title,headRefName,author
gh pr view <PR 号> -R gogf/gf --json title,state,mergeable,mergeStateStatus,files,statusCheckRollup
gh pr merge <PR 号> -R gogf/gf --squash --delete-branch
```

合并前确认：标题与头分支对应本次 tag、作者是机器人、CI 已通过（`statusCheckRollup`无失败项）、没有冲突；`files`里不应出现`cmd/gf/go.mod`的版本回退。

合并后核对`master`：

```bash
git checkout master && git pull origin master && git log --oneline -3
grep -n "github.com/gogf/gf" cmd/gf/go.mod | head
grep -c "github.com/gogf/gf.*vX.Y.Z" cmd/gf/go.sum
```

`cmd/gf/go.mod`中`github.com/gogf/gf/v2`以及`contrib/**`依赖都应为`vX.Y.Z`，`cmd/gf/go.sum`里每个子模块都已有`vX.Y.Z`的条目（本次参考量为 16 条，随模块数量变化），`master`顶部出现该机器人提交。PR 没出现时先按步骤 7 确认`tag.yml`成功，再等机器人推送；始终不要手工改`cmd/gf/go.sum`或手工建这个 PR。

### 9. 生成并整理 Release Notes

Release 页面默认只有自动生成的原始变更列表，需要补一版可读的说明。这一步由技能执行者引导操作人员完成，不要自己直接改页面内容。

1. 通知操作人员：等步骤 7 的`release.yml`完成后，打开`https://github.com/gogf/gf/releases/tag/vX.Y.Z`，点击编辑，在正文处点击`Generate release notes`（自动生成发布说明），确认对比范围是上一版本`v<上一版本>`到本次 tag。
2. 请操作人员把生成出来的正文整段贴回对话。可用同样的数据源自检，但最终以页面生成的结果为准：

```bash
gh api repos/gogf/gf/releases/generate-notes -f tag_name=vX.Y.Z -f previous_tag_name=v<上一版本> --jq '.body'
```

3. 对贴回的内容做分类整理。不新增、不删除条目，不改写作者：
   - `Features`：`feat`类型的 PR。
   - `Bug Fixes`：`fix`类型的 PR，按组件分组（`net/ghttp`、`database/gdb`、`os/gtime`等），组件名写完整路径。
   - `Performance`：`perf`类型。
   - `Refactor`：`refactor`类型。
   - `Docs` / `Tests`：`docs`、`test`类型，条目少时可以合并。
   - `CI & Dependencies`：`ci`、`build`、`chore`类型，以及`fix: update gf cli to vX.Y.Z`这类机器人提交。
   - 安全相关修复（路径穿越、越权读取、panic、数据竞争等）在条目中标注`Security`，不要漏。
   - 每个条目保留`[#NNNN](链接)（@作者）`格式；同一组件内按 PR 号或重要性排序。
4. 输出为纯 markdown，保持英文（与页面现有内容和历史 Release 一致），保留末尾的`**Full Changelog**`链接，删除模板残留与无关内容。
5. 把整理结果交回操作人员，由其在页面替换正文后点击`Update release`；保存后核对：

```bash
gh release view vX.Y.Z -R gogf/gf --json body --jq '.body' | head -40
```

只改正文，不要改 tag、附件和发布状态，不要删除自动生成的真实条目。

### 10. 同步站点发布文档（可选）

`gf-site`仓库可用时（默认与本仓库同级，例如`../gf-site`），把同一批变更整理成站点发布文档，分类口径与步骤 9 保持一致。

1. 选择目录：补丁版本放在`docs/release/v<major>.<minor>.x/`（例如`docs/release/v2.10.x/`）；新子版本放在`docs/release/`根目录，命名为`v<major>.<minor> YYYY-MM-DD.md`。
2. 文件名用`vX.Y.Z YYYY-MM-DD.md`，日期取正式发布日（与 Release 的`publishedAt`一致），不要写草稿日期。
3. `front matter`按同目录既有文件填写：`slug: '/release/vX.Y.Z'`、`title: 'vX.Y.Z YYYY-MM-DD'`、`sidebar_position`（续同目录序号）、`hide_title: true`、`keywords`、`description`。
4. 正文沿用同目录结构：导语 H2（概括本版本重点）→ 完整变更与代码差异链接 → `## 新功能` → `## 问题修复`（按组件分`###`）→ `## 性能优化` → `## 代码重构` → `## CI 与依赖` → `## 新贡献者` → `## 升级建议`（升级命令固定写`go get -v github.com/gogf/gf/v2@vX.Y.Z`，不要带`-u`，避免顺手升级无关的传递依赖；后面再列需要重点验证的能力）。
5. 按站点规范补英文镜像：`i18n/en/docusaurus-plugin-content-docs/current/release/`下保持相同相对路径，翻译正文与`keywords`、`description`，markdown 链接保留中文目录名。
6. 在`gf-site`执行`make check`确认中英文文档一一对应，再提交 PR。只新增本版本文档，不要顺手修改历史版本文档。

## 常见问题

| 现象 | 处理 |
| --- | --- |
| `make version`在`gf pack`阶段报`gf: command not found` | `cd cmd/gf && go install .`，或本次无需刷新模板时直接跑`./.make_version.sh ./ vX.Y.Z` |
| 用户只说「发版」但没有指定版本类型 | 先给提交统计并用结构化选项询问，选项文案写成真实版本号（如`修正版本 v2.10.6`），未选中前不要建分支、不要改版本号 |
| 某个`go.mod`没有升到新版本 | 检查该模块是否有`replace github.com/gogf/gf/v2 => ...`以及`go mod tidy`报错，修好后在该目录单独重跑 |
| `go.sum`出现大量无关升级 | 回退该模块改动单独 tidy 复查，可疑变更不要带进发版提交 |
| tag 推送后没有 Release | 检查`release.yml`运行记录，确认 tag 以`v`开头且指向`master`上的提交 |
| 子模块 tag 缺失 | `tag.yml`只在 tag push 时触发，可重跑该 workflow，但不要手工给其他提交补 tag |
| 机器人 CLI 依赖 PR 迟迟不出现 | 先看`tag.yml`是否成功；成功后仍无 PR 时重跑该 workflow，不要手工建 PR 或手工改`cmd/gf/go.sum` |
| 机器人 CLI 依赖 PR 有冲突或 CI 失败 | 在该头分支上重跑`.github/workflows/scripts/update_version.sh ./cmd/gf vX.Y.Z`后推送，或重跑`tag.yml`，修好再合并，不要跳过一个有冲突的 PR |
| `gh run watch`长时间无输出 | 改用`gh run list`轮询`status`与`conclusion`，不要在流水线未完成时编辑 Release |
| 自动生成的发布说明对比范围不对 | 指定`previous_tag_name`为上一个正式 tag 重新生成，不要手工拼 PR 列表 |
| 操作人员贴回的发布说明缺少作者或链接 | 用`generate-notes`接口结果对照补全，仍对不上就请操作人员重新生成 |
| 站点`make check`报缺英文文档 | 在`i18n/en/docusaurus-plugin-content-docs/current/release/`下补同路径英文镜像 |
| 站点文档日期与实际发布日不一致 | 同时改文件名、`slug`与`title`，保持三者一致 |
| 版本号发错 | tag 尚未发布内容且无人升级时，与维护者确认后可删远端 tag 并修正；已经发布则只能发新版本 |

## 输出约定

完成后报告：确定的版本号与依据、发版分支名、bump 提交与 PR URL、tag 名、`release.yml`与`tag.yml`的结论、Release 与子模块 tag 的核对结果、机器人 CLI 依赖 PR 的编号与合并结果、整理后的 Release Notes 分类结果、站点发布文档的中英文路径（如有），以及所有未完成或需要人工确认的事项。

不要输出 token 或完整 diff；没有实际执行成功的命令不要描述成已完成。
