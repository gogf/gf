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
6. PR 正文必须应用模板规则，但发出去的正文里不得再保留模板中的说明清单。正文是写给没参加这次对话的社区维护者看的审查材料，必须把背景、方案、改动和关联讨论写清楚，不能只留一两句摘要。写出的 Markdown 必须遵守`.agents/instructions/markdown-format.instructions.md`，发出去前读该文件并自检，不要凭记忆。
7. 标题、正文和本次新产生的 commit subject 都根据实际 diff、会话上下文和查到的关联讨论生成，不要只根据用户一句话臆造，也不要编造没出现过的 issue / PR / 评论。
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

- `.github/PULL_REQUEST_TEMPLATE.MD`（标题格式和`Fixes`/`Updates`的来源；详细正文结构以本技能「PR 正文」为准）
- `CONTRIBUTING.md`（PR 应对着`master`打开）
- `.agents/instructions/markdown-format.instructions.md`（PR 正文 Markdown 格式的唯一来源）

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

### 6. 收集审查背景并撰写正文

推送成功后、调用`gh pr create`之前，先为维护者写好完整正文。维护者通常没看过这次对话，看不到你为什么改、为什么选这个方案、对应哪条评论。缺这些信息时，审查只能对着 diff 猜意图。

按「PR 正文」收集材料并写成`$BODY_FILE`。材料不够就补查，不要用空泛的「修复问题」「更新代码」交差。写成后对照`.agents/instructions/markdown-format.instructions.md`改到合规，再进入下一步。

已有同一 head 的开放 PR、这次只是推送更新时：若现有正文缺少「PR 正文」要求的章节，用`gh pr edit --body-file`补全后再报告 URL，不要为此再开一条 PR。

### 7. 创建 Pull Request

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

### 8. 添加审查人

步骤 7 成功之后（新建或复用开放 PR 都算），把目标仓库最近活跃的 5 名可审查社区成员加为 reviewer。这一步失败不得回滚已经创建的 PR。

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

正文的读者是社区维护者，不是这次会话里的你和用户。写完后自检：一个没看过对话、只打开 Files changed 的人，能否明白痛点、方案、改了什么、以及该去哪条 issue / PR / 评论核对。不能的话就补材料，不要发出去。

应用 GitHub 模板后删除模板里的说明清单。标题仍保持英文。正文语言必须和步骤 3 的选择一致，不要中英混写两套说明。不要贴完整 diff。

### Markdown 格式

正文 Markdown 以`.agents/instructions/markdown-format.instructions.md`为准。写出后、发出去前读该文件并自检，不要凭记忆。该文件未覆盖的 GitHub 约定：

- `Fixes`/`Updates`、`#1234`、完整 URL 保持英文半角，以便 GitHub 识别。
- `@用户名`保持可通知的提及，不要包进反引号。
- 中文正文用该文件的中文全角标点；英文正文用英文标点。

### 撰写前收集

按下面来源补齐事实，再动笔。只写查到的内容。整节没有材料就省略该节，不要留空标题，也不要拿「无」凑字。

1. **这次对话**：用户原话、探索结论、被否掉的方案、审查意见、明确的 follow-up。
2. **实际改动**：`git diff` / `git log`，前后行为差异，受影响的包和公开 API，测试覆盖了哪条路径。
3. **本地设计材料**：若工作区有对应的`openspec/changes/<name>/proposal.md`或`design.md`，用来还原动机和取舍，但不要把 OpenSpec 文件名或工作流术语写进正文。
4. **关联讨论**：对话、分支名、commit body 里出现的 issue / PR 编号，用`gh issue view`、`gh pr view`、`gh api repos/<owner>/<repo>/issues/<n>/comments`或行内评论接口核对标题和关键意见。是 follow-up 时，把源 PR 里仍适用的评论要点写进「关联」，并附链接。
5. **没有线索就不要搜一圈硬凑。** 没有对应 issue 不要编造编号。

### 正文结构

用步骤 3 所选语言的标题，按这个顺序写。某节确实没有材料时，不要留空标题；「关联」里没有 issue 就不要写`Fixes`。

中文标题：`背景`、`痛点与场景`、`技术方案`、`改动内容`、`关联`、`验证`。  
英文标题：`Background`、`Problem`、`Approach`、`Changes`、`Related`、`Verification`。

1. **背景**  
   这段改动从哪来：既有行为、哪次讨论、哪条 PR 的后续、用户或社区遇到的什么情况。把时间线和动机讲清楚。

2. **痛点与场景**  
   不改会怎样。用可核对的现象写：错误结果、错误 SQL、缺字段的批量写入、调用方必须补额外参数等。有业务场景就写谁在什么操作下会踩到。

3. **技术方案**  
   选了什么、为什么选、和哪些备选做过取舍，以及明确不做的范围。维护者要能判断方案是否过大或过小。

4. **改动内容**  
   按包或驱动列出行为变化，写清改前 / 改后，而不是文件名清单。点出公开 API、错误码、测试补了哪类用例。不要贴 hunk。

5. **关联**  
   列出查到的 issue、PR、关键评论，带链接或`#编号`。完整修复用单独一行`Fixes #1234`，尚未完整修复用`Updates #1234`。这两个关键字保持英文。没有关联就整节省略，不要写「Fixes #0」。

6. **验证**  
   本地或 CI 实际跑过的命令、结果，以及没跑的部分（例如需要真实库的驱动测试留给 CI）。没跑的不要写成已经验证。

### 示例

中文：

```markdown
## 背景

#4766让 SQLite 的`Save`在未指定`OnConflict`时从主键推断冲突列，并且要求每一行都带齐主键字段。合并后评论里指出`pgsql`、`dm`、`oracle`等驱动仍是「复合主键里只要出现一个字段就过」，和 SQLite 不一致，需要另开`PR`对齐。

## 痛点与场景

调用方对复合主键表做`Save`、且未手写`OnConflict`时，这些驱动只要`list[0]`里有任意一个主键列，就会把整组主键写进`ON CONFLICT`或`MERGE ON`。后面的行如果缺列，SQL 仍然按完整唯一索引生成，更新对不上或直接执行失败。批量 upsert 尤其容易踩。

## 技术方案

把「每一行是否包含全部主键列」抽成`gdb.HasPrimaryKeys`，放在现有的`database/gdb/gdb_core_utility.go`（紧挨`GetPrimaryKeys`），各 upsert 驱动共用。`sqlite`、`sqlitecgo`也改走同一实现，避免两套校验再漂移。`gaussdb`的`InsertIgnore`在主键不完整时仍走原来的降级插入；本次不改 MySQL 系的`ON DUPLICATE KEY UPDATE`。

## 改动内容

- `gdb.HasPrimaryKeys`：空 list、空主键、大小写、复合主键缺列、后续行缺列都会判失败。
- `pgsql`、`gaussdb`、`dm`、`oracle`、`mssql`：推断`OnConflict`前改为检查全部主键列，错误文案改为要求 save data 里带齐主键值。
- `sqlite`、`sqlitecgo`：删除本地副本，改为调用`gdb.HasPrimaryKeys`。
- 各相关驱动补了复合主键`Save`成功、缺列失败、批量后续行缺列失败的用例。

## 关联

- Follow-up of https://github.com/gogf/gf/pull/4766
- #4766中@LanceAdd指出 sqlite 对复合主键更严，`pgsql`、`dm`、`oracle`仍是遇到其中一个字段即可：https://github.com/gogf/gf/pull/4766#issuecomment-4404490806

## 验证

- `go test ./database/gdb -run 'Test_HasPrimaryKeys|Test_mapHasKey'`通过
- `sqlite`、`sqlitecgo`的`Save`与复合主键用例通过
- `pgsql`、`gaussdb`、`dm`、`oracle`、`mssql`的同类用例已添加，需要对应数据库，留给 CI
```

英文用同一六段，标题换成`Background` / `Problem` / `Approach` / `Changes` / `Related` / `Verification`。关联关键字仍是`Fixes`或`Updates`。

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

