**Please ensure you adhere to every item in this list.**
+ The PR title is formatted as follows: `<type>[optional scope]: <description>` For example, `fix(os/gtime): fix time zone issue`
  + `<type>` is mandatory and can be one of `fix`, `feat`, `build`, `ci`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
    + fix: Used when a bug has been fixed.
    + feat: Used when a new feature has been added.
    + build: Used for modifications to the project build system, such as changes to dependencies, external interfaces, or upgrading Node version.
    + ci: Used for modifications to continuous integration processes, such as changes to Travis, Jenkins workflow configurations.
    + docs: Used for modifications to documentation, such as changes to README files, API documentation, etc.
    + style: Used for changes to code style, such as adjustments to indentation, spaces, blank lines, etc.
    + refactor: Used for code refactoring, such as changes to code structure, variable names, function names, without altering functionality.
    + perf: Used for performance optimization, such as improving code performance, reducing memory usage, etc.
    + test: Used for modifications to test cases, such as adding, deleting, or modifying test cases for code.
    + chore: Used for modifications to non-business-related code, such as changes to build processes or tool configurations.
  + After `<type>`, specify the affected package name or scope in parentheses, for example, `(os/gtime)`.
  + The part after the colon uses the verb tense + phrase that completes the blank in
  + Lowercase verb after the colon
  + No trailing period
  + Keep the title as short as possible. ideally under 76 characters or shorter
  + [Reference Documentation](https://www.conventionalcommits.org/en/v1.0.0/)
+ If there is a corresponding issue, add either `Fixes #1234` or `Updates #1234`
  (the latter if this is not a complete fix) to this comment
+ Delete these instructions once you have read and applied them

**提交前请遵守每个事项，感谢！**
+ PR 标题格式如下：`<类型>[可选 范围]: <描述>` 例如 `fix(os/gtime): fix time zone issue`
  + `<类型>`是必须的，可以是 `fix`、`feat`、`build`、`ci`、`docs`、`style`、`refactor`、`perf`、`test`、`chore` 中的一个
    + fix: 用于修复了一个 bug
    + feat: 用于新增了一个功能
    + build: 用于修改项目构建系统，例如修改依赖库、外部接口或者升级 Node 版本等
    + ci: 用于修改持续集成流程，例如修改 Travis、Jenkins 等工作流配置
    + docs: 用于修改文档，例如修改 README 文件、API 文档等
    + style: 用于修改代码的样式，例如调整缩进、空格、空行等
    + refactor: 用于重构代码，例如修改代码结构、变量名、函数名等但不修改功能逻辑
    + perf: 用于优化性能，例如提升代码的性能、减少内存占用等
    + test: 用于修改测试用例，例如添加、删除、修改代码的测试用例等
    + chore: 用于对非业务性代码进行修改，例如修改构建流程或者工具配置等
  + `<类型>`后在括号中填写受影响的包名或范围，例如 `(os/gtime)`
  + 冒号后使用动词时态 + 短语
  + 冒号后的动词小写
  + 不要有结尾句号
  + 标题尽量保持简短，最好在 76 个字符或更短
  + [参考文档](https://www.conventionalcommits.org/zh-hans/v1.0.0/)
+ 如果有对应的 issue，请在此评论中添加 `Fixes #1234`，如果不是完全修复则添加 `Updates #1234`
+ 应用这些规则后删除所有的说明
