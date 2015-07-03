#!/bin/bash
#
#   /etc/rc.d/init.d/githookproxy
#
#  githookproxy service
#   
### BEGIN INIT INFO
# Provides:          githookproxy
# Required-Start:    $network $remote_fs $named
# Required-Stop:     $network $remote_fs $named
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Starts Git Hook Proxy
# Description:       Starts githookproxy using daemon function
### END INIT INFO

# Only configuration needed here
NAME="githookproxy"
DESC="Git Hook Proxy"
PORT=9999
IFACE=0.0.0.0
USER="jenkins"
LOGFILE="/var/log/$NAME.log"
PID_FILE="/var/run/$NAME.pid"
DAEMON="/usr/local/share/githookproxy/githookproxy"
DAEMON_OPTS="-i $IFACE:$PORT -l $LOGFILE -p $PID_FILE"

PATH=/bin:/usr/bin:/sbin:/usr/sbin

# Source function library.
. /etc/rc.d/init.d/functions

# Variables USER, GROUP, PID_FILE, DAEMON and DAEMON_OPTS are stored in $VARIABLES
start() {
    echo -n "Starting $NAME: "
    touch $PID_FILE
    chown $USER:$USER $PID_FILE
    daemon --user $USER --pidfile $PID_FILE "$DAEMON $DAEMON_OPTS"
    echo
    touch /var/lock/subsys/$NAME
}   

# Variable PID_FILE is stored in $VARIABLES
stop() {
    echo -n "Shutting down $NAME: "
    killproc -p $PID_FILE $NAME
    local exit_status=$?
    echo
    rm -f /var/lock/subsys/$NAME
    return $exit_status
}

restart(){
    stop
    sleep 1
    start
}

usage(){
    echo "Usage: $NAME {start|stop|status|restart|reload|force-reload|condrestart}"
    exit 1
}

case "$1" in
    start)                          start ;;
    stop)                           stop ;;
    status)                         status -p $PID_FILE $NAME ;;
    restart|reload|force-reload)    restart ;;
    *)                              usage ;;
esac
exit $?
