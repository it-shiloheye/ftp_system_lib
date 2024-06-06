# Commit Notes

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