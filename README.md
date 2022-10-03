# confctl
sqlconf editor

## Usage
<code>

git clone https://github.com/harryzhu/confctl

cd confctl

./build.sh

</code>

### set KEY=VAL into conf database
<code>
./confctl set --file="./conf.db" --name=appname --val=s3uploader
</code>

### delete KEY from conf database
<code>
./confctl delete --file="./conf.db" --name=appname
</code>