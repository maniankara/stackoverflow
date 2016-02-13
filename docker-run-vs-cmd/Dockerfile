FROM nginx:1.9.9

RUN apt-get update
RUN apt-get -y install php5-fpm php5-mysql php-apc php5-imagick\
        php5-imap php5-mcrypt php5-curl php5-cli php5-gd php5-pgsql\
        php5-common php-pear curl php5-json

ADD index.html /usr/share/nginx/html/index.html
ADD index.php /usr/share/nginx/html/index.php
ADD default /etc/nginx/sites-available/default

RUN /etc/init.d/php5-fpm restart
