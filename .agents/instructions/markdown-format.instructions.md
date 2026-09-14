---
name: "Standardize markdown document formatting"
description: "Standardize the formatting of all markdown documents to keep structure clear, content readable, and the overall quality and user experience consistent. This document explains requirements for heading levels, paragraph formatting, code block usage, list formatting, and image and link insertion so authors can follow a unified style that is easier to read and maintain."
applyTo: "*.{md,MD}"
---

# Primary Formatting Requirements

- Keywords or specialized terms in the document must be formatted with inline code, for example `RuntimeClass`, `containerd`, `GPU`, and `AI`.
- In Chinese text, do not add spaces around inline code.
- For technical articles, review the generated content before finalizing it to ensure the material is technically accurate and contains no incorrect technical descriptions.
- When the generated content is too large, split it into multiple tasks to avoid exceeding model output limits and causing the workflow to fail.


# Detailed Content Requirements

- When documenting parameters or configuration items for a component or project, prefer tables when practical, and keep tables short enough to avoid horizontal scrolling during normal reading.
- In Chinese paragraphs, use full-width punctuation rather than half-width punctuation, for example `，` instead of `,` and `；` instead of `;`.
- Use `mermaid` for architecture diagrams, flowcharts, and similar visuals. If you need line breaks inside `mermaid`, use `<br/>` instead of `\n`.
- If a code block is not a `mermaid` diagram and instead uses box-drawing characters such as `┌─`, `┐`, `┤`, or `│`, keep the content in English so the layout stays aligned.
- Do not use `---` as a separator between paragraphs.
