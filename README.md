# Stock Analysis Tool in Golang

## Overview

This repository provides the architecture and design of a **Stock Analysis Tool** implemented entirely in **Golang**. The tool fetches real-time stock data, performs technical analysis, applies machine learning insights, sends alerts, and uses job scheduling. It employs technologies like **Gin**, **Goroutines**, **Kafka**, **PostgreSQL**, **Redis**, and **InfluxDB** to achieve this.

## Features

- **Real-Time Stock Data**: Fetches live stock data from external APIs.
- **Technical Analysis**: Implements common technical analysis algorithms (MA, RSI, MACD).
- **Machine Learning Insights**: Applies machine learning models for stock predictions.
- **Job Scheduling**: Schedules recurring tasks using **Go-Cron**.
- **Notifications**: Sends alerts via **email (SMTP)**, **SMS (Twilio)**, and **Push Notifications (Firebase)**.
- **Database**: Utilizes **PostgreSQL** for relational data, **InfluxDB** for time-series data, and **Redis** for caching.

---

## Architecture

### 1. **User Layer**
   - **User Device (Web/Mobile App)**: Users interact with the system through a web or mobile app.
  
### 2. **API Gateway (Gin)**
   - **Gin Framework** handles routing and user requests.
   - Includes **JWT** authentication, **Rate Limiting**, and **Error Handling**.

### 3. **Backend Service Layer**
   - **Stock Data Fetcher**: Fetches stock data from external APIs like **Yahoo Finance** or **Alpha Vantage**.
   - **Technical Analysis Module**: Implements technical analysis indicators such as **Moving Averages (MA)**, **Relative Strength Index (RSI)**, and **MACD**.
   - **Machine Learning Insights**: Applies machine learning algorithms for stock predictions.

### 4. **Data Processing Layer**
   - **Kafka** for real-time data streaming.
   - **Golang Cron Jobs** for task scheduling (e.g., fetching stock data periodically).

### 5. **Storage Layer**
   - **InfluxDB** for time-series data storage.
   - **PostgreSQL** for relational data storage.
   - **Redis** for caching frequently accessed data.

### 6. **Notification & Alerts**
   - Sends notifications via **Twilio** (SMS), **SendGrid** (Email), and **Firebase Cloud Messaging** (Push Notifications).

### 7. **Authentication & Authorization**
   - **JWT/OAuth2** for secure user authentication and authorization.

---

## Technologies Used

- **Backend Framework**: **Gin** (for API Gateway)
- **Real-Time Streaming**: **Kafka** (using **Sarama** Golang client)
- **Data Storage**: **InfluxDB** (Time-series DB), **PostgreSQL** (Relational DB), **Redis** (Cache)
- **Job Scheduling**: **Go-Cron** (for scheduling tasks), **Redis** (for task queueing)
- **Authentication**: **JWT**, **OAuth2**
- **Notification Services**: **Twilio** (SMS), **SendGrid** (Email), **Firebase Cloud Messaging** (Push)

---

## System Flow

1. **User Request**: The user interacts with the app (web/mobile) for stock data, analysis, or alerts.
2. **API Gateway**: Requests are routed through the **Gin API Gateway**, which handles authentication, rate limiting, and forwards the request to the backend.
3. **Stock Data Fetcher**: Fetches real-time stock data from external APIs and stores it in the database.
4. **Technical Analysis**: Performs technical analysis (such as MA, RSI, MACD) on the stock data and stores results in the database.
5. **Job Scheduler**: Periodic tasks like fetching stock data and performing analysis are scheduled using **Go-Cron**.
6. **Notifications**: Alerts are sent to users via **email**, **SMS**, or **push notifications** based on user preferences.

---


---
POSTGRES:
https://neon.tech/
stockanalyser@yopmail.com
Controlstock@123


---
REDIS:
https://console.upstash.com/
stockanalyser@yopmail.com
Controlstock@123

MAILBOX:
https://yopmail.com/en/wm



TABLE COMMAND COMMANDS:
CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

## Architecture Diagram

<img width="858" alt="image" src="https://github.com/user-attachments/assets/8b3a45d0-611d-44c1-8ed1-5ae139c4089a" />



