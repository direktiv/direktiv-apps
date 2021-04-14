FROM openjdk:17-jdk-alpine3.13 as build

WORKDIR /app

COPY ./greeting/Greeting.java ./greeting/Greeting.java
COPY ./manifest.txt ./

# Build json-java.jar
RUN wget https://github.com/stleary/JSON-java/archive/refs/tags/20210307.tar.gz
RUN tar -xvf 20210307.tar.gz
RUN cd ./JSON-java-20210307/src/main/java && javac org/json/*.java
RUN cd ./JSON-java-20210307/src/main/java && jar cf json-java.jar org/json/*.class
RUN cp ./JSON-java-20210307/src/main/java/json-java.jar /

# Compile Greeting.jar
RUN javac -classpath "/json-java.jar" ./greeting/Greeting.java
RUN jar cfm Greeting.jar ./manifest.txt ./greeting/*.class

CMD ["java", "-cp", "./Greeting.jar:./json-java.jar", "greeting/Greeting"]
