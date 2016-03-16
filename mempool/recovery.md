## About the recovery system

### Index

#### Backup

1.  Index was sync to file with an IndexContent struct recording its address in database and its baseaddr
2.  All data of CSBT was sync to image file before

#### Recovery

1.  Read `index.json` to load IndexContent table
2.  For each index read related malloc information in `malloc.json`
3.  Load the raw image file to read all data of CSBT from the root and use a recursion
4.  Allocate a new block to storage all data of CSBT

### DataRow

#### Backup

1.  When any allocate occurs, including Index and DataRow, write malloc table with baseaddr and its file image to `malloc.json`
2.  Each DataRow has a special struct that defines how much size it used
3.  When insert DataRows, we allocate a large block and put them one by one, if there exist a index, tell index the address of the head of this row

#### Recovery

1.  When we need to recovery a table, we recovery its primary index tree and recovery every data
2.  We allocate a new block to storage them