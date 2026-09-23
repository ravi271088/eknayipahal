# PR Review Skill
This skill performs a comprehensive review of the current branch compared to the base branch (default: main).

## Review Process
1. **Diff Analysis**: Run `git diff main...HEAD` to identify all changed lines.
2. **Static Analysis**: Review the changes for:
   - **Correctness**: Logic errors, edge cases, or potential crashes.
   - **Efficiency**: Performance bottlenecks or redundant code.
   - **Style**: Adherence to Go idioms and HTML/CSS best practices.
   - **Security**: Potential vulnerabilities.
3. **Reporting**: For each finding, provide:
   - The file path and line number.
   - A clear explanation of the issue.
   - A suggested fix.

## Execution
The agent should be spawned with a 'high' effort level to ensure thoroughness.
