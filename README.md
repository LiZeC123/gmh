GMH(GO-MY-HTTP)
=================

用GO语言实现所有需要的HTTP工具




启动Shell补全
-------------

本则能少输入几个字母, 就少输入几个字母的原则, 可以开启代码补全, 具体可查看`completion`指令的帮助

```bash
liuzheng@LUCASLLIU-MC1 Documents % gmh completion --help
NAME:
   gmh completion - Output shell completion script for bash, zsh, fish, or Powershell

USAGE:
   gmh completion [options]

DESCRIPTION:
   Output shell completion script for bash, zsh, fish, or Powershell.
   Source the output to enable completion.

   # .bashrc
   source <(gmh completion bash)

   # .zshrc
   source <(gmh completion zsh)

   # fish
   gmh completion fish > ~/.config/fish/completions/gmh.fish

   # Powershell
   Output the script to path/to/autocomplete/gmh.ps1 an run it.
```


对于`bash`和`zsh`来说, 只需要在对应的脚本里面加入一行指令即可.

> 代码里面自动生成的帮助信息比官方文档的教程要更新, 果然文档还是要和代码一起才能保证时效性



依赖项目文档
------------

- [urfave/cli](https://cli.urfave.org/v3/getting-started/)