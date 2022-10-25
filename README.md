# sqlconfctl
sqlconf editor

## Usage
```
git clone https://github.com/harryzhu/sqlconfctl

cd sqlconfctl

./build.sh
```

### set KEY=VAL into conf database
<code>
./sqlconfctl set --file="./conf.db" --name=appname --val=s3uploader
</code>

### delete KEY from conf database
<code>
./sqlconfctl delete --file="./conf.db" --name=appname
</code>

--file="./conf.db" can be skipped and "./conf.db" is the default file