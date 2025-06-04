# 简单的翻译工具

配置文件
LLM openai/ollama
SYSTEM "xxx"

BASE ""
KEY ""
MODEL ""


示例 如果你使用 deepseek，可以如下配置

FROM openai
BASE https://api.deepseek.com
KEY sk-xxx

完毕

可选配置 SYSTEM 提示词 像这样：

SYSTEM """你是一个专业的翻译员, 精通中文与英文的翻译工作, 熟谙各种翻译技巧和方法, 能够准确地将中文翻译成英文, 并且能够理解和传达原文的意思和情感。
	你主要工作在 IT/计算机领域, 需要翻译各种技术文档、代码注释、用户手册等内容。请确保翻译的内容准确、流畅，并且符合目标语言的语法和用词习惯。
    你只能从中文翻译到英文，如用户输入的已经是英文，则直接输出 「无须翻译」。
	不需要输出翻译说明。只需要输出翻译后的内容。"""

所以本质上，这就是个简单的一次性生成工具，完全可以利用 LLM 的能力，移做他用

# TODO:
- ~~copy to clipboard~~
- loading status
- ~~migrate to openai sdk, langchaingo is too big.~~
- ~~multi conf profile~~
- ~~custom title~~
- ~~shortcut~~
- ~~system tray icon~~
- ~~always on the top~~
- screenshot and ocr
