## NilDB: An experimental RDBMS project to learn/practice Database internals and Distributed systems.

### This is my attempt to build a mini distributed RDBMS database system in Go(Golang); I heavily rely on Carnge Mellom University 2020 lecture videos that is available on Youtube [HERE](https://www.youtube.com/watch?v=vdPALZ-GCfI&list=PLSE8ODhjZXjbj8BMuIrRcacnQh20hmY9g). In addition to this, I will be using MIT Open course distributed systems that is also available on Youtube [HERE](https://www.youtube.com/mit)

## FEATURES 

### DATABASE
    |- STORAGE
        |-> IO OPERATOR(files, in-memory)
        |-> CACHE(LRU-K)
        |-> SCHEDULER
        |-> PAGER(BUFFER POOL MANAGER)
        |-> PAGES, TUPLES,...
    |- B+TREE
        |-> B+TREE IMPLEMENTATION
        |-> CONCURRENT B+TREE
        |-> INDEXES WITH B+TREE
    |- SQL PARSER
    |- QUERY ENGINE
    |- TRANSACTION MANAGER
    |- SERVER

### PARTITIONING
    |-> TABLE PARTITIONING(i.e multiple file per db)
    |-> OPTIMIZATIONS

### SHARDING
    |-> DATABASE CONTENTS SPREAD ACROSS MULTIPLE DATABASE
    |-> More info later...

### REPLICATION
    |-> DATA REPLICATION
    |-> RAFT
    |-> QUORUM
    |-> MASTER, SLAVE CONFIGURATIONS

### UI
    |-> BEAUTIFUL UI TO WRITTEN IN REACT

#### I am working on this project because I want to have a full graps of database internals from BTrees, Indexes, Distributed systems, Query Optimization...

For my learning, I use the CMU Database group youtube video and also MIT distributed systems course.

*** I will make significant updates to this readme file as time goes on***