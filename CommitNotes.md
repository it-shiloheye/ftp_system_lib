# Commit Notes

### 17th June 2024, 19:33 PM GMT +3
```sh
1. Improving logging:
    - moving log files into logging struct 
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   logging/logging_struct.go
#

```

### 14th June 2024, 01:36 AM GMT +3
```sh
1. Moving logging into ftp_system_lib 
    for uniformity
2. Want 4 layer logging:
    - day_log
    - err_log
    - latest_log
    - history_log (max file size count)
3. Complex logger structure not yet ready
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	new file:   logging/error_type.go
#	new file:   logging/fake_logger.go
#	new file:   logging/logging_struct.go
#

```

### 08th June 2024, 09:50 AM GMT+3
    1. Fixing bug in lockfile:
        - was returning a nil pointer at .Lock
        - was calling err1.Error() at .Unlock without checking for nil


### 08th June 2024, 09:50 AM GMT+3
    1. Adding MetaData Map to FileHash object

### 08th June 2024, 02:17 PM GMT+3
    1. fixing bugs in modtime
    2. adding lockfile item to lock directory/file

### 06th June 2024, 19:11 PM GMT+3
    1. fixing json error: 'json: cannot unmarshal object into Go struct field Fs of type fs.FileInfo\'

### 06th June 2024, 13:57 PM GMT+3
    
    1. Fixing Run Time Bug (panic on Calling FileBasic.Direntry)
    2. Able to create and write a file

### 06th June 2024, 11:39 AM GMT+3
    
    1. Deleting old FileBasic, old FileTree, old FileHash
    2. New helper functions for operating FileBasic and FileHash

### 06th June 2024, 11:28 AM GMT+3

    1. Simplifying  file basic type
    2. Simplifying byte store type
    3. Want to expose the underlying file
    4. Improving the error message