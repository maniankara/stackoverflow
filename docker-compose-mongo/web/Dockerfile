FROM node
RUN apt-get update -y
RUN apt-get install -y git
RUN apt-get install -y vim
RUN git clone https://github.com/maniankara/quickbuild-docker.git /home/app
WORKDIR /home/app/src
RUN npm install
CMD "npm start"
EXPOSE 3000
