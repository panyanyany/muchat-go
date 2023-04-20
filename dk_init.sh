[ ! -f 配置/config.yml ] && cp 配置/demo.config.yml 配置/config.yml
[ ! -f 配置/seelog.xml ] && cp 配置/demo.seelog.app.xml 配置/seelog.xml
echo 'init done'