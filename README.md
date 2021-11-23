# tasks

tasks is a CLI TODO list manager

*based on John Calhoun's Gophercises*

## How it works?
Basically, I use a key and value db called bolt that stores the buckets (lists) in a file in $HOMEDIR/.tasks. Each file created there refers to a list of tasks created. And its contents are the added tasks.

## Run

```console
git clone https://github.com/lmagdanello/tasks.git
cd tasks/
go install .
tasks
```

## E.g.:

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

```shell
 ~  tasks init --bucket github
Bucket "github" created with success!

 ~  tasks list --bucket github
You have no tasks!

 ~  tasks add --bucket github "New Task"
Bucket: github
Added "New Task" to your task list!

 ~  tasks add --bucket github "Second Task"
Bucket: github
Added "Second Task" to your task list!

 ~  tasks list --bucket github             
Bucket: github
List tasks:
1. New Task
2. Second Task

 ~  tasks done --bucket github 1           
Bucket: github
Task "1" completed!

 ~  tasks list --bucket github  
Bucket: github
List tasks:
1. Second Task
```

Still working on the project, feel free to contribute or clone!

----