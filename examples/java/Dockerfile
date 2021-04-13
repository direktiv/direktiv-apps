FROM openjdk:17-jdk-alpine3.13 as build

WORKDIR /app

COPY ./greeting/Greeting.java ./greeting/Greeting.java
COPY ./json-java.jar ./
COPY ./manifest.txt ./


RUN javac -classpath "./json-java.jar" ./greeting/Greeting.java
RUN jar cfm Greeting.jar ./manifest.txt ./greeting/*.class
RUN jar tf Greeting.jar


CMD ["java", "-cp", "./Greeting.jar:./json-java.jar", "greeting/Greeting"]
