#!/bin/sh

mkdir -p /etc/windmaker-alarmstatuswatcher

echo "### NOT starting on installation, please execute the following statements to configure windmaker-alarmstatuswatcher to start automatically using systemd"
echo " sudo /bin/systemctl daemon-reload"
echo " sudo /bin/systemctl enable windmaker-alarmstatuswatcher"
echo "### You can start windmaker-alarmstatuswatcher by executing"
echo " sudo /bin/systemctl start windmaker-alarmstatuswatcher"
