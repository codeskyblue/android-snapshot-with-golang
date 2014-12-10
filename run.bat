@echo off
adb push cgostudy /data/local/tmp/
adb shell chmod 755 /data/local/tmp/cgostudy
adb shell /data/local/tmp/cgostudy
adb pull /data/local/tmp/s.png
s.png
