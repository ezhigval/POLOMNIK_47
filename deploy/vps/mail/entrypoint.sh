#!/bin/sh
set -eu

: "${POSTGRES_HOST:=postgres}"
: "${POSTGRES_PORT:=5432}"
: "${POSTGRES_DB:?POSTGRES_DB is required}"
: "${POSTGRES_USER:?POSTGRES_USER is required}"
: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}"
: "${SMTP_RELAY_USERNAME:?SMTP_RELAY_USERNAME is required}"
: "${SMTP_RELAY_PASSWORD:?SMTP_RELAY_PASSWORD is required}"

if [ ! -d /var/spool/postfix/maildrop ]; then
  cp -a /var/spool/postfix-default/. /var/spool/postfix/
fi
for queue_dir in active bounce corrupt defer deferred flush hold incoming private saved trace; do
  chown postfix:root "/var/spool/postfix/$queue_dir"
  chmod 0700 "/var/spool/postfix/$queue_dir"
done
chown postfix:postdrop /var/spool/postfix/maildrop /var/spool/postfix/public
chmod 0730 /var/spool/postfix/maildrop
chmod 2710 /var/spool/postfix/public

umask 077
cat > /etc/postfix/pgsql_virtual_aliases.cf <<EOF
hosts = ${POSTGRES_HOST}:${POSTGRES_PORT}
user = ${POSTGRES_USER}
password = ${POSTGRES_PASSWORD}
dbname = ${POSTGRES_DB}
query = SELECT replace(mail_forward_to, ',', ', ') FROM site_settings WHERE id = 1 AND '%s' = 'info@tikhvin-palomnik.ru' AND btrim(mail_forward_to) <> ''
EOF
printf '/.*/ info@tikhvin-palomnik.ru\n' > /etc/postfix/sender_canonical

printf '[smtp.mail.selcloud.ru]:1126 %s:%s\n' \
  "$SMTP_RELAY_USERNAME" "$SMTP_RELAY_PASSWORD" > /etc/postfix/private/sasl_passwd
postmap /etc/postfix/private/sasl_passwd
rm -f /etc/postfix/private/sasl_passwd
chown root:postfix /etc/postfix/pgsql_virtual_aliases.cf /etc/postfix/private/sasl_passwd.db
chmod 0640 /etc/postfix/pgsql_virtual_aliases.cf /etc/postfix/private/sasl_passwd.db

cert_dir="$(find /caddy_data/caddy/certificates -type f -name 'mail.tikhvin-palomnik.ru.crt' -print -quit 2>/dev/null | sed 's#/mail.tikhvin-palomnik.ru.crt$##')"
if [ -n "$cert_dir" ] && [ -f "$cert_dir/mail.tikhvin-palomnik.ru.key" ]; then
  cp "$cert_dir/mail.tikhvin-palomnik.ru.crt" /etc/postfix/tls.crt
  cp "$cert_dir/mail.tikhvin-palomnik.ru.key" /etc/postfix/tls.key
  chmod 0600 /etc/postfix/tls.crt /etc/postfix/tls.key
  postconf -e 'smtpd_tls_cert_file=/etc/postfix/tls.crt' 'smtpd_tls_key_file=/etc/postfix/tls.key'
fi

postfix check
exec /usr/sbin/postfix start-fg
