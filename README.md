# tasks

```shell
~  tasks
tasks is a CLI TODO list manager

Usage:
  tasks [command]

Available Commands:
  add         Add a new task to your task list
  completion  generate the autocompletion script for the specified shell
  done        Marks a task as complete
  help        Help about any command
  list        List tasks

Flags:
  -h, --help   help for tasks

Use "tasks [command] --help" for more information about a command.
```
----

- [X] function chooseBucket should just choose the bucket and not init a new bucket (db file);
- [ ] rm command to remove the tasks and not complete;
- [ ] done cmd tuning (add status to struct (like COMPLETE and PENDING for each task)).
- [ ] personalize help;
- [ ] all functions show which bucket you are currently working on;
- [ ] improve logging errors;
- [ ] all functions missing args drops the help menu

:)