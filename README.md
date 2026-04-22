# My Blog (Hugo + PaperMod)

这套博客已经配置成“只管写文章”的结构，你平时只需要在 `content/posts/` 写 Markdown。

## 目录说明

- `content/posts/`: 文章目录（你主要用这个）
- `content/about.md`: 关于页
- `content/archives.md`: 归档页
- `content/search.md`: 搜索页
- `archetypes/posts.md`: 新文章模板

## 写作流程

1. 新建文章

```bash
hugo new posts/你的文章名.md
```

2. 写文章内容

编辑新文件，重点改这几项：

- `title`
- `tags`
- `categories`
- `draft` 改为 `false`

3. 本地预览

```bash
hugo server -D
```

浏览器打开 `http://localhost:1313/my-blog/`。

4. 发布上线

```bash
git add -A
git commit -m "发布文章: 标题"
git push origin master
```

推送后 GitHub Actions 会自动部署到 GitHub Pages。

## 更省事的命令

```bash
make new name=你的文章名
make dev
make publish
```

- `make new`: 新建文章到 `content/posts/`
- `make dev`: 本地预览
- `make publish`: 提交并推送（触发自动部署）

## 线上地址

- `https://stan0930.github.io/my-blog/`
