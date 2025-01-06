#define _GNU_SOURCE

#include<signal.h>
#include<stdlib.h>
#include<stdio.h>
#include<unistd.h>
#include<linux/unistd.h>
#include<linux/sched.h>
#include<sched.h>
#include<sys/wait.h>
#include<sys/syscall.h>

#include "lib/helper.h"
#include "lib/process.h"
#include "lib/isoproc.h"
#include "lib/parse_proc_spec.h"


void start_process(char* process_yaml_loc, struct Process* p) {
    parse_process_yaml(process_yaml_loc, p);
    
    if ( chdir(p->ContextDir) != 0 ) {
        perror("error changing dir");
        exit(1);
    }

    int clone_flags = SIGCHLD | CLONE_NEWNS | CLONE_NEWUTS | CLONE_NEWUSER;
    char* cmd_stack = malloc(STACKSIZE);

    pid_t pid = clone(isoproc, cmd_stack + STACKSIZE, clone_flags, (void*)p);
    if (pid == -1){
        perror("clone");
        free(cmd_stack);
        exit(EXIT_FAILURE);
    }

    p->Pid = pid;
    p->Stack = cmd_stack;

    if( waitpid(pid, NULL, 0) == -1 ) {
        graceful_exit(p, "waitpid failed", 1);
    }

    graceful_exit(p, "success", 0);
}


int main() {

    struct Process* p = (struct Process*)calloc(1, sizeof(struct Process));
    
    printf("starting process\n");
    start_process("process.yaml", p);
    printf("starting process\n");

    free_process(p);

    return 0;
}