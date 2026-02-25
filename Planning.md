# Planning app

Welcome to `cmdr` (Commander)!

Cmdr is a task runner for your custom cli workflows. Each workflow can store several steps that can be used to build out any command you use day-to-day. 
Beyond that, your workflows can be transferred from project to project, or passed to new employees being on-boarded. Instead of explaing each step
of the build or deploy process, a workflow can handle all situations with ease!

A workflow consists of Steps. 

Each Step runs in your terminal to display an input for text, a select, and a command. 

### Example

```JSON
[
    "name": "Example Workflow",
    "description": "This example can run a command to greet the user.",
    "steps": [
        {
            "type": "input",
            "helpText": "",
            "prompt": "Enter your name.",
            "variable": "user_name",
        },
        {
            "type": "select",
            "prompt": "Select the current time of day."
            "helpText": "This is some help text for choosing the correct time of day!",
            "options": "Morning, Afternoon, Evening",
            "variable": "time_of_day",
        },
        {
            "type": "command",
            "command": "echo \"Good {{time_of_day}}, {{user_name}}!\"",
            "capture_output": false,
        }
    ]
]
```

Or in yaml

```yaml
---
- name: Example Workflow
  description: This example can run a command to greet the user.
  steps:
  - type: input
    prompt: Enter your name.
    helpText: 
    variable: user_name
  - type: select
    prompt: Select the current time of day.
    helpText: This is some help text for choosing the correct time of day!
    options: Morning, Afternoon, Evening
    variable: time_of_day
  - type: command
    command: echo "Good {{time_of_day}}, {{user_name}}!"
    capture_output: false
```

### Definitions

- Workflow
- Step

#### Workflows
Workflows are individual Commands withing Cmdr (Commander)

Each workflow will contain a Name, Description, and a series of Steps. 

Each step can be one of Input, Select, or Command.

workflows:
    name
    description
    steps
        input
        select
        command

- Input
    - prompt 
    - helpText
    - variable

- Select
    - prompt
    - helpText
    - options
    - variable

- Commad
    - command


