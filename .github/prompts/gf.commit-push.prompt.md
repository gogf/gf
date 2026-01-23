---
description: "自动整理当前分支修改内容，生成符合规范的commit message，并执行git add、commit和push操作"
agent: "agent"
tools:
  ['gitkraken/*']
---

# Git 提交和推送任务

## 任务目标

自动分析当前 Git 分支的修改内容，生成符合项目规范的 commit message，并完成代码的提交和推送。

## 执行步骤

### 1. 查看当前修改状态

首先执行 `git status` 命令查看当前工作区的修改状态，了解哪些文件被修改、新增或删除。

### 2. 分析修改内容

根据 `git diff` 和 `git status` 的输出，分析本次修改的主要内容：
- 识别修改的文件和模块
- 理解修改的类型（bug修复、新功能、重构、文档更新等）
- 识别影响的组件或包（如 `os/gtime`、`net/ghttp` 等）
- 总结修改的核心内容

### 3. 生成 Commit Message

根据修改内容，生成符合以下规范的 commit message：

#### Commit Message 格式规范

格式：`<type>[optional scope]: <description>`

例如：`fix(os/gtime): fix time zone issue`

#### Type 类型说明

- `fix`: 修复了一个bug，通常会对应有一个issue
- `feat`: 新增了一个功能，或者对现有组件执行了一些功能改进
- `build`: 修改项目构建系统，例如修改依赖库、外部接口或者升级 Node 版本等
- `ci`: 修改持续集成流程，例如修改 Travis、Jenkins 等工作流配置
- `docs`: 修改文档，例如修改 README 文件、API 文档等
- `style`: 修改代码的样式，例如调整缩进、空格、空行等
- `refactor`: 重构代码，例如修改代码结构、变量名、函数名等
- `perf`: 优化性能，例如提升代码的性能、减少内存占用等
- `test`: 修改测试用例，例如添加、删除、修改代码的测试用例等
- `chore`: 对非业务性代码进行修改，例如修改构建流程或者工具配置等

#### Scope 范围说明

- 在 `<type>` 后的括号中填写受影响的包名或范围
- 例如：`(os/gtime)`、`(net/ghttp)`、`(database/gdb)` 等
- 如果影响多个组件，选择主要影响的组件，或使用更通用的范围

#### Description 描述说明

- 冒号后使用动词时态 + 短语
- 冒号后的动词小写
- 不要有结尾句号
- 标题尽量保持简短，最好在 76 个字符或更短
- 使用英文描述

#### 完整示例

```text
fix(os/gtime): fix time zone issue
feat(net/ghttp): add middleware support for request validation
docs(README): update installation instructions
refactor(database/gdb): improve connection pool management
test(container/garray): add unit tests for sorted array
```

### 4. 执行 Git 操作

按照以下顺序执行 git 操作：

1. **git add -A**
   - 将所有修改的文件添加到暂存区

2. **git commit -m "commit message"**
   - 使用生成的 commit message 提交代码
   - 确保 commit message 符合上述规范

3. **git push**
   - 将本地提交推送到远程仓库
   - 如果是首次推送新分支，使用 `git push -u origin <branch-name>`

## 注意事项

1. **冲突处理**：如果push时遇到冲突，需要先执行 `git pull` 解决冲突后再推送

2. **分支检查**：确认当前在正确的分支上进行操作

3. **提交范围**：确保只提交相关的修改，避免混入无关的文件

4. **Commit Message 质量**：
   - 确保描述准确、简洁
   - 类型选择要正确
   - 范围要明确
   - 描述要有意义，避免模糊的描述如 "update code"

5. **相关 Issue**：
   - 如果有对应的 issue，在 commit message body 中添加 `Fixes #1234` 或 `Updates #1234`
   - 完全修复使用 `Fixes`，部分修复使用 `Updates`

## 输出要求

完成所有操作后，提供以下信息：
- 生成的commit message
- 执行的git命令及其输出
- 推送结果（成功/失败）
- 如有问题或建议，给出相应的提示

## 参考文档

- [项目PR模板](../../.github/PULL_REQUEST_TEMPLATE.MD)
