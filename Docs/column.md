# Column Define Language
A example will be like:
```python
Column('group_name', String(64), unique=True, nullable=False)
Column('idx', Integer, Sequence('idx_seq'), nullable=False, index=True)
Column('parentId', Integer, ForeignKey('parents.id', ondelete='CASCADE', onupdate='CASCADE'), nullable=False)
```
## Column function
Column('name', type, ...params)
## Sequence function
Sequence('name', ...params)
## ForeignKey function
ForeignKey('fieldname',...params)
## default value
Can be a value: string, integer, fload, date, datetime, time, ...
Can be a function call: support be DBMS.