# 流水线说明

## 打Tag以及发布

1. 在要打Tag前，执行`make version to=[version]`命令将主库版本文件和子模块依赖预升级
2. 执行命令`git tag [version]`打新Tag
3. `git push --tags`推送新Tag到Github
4. 推送新Tag将会触发`tag.yml`中的流水线任务
   1. 给子模块打Tag
   2. 升级工具链库的依赖
   3. 将升级后的依赖变化推送pr到GitHub
   4. 给工具链库打Tag（使go install可用）
5. 人工合并上一步骤推送的pr
6. 合并pr后，将触发`cli.yml`流水线
   1. 构建gf工具链
   2. 发布release
   3. 将构建产出推送到release中