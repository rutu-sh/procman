#include <stdlib.h>

#include "process.h"
#include "helper.h"

void graceful_exit(struct Process* proc, char* msg, int exit_code) {
    free_process(proc);
    perror(msg);
    exit(exit_code);
}