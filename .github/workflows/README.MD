# 流水线说明

## 打Tag以及发布

1. 在要打Tag前，执行`make version to=[version]`命令将主库版本文件和子模块依赖预升级
2. 提交变更并合并到主分支
   > 注意：`一定要逐一查看变更的文件，如果版本号异常则再执行一次命令`
3. 执行命令`git tag [version]`打新Tag
4. `git push --tags`推送新Tag到Github
5. 推送新Tag将触发`tag.yml`中的流水线任务
   1. 给子模块打Tag
   2. 升级工具链库的依赖
   3. 将升级后的依赖变化推送pr到GitHub
   4. 给工具链库打Tag（使go install可用）
6. 人工合并上一步骤推送的pr
7. 合并pr后，将触发`cli.yml`流水线
   1. 构建gf工具链
   2. 发布release
   3. 将构建产出推送到release中