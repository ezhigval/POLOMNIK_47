FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive
RUN printf 'postfix postfix/main_mailer_type select Internet Site\npostfix postfix/mailname string tikhvin-palomnik.ru\n' | debconf-set-selections \
    && apt-get update \
    && apt-get install -y --no-install-recommends postfix postfix-pgsql libsasl2-modules ca-certificates ssl-cert \
    && sed -E -i 's/^([^[:space:]]+[[:space:]]+[^[:space:]]+[[:space:]]+[^[:space:]]+[[:space:]]+[^[:space:]]+[[:space:]]+)y([[:space:]]+)/\1n\2/' /etc/postfix/master.cf \
    && cp -a /var/spool/postfix /var/spool/postfix-default \
    && rm -rf /var/lib/apt/lists/*

COPY main.cf /etc/postfix/main.cf
COPY entrypoint.sh /usr/local/sbin/polomnik-mail-entrypoint
RUN chmod 0755 /usr/local/sbin/polomnik-mail-entrypoint \
    && install -d -o root -g root -m 0700 /etc/postfix/private

EXPOSE 25
ENTRYPOINT ["/usr/local/sbin/polomnik-mail-entrypoint"]
