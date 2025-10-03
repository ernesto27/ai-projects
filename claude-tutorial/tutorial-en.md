# Claude Tutorial & Tips

## Table of Contents

- [Introduction & Definition](#introduction--definition)
- [Installation](#installation)
- [Init Claude Command](#init-claude-command)
- [Project Explanation & Analysis](#project-explanation--analysis)
- [Implement a Feature / Task](#implement-a-feature--task)
- [Reviewing Changes](#reviewing-changes)
- [Using @ for Files](#using--for-files)
- [Fix a BUG](#fix-a-bug)
- [Plan Mode](#plan-mode)
- [Thinking Mode](#thinking-mode)
- [Claude Code Tools](#claude-code-tools)
- [Adding Images to Context](#adding-images-to-context)
- [Using Commands](#using-commands)
- [Resume a Conversation](#resume-a-conversation)
- [Using Sub-Agents](#using-sub-agents)
- [Using MCP - Playwright](#using-mcp---playwright)
- [Output style](#output-style)

# Introduction & Definition

Claude Code is an AI agent that runs from the terminal and allows (among other things) analyzing projects, adding features, fixing bugs, etc. It was released in preview in February 2025 and became available to all users in May.

Over the past few months it has become one of the most popular tools among developers, which led other companies in recent months like Google and OpenAI to release similar tools:

- [Gemini CLI](https://blog.google/technology/developers/introducing-gemini-cli-open-source-ai-agent/)
- [OpenAI Codex CLI](https://help.openai.com/en/articles/11096431-openai-codex-cli-getting-started)

### Advantages of Claude Code

- **Runs in the terminal**: For people who work a lot in the terminal, this is a more natural environment than an IDE panel.

- **Separation of AI interface outside the editor**: One thing I don’t like about editors with AI (VS Code, Cursor) is that the UX gets cluttered with many panels and buttons and you can lose focus on the main editing area.

![Claude tutorial overview](https://github.com/ernesto27/ai-projects/blob/master/claude-tutorial/image01.png)

- **Editor-agnostic and independent**: Although if needed Claude Code can integrate with editors like VS Code or Cursor.

- **Tools**: Claude Code has an excellent tool implementation—for searching, editing files, running shell commands, web search, etc.

# Installation

### Requirements

At the moment you need an Anthropic subscription to use Claude Code:

- **Pro**: Lets you use the Claude Code 4 model, which is excellent for daily development tasks.
- **Max**: Lets you use the Claude Opus model, Anthropic’s most advanced model, better for more complex tasks.

Or you can use an API key—in that case you pay per use, which can increase quickly depending on the task type. Don’t forget that as an AI agent each question can trigger multiple Anthropic API calls.

**NodeJS**: 18 or higher

### Claude Code Installation

```bash
npm install -g @anthropic-ai/claude-code
```

The examples in this tutorial use the Pro model.

**Other requirements**:
- Docker

# Init Claude Command

The first thing to do is run the `claude` command inside the repository or project you want to work with.

Accept the permissions so Claude can modify project files.

The first command to run is:

```bash
/init 
```

This command initializes the agent: it analyzes the current project (structure, files, dependencies, etc.) and generates a configuration file called `CLAUDE.md` at the project root with information about the project. The utility of this file is that Claude sends these instructions every time a new conversation starts so it can maintain "memory" of project structure and practices.

### Types of CLAUDE.md Files

There are three types of `CLAUDE.md` files that can be generated:

1. **Per project**: checked into the repository.
2. **Local only**: with project-specific rules, not committed to the repo.
3. **Global**: shared across all projects where you use Claude Code.

In our case we’ll use option 1.

# Project Explanation & Analysis
One of the most effective uses of Claude Code is analyzing/explaining a project. This can be useful when:

- You’re onboarding to a long-lived project.
- You need to understand a feature.
- You’re reviewing a pull request.

We could use a prompt like:

```bash
Analyze the project: code, structure, dependencies, integrations, etc.
Generate a basic flow diagram and also one in Mermaid format. Save the content in a file named RESEARCH.md
```


Here we see Claude Code in action, generating a TODO list and using different tools to investigate the project.

![Claude tutorial overview](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/explain.png)

Generated result:

- [RESEARCH.md](RESEARCH.md)

# Implement a Feature / Task

When implementing a task, it’s best to be as specific as possible in the prompt. Otherwise Claude Code will infer/assume and the result might be far from what you expect. Like any tool, if used incorrectly it can generate more work than value.

Let’s say we have a backend API and we want to add a MySQL database connection.

### Example of a Bad Prompt

```bash
Add a MySQL database to the API.
```

While this might yield a functional result, Claude Code will make choices you might not want—like using a library you don’t prefer or overengineering the solution.

### A Better Prompt

```bash
I want to add a MySQL database connection to the existing API.

- Add a MySQL service in docker compose using version 8.
- Use the "goose" library for database migrations.
- Use the GORM library for the database connection.
- Create a migration that adds a table named "posts" with the fields: id, title, content
- Place all database-related files in a folder named "db" 
```

By detailing versions, libraries, and paths you get a result closer to expectations. Always take a few minutes to think through the prompt details.

# Reviewing Changes

As Claude Code generates changes, it lets you configure "auto-accept" either by selecting that option when it needs to edit a file or by by pressing "shift+tab" once.

![Claude tutorial overview](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/auto-accept.png)

While auto-accept can be useful in some circumstances (e.g., running the agent in the background), the recommended practice is to review changes line by line—ensuring correctness and iteratively instructing Claude on adjustments, similar to a code review but in real time.

Remember LLMs are non-deterministic—they can generate different results for the same prompt. More importantly, you are also responsible for bugs that ships to production!

# Using @ for Files

Although Claude Code can understand project context and locate the file to change, we can be more specific by telling Claude exactly which file we want to edit. This makes it faster and saves tokens otherwise spent using tools.

![File](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/file.png)

# Fix a BUG
Generally when we have a bug and we have the text of the error/message, copying it into Claude is enough.

![Fix Bug](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/debug.png)

It helps to mention what you already tried to fix it, plus add files, context, etc.

# Plan Mode
Claude has a plan mode, which can be enabled by pressing shift+tab twice.

![Plan](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/plan.png)

This mode is useful for medium/complex tasks because it lets Claude plan structured steps, iterating and validating each one.

For example:

```bash
Integrate our existing Golang API with AWS services: S3 for file uploads and SQS for async task processing
```

Example result:

![Plan result](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/plan2.png)

After that we can accept or refine the plan if is necessary.

# Thinking Mode

Sometimes Claude can’t find a solution due to the complexity of the task. For that, we can enable "thinking mode" by adding `, please think`.

Example prompt:

```bash
In Go, my app sometimes hangs, sometimes crashes with concurrent map read and map write, and sometimes gives wrong results under high concurrency. Why does this happen and how can I debug it?",  please think
```

Result:
![Thinking result](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/think1.png)

This can yield a better result, but keep in mind this mode uses more tokens and is therefore more expensive.

# Claude Code Tools

AI agents like Claude Code, Copilot, Cursor, etc., can use tools to perform specific tasks on the path to solving a problem.

In Claude Code, some of the tools currently available are:

### File & Project Tools

- **Read**: Read project files
- **Write**: Create new files or edit when needed
- **Edit/MultiEdit**: Modify existing files with precise search & replace operations

### Search Tools

- **Glob**: File pattern search (e.g. `**/*.ts`, `src/**/*.tsx`)
- **Grep**: Text search similar to the shell command
- **LS**: List directory contents

### System Tools

- **Bash**: Execute shell commands with proper quoting and background execution support

### Web Tools

- **WebFetch**: Retrieve and analyze web content
- **WebSearch**: Search the web for current information

# Adding Images to Context

To add images to the conversation, copy the image path (e.g., on Ubuntu right-click the image and select "copy") and then paste with `ctrl+shift+v` in the terminal.

This can help analyze diagrams, designs, errors, etc.

![Image from clipboard](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/image-command.png)

# Using Commands

At some point you may find yourself repeating the same prompt/task (project explanation, update, review, etc.). To avoid that we can create Markdown files in the `.claude/commands` folder.

Example: a command for a security review.

security-audit.md
```markdown
Perform a security audit on this code. Look for:
- SQL injection vulnerabilities
- XSS risks
- Authentication/authorization issues
- Sensitive data exposure
- Input validation problems
```

Another example: command to update the README.md with latest changes.

update-readme.md
```markdown
Update README.md with the recent changes:

- Analyze modified files and detect: new features, dependencies, configurations, structural changes
- Update affected sections: Features, Installation, Usage, Configuration, API docs
- Add or update the Changelog using headings: ### Added / Changed / Fixed / Removed
- Keep existing format and language; do not remove valid content
- List the proposed changes before applying them
```

These commands are used from Claude Code like:

```bash
/security-audit
```

```bash
/update-readme
```

# Resume a Conversation

If we need to resume a previous conversation—for instance, because a task wasn’t finished or we hit `ctrl+c` by mistake—we can use the `/resume` command.

```bash
/resume
```

It will show a list of recent conversations; select the one to resume and continue.

![Resume](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/resume.png)

# Using Sub-Agents
Claude lets us create sub-agents to perform various tasks. This is similar to commands but with important differences:

- The agent has its own context window (doesn’t share conversation context with Claude like commands do)
- It has its own system prompt
- Different tools can be enabled/disabled

To create a sub-agent use the command:

```bash
/agent
```
Create new agent

![Sub-agent](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/agent1.png)

We can generate an agent specifically for the current project or a global one.

Select the option:

1. Generate with Claude (recommended)

Next define the agent’s role.

![Sub-agent description](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/agent2.png)

Here we’ll create an agent to generate a code review.

Example:

```markdown
You are a senior Golang specialist performing a comprehensive code review. Analyze this Go code with focus on:

GOLANG BEST PRACTICES:
- Idiomatic Go patterns and conventions
- Effective use of goroutines and channels
- Proper error handling (not just if err != nil)
- Context usage and cancellation
- Defer statements placement and usage

PERFORMANCE & MEMORY:
- Memory allocations and potential leaks
- Goroutine leaks
- Efficient use of slices vs arrays
- String concatenation optimization
- Sync.Pool usage where appropriate
- Benchmark-worthy code sections

CONCURRENCY SAFETY:
- Race conditions
- Proper mutex usage (sync.Mutex vs sync.RWMutex)
- Channel deadlocks
- WaitGroup patterns
- Atomic operations where needed

CODE STRUCTURE:
- Interface design and composition
- Package organization and naming
- Exported vs unexported identifiers
- Embedded types usage
- Method receivers (pointer vs value)

TESTING & RELIABILITY:
- Table-driven test completeness
- Benchmark tests for critical paths
- Race detector compatibility
- Mock interfaces design
- Test coverage gaps

SECURITY:
- SQL injection in database/sql usage
- Command injection in os/exec
- Path traversal vulnerabilities
- Proper crypto package usage
- Sensitive data in logs

Provide specific line-by-line feedback with severity levels (Critical/Major/Minor) and code examples for improvements. Focus on Go-specific issues that linters might miss.
```

We can select which tools to enable—here we’ll leave defaults.

![Sub-agent tools](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/agent4.png)

If we inspect the project, Claude created a new file in `.claude/agents` named `golang-code-reviewer`.

We can invoke the agent like:

```bash
use golang-code-reviewer
```

## Using MCP - Playwright

In Claude Code we can configure different MCP services to extend the agent's capabilities. One of the best known is Playwright, which lets us interact with a web browser to perform various actions—like automated testing.

Repository:
https://github.com/microsoft/playwright-mcp

### Installation

```bash
claude mcp add playwright npx @playwright/mcp@latest
```

Once installed, we can verify that Playwright is available by running:

```bash
/mcp
```

![MCP](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/mcp.png)

### Usage Example

We can give a prompt like:

```bash
open this url using playwright http://localhost:8000/ and try to login using username "admin" and password "password" if there is an error take a screenshot
```

Here we see Playwright in action opening a browser:

![Playwright](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/playwright1.png)

If it cannot log in for some reason, it will take a screenshot and save it in the `playwright-mcp` folder:

![Playwright screenshot](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/playwright2.png)


# Output style

You can change the style in which Claude interacts with you.

Run this command:

```bash
/output-style
```

![Output style](https://raw.githubusercontent.com/ernesto27/ai-projects/master/claude-tutorial/output-style.png)

I highly recommend using the "Explanatory" style,  with this option Claude explain about why decided to do something in a certain way,  great for learn a new language and architecture.



