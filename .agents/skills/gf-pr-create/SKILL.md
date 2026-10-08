---
name: gf-pr-create
description: >-
  把当前仓库的本地修改放到独立功能分支上，按 GoFrame PR 模板创建并打开合并到主仓库的 GitHub Pull Request。
  必须由用户手动触发，禁止自动运行。
compatibility: 需要 git，以及已登录的 GitHub CLI gh（优先）或可用的 GitHub MCP。创建 PR 需要向 head 仓库推送分支、向目标仓库开 PR 的权限。
---

# GF PR Create

把当前工作区的改动提交成一条指向主仓库的 GitHub Pull Request。这是执行型技能：触发后要真的建分支、提交、推送并开 PR，而不是只给命令建议。

## 核心规则

1. 默认目标仓库是`gogf/gf`，默认目标分支是`master`。用户另有指定时以用户为准。
2. 禁止在主分支上直接改内容、commit、amend、rebase 或 push。主分支包括`master`和`main`。
3. 提交前，当前改动必须位于一条独立的非主分支上。若当前就在主分支，先切出新分支再提交。
4. 提交前必须用**结构化选项**询问用户：本次 PR **正文**用英文还是中文。该选择只影响正文，不影响标题。未得到选择前，不得生成正文或创建 PR。
5. PR 标题始终遵守`.github/PULL_REQUEST_TEMPLATE.MD`，始终用英文书写，不随步骤 3 的语言选项变化。每次执行都要重新读取该文件，不要靠记忆。
6. PR 正文必须应用模板规则，但发出去的正文里不得再保留模板中的说明清单。
7. 标题、正文和本次新产生的 commit subject 都根据实际 diff 生成，不要只根据用户一句话臆造。
8. 不要 force push，不要改写已经推送到远程的历史。不要调用 Copilot 代写 PR 的工具。
9. 创建或复用 PR 之后，把目标仓库最近活跃的 5 名可审查社区成员加为 reviewer。加不上时保留已创建的 PR，不要假装已经请求成功。

## 工作流

按顺序执行。任何一步无法安全继续时，说明原因并停止，不要假装已经开出 PR。

### 1. 只读检查

```bash
git rev-parse --show-toplevel
git branch --show-current
git status --short --branch
git remote -v
git diff --stat
git diff --cached --stat
git log --oneline --decorate -n 20
```

同时读取：

- `.github/PULL_REQUEST_TEMPLATE.MD`（标题和正文规则的唯一来源）
- `CONTRIBUTING.md`（PR 应对着`master`打开）

确认 GitHub 访问：

```bash
gh auth status
gh api user --jq .login
```

`gh`不可用时，改用 GitHub MCP 完成只读确认和后续开 PR。权限不够就停在已确认的部分，不要假装 PR 已创建。

识别目标仓库：

```bash
gh repo view --json nameWithOwner,isFork,parent,defaultBranchRef
```

- 当前仓库是`gogf/gf`的 fork：目标仓库用`parent`（通常是`gogf/gf`）
- 否则：目标仓库用`origin`对应的 GitHub 仓库；在本项目中通常就是`gogf/gf`
- 目标分支默认`master`；只有用户明确要求时才改成其他分支

没有可提交内容时停止：工作区干净，且当前分支相对目标分支也没有超前提交。

`HEAD` 处于 detached 状态时，先进入步骤 2 建分支，不要在 detached HEAD 上 commit 或 push。

### 2. 独立分支

受保护的主分支：`master`、`main`。

**当前在主分支，或 detached HEAD：**

1. 根据实际改动生成短 ASCII kebab-case 分支名，例如`fix/os-gtime-timezone`、`docs/readme`、`chore/remove-sonar-properties`。
2. 名称已存在则追加`-2`、`-3`。不要使用`master`、`main`、`HEAD`。
3. 从当前工作区切出新分支，带走未提交改动：

```bash
git switch -c "$branch_name"
```

4. 若刚才离开的是主分支，且该主分支相对`origin/<main>`还有本地提交，把主分支指针恢复到远程主分支，避免改动留在主分支上。此时必须已经在新分支上：

```bash
git branch -f master origin/master
```

只移动主分支指针，不要`reset --hard`，不要丢弃未提交改动。新分支已经指向原来的提交，这些提交不会丢。

**当前已经在非主分支：**

沿用该分支，不要再套一层。用户指定了新分支名时除外。

此后所有 commit 和 push 只针对这条功能分支。

禁止：

```bash
git commit ...          # 当仍在 master/main 上
git push origin master
git push origin HEAD:master
git push origin main
```

不要调用`git-commit-push`。

### 3. 询问正文语言

分支就绪、确认有内容可提交之后，**立刻**用结构化选择题询问，然后等待回答。

使用当前环境的选择题工具：Grok 用`ask_user_question`，Claude 用`AskUserQuestion`。不要用开放式问题代替选项。工具不可用时，列出下面两个选项并等待用户回复编号。

问题：`本次 PR 正文使用哪种语言？`

选项：

1. `English` — PR 正文使用英文
2. `中文` — PR 正文使用中文

当前对话主要是中文时，把`中文`放第一项并标为推荐；主要是英文时，把`English`放第一项并标为推荐。

即使触发语里已经出现「用中文」或「use English」，仍然要弹出这两个选项；可以把用户已表达的偏好标成推荐项。未收到选项结果前，不要写 PR 正文、不要开 PR。

语言**只影响 PR 正文**。标题和本次新产生的 commit subject 始终按模板用英文书写，不受该选项影响。正文里的`Fixes`/`Updates`关键字始终保持英文，GitHub 靠它们关联 issue。

### 4. 提交本地改动

有未提交改动时，在功能分支上暂存并提交：

```bash
git add -A
git commit -m "<title>"
```

Commit subject 与即将使用的 PR 标题相同，始终用英文，并遵守同一套标题规则。不要提交空 commit。用户明确排除的文件不要加进去。

已经有提交、工作区干净时，不要为了开 PR 再造一个空 commit。有对应 issue 且需要写进 commit body 时，把`Fixes #1234`或`Updates #1234`放在 body 里。

### 5. 推送功能分支

```bash
git push -u origin HEAD
```

推送失败时报告原始错误并停止。远程分支被别人更新时，不要自动 rebase、merge 或 force push。

### 6. 创建 Pull Request

先查是否已有同一 head 的开放 PR，避免重复开：

```bash
gh pr list -R "$TARGET_REPO" --head "$HEAD_FILTER" --state open --json number,title,url,baseRefName
```

同一仓库推送：`$HEAD_FILTER`为分支名。从 fork 往主仓库开：`$HEAD_FILTER`为`owner:branch`。

已有开放 PR 时，推送更新后报告已有 URL，不要再开一条，除非用户明确要求新开。

新建：

```bash
gh pr create \
  -R "$TARGET_REPO" \
  --base master \
  --head "$HEAD_REF" \
  --title "$TITLE" \
  --body-file "$BODY_FILE"
```

同一仓库：`--head`为分支名。Fork：`--head owner:branch`。

`gh`不可用时，用 GitHub MCP 的`create_pull_request`（先查当前环境里的工具 schema，不要猜参数名）：

- `owner` / `repo`：目标仓库，例如`gogf` / `gf`
- `base`：`master`
- `head`：同一仓库用分支名；跨 fork 用`owner:branch`
- `title` / `body`：按下方规则
- `draft`：默认 false
- `maintainer_can_modify`：true

不要使用会把实现交给 Copilot 的创建 PR 工具。

### 7. 添加审查人

步骤 6 成功之后（新建或复用开放 PR 都算），把目标仓库最近活跃的 5 名可审查社区成员加为 reviewer。这一步失败不得回滚已经创建的 PR。

1. 读取当前审查请求和 PR 作者：

```bash
gh pr view "$PR_NUMBER" -R "$TARGET_REPO" --json author,reviewRequests
```

个人 reviewer 已有 5 名或以上时跳过。不要清掉已经请求过的人。不要请求 team。

2. 收集有 push 权限的协作者。GitHub 对没有写权限的用户会静默丢弃审查请求，所以先过滤，不要只按最近提交作者硬加：

```bash
gh api "repos/$TARGET_REPO/collaborators?per_page=100" --paginate \
  --jq '.[] | select(.permissions.push == true) | .login'
```

3. 从目标默认分支提交里按时间从新到旧收集`author.login`并去重：

```bash
gh api "repos/$TARGET_REPO/commits?sha=$DEFAULT_BRANCH&per_page=100&page=$PAGE"
```

排除：PR 作者；login 以`[bot]`结尾；`github-actions`、`web-flow`、`Copilot`；login 含`bot`（大小写不敏感）；已经在`reviewRequests`里的人。只保留第 2 步中有 push 权限的登录名。不够 5 人就翻页，直到凑齐或提交历史用尽。

4. 一次加上候选，然后必须再读`reviewRequests`核对实际生效的人：

```bash
gh pr edit "$PR_NUMBER" -R "$TARGET_REPO" --add-reviewer "a,b,c,d,e"
gh pr view "$PR_NUMBER" -R "$TARGET_REPO" --json reviewRequests
```

不够 5 人就继续从提交历史补人。协作者列表拿不到时，按提交顺序逐个`--add-reviewer`并核对，跳过未被接受的人。

`gh`不可用时，先查 GitHub MCP 工具 schema，用带`reviewers`参数的更新 PR 工具。不要调用 Copilot 代审。

凑不满 5 人也保持已创建的 PR，报告实际加上的登录名和缺口原因。

## PR 标题

每次都以`.github/PULL_REQUEST_TEMPLATE.MD`为准。标题始终用英文，不随后续正文语言选项变化。生成后自检，不合格就改到合格再提交。

格式：`<type>[optional scope]: <description>`

- `<type>`必须是模板列出的一种：`fix`、`feat`、`build`、`ci`、`docs`、`style`、`refactor`、`perf`、`test`、`chore`
- 类型、scope、冒号后的描述全部用英文，不要写成`修复(os/gtime): 修复时区问题`
- 有明确受影响包或范围时加 scope，例如`(os/gtime)`、`(net/ghttp)`、`(contrib/drivers/mysql)`
- 多个包被改、没有单一主范围时，可以省略 scope，不要堆一串 scope
- 冒号后是小写动词 + 短语，例如`fix`、`add`、`remove`，不要写成`Fix`
- 不要结尾句号
- 尽量短，最好不超过 76 个字符
- 从实际改动的文件和 diff 归纳，不要用`update`、`update code`、`fix bug`这种空描述

示例：

```text
fix(os/gtime): fix time zone issue
feat(net/ghttp): add request timeout option
docs: update contributing guide for agent setup
chore: remove unused sonar project properties
```

## PR 正文

应用模板后删除说明文字。发出去的正文只保留对这次改动有用的内容。

结构：

1. 用步骤 3 所选语言写清楚改了什么、为什么改。可以分短句或要点，不要贴完整 diff。标题仍保持英文。
2. 有对应 issue 时，单独一行写`Fixes #1234`（已完整修复）或`Updates #1234`（尚未完整修复）。这两个关键字保持英文，GitHub 靠它们关闭或关联 issue。
3. 没有对应 issue 就不要编造编号。
4. 不要把模板里的类型说明、参考链接或「删除这些说明」复制进正文。
5. 不要同时堆中英两套说明。正文语言必须和步骤 3 的选择一致。

英文示例：

```markdown
Remove the unused Sonar project properties file and refresh the root README badges.

Fixes #1234
```

中文示例：

```markdown
删除未使用的 Sonar 项目配置，并更新根目录 README 中的徽章。

Fixes #1234
```

## 推送与目标仓库

| 本地`origin` | `--head` | `-R` / owner/repo |
| --- | --- | --- |
| `gogf/gf` | 功能分支名 | `gogf/gf` |
| 自己的 fork | `fork-owner:branch` | `gogf/gf`（或 fork 的 parent） |

`origin`指向 Gitee 或其他非 GitHub 远程、但 GitHub 主仓库仍是`gogf/gf`时，把功能分支推到可开 GitHub PR 的 GitHub 远程（常见是`origin`或名为`github`的 remote），不要推到主分支。

## 输出约定

完成后向用户报告：

- 功能分支名，以及是否从主分支新建
- 用户选择的正文语言（标题始终为英文）
- 最终 PR 标题
- 目标仓库和目标分支
- PR URL 和编号
- 已请求的 reviewer 登录名；不足 5 人时说明原因
- 若没有开新 PR：说明是复用已有 PR，还是因无改动、权限不足、推送失败而停止

不要输出 token、密码或完整无用 diff。不要声称已经创建 PR，除非命令或 MCP 调用确实成功并返回了 URL。

