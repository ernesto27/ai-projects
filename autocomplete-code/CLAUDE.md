# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

Build and compile:
- `npm run compile` - Compile using webpack (development mode)
- `npm run watch` - Watch mode for development 
- `npm run package` - Production build with optimized bundle

Testing:
- `npm run test` - Run all tests (includes lint, compile, and test execution)
- `npm run compile-tests` - Compile TypeScript tests to `out/` directory
- `npm run watch-tests` - Watch mode for test compilation

Linting and code quality:
- `npm run lint` - Run ESLint on source files
- `npm run pretest` - Run all quality checks (compile-tests, compile, lint)

## Project Architecture

This is a VS Code extension for AI-powered code completion built with TypeScript and webpack.

### Key Components

- **Entry Point**: `src/extension.ts` - Contains `activate()` and `deactivate()` functions
- **Build Output**: `dist/extension.js` - Webpack bundles everything into this single file
- **Extension Manifest**: `package.json` - Defines commands, activation events, and VS Code API version (^1.102.0)

### Build System

- Uses webpack with ts-loader for TypeScript compilation
- Target: Node.js (VS Code extension host environment)  
- Source maps enabled for debugging
- Production builds use hidden source maps and minification

### Code Standards

- TypeScript with strict mode enabled
- ESLint with TypeScript plugin
- Naming conventions: camelCase/PascalCase for imports
- Requires semicolons, curly braces, and strict equality

### Extension Structure

Currently implements a basic "Hello World" command (`ai-autocomplete.helloWorld`) as a starting point. The extension activates automatically and registers the command through VS Code's command palette.

## Testing Setup

Uses `@vscode/test-cli` and `@vscode/test-electron` for VS Code extension testing. Tests are compiled separately to the `out/` directory.