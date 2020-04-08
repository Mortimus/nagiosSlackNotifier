Add to misccommands.cfg

define command {
command_name                    host-notify-by-slack
command_line                    /usr/bin/slackNagiosBot "HOST" "COMPANY" "#testing" '$NOTIFICATIONTYPE$' '$SERVICEDESC$' '$HOSTALIAS$' '$HOSTADDRESS$' '$SERVICESTATE$' '$LONGDATETIME$' '$HOSTOUTPUT$' '$NOTIFICATIONCOMMENTS$'
}
​
define command {
command_name                    service-notify-by-slack
command_line                    /usr/bin/slackNagiosBot "SERVICE" "COMPANY" "#testing" '$NOTIFICATIONTYPE$' '$SERVICEDESC$' '$HOSTALIAS$' '$HOSTADDRESS$' '$SERVICESTATE$' '$LONGDATETIME$' '$SERVICEOUTPUT$'
}

then a contact is just needed for slack and notification commands are host-notify-by-slack for hosts and service-notify-by-slack for services, and binary added to /usr/bin/slackNagiosBot and chmod +x

move the config file to /opt/NagiosBot/config.json
make sure the log file in the config is writable by the nagios user